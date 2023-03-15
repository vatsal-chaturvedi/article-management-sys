package handler

import (
	"encoding/json"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/logic"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource"

	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal-chaturvedi/article-management-sys/internal/handler ArticleManagementHandlerI

type ArticleManagementHandlerI interface {
	MethodNotAllowed(http.ResponseWriter, *http.Request)
	RouteNotFound(http.ResponseWriter, *http.Request)
}

type articleManagement struct {
	logic logic.ArticleManagementLogicI
}

func NewArticleManagementHandlerI(ds datasource.DataSourceI) ArticleManagementHandlerI {
	svc := &articleManagement{
		logic: logic.NewArticleManagementLogicI(ds),
	}
	return svc
}

func (a articleManagement) MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	_ = json.NewEncoder(w).Encode(&model.Response{
		Status:  http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
		Data:    nil,
	})
}

func (a articleManagement) RouteNotFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(&model.Response{
		Status:  http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
		Data:    nil,
	})
}
