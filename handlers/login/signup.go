package login

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/handlers/app"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var users []User
		err := a.DB.Where("first_name = ?", req.Username).Find(&users).Error
		if err != nil {
			http.Error(w, "Error quering users", http.StatusInternalServerError)
			return
		}

		if len(users) != 0 {
			http.Error(w, "User with username already exits", http.StatusInternalServerError)
			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		fmt.Println(password)

		err = a.DB.Create(&User{FirstName: req.Username}).Error
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Created Successfully!"))

	}
}
