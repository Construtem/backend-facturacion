package models

import "time"

type Cotizacion struct {
	ID           uint      `json:"id"`
	FechaEmision time.Time `json:"fecha_emision"`
	CotizacionID int
	Subtotal     float64 `json:"subtotal"`
	Impuesto     float64 `json:"impuesto"`
	Total        float64 `json:"total"`
	ItemsJSON    string  `json:"-" gorm:"column:items_json"`
}

func (Cotizacion) TableName() string {
	return "cotizaciones"
}
