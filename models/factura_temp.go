package models

import (
	"encoding/json"
	"time"
)

type FacturaTemp struct {
	ID                  uint `gorm:"primaryKey"`
	CotizacionID        int
	QuotePreviewID      int
	RutCliente          string
	TipoDocumento       string
	Folio               string
	FechaEmision        time.Time
	FechaVencimiento    time.Time
	TimbreElectronico   string
	SiiIndicacion       string
	FraseLegal          string
	RutEmisor           string
	RazonSocialEmisor   string
	GiroEmisor          string
	DireccionEmisor     string
	ComunaEmisor        string
	CiudadEmisor        string
	TelefonoEmisor      string
	EmailEmisor         string
	RutReceptor         string
	RazonSocialReceptor string
	GiroReceptor        string
	DireccionReceptor   string
	ComunaReceptor      string
	CiudadReceptor      string
	ContactoReceptor    string
	Items               json.RawMessage `gorm:"type:jsonb"`
	SubtotalNeto        float64
	Iva19               float64
	IvaRetenido         float64
	TotalFinal          float64
	UrlPdf              string
	UrlVerificacion     string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (FacturaTemp) TableName() string {
	return "factura_temp"
}
