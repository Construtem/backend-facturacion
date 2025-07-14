package models

import (
	"fmt"
	"regexp"
	"strings"
)

func (f *Factura) ValidateForPDF() error {
	var errors []string

	if f.RutEmisor == "" {
		errors = append(errors, "RUT emisor es requerido")
	}
	if f.RazonSocialEmisor == "" {
		errors = append(errors, "Razón social emisor es requerida")
	}
	if f.RutReceptor == "" {
		errors = append(errors, "RUT receptor es requerido")
	}
	if f.RazonSocialReceptor == "" {
		errors = append(errors, "Razón social receptor es requerida")
	}
	if f.Folio == "" {
		errors = append(errors, "Folio es requerido")
	}

	// FIX: Validar detalles correctamente
	detalles, err := f.GetDetalles()
	if err != nil {
		errors = append(errors, "error al deserializar items: "+err.Error())
	} else if len(detalles) == 0 {
		errors = append(errors, "items son requeridos")
	}

	// Validar formato de RUTs
	if f.RutEmisor != "" && !ValidateRUT(f.RutEmisor) {
		errors = append(errors, "formato de RUT emisor inválido")
	}
	if f.RutReceptor != "" && !ValidateRUT(f.RutReceptor) {
		errors = append(errors, "formato de RUT receptor inválido")
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores de validación: %s", strings.Join(errors, ", "))
	}

	return nil
}

// ValidateRUT valida formato de RUT chileno
func ValidateRUT(rut string) bool {
	if rut == "" {
		return false
	}
	// Formato más flexible: XX.XXX.XXX-X o XXXXXXXX-X
	re := regexp.MustCompile(`^\d{1,2}\.\d{3}\.\d{3}-[\dkK]$|^\d{7,8}-[\dkK]$`)
	return re.MatchString(rut)
}
