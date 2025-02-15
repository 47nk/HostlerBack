package announcement

import (
	"encoding/json"
	"errors"
	"hostlerBackend/app"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func CreateChannel(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req    CreateChannelReq
			entity Entity
		)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		//validations
		if req.EntityId <= 0 {
			http.Error(w, "Invalid Entity", http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "Please Provide a Channel Name", http.StatusBadRequest)
			return
		}
		if req.Type == "" {
			http.Error(w, "Please Choose type of Channel", http.StatusBadRequest)
			return
		}

		//check if entity exists
		err = a.DB.Preload("Channels", func(db *gorm.DB) *gorm.DB {
			return db.Where("entity_id = ?", req.EntityId)
		}).Where("id = ? AND active = true", req.EntityId).First(&entity).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Entity not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error retrieving entity details", http.StatusInternalServerError)
			}
			return
		}

		//avoiding channels with same name in entity
		for _, channel := range entity.Channels {
			if channel.Name == req.Name && channel.Active {
				http.Error(w, "Channel with this Name already Exists", http.StatusConflict)
				return
			}
		}

		//create channel
		err = a.DB.Create(&Channel{
			CreatedBy:   uint(userId),
			UpdatedBy:   uint(userId),
			EntityID:    req.EntityId,
			Name:        req.Name,
			Type:        req.Type,
			Description: req.Description,
		}).Error
		if err != nil {
			http.Error(w, "Internal Error Creating Channel", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Channel created successfully"})
	}
}

func GetChannels(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			entities []Entity
		)

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}

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
