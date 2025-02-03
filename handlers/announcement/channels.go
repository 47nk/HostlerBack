package announcement

import (
	"encoding/json"
	"hostlerBackend/handlers/app"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func CreateChannel(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Channel
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func GetChannels(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			entities  []Entity
			userIdStr = r.Context().Value("user_id").(string)
		)

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = a.DB.
			Preload("Channels", func(db *gorm.DB) *gorm.DB {
				// Join user_channel table and filter by user_id and active status
				return db.Joins("JOIN user_channel uc ON uc.channel_id = channels.id").
					Where("uc.user_id = ? AND uc.active = true", userId)
			}).
			Where("id IN (?) AND active = true", a.DB.Table("user_entity").Select("entity_id").Where("user_id = ? AND active = true", userId)).
			Find(&entities).Error
		if err != nil {
			log.Printf("Error finding channels for user %d: %v", 1, err)
			http.Error(w, "Error finding channels", http.StatusInternalServerError)
			return
		}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Encode and return the bills
		if err := json.NewEncoder(w).Encode(entities); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
