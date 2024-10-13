package main

import (
	"hostlerBackend/db"
	"hostlerBackend/handlers/app"
	"hostlerBackend/handlers/login"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db, err := db.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	app := &app.App{DB: db}

	r := mux.NewRouter()
	userGroup := r.PathPrefix("/users").Subrouter()
	userGroup.HandleFunc("/{id}", login.UpdateUser(app)).Methods("PUT")
	userGroup.HandleFunc("/complex-query", login.ComplexQuery(app)).Methods("GET")
	userGroup.HandleFunc("/login", login.Login(app)).Methods("POST")
	corsHandler := cors.Default().Handler(r)
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsHandler)
}
