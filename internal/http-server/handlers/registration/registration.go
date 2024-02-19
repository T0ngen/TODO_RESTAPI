package registration

import (
	"log"
	"net/http"

	resp "TodoRESTAPI/internal/lib/response"
	"TodoRESTAPI/internal/lib/response/passhash"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
)


type NewUser struct{
	Username string `json:"username"`
	Password string `json:"password"`


}

type BadResponse struct{
	Error error `json:"error"`
	Number int `json:"number"`
}

type RegInterface interface{
	CheckUsername(username string)(bool, error)
	CreateNewUser(username string, password string) (bool, error)
}



func (u NewUser) Validate() error {
	return validation.ValidateStruct(&u,
		
		validation.Field(&u.Username, validation.Required, validation.Length(4, 12).Error("The username length must be from 4 to 12 characters")),
		validation.Field(&u.Password, validation.Required, validation.Length(6,15).Error("The password length must be from 6 to 15 characters")),
		
		
	)
}

func New(reg RegInterface)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
	
	const op = "handlers.registration.registration.New"
	var user NewUser

	err := render.DecodeJSON(r.Body, &user)
	if err != nil{
		log.Println("user registration body is empty")

		render.JSON(w, r, resp.Error("empty request"))

		return
	}
	err = user.Validate()
	if err != nil{
		badResponse := BadResponse{Error: err, Number:  500}
		w.WriteHeader(badResponse.Number)
		render.JSON(w,r, badResponse)
		return
	}
	ok, err :=reg.CheckUsername(user.Username)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if ok{
		w.WriteHeader(500)
		render.JSON(w, r, resp.Error("This username already exists"))
		return
	}
	HashPass := passhash.HashedPassword(user.Password)
	ok, err = reg.CreateNewUser(user.Username, HashPass)
	if err != nil{
		log.Fatalf("%v", op)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok{
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	render.JSON(w, r, resp.RespOK("Account created"))
		

}}