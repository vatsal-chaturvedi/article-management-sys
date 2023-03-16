package logic

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/codes"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/repo/datasource"
	"github.com/vatsal-chaturvedi/article-management-sys/pkg/mock"
	"net/http"
	"reflect"
	"testing"
)

func TestArticleManagementLogic_InsertArticle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() datasource.DataSourceI
		give  *model.Article
		want  *model.Response
	}{
		{
			name: "Success",
			setup: func() datasource.DataSourceI {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				x := model.ArticleDs{
					Id:      "1",
					Title:   "title",
					Content: "content",
					Author:  "author",
				}
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).
					DoAndReturn(func(article model.ArticleDs) *model.Response {
						article.Id = "1"
						if !reflect.DeepEqual(article, x) {
							t.Logf("Want: %v, Got: %v", x, article)
							t.Fail()
						}
						return nil
					})
				return mockDs
			},
			give: &model.Article{
				Title:   "title",
				Content: "content",
				Author:  "author",
			},
			want: &model.Response{
				Status:  http.StatusCreated,
				Message: "Success",
				Data:    map[string]interface{}{"id": "1"},
			},
		},
		{
			name: "Failure:: Datasource Error",
			setup: func() datasource.DataSourceI {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				x := model.ArticleDs{
					Title:   "title",
					Content: "content",
					Author:  "author",
				}
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).
					DoAndReturn(func(article model.ArticleDs) *model.Response {
						article.Id = ""
						if !reflect.DeepEqual(article, x) {
							t.Logf("Want: %v, Got: %v", x, article)
							t.Fail()
						}
						return nil
					}).Return(errors.New(codes.GetErr(codes.ErrDataSource)))
				return mockDs
			},
			give: &model.Article{
				Title:   "title",
				Content: "content",
				Author:  "author",
			},
			want: &model.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrDataSource),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewArticleManagementLogicI(tt.setup())
			got := rec.InsertArticle(tt.give)
			if !reflect.DeepEqual(got.Status, tt.want.Status) {
				t.Logf("Want: %v, Got: %v", tt.want.Status, got.Status)
				t.Fail()
			}
			if !reflect.DeepEqual(got.Message, tt.want.Message) {
				t.Logf("Want: %v, Got: %v", tt.want.Message, got.Message)
				t.Fail()
			}
		})
	}
}

func TestArticleManagementLogic_GetArticle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() datasource.DataSourceI
		give  *model.Article
		want  *model.Response
	}{
		{
			name: "Success",
			setup: func() datasource.DataSourceI {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				x := model.ArticleDs{
					Id:      "1",
					Title:   "title",
					Content: "content",
					Author:  "author",
				}
				mockDs.EXPECT().Get(map[string]interface{}{"id": "1"}, 1, 0).Times(1).Return([]model.ArticleDs{x}, nil)
				return mockDs
			},
			want: &model.Response{
				Status:  http.StatusOK,
				Message: "Success",
				Data: []model.ArticleDs{{
					Id:      "1",
					Title:   "title",
					Content: "content",
					Author:  "author",
				}},
			},
		},
		{
			name: "Failure:: No article found",
			setup: func() datasource.DataSourceI {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"id": "1"}, 1, 0).Times(1).Return([]model.ArticleDs{}, nil)
				return mockDs
			},
			want: &model.Response{
				Status:  http.StatusBadRequest,
				Message: codes.GetErr(codes.ErrArticleNotFound),
				Data:    nil,
			},
		},
		{
			name: "Failure:: Datasource Error",
			setup: func() datasource.DataSourceI {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"id": "1"}, 1, 0).Times(1).Return(nil, errors.New(""))
				return mockDs
			},
			want: &model.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrDataSource),
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewArticleManagementLogicI(tt.setup())
			got := rec.GetArticle("1")
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("Want: %v, Got: %v", tt.want, got)
				t.Fail()
			}
		})
	}
}

func TestArticleManagementLogic_GetAllArticle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() datasource.DataSourceI
		give  *model.Article
		want  *model.Response
	}{
		{
			name: "Success",
			setup: func() datasource.DataSourceI {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				x := model.ArticleDs{
					Id:      "1",
					Title:   "title",
					Content: "content",
					Author:  "author",
				}
				y := model.ArticleDs{
					Id:      "2",
					Title:   "title",
					Content: "content",
					Author:  "author",
				}
				mockDs.EXPECT().Get(nil, 5, 0).Times(1).Return([]model.ArticleDs{x, y}, nil)
				return mockDs
			},
			want: &model.Response{
				Status:  http.StatusOK,
				Message: "Success",
				Data: []model.ArticleDs{{
					Id:      "1",
					Title:   "title",
					Content: "content",
					Author:  "author",
				},
					{
						Id:      "2",
						Title:   "title",
						Content: "content",
						Author:  "author",
					}},
			},
		},
		{
			name: "Failure:: Datasource Error",
			setup: func() datasource.DataSourceI {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(nil, 5, 0).Times(1).Return(nil, errors.New(""))
				return mockDs
			},
			want: &model.Response{
				Status:  http.StatusInternalServerError,
				Message: codes.GetErr(codes.ErrDataSource),
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewArticleManagementLogicI(tt.setup())
			got := rec.GetAllArticle(5, 1)
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("Want: %v, Got: %v", tt.want, got)
				t.Fail()
			}
		})
	}
}
