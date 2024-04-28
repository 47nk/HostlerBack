package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Numbers struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

type SumResponse struct {
	Sum int `json:"sum"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		var numbers Numbers
		if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sum := Numbers.Num1 + Numbers.Num2
		sumResponse := SumResponse{Sum: sum}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sumResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	// CORS configuration
	corsHandler := cors.Default().Handler(r)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsHandler)
}
