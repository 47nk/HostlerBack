package announcement

import (
	"encoding/json"
	"hostlerBackend/handlers/app"
	"log"
	"net/http"
	"strconv"
)

func AddAnnouncement(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req       Announcement
			userIdStr = r.Context().Value("user_id").(string)
		)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil || userId <= 0 {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		//validations
		if req.Title == "" {
			http.Error(w, "Title is Required!", http.StatusUnauthorized)
			return
		}
		if req.Type == "" {
			http.Error(w, "Type is Required!", http.StatusUnauthorized)
			return
		}
		if req.Description == "" {
			http.Error(w, "Description is Required!", http.StatusUnauthorized)
			return
		}
		if req.ChannelId <= 0 {
			http.Error(w, "Channel Id is Required!", http.StatusUnauthorized)
			return
		}

		//create announcement
		req.CreatedBy = uint(userId)
		req.UpdatedBy = uint(userId)
		if err := a.DB.Create(&req).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Announcement created successfully"})
	}
}

func GetAnnouncements(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			channelIdStr = r.URL.Query().Get("channel_id")
			limitStr     = r.URL.Query().Get("limit")
			offsetStr    = r.URL.Query().Get("offset")
			result       []Announcement
		)

		//input validation
		channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
		if err != nil || channelId <= 0 {
			http.Error(w, "Invalid Channel Id!", http.StatusUnauthorized)
			return
		}
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid limit!", http.StatusUnauthorized)
			return
		}
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid offset!", http.StatusUnauthorized)
			return
		}

		//find announcements
		err = a.DB.Where("channel_id = ? and active = true", channelId).Order("created_at desc").Limit(int(limit)).Offset(int(offset)).Find(&result).Error
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
			http.Error(w, "Error fetching announcements", http.StatusInternalServerError)
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
