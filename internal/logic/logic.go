package logic

import (
	"github.com/google/uuid"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/codes"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource"
	"log"
	"net/http"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal-chaturvedi/article-management-sys/internal/logic ArticleManagementLogicI

type ArticleManagementLogicI interface {
	InsertArticle(req *model.Article) *model.Response
}

type ArticleManagementLogic struct {
	DsSvc datasource.DataSourceI
}

func NewArticleManagementLogicI(ds datasource.DataSourceI) ArticleManagementLogicI {
	return &ArticleManagementLogic{
		DsSvc: ds,
	}
}

func (l ArticleManagementLogic) InsertArticle(req *model.Article) *model.Response {
	article := model.ArticleDs{
		Id:      uuid.NewString(),
		Title:   req.Title,
		Author:  req.Author,
		Content: req.Content,
	}
	err := l.DsSvc.Insert(article)
	if err != nil {
		log.Print(codes.GetErr(codes.ErrDataSource), err)
		return &model.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrDataSource),
			Data:    nil,
		}
	}
	return &model.Response{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    map[string]string{"id": article.Id},
	}
}
