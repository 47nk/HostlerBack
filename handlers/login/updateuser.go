package login

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/handlers/app"
	"net/http"

	"github.com/gorilla/mux"
)

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
}

func ComplexQuery(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users User

		// Example complex query

		var user User
		if err := a.DB.Model(&User{}).Where("id = ?", 1).First(&user).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Print(user.FirstName)
		response, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func UpdateUser(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		var user User
		if err := a.DB.First(&user, id).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		a.DB.Model(&user).Updates(User{
			FirstName: req.FirstName,
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User updated successfully"))
	}
}
