package main

import (
	"TodoRESTAPI/internal/config"
	"TodoRESTAPI/internal/http-server/handlers/addtask"
	deletetask "TodoRESTAPI/internal/http-server/handlers/delete"
	"TodoRESTAPI/internal/http-server/handlers/registration"
	"TodoRESTAPI/internal/http-server/handlers/taskbyid"
	"TodoRESTAPI/internal/http-server/handlers/tasks"
	"TodoRESTAPI/internal/http-server/handlers/updatetask"
	"TodoRESTAPI/internal/http-server/middlewareauth"
	"TodoRESTAPI/internal/storage/postgresql"

	"net/http"
	"os"

	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)




func main(){
	cfg := config.MustLoad()
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
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	log.Println("starting server", cfg.Address)
	router.Route("/", func(r chi.Router) {
		r.Use((middlewareauth.BasicAuthFromDB(storage)))
		r.Delete("/tasks/{id}", deletetask.New(storage))
		r.Put("/tasks/{id}", updatetask.New(storage))
		r.Post("/tasks", addtask.New(storage))
		r.Get("/tasks/{id}", taskbyid.ById(storage))
		r.Get("/tasks", tasks.All(storage))
		
		
		
		
	})

	router.Post("/register", registration.New(storage))

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
