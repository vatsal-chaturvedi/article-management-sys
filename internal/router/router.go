package router

import (
	"github.com/gorilla/mux"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/handler"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource"
	"net/http"
)

func Register(svcCfg *config.SvcConfig) *mux.Router {
	m := mux.NewRouter()

	m.StrictSlash(true)
	dataSource := datasource.NewSql(svcCfg.DbSvc, svcCfg.Cfg.DataBase.TableName)
	svc := handler.NewArticleManagementHandlerI(dataSource)

	m.NotFoundHandler = http.HandlerFunc(svc.RouteNotFound)
	m.MethodNotAllowedHandler = http.HandlerFunc(svc.MethodNotAllowed)

	router1 := m.PathPrefix("").Subrouter()
	router1.HandleFunc("/articles", svc.InsertArticle).Methods(http.MethodPost)

	return m
}
