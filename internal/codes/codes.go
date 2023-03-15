package codes

import (
	"fmt"
)

type errCode int

const (
	ErrAssertid errCode = iota + 1000
	ErrReadingReqBody
	ErrUnmarshall
	ErrDataSource
	ErrArticleNotFound
)

var errCodes = map[errCode]string{
	ErrAssertid:        "Unable to assert article id",
	ErrReadingReqBody:  "Unable to read request body",
	ErrUnmarshall:      "Unable to unmarshal request body",
	ErrDataSource:      "DataSource error",
	ErrArticleNotFound: "No article found for specified id",
}

func GetErr(code errCode) string {
	x, ok := errCodes[code]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s", x)
}
