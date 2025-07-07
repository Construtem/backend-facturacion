package models

import "time"

type Payment_intent struct {
	ID                uint `gorm:"primaryKey"`
	QuotePreviewID    int
	PagoID            int
	Status            string
	TransactionAmount float64
	MetodoPago        string
	EventType         string
	CreatedAt         string
	UpdatedAt         time.Time
}

func (Payment_intent) TableName() string {
	return "payment_intent"
}
