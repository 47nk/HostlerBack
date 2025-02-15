package announcement

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/app"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	storage "github.com/supabase-community/storage-go"
)

func AddAnnouncement(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Announcement

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}

		// Decode request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		// Convert user ID to integer
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil || userId <= 0 {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusUnauthorized)
			return
		}

		// Validate input fields
		if req.Title == "" {
			http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
			return
		}
		if req.Type == "" {
			http.Error(w, `{"error": "Type is required"}`, http.StatusBadRequest)
			return
		}
		if req.Description == "" {
			http.Error(w, `{"error": "Description is required"}`, http.StatusBadRequest)
			return
		}
		if req.ChannelId <= 0 {
			http.Error(w, `{"error": "Channel ID is required"}`, http.StatusBadRequest)
			return
		}

		// Set creator and updater IDs
		req.CreatedBy = uint(userId)
		req.UpdatedBy = uint(userId)

		// Create the announcement
		if err := a.DB.Create(&req).Error; err != nil {
			log.Printf("Error creating announcement: %v", err)
			http.Error(w, `{"error": "Failed to create announcement"}`, http.StatusInternalServerError)
			return
		}

		// Success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Announcement created successfully"})
	}
}

func UploadAnnouncementAttachment(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Enforce max request size
		maxFileSize := int64(10 << 20) // 10MB
		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read first 512 bytes for MIME detection
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		detectedType := http.DetectContentType(buffer)

		// Reset file reader after detection
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			http.Error(w, "Failed to reset file", http.StatusInternalServerError)
			return
		}

		allowedTypes := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"video/mp4":       true,
			"text/csv":        true,
			"text/plain":      true,
			"application/zip": true, // XLSX is a ZIP archive
		}
		if !allowedTypes[detectedType] {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		// Sanitize filename
		sanitizedFilename := filepath.Base(fileHeader.Filename)
		fileName := fmt.Sprintf("%s_%s", uuid.New().String(), sanitizedFilename)

		// Upload to Supabase
		bucketName := "announcements"
		supabaseStorageURL := fmt.Sprintf("%s/storage/v1", os.Getenv("SUPABASE_URL"))
		client := storage.NewClient(supabaseStorageURL, os.Getenv("SUPABASE_KEY"), nil)
		_, err = client.UploadFile(bucketName, fileName, file, storage.FileOptions{
			ContentType: &detectedType,
		})
		if err != nil {
			log.Printf("Upload failed: %v", err)
			http.Error(w, "Upload failed", http.StatusInternalServerError)
			return
		}

		fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucketName, fileName)
		response := map[string]string{
			"message":  "Upload successful",
			"file_url": fileURL,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
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

		// Validate channel_id
		channelId, err := strconv.ParseInt(channelIdStr, 10, 64)
		if err != nil || channelId <= 0 {
			http.Error(w, "Invalid channel_id! Must be a positive integer.", http.StatusBadRequest)
			return
		}

		// Set default values for limit and offset if not provided
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit <= 0 {
			limit = 10 // Default limit
		}

		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || offset < 0 {
			offset = 0 // Default offset
		}

		// Fetch announcements with relationships
		err = a.DB.
			Preload("Creator").
			Preload("Attachments", "active = ?", true).
			Where("channel_id = ? AND active = true", channelId).
			Order("created_at DESC").
			Limit(int(limit)).
			Offset(int(offset)).
			Find(&result).Error
		if err != nil {
			log.Printf("Error fetching announcements: %v", err)
			http.Error(w, "Failed to fetch announcements", http.StatusInternalServerError)
			return
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}
