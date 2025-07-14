package handlers

import (
	"github.com/gin-gonic/gin"
)

// Estructura para los items de la factura
type ItemFactura struct {
	Sku            string  `json:"sku"`
	Cantidad       int     `json:"cantidad"`
	Descripcion    string  `json:"descripcion"`
	TotalLinea     float64 `json:"total_linea"`
	RecargoPorc    float64 `json:"recargo_porc"`
	DescuentoPorc  float64 `json:"descuento_porc"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

// Estructura para la respuesta completa
type DatosFactura struct {
	CotizacionID      int           `json:"cotizacion_id"`
	RutCliente        string        `json:"rut_cliente"`
	TipoDocumento     string        `json:"tipo_documento"`
	RutEmisor         string        `json:"rut_emisor"`
	EmailEmisor       string        `json:"email_emisor"`
	RutReceptor       string        `json:"rut_receptor"`
	DireccionReceptor string        `json:"direccion_receptor"`
	ComunaReceptor    string        `json:"comuna_receptor"`
	CiudadReceptor    string        `json:"ciudad_receptor"`
	ContactoReceptor  string        `json:"contacto_receptor"`
	Items             []ItemFactura `json:"items"`
}

func PruebaVentas(c *gin.Context) {
	// Crear datos estáticos para la factura
	items := []ItemFactura{
		{
			Sku:            "ABC123",
			Cantidad:       5,
			Descripcion:    "Producto A",
			TotalLinea:     10000,
			RecargoPorc:    0,
			DescuentoPorc:  0,
			PrecioUnitario: 2000,
		},
		{
			Sku:            "XYZ456",
			Cantidad:       1,
			Descripcion:    "Servicio de Instalacion",
			TotalLinea:     5000,
			RecargoPorc:    0,
			DescuentoPorc:  0,
			PrecioUnitario: 5000,
		},
		{
			Sku:            "DEF789",
			Cantidad:       2,
			Descripcion:    "Articulo C",
			TotalLinea:     7500,
			RecargoPorc:    5,
			DescuentoPorc:  0,
			PrecioUnitario: 3750,
		},
	}

	// Crear la respuesta completa
	datosFactura := DatosFactura{
		CotizacionID:      1,
		RutCliente:        "11111111-1",
		TipoDocumento:     "NULO",
		RutEmisor:         "22222222-2",
		EmailEmisor:       "prueba@prueba.com",
		RutReceptor:       "33333333-3",
		DireccionReceptor: "HOLA MUNDO",
		ComunaReceptor:    "HOLA",
		CiudadReceptor:    "MUNDO",
		ContactoReceptor:  "HOLAMUNDO",
		Items:             items,
	}

	// Devolver los datos en formato JSON
	c.JSON(200, datosFactura)
}
