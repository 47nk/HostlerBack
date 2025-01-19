package dashboard

import (
	"encoding/json"
	"hostlerBackend/handlers/app"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetBills(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve query parameters
		var (
			userIdStr = r.URL.Query().Get("user_id")
			limitStr  = r.URL.Query().Get("limit")
			offsetStr = r.URL.Query().Get("offset")
			bills     []Bill
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

func GetTransactions(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			userIdStr       = r.URL.Query().Get("user_id")
			billingMonthStr = r.URL.Query().Get("billing_month")
			transactions    = []Transaction{}
		)

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

type DueDetails struct {
	TotalDue TotalDue
	MealDue  float64 `json:"meal_due"`
	MiscDue  float64 `json:"misc_due"`
}
type TotalDue struct {
	MiscTotalDue float64 `json:"misc_total_due"`
	MealTotalDue float64 `json:"meal_total_due"`
}

func GetDueDetails(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			userIdStr    = r.URL.Query().Get("user_id")
			unpaidBills  []Bill
			dueDetails   = DueDetails{}
			currentMonth = time.Now().Format("200601")
		)

		// Convert query parameters to appropriate types
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = a.DB.Where("user_id = ? and payment_status != ?", userId, "complete").Order("billing_month desc").Find(&unpaidBills).Error
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
				dueDetails.TotalDue.MealTotalDue += unpaidBills[i].Amount

			default:
				if i == 0 && unpaidBills[i].BillingMonth == currentMonth {
					dueDetails.MiscDue += unpaidBills[i].Amount
				}
				dueDetails.TotalDue.MiscTotalDue += unpaidBills[i].Amount
			}
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
