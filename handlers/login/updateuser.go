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
		fmt.Println("yesssss")
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
