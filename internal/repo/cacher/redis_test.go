package cacher

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name         string
		requestBody  string
		setupFunc    func() (*redis.Client, redismock.ClientMock)
		validateFunc func(string, error)
	}{
		{
			name: "Success::Health",
			setupFunc: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectPing().SetVal("PONG")
				return db, mock
			},
			validateFunc: func(s string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
				if s != "PONG" {
					t.Errorf("want %v got %v", "PONG", s)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb, mockCache := tt.setupFunc()
			mockCacher := NewCacher(config.CacheSvc{Rdb: mockDb})
			x, err := mockCacher.Health()
			if mockCache.ExpectationsWereMet() != nil {
				t.Log(mockCache.ExpectationsWereMet())
				t.Fail()
			}
			tt.validateFunc(x, err)
		})
	}

}

func TestSet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		expiry       time.Duration
		setupFunc    func() (*redis.Client, redismock.ClientMock)
		validateFunc func(error)
	}{
		{
			name:        "Success:: Set",
			requestBody: "ABC",
			setupFunc: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectSet("1", "ABC", 0*time.Second).SetVal("ABC")
				return db, mock
			},
			validateFunc: func(err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Failure:: Set",
			requestBody: "ABC",
			setupFunc: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectSet("1", "ABC", 0*time.Second).SetErr(errors.New("error"))
				return db, mock
			},
			validateFunc: func(err error) {
				if err == nil {
					t.Errorf("want %v got %v", "error", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockCache := tt.setupFunc()
			mockCacher := NewCacher(config.CacheSvc{Rdb: mockDB})
			data := tt.requestBody
			err := mockCacher.Set("1", data, tt.expiry)
			if mockCache.ExpectationsWereMet() != nil {
				t.Log(mockCache.ExpectationsWereMet())
				t.Fail()
			}
			tt.validateFunc(err)
		})
	}

}
func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		expiry       time.Duration
		setupFunc    func() (*redis.Client, redismock.ClientMock)
		validateFunc func(string, error)
	}{
		{
			name:        "Success:: Get",
			requestBody: "ABC",
			setupFunc: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectGet("1").SetVal("ABC")
				return db, mock
			},
			validateFunc: func(data string, err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Failure:: Get",
			requestBody: "ABC",
			setupFunc: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectGet("1").SetErr(errors.New("error"))
				return db, mock
			},
			validateFunc: func(data string, err error) {
				if err == nil {
					t.Errorf("want %v got %v", "error", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockCache := tt.setupFunc()
			mockCacher := NewCacher(config.CacheSvc{Rdb: mockDB})

			data, err := mockCacher.Get("1")
			if mockCache.ExpectationsWereMet() != nil {
				t.Log(mockCache.ExpectationsWereMet())
				t.Fail()
			}
			tt.validateFunc(string(data), err)
		})
	}

}
func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name         string
		requestBody  string
		expiry       time.Duration
		setupFunc    func() (*redis.Client, redismock.ClientMock)
		validateFunc func(error)
	}{
		{
			name:        "Success:: Delete",
			requestBody: "ABC",
			setupFunc: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectDel("1").SetVal(1)
				return db, mock
			},
			validateFunc: func(err error) {
				if err != nil {
					t.Errorf("want %v got %v", nil, err.Error())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockCache := tt.setupFunc()
			mockCacher := NewCacher(config.CacheSvc{Rdb: mockDB})
			err := mockCacher.Delete("1")
			if mockCache.ExpectationsWereMet() != nil {
				t.Log(mockCache.ExpectationsWereMet())
				t.Fail()
			}
			tt.validateFunc(err)
		})
	}

}
