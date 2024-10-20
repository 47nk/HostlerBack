package dashboard

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/handlers/app"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type CreateTransactionReq struct {
	RollNumber      string  `json:"roll_number"`
	TransactionType string  `json:"transaction_type"`
	Items           int     `json:"items"`
	Price           float64 `json:"price"`
	ExtraItems      int     `json:"extra_items"`
	ExtraPrice      float64 `json:"extra_price"`
}

func CreateTransaction(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTransactionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := a.DB.Transaction(func(tx *gorm.DB) error {
			var (
				user         User
				bills        []Bill
				billId       int64
				billingMonth = time.Now().Format("200601") // YYYYMM format
			)

			err := tx.Where("roll_number = ?", req.RollNumber).First(&user).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return fmt.Errorf("user retrieval error: %w", err)
			}
			if user.ID == 0 {
				return fmt.Errorf("user not found")
			}

			// Find pending bills
			err = tx.Where("user_id = ? AND bill_type = ? AND payment_status = ? AND billing_month = ?", user.ID, req.TransactionType, "pending", billingMonth).Find(&bills).Error
			if err != nil {
				return fmt.Errorf("internal error in finding bills: %w", err)
			}

			if len(bills) == 0 {
				// Need to create an entry in bills
				bill := Bill{
					UserId:       user.ID,
					BillType:     req.TransactionType,
					BillingMonth: billingMonth,
				}

				result := tx.Create(&bill)
				if result.Error != nil {
					return fmt.Errorf("internal error in creating bill record: %w", result.Error)
				}
				billId = bill.ID
			} else {
				billId = bills[0].ID
			}

			// Create transaction
			transaction := Transaction{
				BillId:          billId,
				Items:           int64(req.Items),
				Price:           req.Price,
				ExtraItems:      int64(req.ExtraItems),
				ExtraPrice:      req.ExtraPrice,
				TransactionType: req.TransactionType,
			}
			err = tx.Create(&transaction).Error
			if err != nil {
				return fmt.Errorf("internal error in creating transaction: %w", err)
			}

			// Update bill amount
			err = tx.Model(&Bill{ID: billId}).UpdateColumn("amount", gorm.Expr("amount + ?", req.Price)).Error
			if err != nil {
				return fmt.Errorf("internal error in updating bill: %w", err)
			}

			return nil
		})

		if err != nil {
			// Handle the error and return appropriate response to the client
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Success response (if everything went well)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Transaction completed successfully"))
	}
}
