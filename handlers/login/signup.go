package login

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/handlers/app"
	"log"
	"net/http"
	"time"

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

		result := a.DB.Where("username = ?", req.Username).First(&user)
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

func TestAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.Context().Value("user_id").(string))
		w.Write([]byte("working fine!"))
	}
}
