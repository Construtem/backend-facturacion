package models

type PaymentStatus struct {
	Status             string `json:"status"`
	TransactionDetails struct {
		TotalPaidAmount float64 `json:"total_paid_amount"`
	} `json:"transaction_details"`
}

func (PaymentStatus) TableName() string {
	return "paymentstatus"
}
