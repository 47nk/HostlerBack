package login

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/handlers/app"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("%v %v", req.Username, req.Password)
		var user []User // Assuming you have a User struct
		if err := a.DB.Where("username = ? AND password = ?", req.Username, req.Password).Find(&user).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
	}
}
