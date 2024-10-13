package announcement

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/handlers/app"
	"net/http"
)

func AddAnnouncement(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Announcement
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("%v", req.Title)

		if err := a.DB.Create(&req).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
