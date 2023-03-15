package datasource

import (
	"database/sql"
	"fmt"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	"strings"
)

type sqlDs struct {
	sqlSvc *sql.DB
	table  string
}

// NewSql creates a new instance of sqlDs with a given database service and table name.
func NewSql(dbSvc config.DbSvc, tableName string) DataSourceI {
	return &sqlDs{
		sqlSvc: dbSvc.Db,
		table:  tableName,
	}
}

// queryFromMap returns a string containing SQL queries generated from a given map of filters and a join separator.
func queryFromMap(d map[string]interface{}, join string) string {
	var (
		q string
		f []string
	)
	for k, v := range d {
		switch v.(type) {
		case string:
			f = append(f, fmt.Sprintf(`%s = '%s'`, k, v))
		default:
			f = append(f, fmt.Sprintf(`%s = %v`, k, v))
		}
	}
	if len(f) > 0 {
		q = fmt.Sprintf(`%s`, strings.Join(f, ` `+join+` `))
	}
	return q
}

// HealthCheck checks the health of the database service.
func (d sqlDs) HealthCheck() bool {
	err := d.sqlSvc.Ping()
	return err == nil
}

// Get retrieves transactions from the database service based on a given set of filters, limit, and offset.
func (d sqlDs) Get(filter map[string]interface{}, limit int, offset int) ([]model.ArticleDs, int, error) {
	var article model.ArticleDs
	var articles []model.ArticleDs
	var count int
	q := fmt.Sprintf("SELECT id, title, author, content FROM %s", d.table)
	whereQuery := queryFromMap(filter, " AND ")
	if whereQuery != "" {
		whereQuery = " WHERE " + whereQuery
		q += whereQuery
	}
	queryCount := fmt.Sprintf("SELECT COUNT(`id`) FROM %s %s", d.table, whereQuery)
	rowsCount, err := d.sqlSvc.Query(queryCount)
	if err != nil {
		return nil, 0, err
	}
	for rowsCount.Next() {
		err = rowsCount.Scan(&count)
		if err != nil {
			return nil, 0, err
		}
	}
	//sort based on created at
	r := " ORDER BY title ;"
	if limit > 0 {
		r = fmt.Sprintf(" ORDER BY title LIMIT %d OFFSET %d ;", limit, offset)
	}
	q += r
	rows, err := d.sqlSvc.Query(q)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		err = rows.Scan(&article.Id, &article.Title, &article.Author, &article.Content)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, article)
	}
	rows.Close()
	return articles, count, nil
}

// Insert adds a new transaction to the database service.
func (d sqlDs) Insert(article model.ArticleDs) error {
	queryString := fmt.Sprintf("INSERT INTO %s", d.table)
	_, err := d.sqlSvc.Exec(queryString+"(id, title, author, content) VALUES(?,?,?,?)", article.Id, article.Title, article.Author, article.Content)
	if err != nil {
		return err
	}
	return err
}
