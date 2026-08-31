package login

import (
	"encoding/json"
	"errors"
	"hostlerBackend/app"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func writeJSON(w http.ResponseWriter, status int, response SignUpResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func requireAdmin(r *http.Request) (int64, bool) {
	userRole, ok := r.Context().Value("role").(string)
	if !ok || userRole == "" || userRole != "admin" {
		return 0, false
	}

	userIDString, ok := r.Context().Value("user_id").(string)
	if !ok || userIDString == "" {
		return 0, false
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil || userID <= 0 {
		return 0, false
	}

	return userID, true
}

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
		userID, ok := requireAdmin(r)
		if !ok {
			writeJSON(w, http.StatusForbidden, SignUpResponse{Error: "Only an authenticated admin can onboard new users"})
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, SignUpResponse{Error: "Invalid JSON request body"})
			return
		}
		req.Username = strings.TrimSpace(req.Username)
		req.FirstName = strings.TrimSpace(req.FirstName)
		req.LastName = strings.TrimSpace(req.LastName)
		req.MobileNumber = strings.TrimSpace(req.MobileNumber)
		req.Role = strings.TrimSpace(req.Role)
		if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.MobileNumber == "" || req.Role == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, SignUpResponse{Error: "username, first_name, last_name, mobile_num, role, and password are required"})
			return
		}

		result := a.DB.
			Where("username = ?", req.Username).
			First(&user)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Error querying user %q: %v", req.Username, result.Error)
			writeJSON(w, http.StatusInternalServerError, SignUpResponse{Error: "Error querying users"})
			return
		}
		if result.RowsAffected != 0 {
			writeJSON(w, http.StatusConflict, SignUpResponse{Error: "User with username already exists"})
			return
		}

		err := a.DB.Where("role = ? AND active = ?", req.Role, true).First(&userRoleDetails).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSON(w, http.StatusBadRequest, SignUpResponse{Error: "No active role found"})
			return
		}
		if err != nil {
			log.Printf("Error querying role %q: %v", req.Role, err)
			writeJSON(w, http.StatusInternalServerError, SignUpResponse{Error: "Error finding role details"})
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			writeJSON(w, http.StatusInternalServerError, SignUpResponse{Error: "Could not secure password"})
			return
		}
		newUser := User{
			Username:     req.Username,
			RoleId:       int64(userRoleDetails.ID),
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			MobileNumber: req.MobileNumber,
			Password:     string(password),
			CreatedAt:    time.Now(),
			CreatedBy:    userID,
			UpdatedBy:    userID,
		}
		err = a.DB.Create(&newUser).Error
		if err != nil {
			log.Printf("Error creating new user: %v", err)
			writeJSON(w, http.StatusInternalServerError, SignUpResponse{Error: "Error creating new user"})
			return
		}
		writeJSON(w, http.StatusCreated, SignUpResponse{Success: "User created successfully"})
	}
}

func SignUpBulk(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func UpdateUser(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Security: only an authenticated admin should be allowed to
		// change another user's details.
		adminID, ok := requireAdmin(r)
		if !ok {
			writeJSON(w, http.StatusForbidden, SignUpResponse{
				Error: "Only an authenticated admin can update users",
			})
			return
		}

		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, SignUpResponse{
				Error: "Invalid JSON request body",
			})
			return
		}

		// Clean user input before saving it.
		req.FirstName = strings.TrimSpace(req.FirstName)

		if req.FirstName == "" {
			writeJSON(w, http.StatusBadRequest, SignUpResponse{
				Error: "first_name is required",
			})
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		var user User
		if err := a.DB.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeJSON(w, http.StatusNotFound, SignUpResponse{
					Error: "User not found",
				})
				return
			}

			log.Printf("Error finding user %s: %v", id, err)
			writeJSON(w, http.StatusInternalServerError, SignUpResponse{
				Error: "Could not find user",
			})
			return
		}

		// Audit: save which admin made the modification.
		user.FirstName = req.FirstName
		user.UpdatedBy = adminID
		user.UpdatedAt = time.Now()

		if err := a.DB.Save(&user).Error; err != nil {
			log.Printf("Error updating user %s: %v", id, err)
			writeJSON(w, http.StatusInternalServerError, SignUpResponse{
				Error: "Could not update user",
			})
			return
		}

		writeJSON(w, http.StatusOK, SignUpResponse{
			Success: "User updated successfully",
		})
	}
}
