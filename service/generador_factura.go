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

	// 2. Crear una nueva Factura con datos mínimos y constantes
	factura = models.Factura{
		// Campos básicos
		QuotePreviewID: int(ID),
		//RutCliente:    "", // Se establecerá más abajo si viene de otra tabla
		Folio:         fmt.Sprintf("%d", rand.Intn(9000)+1000),
		FechaEmision:  time.Now(),
		TipoDocumento: "Factura Electronica",

		// DATOS DEL EMISOR - Estos se guardan directamente en factura
		RutEmisor:         "76.123.456-7",
		RazonSocialEmisor: "Ferretería Construtem S.A.",
		GiroEmisor:        "Venta al por mayor y menor de artículos de ferreteria",
		DireccionEmisor:   "Av. Siempre Viva 742",
		ComunaEmisor:      "Santiago",
		CiudadEmisor:      "Santiago",
		TelefonoEmisor:    "+56 2 1234 5678",
		EmailEmisor:       "contacto@construtem.cl",

		// DATOS DEL RECEPTOR - Estos vienen de otras celdas/tablas y se guardan directamente
		// TODO: Estos datos deben venir del frontend o de consultas a otras tablas
		RutReceptor:         "12.345.678-9",      //  Viene de otra celda
		RazonSocialReceptor: "Cliente Ejemplo",   //  Viene de otra celda
		GiroReceptor:        "",                  //  Viene de otra celda
		DireccionReceptor:   "Dirección Cliente", //  Viene de otra celda
		ComunaReceptor:      "Comuna Cliente",    //  Viene de otra celda
		CiudadReceptor:      "Ciudad Cliente",    //  Viene de otra celda
		ContactoReceptor:    "Contacto Cliente",  //  Viene de otra celda

		// Datos SII
		TimbreElectronico: "",
		FraseLegal:        "Esta factura es representacion fiel del documento electronico firmado digitalmente segun ley N° 19.799",
	}

	// Detalles de ejemplo
	detalles := []models.DetalleFactura{
		{
			SKU:            "PROD001",          //  Viene de otra celda
			Descripcion:    "Producto Ejemplo", //  Viene de otra celda
			Cantidad:       2.0,                // Viene de otra celda
			PrecioUnitario: 15000.0,            //  Viene de otra celda
			TotalLinea:     30000.0,            //  Calculado
		},
		{
			SKU:            "SERV001",                 //  Viene de otra celda
			Descripcion:    "Servicio de Instalacion", //  Viene de otra celda
			Cantidad:       1.0,                       //  Viene de otra celda
			PrecioUnitario: 25000.0,                   //  Viene de otra celda
			TotalLinea:     25000.0,                   //  Calculado
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

	factura.SubtotalNeto = subtotalNeto //  Calculado
	factura.Iva19 = iva                 //  Calculado
	factura.IvaRetenido = 0             //  Por defecto
	factura.TotalFinal = totalFinal

	// Crear la factura en la base de datos
	if err := utils.DB.Create(&factura).Error; err != nil {
		return nil, fmt.Errorf("error al guardar factura: %w", err)
	}

	return &factura, nil
}

// Helper para formatear montos (ej. 1.234.567)
func FormatMoney(amount float64) string {
	// Implementar lógica de formateo para CLP
	s := fmt.Sprintf("%.0f", amount) // Redondear a entero sin decimales
	n := len(s)
	if n <= 3 {
		return s
	}
	// Añadir puntos para miles
	for i := n - 3; i > 0; i -= 3 {
		s = s[:i] + "." + s[i:]
	}
	return "$ " + s // Con el signo peso
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
