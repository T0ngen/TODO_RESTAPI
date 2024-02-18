package main

import (
	"TodoRESTAPI/internal/config"
	"TodoRESTAPI/internal/storage/postgresql"
	"fmt"
	"os"

	"log"
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
	_ = storage
}