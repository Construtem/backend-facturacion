package models

import (
	"encoding/json"
	"time"
)

type Factura struct {
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
	Envio               float64
	Estado              string
	UrlPdf              string
	UrlVerificacion     string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (Factura) TableName() string {
	return "facturas"
}

type DetalleFactura struct {
	SKU            string  `json:"sku"`
	Nombre         string  `json:"nombre"`
	Cantidad       float64 `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
	Sucursal       string  `json:"sucursal"`
	Descuento      float64 `json:"descuento"`
}

// Método para obtener los detalles desde el JSON
func (f *Factura) GetDetalles() ([]DetalleFactura, error) {
	var detalles []DetalleFactura
	if len(f.Items) == 0 {
		return detalles, nil
	}

	err := json.Unmarshal(f.Items, &detalles)
	if err != nil {
		return nil, err
	}

	return detalles, nil
}

// Método para establecer los detalles en el JSON
func (f *Factura) SetDetalles(detalles []DetalleFactura) error {
	data, err := json.Marshal(detalles)
	if err != nil {
		return err
	}

	f.Items = json.RawMessage(data)
	return nil
}
