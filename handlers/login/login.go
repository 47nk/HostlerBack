package login

import (
	"encoding/json"
	"hostlerBackend/handlers/app"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

		if req.Username == "" || req.Password == "" {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		var user []User
		err := a.DB.Where("roll_number = ?", req.Username).Find(&user).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//validate password

		tokenString, err := GenerateJWT(req.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &cookie)
		w.Write([]byte("Login successful"))
	}
}
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
