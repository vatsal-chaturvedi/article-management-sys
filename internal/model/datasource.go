package model

type ArticleDs struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

const Schema = `
	(
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		content VARCHAR(255) NOT NULL
	);
`
