package middlewareauth

import (
	"TodoRESTAPI/internal/lib/response/passhash"
	"log"
	"net/http"
)



type AuthUser interface{
	CheckUserInDb(username string, password string) (bool, error)
}





func BasicAuthFromDB(auth AuthUser)func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
	  fn := func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		hashedPassword := passhash.HashedPassword(password)
		ok, err :=auth.CheckUserInDb(username, hashedPassword)
		
		if err != nil{
			log.Printf("Error: %v", err) 
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// Проверка, что username и password корректные
		
		if !ok{
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
  
		// Прошли аутентификацию, вызываем следующий обработчик
		next.ServeHTTP(w, r)
	  }
	  return http.HandlerFunc(fn)
	}
  }