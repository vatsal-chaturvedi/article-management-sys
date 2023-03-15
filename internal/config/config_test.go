package config

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestInitSvcConfig(t *testing.T) {
	type args struct {
		cfg Config
	}
	type testCase struct {
		name string
		args func() args
		want func(args) *SvcConfig
	}
	_, mock1, err := sqlmock.NewWithDSN(":@tcp(:)/?charset=utf8mb4&parseTime=True")
	if err != nil {
		t.Log(err)
		t.Fail()

	}
	_, mock2, err := sqlmock.NewWithDSN(":@tcp(:)/newTemp?charset=utf8mb4&parseTime=True")
	if err != nil {
		t.Log()
		t.Fail()
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []testCase{
		{
			name: "Success",
			args: func() args {
				mock1.ExpectPrepare("CREATE SCHEMA IF NOT EXISTS newTemp ;").ExpectExec().WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1))
				mock1.ExpectClose()
				mock2.ExpectExec(regexp.QuoteMeta("create table if not exists ( id VARCHAR(255) NOT NULL PRIMARY KEY, title VARCHAR(255) NOT NULL, author VARCHAR(255) NOT NULL, content VARCHAR(255) NOT NULL );")).WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1))
				return args{
					cfg: Config{
						ServerConfig: ServerConfig{},
						DataBase: DbCfg{
							Driver: "sqlmock",
							DbName: "newTemp",
						},
						Cacher: CacheConfig{KeyExpiry: "1m"},
					},
				}
			},
			want: func(arg args) *SvcConfig {
				required := &SvcConfig{
					Cfg: &Config{
						ServerConfig: ServerConfig{},
						DataBase: DbCfg{
							Driver: "sqlmock",
							DbName: "newTemp",
						},
						Cacher: CacheConfig{KeyExpiry: "1m", KeyExpiryDuration: time.Minute},
					},
					SvrCfg: ServerConfig{},
				}
				return required
			},
		},
		{
			name: "Failure::DB Open 1",
			args: func() args {
				return args{cfg: Config{DataBase: DbCfg{Driver: "", DbName: "newTemp"}}}
			},
		},
		{
			name: "Failure::DB Prepare",
			args: func() args {
				mock1.ExpectPrepare("CREATE SCHEMA IF NOT EXISTS newTemp ;").WillReturnError(errors.New("error "))
				return args{cfg: Config{DataBase: DbCfg{Driver: "sqlmock", DbName: "newTemp"}}}
			},
		},
		{
			name: "Failure:DB Exec",
			args: func() args {
				mock1.ExpectPrepare("CREATE SCHEMA IF NOT EXISTS newTemp ;").ExpectExec().WillReturnError(errors.New("error")).WillReturnResult(sqlmock.NewResult(1, 1))
				return args{cfg: Config{DataBase: DbCfg{Driver: "sqlmock", DbName: "newTemp"}}}
			},
		},
		{
			name: "Failure:: Exec err 2",
			args: func() args {
				mock1.ExpectPrepare("CREATE SCHEMA IF NOT EXISTS newTemp ;").ExpectExec().WillReturnError(nil).WillReturnResult(sqlmock.NewResult(1, 1))
				mock2.ExpectExec(regexp.QuoteMeta("create table if not exists ( id VARCHAR(255) NOT NULL PRIMARY KEY, title VARCHAR(255) NOT NULL, author VARCHAR(255) NOT NULL, content VARCHAR(255) NOT NULL );")).WillReturnError(errors.New("error exec")).WillReturnResult(sqlmock.NewResult(1, 1))
				return args{cfg: Config{DataBase: DbCfg{Driver: "sqlmock", DbName: "newTemp"}}}
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				a := recover()
				if a != nil {
					t.Log("RECOVER"+tt.name, a)
				}
			}()
			s := tt.args()
			got := InitSvcConfig(s.cfg)
			got.DbSvc.Db = nil
			got.CacherSvc.Rdb = nil
			want := tt.want(s)
			if !reflect.DeepEqual(&got, &want) {
				t.Errorf("Want: %v, Got: %v", &want, &got)
			}
		})

	}
}
