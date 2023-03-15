package middleware

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/codes"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"github.com/vatsal-chaturvedi/article-management-sys/pkg/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func test(hit *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*hit = true
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&model.Response{
			Status:  http.StatusOK,
			Message: "passed",
			Data:    nil,
		})
	}
}

func TestMiddleware_Cacher(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name      string
		config    config.Config
		setupFunc func() (*http.Request, *mock.MockCacherI)
		validator func(*httptest.ResponseRecorder, *bool)
	}{
		{
			name:   "SUCCESS::Cacher::Cached Response",
			config: config.Config{},
			setupFunc: func() (*http.Request, *mock.MockCacherI) {
				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				mockCacher := mock.NewMockCacherI(mockCtrl)
				cacheResponse := model.CacheResponse{Status: http.StatusOK, Response: "ok", ContentType: "application/json"}
				b, _ := json.Marshal(cacheResponse)
				mockCacher.EXPECT().Get("http://localhost:80").Return(b, nil)
				return req, mockCacher
			},
			validator: func(res *httptest.ResponseRecorder, hit *bool) {
				if *hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
					return
				}
				by, _ := ioutil.ReadAll(res.Body)
				if !reflect.DeepEqual([]byte("ok"), by) {
					t.Errorf("Want: %v, Got: %v", "ok", string(by))
				}
				if !reflect.DeepEqual(http.StatusOK, res.Code) {
					t.Errorf("Want: %v, Got: %v", http.StatusOK, res.Code)
				}
				if !reflect.DeepEqual("application/json", res.Header().Get("Content-Type")) {
					t.Errorf("Want: %v, Got: %v", "application/json", res.Header().Get("Content-Type"))
				}
			},
		},
		{
			name:   "Failure::Cacher::Cached Response::unmarshal error",
			config: config.Config{},
			setupFunc: func() (*http.Request, *mock.MockCacherI) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				mockCacher := mock.NewMockCacherI(mockCtrl)
				mockCacher.EXPECT().Get("http://localhost:80").Return([]byte("123"), nil)
				return req, mockCacher
			},
			validator: func(res *httptest.ResponseRecorder, hit *bool) {
				if *hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
					return
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					return
				}
				expected := &model.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrUnmarshall),
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name:   "SUCCESS::Cacher::Normal Response",
			config: config.Config{Cacher: config.CacheConfig{KeyExpiryDuration: time.Minute}},
			setupFunc: func() (*http.Request, *mock.MockCacherI) {
				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				mockCacher := mock.NewMockCacherI(mockCtrl)
				mockCacher.EXPECT().Get("http://localhost:80").Return(nil, errors.New("error"))
				mockCacher.EXPECT().Set("http://localhost:80", []byte("{\"Status\":200,\"Response\":\"{\\\"status\\\":200,\\\"message\\\":\\\"passed\\\",\\\"data\\\":null}\\n\",\"ContentType\":\"application/json\"}"), time.Minute)
				return req, mockCacher
			},
			validator: func(res *httptest.ResponseRecorder, hit *bool) {
				if *hit != true {
					t.Errorf("Want: %v, Got: %v", true, *hit)
					return
				}
				var resp model.Response
				by, _ := ioutil.ReadAll(res.Body)
				json.Unmarshal(by, &resp)
				expectedResp := model.Response{Status: 200, Message: "passed", Data: nil}
				if !reflect.DeepEqual(expectedResp, resp) {
					t.Errorf("Want: %v, Got: %v", expectedResp, resp)
				}
				if !reflect.DeepEqual("application/json", res.Header().Get("Content-Type")) {
					t.Errorf("Want: %v, Got: %v", "application/json", res.Header().Get("Content-Type"))
				}
			},
		},
		{
			name:   "Failure::Cacher::Normal Response::Redis fail",
			config: config.Config{Cacher: config.CacheConfig{KeyExpiryDuration: time.Minute}},
			setupFunc: func() (*http.Request, *mock.MockCacherI) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				mockCacher := mock.NewMockCacherI(mockCtrl)
				mockCacher.EXPECT().Get("http://localhost:80").Return(nil, errors.New("error"))
				mockCacher.EXPECT().Set("http://localhost:80", []byte("{\"Status\":200,\"Response\":\"{\\\"status\\\":200,\\\"message\\\":\\\"passed\\\",\\\"data\\\":null}\\n\",\"ContentType\":\"application/json\"}"), time.Minute).Return(errors.New("error"))
				return req, mockCacher
			},
			validator: func(res *httptest.ResponseRecorder, hit *bool) {
				if *hit != true {
					t.Errorf("Want: %v, Got: %v", true, *hit)
					return
				}
				var resp model.Response
				by, _ := ioutil.ReadAll(res.Body)
				json.Unmarshal(by, &resp)
				expectedResp := model.Response{Status: 200, Message: "passed", Data: nil}
				if !reflect.DeepEqual(expectedResp, resp) {
					t.Errorf("Want: %v, Got: %v", expectedResp, string(by))
				}
			},
		},
	}

	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// STEP 1: seting up all instances for the specific test case
			res := httptest.NewRecorder()
			req, cacher := tt.setupFunc()
			cfg := config.SvcConfig{Cfg: &tt.config}
			middleware := Middleware{
				cacher: cacher,
				cfg:    cfg.Cfg}
			var hit bool
			testFunc := test(&hit)
			x := middleware.Cacher(testFunc)
			x.ServeHTTP(res, req)

			tt.validator(res, &hit)

		})
	}
}
