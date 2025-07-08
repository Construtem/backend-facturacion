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

	// AGREGAR: Configurar transformación UTF-8
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Arial", "", 10)

	// IMPORTANTE: Agregar una página
	pdf.AddPage()

	pdf.SetAutoPageBreak(false, 0)

	/*
	 logoPath := "./static/construtem_logo.png"
	 if _, err := os.Stat(logoPath); os.IsNotExist(err) {
	 	log.Printf("ADVERTENCIA: Logo no encontrado en %s. Continuando sin logo.", logoPath)
	 } else {
	 	pdf.Image(logoPath, 15, 15, 30, 0, false, "", 0, "") // Posición y tamaño del logo
	 }
	*/
	// Ahora usar tr() para todos los textos con caracteres especiales

	// 3. Datos del Emisor y Título de Factura
	// Guardar la posición Y inicial para alinear elementos
	initialY := pdf.GetY()

	// --- Sección del Título de Factura y RUT (Recuadro Rojo a la Derecha) ---
	// Definir color de borde rojo y relleno blanco
	pdf.SetDrawColor(255, 0, 0)     // Borde rojo
	pdf.SetFillColor(255, 255, 255) // Relleno blanco
	pdf.SetTextColor(0, 0, 0)       // Texto negro
	pdf.SetLineWidth(0.5)           // Grosor del borde

	// Posición y dimensiones del recuadro rojo
	rectRightX := 120.0
	rectRightY := 10.0
	rectRightWidth := 80.0
	rectRightHeight := 30.0

	pdf.Rect(rectRightX, rectRightY, rectRightWidth, rectRightHeight, "D") // "D" para dibujar solo el borde (y no "FD" para rellenar)

	// Contenido dentro del recuadro rojo - USAR DATOS DIRECTOS DE FACTURA
	pdf.SetTextColor(255, 0, 0)

	// R.U.T. centrado
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(rectRightX+5, rectRightY+3)
	pdf.CellFormat(rectRightWidth-10, 5, tr(fmt.Sprintf("R.U.T.: %s", factura.RutEmisor)), "", 0, "C", false, 0, "")

	// "Factura Electrónica" centrado (ya estaba centrado)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(rectRightX+5, rectRightY+10)
	pdf.MultiCell(rectRightWidth-10, 5, tr("FACTURA ELECTRONICA"), "", "C", false)

	// N° folio centrado
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(rectRightX+5, rectRightY+22)
	pdf.CellFormat(rectRightWidth-10, 5, tr(fmt.Sprintf("N° %s", factura.Folio)), "", 0, "C", false, 0, "")

	// RESTABLECER COLOR DE DIBUJO A NEGRO
	pdf.SetDrawColor(0, 0, 0)

	// --- Sección de Datos de la Ferretería (A la Izquierda) - USAR DATOS DE FACTURA ---
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetXY(15, initialY+5)
	pdf.Cell(0, 6, tr(factura.RazonSocialEmisor)) // ✅ Directo de factura
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 9)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(factura.GiroEmisor)) // ✅ Directo de factura
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Dirección: %s, %s", factura.DireccionEmisor, factura.ComunaEmisor)))
	pdf.Ln(4)
	pdf.SetX(15)
	pdf.Cell(0, 4, tr(fmt.Sprintf("Email: %s | Tel: %s", factura.EmailEmisor, factura.TelefonoEmisor))) // ✅ Directo de factura
	pdf.Ln(15)

	// Para el rectángulo naranja
	currentY := pdf.GetY()
	pdf.SetFillColor(255, 102, 0)
	pdf.Rect(10, currentY, 190, 25, "F")

	pdf.SetTextColor(255, 255, 255) // Texto blanco para el rectángulo naranja
	pdf.SetFont("Arial", "B", 10)

	// Primera línea con caracteres especiales
	pdf.SetXY(15, currentY+5)
	pdf.Cell(0, 5, tr(fmt.Sprintf("N° Factura: %s", factura.Folio)))

	pdf.SetXY(110, currentY+5)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Fecha de Emisión: %s", FormatDateChilean(factura.FechaEmision))))

	// Segunda línea
	pdf.SetXY(15, currentY+15)
	pdf.Cell(0, 5, tr(fmt.Sprintf("Cliente: %s", factura.RazonSocialReceptor)))

	pdf.SetXY(110, currentY+15)
	pdf.Cell(0, 5, tr(fmt.Sprintf("RUT Cliente: %s", factura.RutReceptor)))

	// Posicionar después del rectángulo
	pdf.SetY(currentY + 30)
	pdf.SetTextColor(0, 0, 0) // Volver a texto negro
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(5)

	// Línea separadora superior (grosor doble)
	pdf.SetLineWidth(0.6)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.SetLineWidth(0.2) // Restaurar grosor por defecto
	pdf.Ln(5)

	pdf.CellFormat(0, 8, tr("DETALLES DE LA COMPRA"), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Línea separadora inferior (grosor doble)
	pdf.SetLineWidth(0.6)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.SetLineWidth(0.2) // Restaurar grosor por defecto
	pdf.Ln(5)

	// 6. Encabezados de la tabla - Mejorar diseño
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240) // Gris claro como en la referencia
	pdf.SetTextColor(0, 0, 0)       // Texto negro

	// Usar los mismos anchos que la referencia
	pdf.CellFormat(25, 8, tr("Referencia"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, tr("Descripción"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, tr("Cantidad"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, tr("Precio unidad"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, tr("Subtotal"), "1", 1, "C", true, 0, "")

	// 7. Filas de datos
	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)

	// Obtener detalles del JSON
	detalles, err := factura.GetDetalles()
	if err != nil {
		return fmt.Errorf("error al deserializar detalles: %w", err)
	}

	for _, item := range detalles {
		pdf.CellFormat(25, 8, tr(item.SKU), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 8, tr(item.Descripcion), "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%.0f", item.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.0f", item.PrecioUnitario), "1", 0, "R", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%.0f", item.TotalLinea), "1", 1, "R", false, 0, "")
	}

	// Después de las filas de datos
	pdf.Ln(5)

	// Línea bajo la tabla de detalles
	pdf.SetLineWidth(0.6)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.SetLineWidth(0.2) // Restaurar grosor por defecto
	pdf.Ln(5)

	// 8. Totales
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 10)

	// Calcular la posición Y para la imagen del timbre
	timbreY := pdf.GetY() + 2

	// Cargar la imagen del timbre electrónico
	timbrePath := "./static/TimbreElectronico.png"
	var (
		timbreWidth  float64 = 80.0
		timbreHeight float64
	)
	if _, err := os.Stat(timbrePath); os.IsNotExist(err) {
		log.Printf("ADVERTENCIA: Timbre electrónico no encontrado en %s.", timbrePath)
	} else {
		info := pdf.RegisterImage(timbrePath, "")
		if info != nil {
			timbreHeight = timbreWidth * info.Height() / info.Width()
			pdf.Image(timbrePath, 15, timbreY, timbreWidth, timbreHeight, false, "", 0, "")
		} else {
			log.Printf("ADVERTENCIA: No se pudo obtener información de la imagen del timbre en %s.", timbrePath)
		}
	}

	// --- Inicio de la sección de Totales con recuadro grande ---
	// Guardar la posición Y actual para el inicio del recuadro de totales
	startTotalsY := pdf.GetY()

	// Definir las dimensiones y posición del recuadro grande
	rectX := 130.0
	rectWidth := 70.0
	numRows := 4
	rowHeight := 7.0
	rectHeight := float64(numRows) * rowHeight

	// Dibujar el recuadro grande (solo el borde)
	pdf.SetDrawColor(0, 0, 0) // Color del borde (negro)
	pdf.SetLineWidth(0.2)     // Grosor del borde
	pdf.Rect(rectX, startTotalsY, rectWidth, rectHeight, "D")

	// Posicionar el cursor para escribir los textos dentro del recuadro
	pdf.SetY(startTotalsY)

	// Subtotal neto
	pdf.SetX(rectX)
	pdf.CellFormat(35, rowHeight, "Subtotal neto", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeight, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(25, rowHeight, FormatMoneySimple(factura.SubtotalNeto), "", 1, "R", false, 0, "")

	// IVA
	pdf.SetX(rectX)
	pdf.CellFormat(35, rowHeight, "IVA 19%", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeight, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(25, rowHeight, FormatMoneySimple(factura.Iva19), "", 1, "R", false, 0, "")

	// IVA RETENIDO
	pdf.SetX(rectX)
	pdf.CellFormat(35, rowHeight, "IVA Retenido", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeight, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(25, rowHeight, FormatMoneySimple(factura.IvaRetenido), "", 1, "R", false, 0, "")

	// Total
	pdf.SetX(rectX)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(35, rowHeight, "Total", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeight, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(25, rowHeight, FormatMoneySimple(factura.TotalFinal), "", 1, "R", false, 0, "")

	// --- Fin de la sección de Totales con recuadro grande ---

	// Ajustar el espacio después de los totales y el timbre
	pdf.Ln(10 + timbreHeight)

	// 8. Mensaje final
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, "Gracias por su compra", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Obtener las dimensiones de la página
	pageWidth, pageHeight := pdf.GetPageSize()

	// 9. Pie de página con rectángulo naranja MÁS ALTO para incluir todo el texto
	pdf.SetFillColor(255, 102, 0)
	pdf.Rect(0, pageHeight-35, pageWidth, 35, "F")

	// Texto dentro del rectángulo naranja (texto blanco)
	pdf.SetTextColor(255, 255, 255)

	// Frase legal dentro del rectángulo
	pdf.SetFont("Arial", "I", 7)
	pdf.SetXY(10, pageHeight-30)
	pdf.MultiCell(pageWidth-20, 3, tr("Esta factura es representacion fiel del documento electronico firmado digitalmente segun ley N° 19.799"), "", "C", false)

	// S.I.I. - Santiago
	pdf.SetXY(10, pageHeight-22)
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, "S.I.I. - Santiago", "", 1, "C", false, 0, "")

	// Número de factura en la esquina derecha
	pdf.SetXY(pageWidth-30, pageHeight-12)
	pdf.Cell(0, 5, factura.Folio)

	// Al final, escribir el PDF
	return pdf.Output(writer) // Escribe el PDF al writer (ej. http.ResponseWriter)
}
