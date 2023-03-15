package main

import (
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/router"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg := config.Config{}
	err := config.LoadFromJson("./configs/config.json", &cfg)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	svcInitCfg := config.InitSvcConfig(cfg)
	r := router.Register(svcInitCfg)
	log.Println("started server on port 8080")
	http.ListenAndServe(":8080", r)
}
