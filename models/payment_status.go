package models

type PaymentResp struct {
	ID                 int     `json:"id"`
	Status             string  `json:"status"`
	StatusDetail       string  `json:"status_detail"`
	TransactionAmount  float64 `json:"transaction_amount"`
	TransactionDetails struct {
		TotalPaidAmount float64 `json:"total_paid_amount"`
	} `json:"transaction_details"`
	DateCreated     string `json:"date_created"`
	PaymentMethodID string `json:"payment_method_id"`
}

func (PaymentResp) TableName() string {
	return "paymentstatus"
}
