package main

import (
	//"blog/pkg/setting"
	//"blog/routers"
	"github.com/happykirk/blog/pkg/setting"
	"github.com/happykirk/blog/routers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := routers.InitRouter()
	log.Println("port:", setting.HTTPPort)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}
