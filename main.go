package main

import (
	"TodoRESTAPI/internal/config"
	"TodoRESTAPI/internal/http-server/middlewareauth"
	"TodoRESTAPI/internal/storage/postgresql"
	"fmt"
	"io"
	"net/http"
	"os"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)




func main(){
	cfg := config.MustLoad()
	fmt.Println(cfg)
	log.Println("starting TODO")
	storage, err :=postgresql.New()
	if err != nil{
		log.Fatal("failed to init storage")
		os.Exit(1)
	}
	log.Println("init STORAGE")
	_ = storage
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	log.Println("starting server", cfg.Address)
	router.Route("/url", func(r chi.Router) {
		r.Use((middlewareauth.BasicAuthFromDB(storage)))
		r.Get("/", HomePage)
		// TODO: add DELETE /url/{id}
	})

	

	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.Timeout,
	}
	if err:= srv.ListenAndServe(); err != nil{
		log.Fatal("failed to start server")
		os.Exit(1)
	}
}
func HomePage(w http.ResponseWriter, r *http.Request){
	io.WriteString(w,"HEllo")
}