package main

import (
	"hostlerBackend/app"
	"hostlerBackend/auth"
	"hostlerBackend/db"
	"hostlerBackend/handlers/announcement"
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
		userGroup.HandleFunc("/login", login.Login(app)).Methods("POST")
		userGroup.HandleFunc("/{id}", auth.JWTMiddleware(login.UpdateUser(app))).Methods("PUT")
		userGroup.HandleFunc("/signup", auth.JWTMiddleware(login.SignUp(app))).Methods("POST")
	}

	//announcement group
	announcementGroup := r.PathPrefix("/announcements").Subrouter()
	{
		announcementGroup.HandleFunc("/add-announcement", auth.JWTMiddleware(announcement.AddAnnouncement(app))).Methods("POST")
		announcementGroup.HandleFunc("/get-announcements", auth.JWTMiddleware(announcement.GetAnnouncements(app))).Methods("GET")
		announcementGroup.HandleFunc("/add-channel", auth.JWTMiddleware(announcement.CreateChannel(app))).Methods("POST")
		announcementGroup.HandleFunc("/get-channels", auth.JWTMiddleware(announcement.GetChannels(app))).Methods("GET")
	}

	//dashboard group
	dashboardGroup := r.PathPrefix("/dashboard").Subrouter()
	{
		dashboardGroup.HandleFunc("/get-transactions", auth.JWTMiddleware(dashboard.GetTransactions(app))).Methods("GET")
		dashboardGroup.HandleFunc("/get-bills", auth.JWTMiddleware(dashboard.GetBills(app))).Methods("GET")
		dashboardGroup.HandleFunc("/get-dues", auth.JWTMiddleware(dashboard.GetDueDetails(app))).Methods("GET")
		dashboardGroup.HandleFunc("/create-transaction", auth.JWTMiddleware(dashboard.CreateTransaction(app))).Methods("POST")

	}

	//cors
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5317"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", corsHandler)
}
