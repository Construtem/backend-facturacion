package models

import "time"

// PaymentStatus define los estados válidos de pago.
type PaymentStatus string

const (
	Pending         PaymentStatus = "pending"
	InterimApproval PaymentStatus = "interim_approval"
	Approved        PaymentStatus = "approved"
	Rejected        PaymentStatus = "rejected"
) // modificar para que solo se pueda usar estos estados

// QuotePreview representa una cotización preliminar y su estado agregado de pago.
type QuotePreview struct {
	ID                        uint          `gorm:"primaryKey"`
	CotizacionId              int           `gorm:"not null"`
	IssuedAt                  time.Time     `gorm:"not null"`
	Subtotal                  float64       `gorm:"not null"`
	Tax                       float64       `gorm:"not null"`
	Total                     float64       `gorm:"not null"`
	PaymentStatus             PaymentStatus `gorm:"type:varchar(30);not null;default:'pending'"`
	SuccessfulPaymentIntentID string        `gorm:"not null"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

func (QuotePreview) TableName() string {
	return "quote_previews"
}
