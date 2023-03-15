package handler

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/codes"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/logic"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource"
	"io/ioutil"
	"log"
	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal-chaturvedi/article-management-sys/internal/handler ArticleManagementHandlerI

type ArticleManagementHandlerI interface {
	MethodNotAllowed(http.ResponseWriter, *http.Request)
	RouteNotFound(http.ResponseWriter, *http.Request)
	InsertArticle(http.ResponseWriter, *http.Request)
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

func (svc articleManagement) InsertArticle(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&model.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrReadingReqBody),
			Data:    nil,
		})
		return
	}
	var article model.Article
	err = json.Unmarshal(bytes, &article)
	if err != nil {
		log.Print(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&model.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrUnmarshall),
			Data:    nil,
		})
		return
	}
	validate := validator.New()
	if err := validate.Struct(article); err != nil {
		log.Print(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&model.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	resp := svc.logic.InsertArticle(&article)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	_ = json.NewEncoder(w).Encode(&model.Response{
		Status:  resp.Status,
		Message: resp.Message,
		Data:    resp.Data,
	})
}
