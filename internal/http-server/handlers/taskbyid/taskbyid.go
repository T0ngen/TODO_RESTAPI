package taskbyid

import (
	resp "TodoRESTAPI/internal/lib/response"
	"TodoRESTAPI/internal/storage/postgresql"
	"encoding/json"

	"log"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)


type TaskByIdInterface interface{
	CheckTaskById(username string, id string)(*postgresql.TaskById, error)
}



func ById(taskbyidinter TaskByIdInterface)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.taskbyid.taskbyid.ByID"
		if r.Method == http.MethodGet{
			id := chi.URLParam(r, "id")
			if id == ""{
				log.Printf("id is empty") 
				render.JSON(w,r, resp.Error("invalid request"))
				return
			}
			username, _, _ := r.BasicAuth()
			result, err :=taskbyidinter.CheckTaskById(username, id)
			if err != nil{
				log.Printf("Error: %v : %s", err, op) 
			}
			
			if result == nil{
				render.JSON(w,r, resp.Error("invalid request"))
				return
			}
			jsonData, err := json.Marshal(result)
			if err != nil{
				log.Println("cant marshal json", op)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		}
		

	}
}