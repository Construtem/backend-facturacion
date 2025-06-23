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
	ID            uint          `gorm:"primaryKey"`
	IssuedAt      time.Time     `gorm:"not null"`
	Subtotal      float64       `gorm:"not null"`
	Tax           float64       `gorm:"not null"`
	Total         float64       `gorm:"not null"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(30);not null;default:'pending'"`

	// Relación 1:N con PaymentIntent (histórico de webhooks/eventos)
	PaymentIntents []PaymentIntent `gorm:"foreignKey:QuotePreviewID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// PaymentIntent almacena cada evento recibido de la pasarela de pago.
type PaymentIntent struct {
	ID             uint   `gorm:"primaryKey"`
	QuotePreviewID uint   `gorm:"not null;index"`
	ProviderID     string `gorm:"size:100;not null"`
	EventType      string `gorm:"size:50;not null"`
	Status         string `gorm:"size:30;not null"`
	Payload        string `gorm:"type:jsonb"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
