package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/codes"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/logic"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal-chaturvedi/article-management-sys/internal/handler ArticleManagementHandlerI

type ArticleManagementHandlerI interface {
	MethodNotAllowed(http.ResponseWriter, *http.Request)
	RouteNotFound(http.ResponseWriter, *http.Request)
	InsertArticle(http.ResponseWriter, *http.Request)
	GetArticleById(w http.ResponseWriter, r *http.Request)
	GetAllArticle(w http.ResponseWriter, r *http.Request)
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
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&model.Response{
			Status:  http.StatusBadRequest,
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
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&model.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrUnmarshall),
			Data:    nil,
		})
		return
	}
	//removed blank spaces
	article.Title = strings.Trim(article.Title, " ")
	article.Author = strings.Trim(article.Author, " ")
	article.Content = strings.Trim(article.Content, " ")

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

func (svc articleManagement) GetArticleById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&model.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrAssertid),
			Data:    nil,
		})
		return
	}
	resp := svc.logic.GetArticle(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	_ = json.NewEncoder(w).Encode(&model.Response{
		Status:  resp.Status,
		Message: resp.Message,
		Data:    resp.Data,
	})
}

func (svc articleManagement) GetAllArticle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit == 0 {
		log.Print(fmt.Sprintf("setting default limit as %d", 20))
		limit = 20
	}
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page == 1 {
		log.Print(fmt.Sprintf("setting default page as %d", 1))
		page = 1
	}
	resp := svc.logic.GetAllArticle(limit, page)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	_ = json.NewEncoder(w).Encode(&model.Response{
		Status:  resp.Status,
		Message: resp.Message,
		Data:    resp.Data,
	})
}
