package utils

import (
	"backend-facturacion/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Estructuras para mapear la respuesta del endpoint
type Cliente struct {
	Rut      string `json:"rut"`
	Nombre   string `json:"nombre"`
	Telefono string `json:"telefono"`
	Email    string `json:"email"`
}

type Usuario struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	RolID  int    `json:"rol_id"`
}

type Direccion struct {
	Direccion string `json:"direccion"`
	Comuna    string `json:"comuna"`
	Ciudad    string `json:"ciudad"`
}

type ItemFactura struct {
	Sku            string  `json:"sku"`
	Nombre         string  `json:"nombre"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
	Sucursal       string  `json:"sucursal"`
}

// DatosFacturaResponse representa la estructura de la respuesta del endpoint
type DatosFacturaResponse struct {
	ID           int           `json:"id"`
	FechaCrea    string        `json:"fecha_crea"`
	Estado       string        `json:"estado"`
	CostoEnvio   float64       `json:"costo_envio"`
	TipoDespacho string        `json:"tipo_despacho"`
	Cliente      Cliente       `json:"cliente"`
	Usuario      Usuario       `json:"usuario"`
	Direccion    Direccion     `json:"direccion"`
	Items        []ItemFactura `json:"items"`
	SubtotalNeto float64       `json:"subtotal_neto"`
	Iva          float64       `json:"iva"`
	Total        float64       `json:"total"`
}

// CrearFactura obtiene datos del endpoint de prueba y crea una factura
func CrearFactura(cotizacionID int, quotePreviewID int) (*models.Factura, error) {
	// URL del endpoint usando variable de entorno
	baseURL := os.Getenv("BACK_VENTAS_URL")
	url := fmt.Sprintf("%sapi/cotizaciones/checkout/%d", baseURL, cotizacionID)

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

	// Convertir los items a json.RawMessage para guardar en la base de datos
	itemsJSON, err := json.Marshal(datosFactura.Items)
	if err != nil {
		return nil, fmt.Errorf("error al convertir items a JSON: %v", err)
	}

	// Crear la factura con los datos obtenidos
	var factura models.Factura

	// Llenar los datos de la factura usando los datos del endpoint
	factura.CotizacionID = cotizacionID     // Usar el parámetro recibido
	factura.QuotePreviewID = quotePreviewID // Usar el parámetro recibido
	factura.RutCliente = datosFactura.Cliente.Rut
	factura.TipoDocumento = "FACTURA"
	factura.RutEmisor = "76543210-9"
	factura.EmailEmisor = datosFactura.Usuario.Email
	factura.RutReceptor = datosFactura.Cliente.Rut
	factura.DireccionReceptor = datosFactura.Direccion.Direccion
	factura.ComunaReceptor = datosFactura.Direccion.Comuna
	factura.CiudadReceptor = datosFactura.Direccion.Ciudad
	factura.ContactoReceptor = datosFactura.Cliente.Telefono
	factura.Items = json.RawMessage(itemsJSON)

	// Usar los datos del endpoint para subtotal y total
	factura.SubtotalNeto = datosFactura.SubtotalNeto
	factura.TotalFinal = datosFactura.Total
	factura.Iva19 = datosFactura.Iva

	// Agregar datos del usuario (si tienes estos campos en el modelo)
	// factura.UsuarioNombre = datosFactura.Usuario.Nombre
	// factura.UsuarioEmail = datosFactura.Usuario.Email
	// factura.UsuarioRolID = datosFactura.Usuario.RolID

	// Campos constantes o valores por defecto
	factura.Folio = fmt.Sprintf("FOLIO%d", datosFactura.ID)
	factura.FechaEmision = time.Now()
	factura.FechaVencimiento = time.Now().AddDate(0, 0, 30)
	factura.TimbreElectronico = "TIMBRE"
	factura.SiiIndicacion = "SII"
	factura.FraseLegal = "Frase legal"
	factura.RazonSocialEmisor = "Empresa S.A."
	factura.GiroEmisor = "Servicios"
	factura.DireccionEmisor = "Calle Falsa 123"
	factura.ComunaEmisor = "Comuna"
	factura.CiudadEmisor = "Ciudad"
	factura.TelefonoEmisor = "123456789"
	factura.RazonSocialReceptor = datosFactura.Cliente.Nombre
	factura.GiroReceptor = "Comercio"
	factura.IvaRetenido = 0
	factura.UrlPdf = ""
	factura.UrlVerificacion = ""
	factura.CreatedAt = time.Now()
	factura.UpdatedAt = time.Now()
	factura.Estado = "pending"

	// Guardar la factura en la base de datos
	db := GetDB()
	if err := db.Create(&factura).Error; err != nil {
		return nil, fmt.Errorf("error al guardar la factura: %v", err)
	}

	return &factura, nil
}
