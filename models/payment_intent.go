package models

type Payment_intent struct {
	PagoID            int
	QuotePreviewID    int
	Status            string
	TransactionAmount float64
	FechaCreacion     string
}

func (Payment_intent) TableName() string {
	return "payment_intent"
}
