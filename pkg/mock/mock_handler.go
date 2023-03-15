// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vatsal-chaturvedi/article-management-sys/internal/handler (interfaces: ArticleManagementHandlerI)

// Package mock is a generated GoMock package.
package mock

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockArticleManagementHandlerI is a mock of ArticleManagementHandlerI interface.
type MockArticleManagementHandlerI struct {
	ctrl     *gomock.Controller
	recorder *MockArticleManagementHandlerIMockRecorder
}

// MockArticleManagementHandlerIMockRecorder is the mock recorder for MockArticleManagementHandlerI.
type MockArticleManagementHandlerIMockRecorder struct {
	mock *MockArticleManagementHandlerI
}

// NewMockArticleManagementHandlerI creates a new mock instance.
func NewMockArticleManagementHandlerI(ctrl *gomock.Controller) *MockArticleManagementHandlerI {
	mock := &MockArticleManagementHandlerI{ctrl: ctrl}
	mock.recorder = &MockArticleManagementHandlerIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleManagementHandlerI) EXPECT() *MockArticleManagementHandlerIMockRecorder {
	return m.recorder
}

// GetAllArticle mocks base method.
func (m *MockArticleManagementHandlerI) GetAllArticle(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetAllArticle", arg0, arg1)
}

// GetAllArticle indicates an expected call of GetAllArticle.
func (mr *MockArticleManagementHandlerIMockRecorder) GetAllArticle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllArticle", reflect.TypeOf((*MockArticleManagementHandlerI)(nil).GetAllArticle), arg0, arg1)
}

// GetArticleById mocks base method.
func (m *MockArticleManagementHandlerI) GetArticleById(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetArticleById", arg0, arg1)
}

// GetArticleById indicates an expected call of GetArticleById.
func (mr *MockArticleManagementHandlerIMockRecorder) GetArticleById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArticleById", reflect.TypeOf((*MockArticleManagementHandlerI)(nil).GetArticleById), arg0, arg1)
}

// InsertArticle mocks base method.
func (m *MockArticleManagementHandlerI) InsertArticle(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InsertArticle", arg0, arg1)
}

// InsertArticle indicates an expected call of InsertArticle.
func (mr *MockArticleManagementHandlerIMockRecorder) InsertArticle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertArticle", reflect.TypeOf((*MockArticleManagementHandlerI)(nil).InsertArticle), arg0, arg1)
}

// MethodNotAllowed mocks base method.
func (m *MockArticleManagementHandlerI) MethodNotAllowed(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MethodNotAllowed", arg0, arg1)
}

// MethodNotAllowed indicates an expected call of MethodNotAllowed.
func (mr *MockArticleManagementHandlerIMockRecorder) MethodNotAllowed(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MethodNotAllowed", reflect.TypeOf((*MockArticleManagementHandlerI)(nil).MethodNotAllowed), arg0, arg1)
}

// RouteNotFound mocks base method.
func (m *MockArticleManagementHandlerI) RouteNotFound(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RouteNotFound", arg0, arg1)
}

// RouteNotFound indicates an expected call of RouteNotFound.
func (mr *MockArticleManagementHandlerIMockRecorder) RouteNotFound(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RouteNotFound", reflect.TypeOf((*MockArticleManagementHandlerI)(nil).RouteNotFound), arg0, arg1)
}
