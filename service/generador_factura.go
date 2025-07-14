package services

import (
	"backend-facturacion/models"
	"backend-facturacion/utils"
	"fmt"
	"math/rand"
	"time"
)

// GetInvoiceData obtiene los datos de una factura existente o devuelve un error si no existe
func GetInvoiceData(quotePreviewID uint) (*models.Factura, error) {
	var factura models.Factura

	if err := utils.DB.Where("quote_preview_id = ?", quotePreviewID).First(&factura).Error; err != nil {
		return nil, fmt.Errorf("factura no encontrada para quote_preview_id %d. Debe crearse primero con datos reales", quotePreviewID)
	}

	return &factura, nil
}

// GetOrCreateInvoiceData obtiene o crea los datos completos de la factura
func GetOrCreateInvoiceData(ID uint) (*models.Factura, error) {
	var factura models.Factura

	// 1. Intentar encontrar una Factura existente para esta QuotePreview
	if err := utils.DB.Where("quote_preview_id = ?", ID).First(&factura).Error; err == nil {
		return &factura, nil // Factura encontrada, devolver
	}

	// Generar folio único
	rand.Seed(time.Now().UnixNano())

	// Cargar configuración de empresa desde variables de entorno
	config := LoadCompanyConfig()

	// 2. Crear una nueva Factura con datos de configuración
	factura = models.Factura{
		// Campos básicos
		QuotePreviewID: int(ID),
		Folio:          fmt.Sprintf("%d", rand.Intn(9000)+1000),
		FechaEmision:   time.Now(),
		TipoDocumento:  "Factura Electronica",

		// DATOS DEL EMISOR - Desde configuración
		RutEmisor:         config.RUT,
		RazonSocialEmisor: config.RazonSocial,
		GiroEmisor:        config.Giro,
		DireccionEmisor:   config.Direccion,
		ComunaEmisor:      config.Comuna,
		CiudadEmisor:      config.Ciudad,
		TelefonoEmisor:    config.Telefono,
		EmailEmisor:       config.Email,

		// DATOS DEL RECEPTOR - TODO: Obtener de base de datos real
		RutReceptor:         "abcd1fg-l",       // 12.345.678-9
		RazonSocialReceptor: "Cliente Ejemplo", // Cliente Ejemplo
		GiroReceptor:        "Comercio",
		DireccionReceptor:   "Dirección Cliente",
		ComunaReceptor:      "Comuna Cliente",
		CiudadReceptor:      "Ciudad Cliente",
		ContactoReceptor:    "Contacto Cliente",

		// Datos SII
		TimbreElectronico: "",
		FraseLegal:        "Esta factura es representacion fiel del documento electronico firmado digitalmente segun ley N° 19.799",
	}

	// Detalles de ejemplo
	detalles := []models.DetalleFactura{
		{
			SKU:            "PROD001",
			Descripcion:    "Producto Ejemplo",
			Cantidad:       2.0,
			PrecioUnitario: 15000.0,
			TotalLinea:     30000.0,
		},
		{
			SKU:            "SERV001",
			Descripcion:    "Servicio de Instalacion",
			Cantidad:       1.0,
			PrecioUnitario: 25000.0,
			TotalLinea:     25000.0,
		},
	}

	// Calcular totales basados en los detalles
	var subtotalNeto float64 = 0
	for _, detalle := range detalles {
		subtotalNeto += detalle.TotalLinea
	}

	// Establecer los detalles en el campo JSON
	if err := factura.SetDetalles(detalles); err != nil {
		return nil, fmt.Errorf("error al serializar detalles: %w", err)
	}

	// Calcular totales
	iva := subtotalNeto * 0.19
	totalFinal := subtotalNeto + iva

	factura.SubtotalNeto = subtotalNeto
	factura.Iva19 = iva
	factura.IvaRetenido = 0
	factura.TotalFinal = totalFinal

	// validacion antes de guardar en base de datos
	if err := factura.ValidateForPDF(); err != nil {
		return nil, fmt.Errorf("datos de factura inválidos: %w", err)
	}

	// Crear la factura en la base de datos (solo si pasa validación)
	if err := utils.DB.Create(&factura).Error; err != nil {
		return nil, fmt.Errorf("error al guardar factura: %w", err)
	}

	return &factura, nil
}

// Helper para formatear montos sin símbolo $ (solo números)
func FormatMoneySimple(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	n := len(s)
	if n <= 3 {
		return s
	}
	for i := n - 3; i > 0; i -= 3 {
		s = s[:i] + " " + s[i:] // Usar espacio en lugar de punto
	}
	return s
}

// Helper para formatear fechas (ej. 18 de enero, 2026)
func FormatDateChilean(t time.Time) string {
	meses := []string{
		"enero", "febrero", "marzo", "abril", "mayo", "junio",
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
	}

	return fmt.Sprintf("%d de %s, %d",
		t.Day(),
		meses[t.Month()-1],
		t.Year())
}
