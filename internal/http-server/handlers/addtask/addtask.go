package addtask

import (
	resp "TodoRESTAPI/internal/lib/response"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/render"
)


type Request struct {
	Note   string `json:"note"`
	Importance int `json:"importance"`
}



//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=AddTaskInterface
type AddTaskInterface interface{
	AddNewTask(username string, note string, importance int)(bool, error)
}



func New(addTask AddTaskInterface)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.addtask.New"
		fmt.Println(r.Body)
		var req Request
		
		err := render.DecodeJSON(r.Body, &req)

		if err != nil{
			log.Println("request body is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		username, _, _ := r.BasicAuth()
		
		ok, err :=addTask.AddNewTask(username, req.Note, req.Importance)

		if err != nil{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !ok{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.RespOK("added to notes"))
		
		
	}}

	