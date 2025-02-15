package login

import (
	"encoding/json"
	"hostlerBackend/handlers/app"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignupRequest struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	RoleId       int64  `json:"role_id"`
	MobileNumber string `json:"mobile_num"`
	Password     string `json:"password"`
}

func TestAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("working fine!"))
	}
}

func SignUp(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req  SignupRequest
			user User
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.MobileNumber == "" || req.RoleId == 0 || req.Password == "" {
			http.Error(w, "Invalid Request Payload", http.StatusInternalServerError)
			return
		}

		result := a.DB.
			Where("username = ?", req.Username).
			First(&user)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			http.Error(w, "Error querying users", http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 0 {
			http.Error(w, "User with username already exits", http.StatusInternalServerError)
			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		newUser := User{
			Username:     req.Username,
			RoleId:       req.RoleId,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			MobileNumber: req.MobileNumber,
			Password:     string(password),
			CreatedAt:    time.Now(),
		}

		err := a.DB.Create(&newUser).Error
		if err != nil {
			log.Printf("Error creating new user: %v", err)
			http.Error(w, "Error creating new user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Created Successfully!"))
	}
}

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
