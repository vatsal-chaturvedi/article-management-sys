package datasource

import "github.com/vatsal-chaturvedi/article-management-sys/internal/model"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../../pkg/mock/mock_datasource.go --package=mock github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource DataSourceI

type DataSourceI interface {
	Get(filter map[string]interface{}, limit int, offset int) ([]model.ArticleDs, error)
	Insert(user model.ArticleDs) error
}
