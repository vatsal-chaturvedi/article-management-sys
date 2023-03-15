package router

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name     string
		setup    func() *config.SvcConfig
		validate func(http.ResponseWriter)
		give     *http.Request
	}{
		{
			name: "Success::No route found Endpoint",
			setup: func() *config.SvcConfig {
				return &config.SvcConfig{
					Cfg: &config.Config{
						DataBase: config.DbCfg{
							Driver: "mysql",
						},
					},
					DbSvc: config.DbSvc{}}
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
			give: httptest.NewRequest(http.MethodGet, "/no-route", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			sqlmock.MonitorPingsOption(true)
			c := tt.setup()
			c.DbSvc.Db = db
			r := Register(c)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tt.give)
			tt.validate(w)
		})
	}
}
