package deletetask

import (
	resp "TodoRESTAPI/internal/lib/response"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)




type DeleteInterface interface{
	DeleteTaskById(username string, id string) (bool, error)
}



func New(deleteInter DeleteInterface)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete{
			const op = "handlers.delete.deleteTask.New"
		
			id := chi.URLParam(r, "id")
			username, _, _ := r.BasicAuth()
			fmt.Println("Extracted ID:", id)
			if id == ""{
				log.Printf("id is empty") 
				render.JSON(w,r, resp.Error("invalid request"))
				return
			}
			ok, err:=deleteInter.DeleteTaskById(username, id)
			if err != nil{
				log.Println("error with delete", op)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if !ok{
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w,r, resp.Error("invalid request"))
				return
			}
			w.WriteHeader(http.StatusOK)
			
		
		}
		


		
	}}