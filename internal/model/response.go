package model

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CacheResponse struct {
	Status      int    // Status code of the cached response
	Response    string // Response body of the cached response
	ContentType string // Content type of the cached response
}
