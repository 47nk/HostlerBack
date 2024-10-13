package announcement

import (
	"encoding/json"
	"hostlerBackend/handlers/app"
	"log"
	"net/http"
)

type GetAnnouncementsReq struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func GetAnnouncements(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GetAnnouncementsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result []Announcement
		err := a.DB.Where("is_active = ?", true).Order("created_at desc").Limit(req.Limit).Offset(req.Offset).Find(&result).Error
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
			http.Error(w, "Error fetching announcements", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode result as JSON and send to client
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("Error encoding response to JSON: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}
