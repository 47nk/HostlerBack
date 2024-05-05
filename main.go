package main

import (
	"fmt"
	"net/http"

	"hostlerBackend/db"
	"hostlerBackend/handlers/login"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {

	db, err := db.InitializeDB()
	if err != nil {
		panic(err)
	}
	defer func() {
		// Close the underlying database connection
		dbSQL, err := db.DB()
		if err != nil {
			panic("failed to get database connection")
		}
		dbSQL.Close()
	}()

	fmt.Println("Successfully connected to SQLite database")

	r := mux.NewRouter()
	loginGroup := r.PathPrefix("/handlers").Subrouter()
	loginGroup.HandleFunc("/login", login.Login).Methods("POST")
	// CORS configuration
	corsHandler := cors.Default().Handler(r)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsHandler)
}
