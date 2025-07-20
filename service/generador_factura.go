package services

import (
	"backend-facturacion/models"
	"backend-facturacion/utils"
	"fmt"
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

// GetOrCreateInvoiceData obtiene los datos de una factura existente - NO crea datos de ejemplo
func GetOrCreateInvoiceData(ID uint) (*models.Factura, error) {
	var factura models.Factura

	// Intentar encontrar una Factura existente para esta QuotePreview
	if err := utils.DB.Where("quote_preview_id = ?", ID).First(&factura).Error; err != nil {
		return nil, fmt.Errorf("factura no encontrada para quote_preview_id %d. Debe crearse primero con datos reales desde el sistema de cotizaciones", ID)
	}

	// Completar datos faltantes desde variables de entorno antes de validar
	if err := completeInvoiceDataFromEnv(&factura); err != nil {
		return nil, fmt.Errorf("error al completar datos de factura: %w", err)
	}

	// Validar que los datos sean suficientes para generar PDF
	if err := factura.ValidateForPDF(); err != nil {
		return nil, fmt.Errorf("datos de factura incompletos para quote_preview_id %d: %w", ID, err)
	}

	return &factura, nil
}

// completeInvoiceDataFromEnv completa los campos faltantes con datos del .env
func completeInvoiceDataFromEnv(factura *models.Factura) error {
	config := LoadCompanyConfig()
	updated := false

	// Completar datos del emisor si están vacíos
	if factura.TipoDocumento == "" {
		factura.TipoDocumento = "FACTURA ELECTRÓNICA"
		updated = true
	}

	if factura.RutEmisor == "" {
		factura.RutEmisor = config.RUT
		updated = true
	}

	if factura.RazonSocialEmisor == "" {
		factura.RazonSocialEmisor = config.RazonSocial
		updated = true
	}

	if factura.GiroEmisor == "" {
		factura.GiroEmisor = config.Giro
		updated = true
	}

	if factura.DireccionEmisor == "" {
		factura.DireccionEmisor = config.Direccion
		updated = true
	}

	if factura.ComunaEmisor == "" {
		factura.ComunaEmisor = config.Comuna
		updated = true
	}

	if factura.CiudadEmisor == "" {
		factura.CiudadEmisor = config.Ciudad
		updated = true
	}

	if factura.TelefonoEmisor == "" {
		factura.TelefonoEmisor = config.Telefono
		updated = true
	}

	if factura.EmailEmisor == "" {
		factura.EmailEmisor = config.Email
		updated = true
	}

	// Completar datos fijos
	if factura.TimbreElectronico == "" {
		factura.TimbreElectronico = "TIMBRE_ELECTRONICO"
		updated = true
	}

	if factura.SiiIndicacion == "" {
		factura.SiiIndicacion = "SII - Santiago"
		updated = true
	}

	if factura.FraseLegal == "" {
		factura.FraseLegal = "Esta factura es representacion fiel del documento electronico firmado digitalmente segun ley N° 19.799"
		updated = true
	}

	// Generar folio si está vacío (usando el ID de la factura)
	if factura.Folio == "" {
		factura.Folio = fmt.Sprintf("F%06d", factura.ID)
		updated = true
	}

	// Si no hay fecha de emisión, usar fecha actual
	if factura.FechaEmision.IsZero() {
		factura.FechaEmision = time.Now()
		updated = true
	}

	// Completar giro del receptor si está vacío (valor por defecto)
	if factura.GiroReceptor == "" {
		factura.GiroReceptor = "Comercio"
		updated = true
	}

	// Guardar cambios en la base de datos si hubo actualizaciones
	if updated {
		if err := utils.DB.Save(factura).Error; err != nil {
			return fmt.Errorf("error al guardar datos completados en BD: %w", err)
		}
	}

	return nil
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
