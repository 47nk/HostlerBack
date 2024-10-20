package main

import (
	"hostlerBackend/db"
	"hostlerBackend/handlers/announcement"
	"hostlerBackend/handlers/app"
	"hostlerBackend/handlers/dashboard"
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
	//users group
	userGroup := r.PathPrefix("/users").Subrouter()
	{
		userGroup.HandleFunc("/{id}", login.UpdateUser(app)).Methods("PUT")
		userGroup.HandleFunc("/login", login.Login(app)).Methods("POST")
		userGroup.HandleFunc("/signup", login.SignUp(app)).Methods("POST")
	}

	//announcement group
	announcementGroup := r.PathPrefix("/announcements").Subrouter()
	{
		announcementGroup.HandleFunc("/add", announcement.AddAnnouncement(app)).Methods("POST")
		announcementGroup.HandleFunc("/get", announcement.GetAnnouncements(app)).Methods("GET")
	}
	//dashboard group
	dashboardGroup := r.PathPrefix("/dashboard").Subrouter()
	{
		dashboardGroup.HandleFunc("/transaction", dashboard.CreateTransaction(app)).Methods("POST")

	}

	//cors
	corsHandler := cors.Default().Handler(r)
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsHandler)
}
