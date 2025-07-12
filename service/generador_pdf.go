package services

import (
	"fmt"
	"io"
	"log"
	"os"

	"backend-facturacion/models"

	"github.com/jung-kurt/gofpdf"
)

// Función para generar el PDF de la factura
func GenerateInvoicePDF(factura *models.Factura, writer io.Writer) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Arial", "", 10)
	pdf.AddPage()
	pdf.SetAutoPageBreak(false, 0)

	// Constantes para control de paginación
	const (
		rowHeight    = 8.0
		pageHeight   = 297.0
		marginBottom = 80.0 //120 es el margen inferior
		tableX       = 20.0
	)

	// Función interna para crear encabezado
	addHeader := func(includeClientInfo bool) {
		initialY := pdf.GetY()

		// Recuadro rojo con información de factura
		pdf.SetDrawColor(255, 0, 0)
		pdf.SetFillColor(255, 255, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetLineWidth(1)

		rectRightX := 120.0
		rectRightY := 5.0
		rectRightWidth := 80.0
		rectRightHeight := 30.0

		pdf.Rect(rectRightX, rectRightY, rectRightWidth, rectRightHeight, "D")

		// Contenido del recuadro rojo
		pdf.SetTextColor(255, 0, 0)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(rectRightX+5, rectRightY+5)
		pdf.CellFormat(rectRightWidth-10, 5, tr(fmt.Sprintf("R.U.T.: %s", factura.RutEmisor)), "", 0, "C", false, 0, "")

		pdf.SetFont("Arial", "B", 12)
		pdf.SetXY(rectRightX+5, rectRightY+10)
		pdf.MultiCell(rectRightWidth-10, 5, tr("FACTURA ELECTRONICA"), "", "C", false)

		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(rectRightX+5, rectRightY+18)
		pdf.CellFormat(rectRightWidth-10, 5, tr(fmt.Sprintf("N° %s", factura.Folio)), "", 0, "C", false, 0, "")

		pdf.SetDrawColor(0, 0, 0)

		// Datos del emisor
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "B", 14)
		pdf.SetXY(15, initialY+5)
		pdf.Cell(0, 6, tr(factura.RazonSocialEmisor))
		pdf.Ln(6)

		pdf.SetFont("Arial", "", 9)
		pdf.SetX(15)
		pdf.Cell(0, 4, tr(factura.GiroEmisor))
		pdf.Ln(4)
		pdf.SetX(15)
		pdf.Cell(0, 4, tr(fmt.Sprintf("Dirección: %s, %s", factura.DireccionEmisor, factura.ComunaEmisor)))
		pdf.Ln(4)
		pdf.SetX(15)
		pdf.Cell(0, 4, tr(fmt.Sprintf("Email: %s | Tel: %s", factura.EmailEmisor, factura.TelefonoEmisor)))
		pdf.Ln(7)
		pdf.SetX(135)
		pdf.Cell(0, 4, tr(fmt.Sprintf("Fecha de Emisión: %s", FormatDateChilean(factura.FechaEmision))))
		pdf.Ln(6)

		// Información del cliente (solo en primera página)
		if includeClientInfo {
			currentY := pdf.GetY()
			pdf.SetFillColor(255, 102, 0)
			pdf.Rect(10, currentY, 190, 25, "F")

			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("Arial", "B", 10)

			pdf.SetXY(15, currentY+5)
			pdf.Cell(0, 5, tr(fmt.Sprintf("Cliente: %s", factura.RazonSocialReceptor)))

			pdf.SetXY(100, currentY+5)
			pdf.Cell(0, 5, tr(fmt.Sprintf("Dirección: %s, %s", factura.DireccionReceptor, factura.ComunaReceptor)))

			pdf.SetXY(15, currentY+10)
			pdf.Cell(0, 5, tr(fmt.Sprintf("RUT Cliente: %s", factura.RutReceptor)))

			pdf.SetXY(100, currentY+10)
			pdf.Cell(0, 5, tr(fmt.Sprintf("Contacto: %s", factura.ContactoReceptor)))

			pdf.SetXY(15, currentY+15)
			pdf.Cell(0, 5, tr(fmt.Sprintf("Giro: %s", factura.GiroReceptor)))

			pdf.SetY(currentY + 30)
			pdf.SetTextColor(0, 0, 0)
		}
	}

	// Función interna para encabezado de tabla
	addTableHeader := func() {
		pdf.SetFont("Arial", "B", 12)
		pdf.Ln(5)

		pdf.SetLineWidth(0.6)
		pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
		pdf.SetLineWidth(0.2)
		pdf.Ln(5)

		pdf.CellFormat(0, 8, tr("DETALLES DE LA COMPRA"), "", 1, "C", false, 0, "")
		pdf.Ln(5)

		pdf.SetLineWidth(0.6)
		pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
		pdf.SetLineWidth(0.2)
		pdf.Ln(5)

		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(240, 240, 240)
		pdf.SetTextColor(0, 0, 0)

		pdf.SetX(tableX)
		pdf.CellFormat(25, 8, tr("Referencia"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(60, 8, tr("Descripción"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, tr("Cantidad"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 8, tr("Precio unidad"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, tr("Subtotal"), "1", 1, "C", true, 0, "")
	}

	// Función para generar el pie de página
	// Esta función se pasa como callback a pdf.SetFooterFunc
	generateFooter := func() {
		// Guardar el estado actual para restaurarlo después de dibujar el pie de página
		pdf.SetY(-25) // Mueve la posición Y al final de la página (25mm desde abajo)
		pageW, pageH := pdf.GetPageSize()
		tr := pdf.UnicodeTranslatorFromDescriptor("")

		// Pie de página con rectángulo naranja
		pdf.SetFillColor(255, 102, 0)
		pdf.Rect(0, pageH-25, pageW, 25, "F")

		// Texto dentro del rectángulo naranja (texto blanco)
		pdf.SetTextColor(255, 255, 255)

		// Frase legal dentro del rectángulo
		pdf.SetFont("Arial", "I", 7)
		pdf.SetXY(10, pageH-20)
		pdf.MultiCell(pageW-20, 3, tr("Esta factura es representacion fiel del documento electronico firmado digitalmente segun ley N° 19.799"), "", "C", false)

		// S.I.I. - Santiago
		pdf.SetXY(10, pageH-12)
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(0, 4, "S.I.I. - Santiago", "", 1, "C", false, 0, "")

		// Número de factura en la esquina derecha
		pdf.SetXY(pageW-30, pageH-10)
		pdf.Cell(0, 5, factura.Folio)

		// Restaura el color de texto a negro para el contenido principal
		pdf.SetTextColor(0, 0, 0)
	}

	// Asignar la función de pie de página al PDF
	pdf.SetFooterFunc(generateFooter)

	// GENERAR PRIMERA PÁGINA
	addHeader(true)
	addTableHeader()

	// PROCESAR DETALLES CON PAGINACIÓN
	detalles, err := factura.GetDetalles()
	if err != nil {
		return fmt.Errorf("error al deserializar detalles: %w", err)
	}

	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)

	for _, item := range detalles {
		currentY := pdf.GetY()

		// Verificar si necesitamos nueva página para el detalle
		if currentY+rowHeight > pageHeight-marginBottom {
			pdf.AddPage()
			addHeader(false) // Sin información del cliente
			addTableHeader()
			pdf.SetFont("Arial", "", 9)
			pdf.SetFillColor(255, 255, 255)
		}

		// Agregar fila
		pdf.SetX(tableX)
		pdf.CellFormat(25, rowHeight, tr(item.SKU), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, rowHeight, tr(item.Descripcion), "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, rowHeight, fmt.Sprintf("%.0f", item.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, rowHeight, fmt.Sprintf("%.0f", item.PrecioUnitario), "1", 0, "R", false, 0, "")
		pdf.CellFormat(25, rowHeight, fmt.Sprintf("%.0f", item.TotalLinea), "1", 1, "R", false, 0, "")
	}

	pdf.Ln(5)

	// Verificar espacio para totales y la línea final de la tabla
	// Asegúrate de que la línea final de la tabla tenga espacio suficiente *antes* de dibujarla
	requiredSpaceForTotalsAndLine := 80.0 // Espacio para totales + espacio para la línea
	if pdf.GetY()+requiredSpaceForTotalsAndLine > pageHeight-50 {
		pdf.AddPage()
		addHeader(false) // Sin información del cliente
		// Puedes añadir un pequeño espacio aquí si lo necesitas entre el encabezado y la línea.
		pdf.SetY(pdf.GetY() + 10) // Mueve la posición Y para dejar espacio
	}

	// Línea final de tabla (ahora se dibujará en la nueva página si hubo un salto)
	pdf.SetLineWidth(0.6)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.SetLineWidth(0.2)
	pdf.Ln(5) // Espacio después de la línea

	pdf.Ln(5)
	pdf.SetFont("Arial", "", 10)

	// TIMBRE ELECTRÓNICO
	timbreY := pdf.GetY() - 2
	timbrePath := "./static/TimbreElectronico.png"
	var (
		timbreWidth  float64 = 90.0
		timbreHeight float64
	)
	if _, err := os.Stat(timbrePath); os.IsNotExist(err) {
		log.Printf("ADVERTENCIA: Timbre electrónico no encontrado en %s.", timbrePath)
	} else {
		info := pdf.RegisterImage(timbrePath, "")
		if info != nil {
			timbreHeight = timbreWidth * info.Height() / info.Width()
			pdf.Image(timbrePath, 15, timbreY, timbreWidth, timbreHeight, false, "", 0, "")

			// Agregar frase de verificación debajo del timbre
			pdf.SetFont("Arial", "", 9)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(35, timbreY+timbreHeight+1)
			pdf.Cell(timbreWidth, 4, tr("Verifique documento: www.sii.cl"))
		} else {
			log.Printf("ADVERTENCIA: No se pudo obtener información de la imagen del timbre en %s.", timbrePath)
		}
	}

	// TOTALES - Fijar posición a la misma altura que el timbre
	startTotalsY := timbreY // Usar la misma Y que el timbre
	rectX := 130.0
	rectWidth := 60.0
	numRows := 4
	rowHeightTotals := 7.0
	rectHeight := float64(numRows) * rowHeightTotals

	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.2)
	pdf.Rect(rectX, startTotalsY, rectWidth, rectHeight, "D")

	pdf.SetY(startTotalsY)

	pdf.SetX(rectX)
	pdf.CellFormat(40, rowHeightTotals, "Subtotal neto", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeightTotals, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(15, rowHeightTotals, FormatMoneySimple(factura.SubtotalNeto), "", 1, "R", false, 0, "")

	pdf.SetX(rectX)
	pdf.CellFormat(40, rowHeightTotals, "IVA 19%", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeightTotals, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(15, rowHeightTotals, FormatMoneySimple(factura.Iva19), "", 1, "R", false, 0, "")

	pdf.SetX(rectX)
	pdf.CellFormat(40, rowHeightTotals, "IVA Retenido", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeightTotals, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(15, rowHeightTotals, FormatMoneySimple(factura.IvaRetenido), "", 1, "R", false, 0, "")

	pdf.SetX(rectX)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, rowHeightTotals, "Total", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeightTotals, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(15, rowHeightTotals, FormatMoneySimple(factura.TotalFinal), "", 1, "R", false, 0, "")

	// --- Fin de la sección de Totales con recuadro grande ---

	// Ajustar el espacio después de los totales y el timbre
	pdf.Ln(20) // Un pequeño espacio para la frase de agradecimiento

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, "Gracias por su compra", "", 1, "C", false, 0, "")

	// Al final, escribir el PDF
	return pdf.Output(writer)
}
