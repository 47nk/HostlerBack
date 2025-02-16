package announcement

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/app"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	storage "github.com/supabase-community/storage-go"
	"gorm.io/gorm"
)

var (
	channelClients = make(map[int]map[*SSEClient]bool) // Channel-specific clients
	clientsMu      sync.Mutex
)

func AddAnnouncement(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form data (max 10MB per file)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, `{"error": "Failed to parse form data"}`, http.StatusBadRequest)
			return
		}

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}

		// Convert user ID to integer
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil || userId <= 0 {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusUnauthorized)
			return
		}

		title := r.FormValue("title")
		type_ := r.FormValue("type")
		description := r.FormValue("description")
		channelIdStr := r.FormValue("channel_id")

		// Validate input fields
		if title == "" {
			http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
			return
		}
		if type_ == "" {
			http.Error(w, `{"error": "Type is required"}`, http.StatusBadRequest)
			return
		}
		if description == "" {
			http.Error(w, `{"error": "Description is required"}`, http.StatusBadRequest)
			return
		}
		if channelIdStr == "" {
			http.Error(w, `{"error": "Channel ID is required"}`, http.StatusBadRequest)
			return
		}

		channelId, err := strconv.Atoi(channelIdStr)
		if err != nil || channelId <= 0 {
			http.Error(w, `{"error": "Invalid channel ID"}`, http.StatusBadRequest)
			return
		}

		var attachments []AnnouncementAttachment
		files := r.MultipartForm.File["attachments"]
		for _, fileHeader := range files {
			attachment, err := uploadAnnouncementAttachment(a, fileHeader)
			if err != nil {
				http.Error(w, `{"error": "Failed to upload attachments"}`, http.StatusBadRequest)
				return
			}
			attachments = append(attachments, attachment)
		}

		// Create the announcement record
		announcement := Announcement{
			Title:       title,
			Type:        type_,
			Description: description,
			ChannelId:   channelId,
			CreatedBy:   uint(userId),
			UpdatedBy:   uint(userId),
			Attachments: attachments,
		}

		err = a.DB.Transaction(func(tx *gorm.DB) error {
			// Create the announcement
			if err := tx.Create(&announcement).Error; err != nil {
				log.Printf("Error creating announcement: %v", err)
				return err
			}

			// Fetch the announcement with the creator details
			if err := tx.Preload("Creator").First(&announcement, announcement.ID).Error; err != nil {
				log.Printf("Error Finding Creator: %v", err)
				return err
			}

			return nil
		})

		if err != nil {
			http.Error(w, `{"error": "Failed to process announcement"}`, http.StatusInternalServerError)
			return
		}

		// Broadcast to users
		broadcastAnnouncement(announcement)

		// Success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Announcement created successfully"})
	}
}

func uploadAnnouncementAttachment(a *app.App, fileHeader *multipart.FileHeader) (AnnouncementAttachment, error) {
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return AnnouncementAttachment{}, err
	}
	defer file.Close()

	// Read first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		log.Printf("Failed to read file: %v", err)
		return AnnouncementAttachment{}, err
	}
	detectedType := http.DetectContentType(buffer)

	// Reset file reader after detection
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Printf("Failed to reset file: %v", err)
		return AnnouncementAttachment{}, err
	}

	allowedTypes := map[string]bool{
		// Images
		"image/jpeg": true,
		"image/png":  true,

		// Videos
		"video/mp4": true,

		// Text & CSV
		"text/csv":   true,
		"text/plain": true,

		// Archives
		"application/zip": true, // XLSX, DOCX, PPTX are ZIP-based formats

		// PDF
		"application/pdf": true,

		// Microsoft Word
		"application/msword": true, // DOC
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // DOCX

		// Microsoft Excel
		"application/vnd.ms-excel": true, // XLS
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true, // XLSX

		// Microsoft PowerPoint
		"application/vnd.ms-powerpoint":                                             true, // PPT
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // PPTX
	}

	if !allowedTypes[detectedType] {
		log.Printf("Invalid file type: %v", err)
		return AnnouncementAttachment{}, err
	}

	// Sanitize filename
	sanitizedFilename := filepath.Base(fileHeader.Filename)
	sanitizedFilename = strings.ReplaceAll(sanitizedFilename, " ", "_") // Replace spaces with underscores
	sanitizedFilename = strings.ReplaceAll(sanitizedFilename, " ", "_") // Replace non-breaking spaces with underscores
	sanitizedFilename = strings.ToLower(sanitizedFilename)
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
		return AnnouncementAttachment{}, err
	}

	fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucketName, fileName)
	response := AnnouncementAttachment{
		FileType: detectedType,
		FilePath: fileURL,
		FileSize: fileHeader.Size,
	}

	return response, nil

}

func AnnouncementsSSE(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Get channel_id from query parameters
		vars := mux.Vars(r)
		channelIDStr := vars["channel_id"]
		channelID, err := strconv.Atoi(channelIDStr)
		if err != nil || channelID <= 0 {
			http.Error(w, `{"error": "Invalid channel ID"}`, http.StatusBadRequest)
			return
		}

		// Create a new SSE client
		client := &SSEClient{Channel: make(chan string, 10)}

		// Register the client to the specific channel
		clientsMu.Lock()
		if _, exists := channelClients[channelID]; !exists {
			channelClients[channelID] = make(map[*SSEClient]bool)
		}
		channelClients[channelID][client] = true
		clientsMu.Unlock()

		// Remove client when disconnected
		defer func() {
			clientsMu.Lock()
			delete(channelClients[channelID], client)
			if len(channelClients[channelID]) == 0 {
				delete(channelClients, channelID) // Cleanup if no clients
			}
			clientsMu.Unlock()
			close(client.Channel)
		}()

		// Keep connection alive
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-client.Channel:
				if !ok {
					return // Channel closed, stop listening
				}
				fmt.Fprintf(w, "data: %s\n\n", msg)
				w.(http.Flusher).Flush() // Push update to client

			case <-ticker.C:
				fmt.Fprintf(w, ":\n\n") // Heartbeat
				w.(http.Flusher).Flush()
			}
		}
	}
}

func broadcastAnnouncement(announcement Announcement) {
	jsonData, err := json.Marshal(announcement)
	if err != nil {
		log.Println("Error encoding announcement:", err)
		return
	}

	// Send only to users subscribed to the specific channel
	clientsMu.Lock()
	for client := range channelClients[announcement.ChannelId] {
		client.Channel <- string(jsonData)
	}
	clientsMu.Unlock()
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
