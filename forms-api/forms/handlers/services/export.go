package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"forms/models"

	"github.com/xuri/excelize/v2"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
)

// ExportSubmissions exporta las respuestas a Excel
func ExportSubmissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	var filter models.FilterRequest
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	// Usar el repository que llama a sp_export_submissions (obtiene datos desde answer_details)
	submissions, err := repositories.ExportRepo.ExportSubmissions(
		filter.FormID,
		filter.StartDate,
		filter.EndDate,
		filter.Province,
		filter.SearchTerm,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		log.Printf("Error obteniendo datos de exportación: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo datos: %v", err))
		return
	}

	// Log para depuración
	log.Printf("Exportación: Se obtuvieron %d submissions del repository", len(submissions))
	if len(submissions) == 0 {
		log.Printf("Advertencia: No se encontraron submissions para exportar. Filtros: formID=%v, startDate=%v, endDate=%v, province=%v, searchTerm=%v, limit=%d, offset=%d",
			filter.FormID, filter.StartDate, filter.EndDate, filter.Province, filter.SearchTerm, filter.Limit, filter.Offset)
		helpers.RespondError(w, http.StatusNotFound, "No se encontraron datos para exportar con los filtros especificados")
		return
	}

	// Verificar si hay datos válidos (al menos un campo no vacío)
	validSubmissions := 0
	for i, subMap := range submissions {
		hasData := false
		if firstName := getStringFromMap(subMap, "first_name"); firstName != "" {
			hasData = true
		}
		if lastName := getStringFromMap(subMap, "last_name"); lastName != "" {
			hasData = true
		}
		if email := getStringFromMap(subMap, "email"); email != "" {
			hasData = true
		}
		if !hasData {
			log.Printf("Advertencia: Submission %d (ID: %v) tiene todos los campos principales vacíos", i, subMap["submission_id"])
		} else {
			validSubmissions++
		}
	}

	log.Printf("Exportación: %d de %d submissions tienen datos válidos", validSubmissions, len(submissions))

	// Crear archivo Excel
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Respuestas"
	f.SetSheetName("Sheet1", sheetName)

	// Encabezados
	baseHeaders := []string{
		"ID Respuesta", "Formulario", "Fecha de Envío",
		"Nombre", "Apellido", "Cédula", "Email", "Teléfono", "Provincia",
		"Género", "Edad", "Nombre del Negocio", "Registro Empresarial",
		"Instagram", "Facebook", "TikTok", "Twitter",
		"Consentimiento de Datos", "Newsletter",
	}

	// Campos que ya están cubiertos por los encabezados base (o que no queremos exportar como columnas dinámicas).
	knownKeys := map[string]struct{}{
		"submission_id": {}, "form_id": {}, "form_name": {},
		"submitted_at": {},
		"first_name": {}, "last_name": {}, "id_number": {}, "email": {}, "phone": {}, "province": {},
		"gender": {}, "age": {}, "business_name": {}, "business_registration": {},
		"instagram": {}, "facebook": {}, "tiktok": {}, "twitter": {},
		"data_consent": {}, "newsletter": {},
	}

	// Descubrir columnas dinámicas directamente desde los maps (table-per-form).
	dynamicKeysSet := make(map[string]struct{})
	for _, subMap := range submissions {
		for k := range subMap {
			if _, ok := knownKeys[k]; ok {
				continue
			}
			if k == "" {
				continue
			}
			dynamicKeysSet[k] = struct{}{}
		}
	}
	dynamicKeys := make([]string, 0, len(dynamicKeysSet))
	for k := range dynamicKeysSet {
		dynamicKeys = append(dynamicKeys, k)
	}
	sort.Strings(dynamicKeys)

	headers := append(append([]string{}, baseHeaders...), dynamicKeys...)

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s1", colName)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, styleHeader)
	}

	// Datos
	rowNum := 2

	for _, subMap := range submissions {
		// Obtener valores del mapa con valores por defecto seguros
		submissionID, ok := subMap["submission_id"].(int64)
		if !ok {
			log.Printf("Error: submission_id no es int64, valor: %v", subMap["submission_id"])
			continue
		}

		formName := getStringFromMap(subMap, "form_name")
		submittedAt := subMap["submitted_at"]
		firstName := getStringFromMap(subMap, "first_name")
		lastName := getStringFromMap(subMap, "last_name")
		idNumber := getStringFromMap(subMap, "id_number")
		email := getStringFromMap(subMap, "email")

		phone := getStringFromMap(subMap, "phone")
		province := getStringFromMap(subMap, "province")
		gender := getStringFromMap(subMap, "gender")
		age := getIntFromMap(subMap, "age")
		businessName := getStringFromMap(subMap, "business_name")
		businessReg := getStringFromMap(subMap, "business_registration")
		instagram := getStringFromMap(subMap, "instagram")
		facebook := getStringFromMap(subMap, "facebook")
		tiktok := getStringFromMap(subMap, "tiktok")
		twitter := getStringFromMap(subMap, "twitter")
		dataConsent := getBoolFromMap(subMap, "data_consent")
		newsletter := getBoolFromMap(subMap, "newsletter")

		// Convertir submittedAt a time.Time si es necesario
		var submittedAtTime time.Time
		switch v := submittedAt.(type) {
		case time.Time:
			submittedAtTime = v
		case []uint8:
			// MySQL devuelve DATETIME como []uint8
			if t, err := time.Parse("2006-01-02 15:04:05", string(v)); err == nil {
				submittedAtTime = t
			}
		case string:
			if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
				submittedAtTime = t
			}
		}

		values := []interface{}{
			submissionID,
			formName,
			submittedAtTime.Format("2006-01-02 15:04:05"),
			firstName,
			lastName,
			idNumber,
			email,
			phone,
			province,
			gender,
			age,
			businessName,
			businessReg,
			instagram,
			facebook,
			tiktok,
			twitter,
			getBoolString(dataConsent),
			getBoolString(newsletter),
		}

		// Valores base
		for i, value := range values {
			colName, _ := excelize.ColumnNumberToName(i + 1)
			cell := fmt.Sprintf("%s%d", colName, rowNum)
			f.SetCellValue(sheetName, cell, value)
		}

		// Valores dinámicos (columnas de la tabla form_X_responses)
		for i, key := range dynamicKeys {
			colIdx := len(baseHeaders) + i + 1 // 1-indexed
			colName, _ := excelize.ColumnNumberToName(colIdx)
			cell := fmt.Sprintf("%s%d", colName, rowNum)
			raw, ok := subMap[key]
			if !ok || raw == nil {
				continue
			}
			if b, ok := raw.([]byte); ok {
				f.SetCellValue(sheetName, cell, string(b))
				continue
			}
			f.SetCellValue(sheetName, cell, raw)
		}

		rowNum++
	}

	// Ajustar ancho de columnas
	totalCols := len(headers)
	for i := 1; i <= totalCols; i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// Configurar respuesta HTTP
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=respuestas_%s.xlsx", time.Now().Format("20060102_150405")))

	if err := f.Write(w); err != nil {
		log.Printf("Error escribiendo archivo Excel: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error generando archivo Excel")
		return
	}
}

// getBoolString convierte un bool a string "Sí"/"No"
func getBoolString(b bool) string {
	if b {
		return "Sí"
	}
	return "No"
}

// getStringFromMap obtiene un string del mapa con valor por defecto vacío
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getIntFromMap obtiene un int del mapa con valor por defecto 0
func getIntFromMap(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int32:
			return int(v)
		case int64:
			return int(v)
		}
	}
	return 0
}

// getBoolFromMap obtiene un bool del mapa con valor por defecto false
func getBoolFromMap(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
