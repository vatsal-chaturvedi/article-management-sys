package datasource

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestSqlDs_Get(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (sqlDs, sqlmock.Sqlmock)
		filter    map[string]interface{}
		validator func([]model.ArticleDs, error, sqlmock.Sqlmock)
	}{
		{
			name: "SUCCESS::Get",
			filter: map[string]interface{}{
				"id": "1234",
			},
			setupFunc: func() (sqlDs, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fail()
				}
				dB := sqlDs{
					sqlSvc: db,
					table:  "newTemp",
				}
				mock.ExpectQuery("SELECT id, title, author, content FROM newTemp WHERE id = '1234' ORDER BY title LIMIT 1 OFFSET 2 ").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "content"}).AddRow("1", "TITLE", "AUTHOR", "CONTENT"))
				return dB, mock
			},
			validator: func(rows []model.ArticleDs, err error, mock sqlmock.Sqlmock) {
				temp := []model.ArticleDs{{
					Id:      "1",
					Title:   "TITLE",
					Author:  "AUTHOR",
					Content: "CONTENT",
				}}
				if mock.ExpectationsWereMet() != nil {
					t.Errorf("Want: %v, Got: %v", nil, mock.ExpectationsWereMet())
					return
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err)
					return
				}
				if !reflect.DeepEqual(rows, temp) {
					t.Errorf("Want: %v, Got: %v", temp, rows)
					return
				}
			},
		},
		{
			name:   "FAILURE::Get:: get rows query error",
			filter: map[string]interface{}{"userid": "1234"},
			setupFunc: func() (sqlDs, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fail()
				}
				dB := sqlDs{
					sqlSvc: db,
					table:  "newTemp",
				}
				mock.ExpectQuery("SELECT id, title, author, content FROM newTemp WHERE userid = '1234' ORDER BY title LIMIT 1 OFFSET 2 ;").WillReturnError(errors.New("Unknown column"))
				return dB, mock
			},
			validator: func(rows []model.ArticleDs, err error, mock sqlmock.Sqlmock) {
				if mock.ExpectationsWereMet() != nil {
					t.Errorf("Want: %v, Got: %v", nil, mock.ExpectationsWereMet())
					return
				}
				if !strings.Contains(err.Error(), "Unknown column") {
					t.Errorf("Want: %v, Got: %v", "Unknown column", err)
				}
			},
		},
	}

	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// STEP 1: seting up all instances for the specific test case
			db, mock := tt.setupFunc()
			// STEP 2: call the test function
			rows, err := db.Get(tt.filter, 1, 2)

			// STEP 3: validation of output
			if tt.validator != nil {
				tt.validator(rows, err, mock)
			}
		})
	}
}

//
func TestSqlDs_Insert(t *testing.T) {
	// table driven tests
	tests := []struct {
		name      string
		data      model.ArticleDs
		setupFunc func() (sqlDs, sqlmock.Sqlmock)
		validator func(sqlmock.Sqlmock, error)
	}{
		{
			name: "SUCCESS:: Insert Transaction",
			data: model.ArticleDs{
				Title:   "TITLE",
				Author:  "AUTHOR",
				Content: "CONTENT",
			},
			setupFunc: func() (sqlDs, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fail()
				}
				dB := sqlDs{
					sqlSvc: db,
					table:  "newTemp",
				}
				m := mock.ExpectExec(regexp.QuoteMeta("INSERT INTO newTemp(id, title, author, content) VALUES(?,?,?,?)")).WithArgs(sqlmock.AnyArg(), "TITLE", "AUTHOR", "CONTENT")
				m.WillReturnError(nil)
				m.WillReturnResult(sqlmock.NewResult(1, 1))
				return dB, mock
			},
			validator: func(mock sqlmock.Sqlmock, err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
					return
				}
				if mock.ExpectationsWereMet() != nil {
					t.Errorf("Want: %v, Got: %v", nil, mock.ExpectationsWereMet())
					return
				}

			},
		},
		{
			name: "FAILURE:: insert :: sql error",
			data: model.ArticleDs{
				Title:   "TITLE",
				Author:  "AUTHOR",
				Content: "CONTENT",
			},
			setupFunc: func() (sqlDs, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fail()
				}
				dB := sqlDs{
					sqlSvc: db,
					table:  "newTemp",
				}
				m := mock.ExpectExec(regexp.QuoteMeta("INSERT INTO newTemp(id, title, author, content) VALUES(?,?,?,?)")).WithArgs(sqlmock.AnyArg(), "TITLE", "AUTHOR", "CONTENT")
				m.WillReturnError(errors.New("sql error"))
				m.WillReturnResult(sqlmock.NewResult(1, 1))
				return dB, mock
			},
			validator: func(mock sqlmock.Sqlmock, err error) {
				if mock.ExpectationsWereMet() != nil {
					t.Errorf("Want: %v, Got: %v", nil, mock.ExpectationsWereMet())
					return
				}
				if err.Error() != errors.New("sql error").Error() {
					t.Errorf("Want: %v, Got: %v", "sql error", err.Error())
					return
				}
			},
		},
	}
	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := tt.setupFunc()
			// STEP 2: call the test function
			err := db.Insert(tt.data)
			// STEP 3: validation of output
			if tt.validator != nil {
				tt.validator(mock, err)
			}
		})
	}
}
