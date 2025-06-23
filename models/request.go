package models

type Request struct {
	TransactionAmount float64 `json:"transaction_amount"`
	PaymentMethodID   string  `json:"payment_method_id"`
	Email             string  `json:"email"`
	CardToken         string  `json:"card_token"`
	CotizacionID      int     `json:"id"`
}

func (Request) TableName() string {
	return "request"
}
