package login

import (
	"encoding/json"
	"hostlerBackend/app"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("working fine!"))
	}
}

func SignUp(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req             SignupRequest
			user            User
			userRoleDetails Role
		)
		//verify user role
		userRole, ok := r.Context().Value("role").(string)
		if !ok {
			http.Error(w, `{"error": "User Role missing or invalid"}`, http.StatusUnauthorized)
			return
		}
		if userRole != "admin" {
			http.Error(w, `{"error": "Only admin can onboard new user"}`, http.StatusUnauthorized)
			return
		}

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
			return
		}

		//decode request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.MobileNumber == "" || req.Role == "" || req.Password == "" {
			http.Error(w, `{"error": "Invalid Request Payload"}`, http.StatusInternalServerError)
			return
		}

		//check if user already exits
		result := a.DB.
			Where("username = ?", req.Username).
			First(&user)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			http.Error(w, `{"error": "Error querying users!"}`, http.StatusInternalServerError)
			return
		}
		if result.RowsAffected != 0 {
			http.Error(w, `{"error": "User with username already exists!"}`, http.StatusInternalServerError)
			return
		}

		//check if role exits
		err = a.DB.Where("role = ? and active = true", req.Role).Find(&userRoleDetails).Error
		if err != nil {
			http.Error(w, `{"error": "Internal Error finding role details!"}`, http.StatusInternalServerError)
			return
		}
		if userRoleDetails.ID == 0 {
			http.Error(w, `{"error" : "No such Role found!"}`, http.StatusInternalServerError)
			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		newUser := User{
			Username:     req.Username,
			RoleId:       int64(userRoleDetails.ID),
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			MobileNumber: req.MobileNumber,
			Password:     string(password),
			CreatedAt:    time.Now(),
			CreatedBy:    userId,
			UpdatedBy:    userId,
		}
		//create new user
		err = a.DB.Create(&newUser).Error
		if err != nil {
			log.Printf("Error creating new user: %v", err)
			http.Error(w, `{"error" : "Internal Error creating new user"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success" : "User Created Successfully!"}`))
	}
}

func SignUpBulk(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
