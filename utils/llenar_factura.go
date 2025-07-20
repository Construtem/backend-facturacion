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

// Estructuras actualizadas para mapear la respuesta del endpoint
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

type DireccionCliente struct {
	Direccion string `json:"direccion"`
	Comuna    string `json:"comuna"`
	Ciudad    string `json:"ciudad"`
	Region    string `json:"region"`
}

type CheckoutItemDTO struct {
	SKU        string  `json:"sku"`
	Nombre     string  `json:"nombre"`
	Cantidad   int     `json:"cantidad"`
	PrecioUnit float64 `json:"precio_unit"`
	Subtotal   float64 `json:"subtotal"`
	Descuento  int     `json:"descuento"` // porcentaje entero 0-100
	Sucursal   string  `json:"sucursal"`
}

// DatosFacturaResponse actualizada para coincidir con CheckoutCotizacionResponse
type DatosFacturaResponse struct {
	ID             int               `json:"id"`
	FechaCrea      string            `json:"fecha_crea"`
	Estado         string            `json:"estado"`
	TipoDespacho   string            `json:"tipo_despacho"`
	EstadoPago     string            `json:"estado_pago"`
	Cliente        Cliente           `json:"cliente"`
	Usuario        Usuario           `json:"usuario"`
	Direccion      DireccionCliente  `json:"direccion"`
	Items          []CheckoutItemDTO `json:"items"`
	CostoEnvio     float64           `json:"costo_envio"`
	SubtotalNeto   float64           `json:"subtotal_neto"`
	DescuentoTotal float64           `json:"descuento_total"`
	IVA            float64           `json:"iva"`
	Total          float64           `json:"total"`
	PreviewID      *int              `json:"preview_id,omitempty"`
}

// CrearFactura obtiene datos del endpoint de prueba y crea una factura
func CrearFactura(cotizacionID int, quotePreviewID int) (*models.Factura, error) {
	// URL del endpoint usando variable de entorno
	baseURL := os.Getenv("BACK_VENTAS_URL")
	url := fmt.Sprintf("%sapi/cotizaciones/checkout/%d", baseURL, cotizacionID)

	fmt.Printf("🚀 Realizando petición a: %s\n", url)

	// Realizar la petición HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al llamar al endpoint: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("📊 Status Code: %d\n", resp.StatusCode)

	// Leer la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Mostrar la respuesta en formato JSON pretty (solo si DEBUG está habilitado)
	if os.Getenv("DEBUG_ENDPOINT") == "true" {
		fmt.Println("📋 Respuesta del endpoint:")
		var prettyData interface{}
		if err := json.Unmarshal(body, &prettyData); err == nil {
			if prettyJSON, err := json.MarshalIndent(prettyData, "", "  "); err == nil {
				fmt.Println(string(prettyJSON))
			}
		}
		fmt.Println("═══════════════════════════════════════════════════════════")
	}

	// Decodificar la respuesta JSON
	var datosFactura DatosFacturaResponse
	if err := json.Unmarshal(body, &datosFactura); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON: %v", err)
	}

	// Mostrar los datos decodificados en consola
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("📦 Datos decodificados del endpoint:")
	if prettyJSON, err := json.MarshalIndent(datosFactura, "", "  "); err == nil {
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Printf("Error al mostrar datos decodificados: %v\n", err)
	}
	fmt.Println("═══════════════════════════════════════════════════════════")

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

	// Asignar los datos del endpoint directamente, sin cálculos manuales
	factura.SubtotalNeto = datosFactura.SubtotalNeto
	factura.TotalFinal = datosFactura.Total
	factura.Iva19 = datosFactura.IVA
	// No considerar CostoEnvio para ningún cálculo
	// Guardar el descuento total si el modelo lo soporta
	// factura.DescuentoTotal = datosFactura.DescuentoTotal

	// Campos adicionales que ahora están disponibles
	factura.Estado = datosFactura.EstadoPago

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
