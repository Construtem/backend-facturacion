package utils

import (
	"backend-facturacion/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// DatosFacturaResponse representa la estructura de la respuesta del endpoint
type DatosFacturaResponse struct {
	CotizacionID      int             `json:"cotizacion_id"`
	RutCliente        string          `json:"rut_cliente"`
	TipoDocumento     string          `json:"tipo_documento"`
	RutEmisor         string          `json:"rut_emisor"`
	EmailEmisor       string          `json:"email_emisor"`
	RutReceptor       string          `json:"rut_receptor"`
	DireccionReceptor string          `json:"direccion_receptor"`
	ComunaReceptor    string          `json:"comuna_receptor"`
	CiudadReceptor    string          `json:"ciudad_receptor"`
	ContactoReceptor  string          `json:"contacto_receptor"`
	Items             json.RawMessage `json:"items"`
}

// CrearFactura obtiene datos del endpoint de prueba y crea una factura
func CrearFactura(quotePreviewID int) (*models.Factura, error) {
	// URL del endpoint de prueba (ajusta según tu configuración)
	url := "http://localhost:8080/api/prueba-ventas"

	// Realizar la petición HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al llamar al endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Decodificar la respuesta JSON
	var datosFactura DatosFacturaResponse
	if err := json.Unmarshal(body, &datosFactura); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON: %v", err)
	}

	// Crear la factura con los datos obtenidos
	var factura models.Factura

	// Llenar los datos de la factura
	factura.CotizacionID = datosFactura.CotizacionID
	factura.QuotePreviewID = quotePreviewID
	factura.RutCliente = datosFactura.RutCliente
	factura.TipoDocumento = datosFactura.TipoDocumento
	factura.RutEmisor = datosFactura.RutEmisor
	factura.EmailEmisor = datosFactura.EmailEmisor
	factura.RutReceptor = datosFactura.RutReceptor
	factura.DireccionReceptor = datosFactura.DireccionReceptor
	factura.ComunaReceptor = datosFactura.ComunaReceptor
	factura.CiudadReceptor = datosFactura.CiudadReceptor
	factura.ContactoReceptor = datosFactura.ContactoReceptor
	factura.Items = datosFactura.Items

	// Obtener una referencia a QuotePreview para el subtotal y total
	db := GetDB()
	var preview models.QuotePreview
	if err := db.First(&preview, quotePreviewID).Error; err != nil {
		return nil, fmt.Errorf("error al buscar QuotePreview: %v", err)
	}

	factura.SubtotalNeto = preview.Subtotal
	factura.TotalFinal = preview.Total

	// Campos constantes o inventados
	factura.Folio = "FOLIO123"
	factura.FechaEmision = time.Now()
	factura.FechaVencimiento = time.Now().AddDate(0, 0, 30) // 30 días de vencimiento
	factura.TimbreElectronico = "TIMBRE"
	factura.SiiIndicacion = "SII"
	factura.FraseLegal = "Frase legal"
	factura.RazonSocialEmisor = "Empresa S.A."
	factura.GiroEmisor = "Servicios"
	factura.DireccionEmisor = "Calle Falsa 123"
	factura.ComunaEmisor = "Comuna"
	factura.CiudadEmisor = "Ciudad"
	factura.TelefonoEmisor = "123456789"
	factura.RazonSocialReceptor = "Cliente S.A."
	factura.GiroReceptor = "Comercio"
	factura.Iva19 = factura.SubtotalNeto * 0.19
	factura.IvaRetenido = 0
	factura.UrlPdf = ""
	factura.UrlVerificacion = ""
	factura.CreatedAt = time.Now()
	factura.UpdatedAt = time.Now()
	factura.Estado = "pending"

	// Guardar la factura en la base de datos
	if err := db.Create(&factura).Error; err != nil {
		return nil, fmt.Errorf("error al guardar la factura: %v", err)
	}

	return &factura, nil
}
