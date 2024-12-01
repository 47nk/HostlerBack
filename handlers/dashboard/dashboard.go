package dashboard

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/handlers/app"
	"net/http"
	"strconv"
)

func GetBills(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve query parameters
		var (
			userIdStr = r.URL.Query().Get("user_id")
			limitStr  = r.URL.Query().Get("limit")
			offsetStr = r.URL.Query().Get("offset")
		)

		// Convert query parameters to appropriate types
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 10
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		var bills []Bill
		fmt.Print(userIdStr)

		// Fetch bills from the database based on user ID, limit, and offset
		err = a.DB.Where("user_id = ?", userId).Limit(limit).Offset(offset).Order("billing_month desc").Find(&bills).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Encode and return the bills
		if err := json.NewEncoder(w).Encode(bills); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
