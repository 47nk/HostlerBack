package dashboard

import (
	"encoding/json"
	"fmt"
	"hostlerBackend/app"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

const YYYYMM = "200601"

func CreateTransaction(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTransactionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Username == "" {
			http.Error(w, `{"error": "invalid user!"}`, http.StatusBadRequest)
			return
		}
		if req.Items <= 0 {
			http.Error(w, `{"error": "insufficient item quantity!"}`, http.StatusBadRequest)
			return
		}
		err := a.DB.Transaction(func(tx *gorm.DB) error {
			var (
				user         User
				bills        []Bill
				billId       int64
				billingMonth = time.Now().Format(YYYYMM)
			)

			err := tx.Where("username = ?", req.Username).First(&user).Error
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
				Items:           req.Items,
				Price:           req.Price,
				ExtraItems:      req.ExtraItems,
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
			log.Printf("Internal Error Creating Transaction: %v", err.Error())
			http.Error(w, `{"error": "Internal Error Creating Transaction"}`, http.StatusInternalServerError)
			return
		}
		// Success response (if everything went well)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Transaction completed successfully"))
	}
}
