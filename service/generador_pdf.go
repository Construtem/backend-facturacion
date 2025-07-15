package services

import (
	"backend-facturacion/models"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/jung-kurt/gofpdf"
)

// CompanyConfig estructura para configuración de empresa
type CompanyConfig struct {
	RUT         string
	RazonSocial string
	Giro        string
	Direccion   string
	Comuna      string
	Ciudad      string
	Telefono    string
	Email       string
	LogoPath    string
}

// LoadCompanyConfig carga configuración desde variables de entorno
// Primer parametro: busca variable de entorno, si no existe usa valor por defecto
// Segundo parametro: valor por defecto si la variable de entorno no está definida
// Retorna un puntero a CompanyConfig con los valores cargados
func LoadCompanyConfig() *CompanyConfig {
	return &CompanyConfig{
		RUT:         getEnv("EMPRESA_RUT", "76.123.456-7"),
		RazonSocial: getEnv("EMPRESA_RAZON_SOCIAL", "Ferretería Construtem S.A."),
		Giro:        getEnv("EMPRESA_GIRO", "Venta al por mayor y menor de artículos de ferreteria"),
		Direccion:   getEnv("EMPRESA_DIRECCION", "Av. Siempre Viva 742"),
		Comuna:      getEnv("EMPRESA_COMUNA", "Santiago"),
		Ciudad:      getEnv("EMPRESA_CIUDAD", "Santiago"),
		Telefono:    getEnv("EMPRESA_TELEFONO", "+56 2 1234 5678"),
		Email:       getEnv("EMPRESA_EMAIL", "contacto@construtem.cl"),
		LogoPath:    getEnv("EMPRESA_LOGO_PATH", "./static/construtem.png"),
	}
}

// getEnv obtiene variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Función para generar el PDF de la factura
func GenerateInvoicePDF(factura *models.Factura, writer io.Writer) error {
	// Validar antes de generar
	if err := factura.ValidateForPDF(); err != nil {
		return fmt.Errorf("validación fallida: %w", err)
	}

	// Cargar configuración
	config := LoadCompanyConfig()

	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Arial", "", 10)
	pdf.AddPage()
	pdf.SetAutoPageBreak(false, 0)

	totalPages := 1 // Empezar con 1 porque ya tenemos la primera página

	// Constantes para control de paginación
	const (
		rowHeight    = 8.0
		pageHeight   = 297.0
		marginBottom = 59.0
		tableX       = 20.0
	)

	// Función interna para crear encabezado
	addHeader := func(includeClientInfo bool) {
		initialY := pdf.GetY()

		// LOGO DE LA EMPRESA (lado izquierdo)
		logoPath := config.LogoPath
		logoX := 15.0
		logoY := 2.5
		logoMaxWidth := 60.0
		logoMaxHeight := 25.0

		if _, err := os.Stat(logoPath); err == nil {
			info := pdf.RegisterImage(logoPath, "")
			if info != nil {
				// Calcular dimensiones manteniendo proporción
				logoWidth := logoMaxWidth
				logoHeight := logoWidth * info.Height() / info.Width()

				// Ajustar si la altura excede el máximo
				if logoHeight > logoMaxHeight {
					logoHeight = logoMaxHeight
					logoWidth = logoHeight * info.Width() / info.Height()
				}

				pdf.Image(logoPath, logoX, logoY, logoWidth, logoHeight, false, "", 0, "")
				log.Printf("INFO: Logo cargado correctamente desde %s", logoPath)
			} else {
				log.Printf("ADVERTENCIA: No se pudo procesar la imagen del logo en %s", logoPath)
			}
		} else {
			log.Printf("ADVERTENCIA: Logo no encontrado en %s", logoPath)
		}

		// Recuadro rojo con información de factura (lado derecho)
		pdf.SetDrawColor(255, 0, 0)
		pdf.SetFillColor(255, 255, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetLineWidth(1)

		rectRightX := 120.0
		rectRightY := 5.0
		rectRightWidth := 80.0
		rectRightHeight := 30.0

		pdf.Rect(rectRightX, rectRightY, rectRightWidth, rectRightHeight, "D")

		// Contenido del recuadro rojo - usar config como fallback
		rutEmisor := factura.RutEmisor
		if rutEmisor == "" {
			rutEmisor = config.RUT
		}

		pdf.SetTextColor(255, 0, 0)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(rectRightX+5, rectRightY+5)
		pdf.CellFormat(rectRightWidth-10, 5, tr(fmt.Sprintf("R.U.T.: %s", rutEmisor)), "", 0, "C", false, 0, "")

		pdf.SetFont("Arial", "B", 12)
		pdf.SetXY(rectRightX+5, rectRightY+10)
		pdf.MultiCell(rectRightWidth-10, 5, tr("FACTURA ELECTRONICA"), "", "C", false)

		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(rectRightX+5, rectRightY+18)
		pdf.CellFormat(rectRightWidth-10, 5, tr(fmt.Sprintf("N° %d", factura.ID)), "", 0, "C", false, 0, "")

		pdf.SetDrawColor(0, 0, 0)

		// Datos del emisor - usar config como fallback
		// Posicionar debajo del logo (Y = 35)
		pdf.SetY(40)

		razonSocial := factura.RazonSocialEmisor
		if razonSocial == "" {
			razonSocial = config.RazonSocial
		}

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "B", 14)
		pdf.SetXY(15, initialY+15)
		pdf.Cell(0, 6, tr(razonSocial))
		pdf.Ln(6)

		giro := factura.GiroEmisor
		if giro == "" {
			giro = config.Giro
		}

		pdf.SetFont("Arial", "", 9)
		pdf.SetX(15)
		pdf.Cell(0, 4, tr(giro))
		pdf.Ln(4)

		direccion := factura.DireccionEmisor
		comuna := factura.ComunaEmisor
		if direccion == "" {
			direccion = config.Direccion
		}
		if comuna == "" {
			comuna = config.Comuna
		}

		pdf.SetX(15)
		pdf.Cell(0, 5, tr(fmt.Sprintf("Dirección: %s, %s", direccion, comuna)))
		pdf.Ln(4)

		email := factura.EmailEmisor
		telefono := factura.TelefonoEmisor
		if email == "" {
			email = config.Email
		}
		if telefono == "" {
			telefono = config.Telefono
		}

		pdf.SetX(15)
		pdf.Cell(0, 4, tr(fmt.Sprintf("Email: %s | Tel: %s", email, telefono)))
		pdf.SetX(150)
		pdf.SetTextColor(255, 0, 0)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 4, tr(fmt.Sprintf("S.I.I. - %s", config.Ciudad)))
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 9)
		pdf.Ln(7)

		// Fecha de emisión (alineada con el recuadro rojo)
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
		pdf.SetFillColor(255, 102, 0)
		pdf.SetTextColor(255, 255, 255) // (0, 0, 0) para texto negro
		pdf.SetDrawColor(255, 102, 0)

		pdf.SetX(tableX)
		pdf.CellFormat(25, 8, tr("Referencia"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(60, 8, tr("Descripción"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, tr("Cantidad"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 8, tr("Precio unidad"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, tr("Subtotal"), "1", 1, "C", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetDrawColor(0, 0, 0)
	}

	// funcion para generar el pie de página
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

		// Frase legal
		pdf.SetFont("Arial", "I", 7)
		pdf.SetXY(10, pageH-20)
		pdf.MultiCell(pageW-20, 3, tr("Esta factura es representacion fiel del documento electronico firmado digitalmente segun ley N° 19.799"), "", "C", false)

		// S.I.I. - Santiago
		pdf.SetXY(10, pageH-12)
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(0, 4, "S.I.I. - Santiago", "", 1, "C", false, 0, "")

		// Número de página en la esquina derecha
		pdf.SetXY(pageW-30, pageH-12)
		pdf.Cell(0, 5, tr(fmt.Sprintf("Página %d de %d", pdf.PageNo(), totalPages)))

		// Número de factura en la esquina derecha
		// 	pdf.SetXY(pageW-30, pageH-10)
		// 	pdf.Cell(0, 5, factura.Folio)

		pdf.SetTextColor(0, 0, 0)
	}

	// Asignar la función de pie de página al PDF
	pdf.SetFooterFunc(generateFooter)

	addHeader(true)
	addTableHeader()

	// Dibujar los detalles de la factura
	detalles, err := factura.GetDetalles()
	if err != nil {
		return fmt.Errorf("error al deserializar detalles: %w", err)
	}

	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)

	for i, item := range detalles {
		descHeight := calculateRowHeight(pdf, item.Descripcion, 60.0)
		finalRowHeight := math.Max(rowHeight, descHeight)

		currentY := pdf.GetY()

		// Verificar si necesitamos nueva página
		if currentY+finalRowHeight > pageHeight-marginBottom {
			pdf.AddPage()
			totalPages = pdf.PageNo()

			addHeader(false)
			addTableHeader()
			pdf.SetFont("Arial", "", 9)
		}

		startY := pdf.GetY()

		// Color de los bordes
		pdf.SetDrawColor(255, 255, 255) // Blanco para bordes
		// pdf.SetDrawColor(240, 240, 240)    // Gris claro
		// pdf.SetDrawColor(0, 102, 204)     // Azul corporativo
		// pdf.SetDrawColor(255, 102, 0)    // Naranja (color del header)
		// pdf.SetDrawColor(0, 0, 0)       // Negro (original)

		// Colores de fondo alternados
		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255) // Blanco
		} else {
			pdf.SetFillColor(240, 240, 240) // Gris claro
		}

		pdf.SetX(tableX)

		// Columna 1: Referencia
		pdf.CellFormat(25, finalRowHeight, tr(item.SKU), "1", 0, "C", true, 0, "")

		// Columna 2: Descripción (con MultiCell)
		descX := pdf.GetX()
		descY := pdf.GetY()

		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255) // Blanco
		} else {
			pdf.SetFillColor(240, 240, 240) // Gris claro
		}
		pdf.Rect(descX, descY, 60, finalRowHeight, "DF")

		pdf.SetDrawColor(255, 255, 255) // Mismo color que las otras celdas
		pdf.Rect(descX, descY, 60, finalRowHeight, "D")

		// Posicionar para MultiCell (con padding interno)
		pdf.SetXY(descX+1, descY+1)
		pdf.MultiCell(58, 4, tr(item.Descripcion), "", "L", false)

		// Restaurar posición para las siguientes columnas
		pdf.SetXY(descX+60, startY)

		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255) // Blanco
		} else {
			pdf.SetFillColor(240, 240, 240) // Gris claro
		}

		// Columna 3: Cantidad
		pdf.CellFormat(25, finalRowHeight, fmt.Sprintf("%.0f", item.Cantidad), "1", 0, "C", true, 0, "")

		// Columna 4: Precio unitario
		pdf.CellFormat(35, finalRowHeight, FormatMoneySimple(item.PrecioUnitario), "1", 0, "R", true, 0, "")

		// Columna 5: Subtotal
		pdf.CellFormat(25, finalRowHeight, FormatMoneySimple(item.TotalLinea), "1", 0, "R", true, 0, "")

		// mover a la siguiente línea
		pdf.SetXY(tableX, startY+finalRowHeight)
	}

	pdf.Ln(5)

	// Verificar espacio para totales y la línea final de la tabla
	// Asegurar que la línea final de la tabla tenga espacio suficiente antes de dibujarla
	requiredSpaceForTotalsAndLine := 60.0 // Espacio para totales + espacio para la línea / 80 esta predefinido
	if pdf.GetY()+requiredSpaceForTotalsAndLine > pageHeight-50 {
		pdf.AddPage()
		// Actualizar el número de páginas después de agregar una nueva página
		totalPages = pdf.PageNo()

		addHeader(false)
		// Puedes añadir un pequeño espacio aquí si lo necesitas entre el encabezado y la línea.
		pdf.SetY(pdf.GetY() + 10) // Mueve la posición Y para dejar espacio
	}

	// asegurar que la línea final de la tabla se dibuje en la posición correcta
	totalPages = pdf.PageNo()

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
			pdf.SetXY(25, timbreY+timbreHeight+1)
			pdf.Cell(timbreWidth, 4, tr("Res.99 de 2014 Verifique documento: www.sii.cl"))
		} else {
			log.Printf("ADVERTENCIA: No se pudo obtener información de la imagen del timbre en %s.", timbrePath)
		}
	}

	// TOTALES - Fijar posición a la misma altura que el timbre
	startTotalsY := timbreY // Usar la misma Y que el timbre
	rectX := 130.0
	rectWidth := 60.0
	numRows := 5 // Número de filas para los totales
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

	// agragar despacho y descuento.

	pdf.SetX(rectX)
	pdf.CellFormat(40, rowHeightTotals, "Despacho", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, rowHeightTotals, "$", "", 0, "L", false, 0, "")
	pdf.CellFormat(15, rowHeightTotals, FormatMoneySimple(0), "", 1, "R", false, 0, "")

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

// Calcular altura necesaria para texto largo
func calculateRowHeight(pdf *gofpdf.Fpdf, text string, cellWidth float64) float64 {
	const lineHeight = 4.0 // Altura de cada línea de texto para MultiCell
	const minHeight = 8.0  // Altura mínima de fila
	const padding = 2.0    // Padding interno de la celda

	// Verificar texto vacío
	if text == "" {
		return minHeight
	}

	// Ancho efectivo de la celda (restando padding)
	effectiveCellWidth := cellWidth - padding

	// Simular MultiCell para obtener altura real
	// Guardar posición actual
	currentX, currentY := pdf.GetX(), pdf.GetY()

	// Mover a una posición temporal fuera del área visible
	pdf.SetXY(300, 300)

	// Realizar MultiCell temporal para medir altura
	startY := pdf.GetY()
	pdf.MultiCell(effectiveCellWidth, lineHeight, text, "", "L", false)
	endY := pdf.GetY()

	// Calcular altura real
	actualHeight := endY - startY + padding

	// Restaurar posición original
	pdf.SetXY(currentX, currentY)

	// Retornar la mayor entre altura mínima y altura calculada
	return math.Max(minHeight, actualHeight)
}
