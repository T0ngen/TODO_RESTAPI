package taskbyid

import (
	resp "TodoRESTAPI/internal/lib/response"
	"TodoRESTAPI/internal/storage/postgresql"
	
	"fmt"

	"encoding/json"

	"log"

	"net/http"

	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=TaskByIdInterface
type TaskByIdInterface interface {
	CheckTaskById(username string, id string) (*postgresql.TaskById, error)
}


func ById(taskbyidinter TaskByIdInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.taskbyid.taskbyid.ByID"
		if r.Method == http.MethodGet {
			

			idS := r.URL.Query().Get("id")
			fmt.Println("id",idS)
			if idS == "" {
				fmt.Println("empty id")
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("empty id"))
				return
			}
			username, _, _ := r.BasicAuth()
			// idS = strconv.Itoa(id)
			result, err := taskbyidinter.CheckTaskById(username, idS)
			if err != nil {
				log.Printf("Error: %v : %s", err, op)
				w.WriteHeader(http.StatusBadRequest)
				
				return
			}

			if result == nil {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("invalid request"))
				return
			}
			jsonData, err := json.Marshal(result)
			if err != nil {
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
