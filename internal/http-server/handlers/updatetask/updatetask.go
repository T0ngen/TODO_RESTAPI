package updatetask

import (
	resp "TodoRESTAPI/internal/lib/response"
	
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)


type Request struct {
	Note   string `json:"note"`
	Importance int `json:"importance"`
}


type UpdateInterface interface{
	UpdateTask(username string, id string, note string, importance int) (bool, error)
}




func New(upd UpdateInterface)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		id := chi.URLParam(r, "id")

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil{
			log.Println("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		ok, err := upd.UpdateTask(username, id, req.Note, req.Importance)
		if err != nil{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !ok{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp.RespOK("note updated"))
	}
}