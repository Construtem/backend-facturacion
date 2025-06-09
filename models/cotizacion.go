package models

import "time"

type Cotizacion struct {
	ID           uint      `json:"id"`
	FechaEmision time.Time `json:"fecha_emision"`
	Subtotal     float64   `json:"subtotal"`
	Impuesto     float64   `json:"impuesto"`
	Total        float64   `json:"total"`
	ItemsJSON    string    `json:"-" gorm:"column:items_json"`
}
