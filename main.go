package main

import (
	"hostlerBackend/db"
	"hostlerBackend/handlers/announcement"
	"hostlerBackend/handlers/app"
	"hostlerBackend/handlers/dashboard"
	"hostlerBackend/handlers/login"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func init() {
	// Load the .env file only in local development
	if os.Getenv("RENDER") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found, skipping...")
		}
	}
}
func main() {
	db, err := db.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	app := &app.App{DB: db}

	r := mux.NewRouter()
	r.HandleFunc("/test", login.TestAPI()).Methods("GET")
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
		dashboardGroup.HandleFunc("/get-bills", dashboard.GetBills(app)).Methods("GET")
		dashboardGroup.HandleFunc("/transaction", dashboard.CreateTransaction(app)).Methods("POST")

	}

	//cors
	corsHandler := cors.Default().Handler(r)
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsHandler)
}
