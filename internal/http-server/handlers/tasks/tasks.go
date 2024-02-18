package tasks

import (
	"TodoRESTAPI/internal/storage/postgresql"
	"encoding/json"
	
	"log"
	"net/http"
)




type AllTasksk interface{
	CheckAllUserTasks(username string) ([]postgresql.Task, error)
}


func New(allTasks AllTasksk)http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		allTasks,err  :=allTasks.CheckAllUserTasks(username)
		if err!= nil{
			log.Printf("Error: %v", err) 
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
		var tasks []postgresql.Task
		for _, task := range allTasks {
			tasks = append(tasks, postgresql.Task{Id: task.Id, Text: task.Text, Importance: task.Importance})
		   }
		 
		jsonData, err := json.Marshal(tasks)
		if err != nil {
			log.Printf("Error encoding JSON: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		   }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}