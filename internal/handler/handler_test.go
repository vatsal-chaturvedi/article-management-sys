package handler

import (
	"encoding/json"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"

	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_common_MethodNotAllowed(t *testing.T) {
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

func Test_common_RouteNotFound(t *testing.T) {
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
