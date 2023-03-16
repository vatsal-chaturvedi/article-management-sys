package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/codes"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"github.com/vatsal-chaturvedi/article-management-sys/pkg/mock"
	"io/ioutil"

	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Reader string

func (Reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
func Test_ArticleManagement_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() ArticleManagementHandlerI
		validate func(http.ResponseWriter)
	}{
		{
			name: "Success",
			setup: func() ArticleManagementHandlerI {
				return &articleManagement{}
			},
			validate: func(w http.ResponseWriter) {
				wIn := w.(*httptest.ResponseRecorder)

				if !reflect.DeepEqual(wIn.Code, http.StatusMethodNotAllowed) {
					t.Errorf("Want: %v, Got: %v", http.StatusMethodNotAllowed, wIn.Code)
				}
				if !reflect.DeepEqual(wIn.Header().Get("Content-Type"), "application/json") {
					t.Errorf("Want: %v, Got: %v", "application/json", wIn.Header().Get("Content-Type"))
				}

				resp := model.Response{}
				err := json.NewDecoder(wIn.Body).Decode(&resp)
				if !reflect.DeepEqual(err, nil) {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}

				if !reflect.DeepEqual(resp, model.Response{
					Status:  http.StatusMethodNotAllowed,
					Message: http.StatusText(http.StatusMethodNotAllowed),
					Data:    nil,
				}) {
					t.Errorf("Want: %v, Got: %v", model.Response{
						Status:  http.StatusMethodNotAllowed,
						Message: http.StatusText(http.StatusMethodNotAllowed),
						Data:    nil,
					}, resp)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.setup()
			w := httptest.NewRecorder()
			handler.MethodNotAllowed(w, nil)
			tt.validate(w)
		})
	}
}

func Test_ArticleManagement_RouteNotFound(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() ArticleManagementHandlerI
		validate func(http.ResponseWriter)
	}{
		{
			name: "Success",
			setup: func() ArticleManagementHandlerI {
				return &articleManagement{}
			},
			validate: func(w http.ResponseWriter) {
				wIn := w.(*httptest.ResponseRecorder)
				if !reflect.DeepEqual(wIn.Code, http.StatusNotFound) {
					t.Errorf("Want: %v, Got: %v", http.StatusNotFound, wIn.Code)
				}

				if !reflect.DeepEqual(wIn.Header().Get("Content-Type"), "application/json") {
					t.Errorf("Want: %v, Got: %v", "application/json", wIn.Header().Get("Content-Type"))
				}
				resp := model.Response{}
				err := json.NewDecoder(wIn.Body).Decode(&resp)
				if !reflect.DeepEqual(err, nil) {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
				if !reflect.DeepEqual(resp, model.Response{
					Status:  http.StatusNotFound,
					Message: http.StatusText(http.StatusNotFound),
					Data:    nil,
				}) {
					t.Errorf("Want: %v, Got: %v", model.Response{
						Status:  http.StatusNotFound,
						Message: http.StatusText(http.StatusNotFound),
						Data:    nil,
					}, resp)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.setup()
			w := httptest.NewRecorder()
			handler.RouteNotFound(w, nil)
			tt.validate(w)
		})
	}
}

func Test_ArticleManagement_InsertArticle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() (ArticleManagementHandlerI, *http.Request)
		want  func(httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() (ArticleManagementHandlerI, *http.Request) {
				mockLogic := mock.NewMockArticleManagementLogicI(mockCtrl)
				article := model.Article{
					Title:   "title",
					Author:  "author",
					Content: "content",
				}
				mockLogic.EXPECT().InsertArticle(&article).
					Return(&model.Response{
						Status:  http.StatusCreated,
						Message: "Success",
						Data:    map[string]string{"id": "1"},
					}).Times(1)

				rec := &articleManagement{
					logic: mockLogic,
				}
				by, err := json.Marshal(article)
				if err != nil {
					return nil, nil
				}
				r, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(by))
				return rec, r
			},
			want: func(recorder httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(recorder.Body)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				var response model.Response
				err = json.Unmarshal(b, &response)
				tempResp := &model.Response{
					Status:  http.StatusCreated,
					Message: "Success",
					Data:    map[string]interface{}{"id": "1"},
				}
				if !reflect.DeepEqual(recorder.Code, http.StatusCreated) {
					t.Errorf("Want: %v, Got: %v", http.StatusCreated, recorder.Code)
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure:: Validate error",
			setup: func() (ArticleManagementHandlerI, *http.Request) {
				mockLogic := mock.NewMockArticleManagementLogicI(mockCtrl)
				article := model.Article{
					Author:  "author",
					Content: "content",
				}

				rec := &articleManagement{
					logic: mockLogic,
				}
				by, err := json.Marshal(article)
				if err != nil {
					return nil, nil
				}
				r, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(by))
				return rec, r
			},
			want: func(recorder httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(recorder.Body)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				var response model.Response
				err = json.Unmarshal(b, &response)
				tempResp := &model.Response{
					Status:  http.StatusBadRequest,
					Message: "Key: 'Article.Title' Error:Field validation for 'Title' failed on the 'required' tag",
					Data:    nil,
				}
				if !reflect.DeepEqual(recorder.Code, http.StatusBadRequest) {
					t.Errorf("Want: %v, Got: %v", http.StatusBadRequest, recorder.Code)
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure::readAll error",
			setup: func() (ArticleManagementHandlerI, *http.Request) {
				mockLogic := mock.NewMockArticleManagementLogicI(mockCtrl)
				rec := &articleManagement{
					logic: mockLogic,
				}
				r, _ := http.NewRequest("POST", "/articles", Reader(""))
				return rec, r
			},
			want: func(recorder httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(recorder.Body)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				var response model.Response
				err = json.Unmarshal(b, &response)
				tempResp := &model.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrReadingReqBody),
					Data:    nil,
				}
				if !reflect.DeepEqual(recorder.Code, http.StatusBadRequest) {
					t.Errorf("Want: %v, Got: %v", http.StatusBadRequest, recorder.Code)
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure::json unmarshall error",
			setup: func() (ArticleManagementHandlerI, *http.Request) {
				mockLogic := mock.NewMockArticleManagementLogicI(mockCtrl)
				rec := &articleManagement{
					logic: mockLogic,
				}
				r, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer([]byte("")))
				return rec, r
			},
			want: func(recorder httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(recorder.Body)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				var response model.Response
				err = json.Unmarshal(b, &response)
				tempResp := &model.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrUnmarshall),
					Data:    nil,
				}
				if !reflect.DeepEqual(recorder.Code, http.StatusBadRequest) {
					t.Errorf("Want: %v, Got: %v", http.StatusBadRequest, recorder.Code)
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.InsertArticle(w, r)
			tt.want(*w)
		})
	}
}

func Test_ArticleManagement_GetArticleById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() (ArticleManagementHandlerI, *http.Request)
		want  func(httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() (ArticleManagementHandlerI, *http.Request) {
				mockLogic := mock.NewMockArticleManagementLogicI(mockCtrl)
				mockLogic.EXPECT().GetArticle("1").
					Return(&model.Response{
						Status:  http.StatusOK,
						Message: "Success",
						Data:    map[string]string{"id": "1"},
					}).Times(1)

				rec := &articleManagement{
					logic: mockLogic,
				}
				r, _ := http.NewRequest("GET", "/articles/:1", nil)
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				return rec, r
			},
			want: func(recorder httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(recorder.Body)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				var response model.Response
				err = json.Unmarshal(b, &response)
				tempResp := &model.Response{
					Status:  http.StatusOK,
					Message: "Success",
					Data:    map[string]interface{}{"id": "1"},
				}
				if !reflect.DeepEqual(recorder.Code, http.StatusOK) {
					t.Errorf("Want: %v, Got: %v", http.StatusOK, recorder.Code)
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure::invalid id",
			setup: func() (ArticleManagementHandlerI, *http.Request) {
				mockLogic := mock.NewMockArticleManagementLogicI(mockCtrl)
				rec := &articleManagement{
					logic: mockLogic,
				}
				r, _ := http.NewRequest("GET", "/articles/:1", nil)
				return rec, r
			},
			want: func(recorder httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(recorder.Body)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				var response model.Response
				err = json.Unmarshal(b, &response)
				tempResp := &model.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrAssertid),
					Data:    nil,
				}
				if !reflect.DeepEqual(recorder.Code, http.StatusBadRequest) {
					t.Errorf("Want: %v, Got: %v", http.StatusBadRequest, recorder.Code)
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.GetArticleById(w, r)
			tt.want(*w)
		})
	}
}

func Test_ArticleManagement_GetAllArticle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() (ArticleManagementHandlerI, *http.Request)
		want  func(httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() (ArticleManagementHandlerI, *http.Request) {
				mockLogic := mock.NewMockArticleManagementLogicI(mockCtrl)
				mockLogic.EXPECT().GetAllArticle(20, 1).
					Return(&model.Response{
						Status:  http.StatusOK,
						Message: "Success",
						Data:    []model.ArticleDs{{Id: "1", Title: "title", Author: "author", Content: "content"}},
					}).Times(1)

				rec := &articleManagement{
					logic: mockLogic,
				}
				r, _ := http.NewRequest("GET", "/articles", nil)
				return rec, r
			},
			want: func(recorder httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(recorder.Body)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				var response model.Response
				err = json.Unmarshal(b, &response)
				tempResp := &model.Response{
					Status:  http.StatusOK,
					Message: "Success",
					Data:    []map[string]interface{}{{"author": "author", "content": "content", "id": "1", "title": "title"}},
				}
				if !reflect.DeepEqual(recorder.Code, http.StatusOK) {
					t.Errorf("Want: %v, Got: %v", http.StatusOK, recorder.Code)
					return
				}
				if !reflect.DeepEqual(response.Status, tempResp.Status) {
					t.Errorf("Want: %v, Got: %v", tempResp.Status, response.Status)
					return
				}
				if !reflect.DeepEqual(response.Message, tempResp.Message) {
					t.Errorf("Want: %v, Got: %v", tempResp.Message, response.Message)
					return
				}
				marshal, _ := json.Marshal(&response.Data)
				marshalExpected, _ := json.Marshal(&tempResp.Data)
				if !reflect.DeepEqual(marshal, marshalExpected) {
					t.Errorf("Want: %v, Got: %v", marshalExpected, marshal)
					return
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.GetAllArticle(w, r)
			tt.want(*w)
		})
	}
}
