package dashboard

import (
	"encoding/json"
	"hostlerBackend/app"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetBills(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve query parameters
		var (
			limitStr  = r.URL.Query().Get("limit")
			offsetStr = r.URL.Query().Get("offset")
			bills     []Bill
		)

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}

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

		// Fetch bills from the database based on user ID, limit, and offset
		err = a.DB.
			Where("user_id = ?", userId).
			Limit(limit).Offset(offset).
			Order("billing_month desc").
			Find(&bills).Error
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

func GetTransactions(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			billingMonthStr = r.URL.Query().Get("billing_month")
			transactions    = []Transaction{}
		)

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}

		// Convert query parameters to appropriate types
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = a.DB.Where("bill_id IN (?)",
			a.DB.Model(&Bill{}).
				Select("id").
				Where("user_id = ? AND billing_month = ?", userId, billingMonthStr),
		).Find(&transactions).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Encode and return the bills
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetDueDetails(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			unpaidBills   []Bill
			dueDetails    = DueDetails{}
			currentMonth  = time.Now().Format("200601")
			dueToValueMap = make(map[string]float64)
		)

		// Get user ID from context
		userIdStr, ok := r.Context().Value("user_id").(string)
		if !ok || userIdStr == "" {
			http.Error(w, `{"error": "User ID missing or invalid"}`, http.StatusUnauthorized)
			return
		}

		// Convert query parameters to appropriate types
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = a.DB.
			Where("user_id = ? and payment_status != ?", userId, "complete").
			Order("billing_month desc").
			Find(&unpaidBills).Error
		if err != nil {
			log.Printf("Error finding unpaid bills for user %d: %v", userId, err)
			http.Error(w, "Error finding Bills", http.StatusInternalServerError)
			return
		}

		for i := range unpaidBills {
			switch unpaidBills[i].BillType {
			case "Daily Meal":
				if i == 0 && unpaidBills[i].BillingMonth == currentMonth {
					dueDetails.MealDue += unpaidBills[i].Amount
				}
			default:
				if i == 0 && unpaidBills[i].BillingMonth == currentMonth {
					dueDetails.MiscDue += unpaidBills[i].Amount
				}
			}
			dueToValueMap[unpaidBills[i].BillType] += unpaidBills[i].Amount
		}

		for dueType, dueValue := range dueToValueMap {
			dueDetails.TotalDueSplit = append(dueDetails.TotalDueSplit, TotalDueSplit{DueType: dueType, DueValue: dueValue})
		}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Encode and return the bills
		if err := json.NewEncoder(w).Encode(dueDetails); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
