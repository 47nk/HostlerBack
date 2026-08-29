package dashboard

type CreateTransactionReq struct {
	Username        string  `json:"username"`
	TransactionType string  `json:"transaction_type"`
	Items           int64   `json:"items"`
	Price           float64 `json:"price"`
	ExtraItems      int64   `json:"extra_items"`
	ExtraPrice      float64 `json:"extra_price"`
}

type DueDetails struct {
	TotalDueSplit []TotalDueSplit
	MealDue       float64 `json:"meal_due"`
	MiscDue       float64 `json:"misc_due"`
}

type TotalDueSplit struct {
	DueType  string  `json:"due_type"`
	DueValue float64 `json:"due_value"`
}
