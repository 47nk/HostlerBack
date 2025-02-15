package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// Middleware to validate JWT and authenticate the user
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := parseJWTFromCookie(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// If token is valid, we can add additional checks or extract claims here
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int64(claims["user_id"].(float64))
			role := claims["role"].(string)
			// Store user info in the context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", fmt.Sprint(userID))
			ctx = context.WithValue(ctx, "role", fmt.Sprint(role))
			// Pass the context with user info to the next handler
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}

// Function to parse JWT from the cookie and validate it
func parseJWTFromCookie(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, fmt.Errorf("cookie not found")
	}

	// Validate and parse the JWT token
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method matches (HS256 for symmetric key)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
