package login

import (
	"encoding/json"
	"hostlerBackend/handlers/app"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomClaims struct {
	Role   string `json:"role"`
	UserID int64  `json:"user_id"`
	jwt.RegisteredClaims
}

func Login(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.Password == "" {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		//fetch user
		var user []User
		err := a.DB.
			Preload("UserRole").
			Where("username = ?", req.Username).
			Find(&user).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(user) == 0 {
			http.Error(w, "OOPS! User Not Found", http.StatusInternalServerError)
			return
		}

		//validate password
		err = bcrypt.CompareHashAndPassword([]byte(user[0].Password), []byte(req.Password))
		if err != nil {
			http.Error(w, "OOPS! Wrong Password", http.StatusInternalServerError)
			return
		}

		tokenString, err := GenerateJWT(user[0].ID, user[0].UserRole.Role)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(w, &cookie)

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(user[0]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GenerateJWT(userId int64, role string) (string, error) {
	claims := CustomClaims{
		UserID: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//sign token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
