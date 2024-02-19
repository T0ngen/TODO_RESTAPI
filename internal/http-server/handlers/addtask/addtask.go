package addtask

import (
	
	"log"
	"net/http"
	resp "TodoRESTAPI/internal/lib/response"
	"github.com/go-chi/render"
)


type Request struct {
	Note   string `json:"note"`
	Importance int `json:"importance"`
}




type AddTaskInterface interface{
	AddNewTask(username string, note string, importance int)(bool, error)

}


func New(addTask AddTaskInterface)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.addtask.New"


		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil{
			log.Println("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		username, _, _ := r.BasicAuth()
		ok, err :=addTask.AddNewTask(username, req.Note, req.Importance)
		log.Println(ok)
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

	