package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/codes"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/model"
	redis "github.com/vatsal-chaturvedi/article-management-sys/internal/repo/cacher"

	"log"
	"net/http"
)

type Middleware struct {
	cfg    *config.Config
	cacher redis.CacherI
}

type respWriterWithStatus struct {
	status   int
	response string
	http.ResponseWriter
}

func (w *respWriterWithStatus) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *respWriterWithStatus) Write(d []byte) (int, error) {
	w.response = string(d)
	return w.ResponseWriter.Write(d)
}

func NewMiddleware(cfg *config.SvcConfig, cacherI redis.CacherI) *Middleware {
	return &Middleware{
		cfg:    cfg.Cfg,
		cacher: cacherI,
	}
}

func (t Middleware) Cacher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key string
		var cacheResponse model.CacheResponse
		key = fmt.Sprint(r.URL.String())

		Cacher := t.cacher
		by, err := Cacher.Get(key)
		if err == nil {
			err = json.Unmarshal(by, &cacheResponse)
			if err != nil {
				log.Print(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(&model.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrUnmarshall),
					Data:    nil,
				})
				return
			}
			w.Header().Set("Content-Type", cacheResponse.ContentType)
			w.Write([]byte(cacheResponse.Response))
			w.WriteHeader(cacheResponse.Status)
			return
		}

		hijackedWriter := &respWriterWithStatus{-1, "", w}
		next.ServeHTTP(hijackedWriter, r)

		if hijackedWriter.status < 200 || hijackedWriter.status >= 300 {
			return
		}

		cacheResponse = model.CacheResponse{
			Status:      hijackedWriter.status,
			Response:    hijackedWriter.response,
			ContentType: w.Header().Get("Content-Type"),
		}
		byt, err := json.Marshal(cacheResponse)
		if err != nil {
			log.Print(err)
			return
		}
		err = Cacher.Set(key, byt, t.cfg.Cacher.KeyExpiryDuration)
		if err != nil {
			log.Print(err)
			return
		}
	})
}
