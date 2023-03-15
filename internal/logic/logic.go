package logic

import (
	"github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/blog/internal/logic BlogLogicIer

type ArticleManagementLogicI interface {
}

type blogLogic struct {
	DsSvc datasource.DataSourceI
}

func NewArticleManagementLogicI(ds datasource.DataSourceI) ArticleManagementLogicI {
	return &blogLogic{
		DsSvc: ds,
	}
}
