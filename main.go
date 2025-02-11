package main

import (
	"handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	handlers.Init()
	fs := http.FileServer(http.Dir("./web/static/css"))
	http.Handle("/static/css/", http.StripPrefix("/static/css/", fs))

	http.HandleFunc("/", handlers.HandleHome)

	server := &http.Server{
		Addr:              ":8080",           // Address of the server (port is for example)
		ReadHeaderTimeout: 10 * time.Second,  // Time allowed to read headers
		WriteTimeout:      10 * time.Second,  // Max time to write response
		IdleTimeout:       120 * time.Second, // Max time between two requests
		MaxHeaderBytes:    1 << 20,           // 1 MB, maximum bytes server will read
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
