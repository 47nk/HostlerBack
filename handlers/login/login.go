package login

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	// Print the received username
	fmt.Println("Received Username:", loginRequest.Username)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Login successful for username: %s", loginRequest.Username)
}
