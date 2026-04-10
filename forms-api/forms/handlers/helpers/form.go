package helpers

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"forms/database"
	"forms/handlers/repositories"
	"forms/models"
)

// FormExists verifica si existe un formulario con el ID dado usando el repository
func FormExists(id int64) (bool, error) {
	return repositories.FormRepo.FormExists(id)
}

var (
	hasIsClosedOnce sync.Once
	hasIsClosedCol  bool
)

// HasIsClosedColumn verifica si existe la columna is_closed en la tabla forms
func HasIsClosedColumn() bool {
	hasIsClosedOnce.Do(func() {
		err := database.DB.QueryRow(`
			SELECT COUNT(*) > 0
			FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE()
			  AND TABLE_NAME = 'forms'
			  AND COLUMN_NAME = 'is_closed'
		`).Scan(&hasIsClosedCol)
		if err != nil {
			log.Printf("Error comprobando columnas is_closed/closed_at: %v", err)
			hasIsClosedCol = false
		}
	})
	return hasIsClosedCol
}

// ScanForm escanea un formulario desde la base de datos
func ScanForm(scanner interface{ Scan(dest ...any) error }, withClosed bool) (models.Form, error) {
	var form models.Form
	var desc sql.NullString
	var expiresAt time.Time

	if withClosed {
		var closedAt sql.NullTime
		if err := scanner.Scan(
			&form.ID,
			&form.FormName,
			&desc,
			&form.CreatedAt,
			&expiresAt,
			&form.IsClosed,
			&closedAt,
		); err != nil {
			return models.Form{}, err
		}
		if closedAt.Valid {
			form.ClosedAt = &closedAt.Time
		}
	} else {
		if err := scanner.Scan(
			&form.ID,
			&form.FormName,
			&desc,
			&form.CreatedAt,
			&expiresAt,
		); err != nil {
			return models.Form{}, err
		}
		form.IsClosed = false
		form.ClosedAt = nil
	}

	if desc.Valid {
		form.Description = desc.String
	}
	form.ExpiresAt = expiresAt
	return form, nil
}

// ValidFieldTypes es un mapa de tipos de campo válidos
var ValidFieldTypes = map[string]bool{
	"text": true, "number": true, "email": true, "select": true,
	"checkbox": true, "textarea": true, "date": true, "tel": true, "radio": true,
}

// MapDBFieldType mapea tipos de campo a tipos de base de datos
func MapDBFieldType(fieldType, fieldName string) string {
	switch fieldType {
	case "tel":
		log.Printf("Mapeando tipo 'tel' a 'text' para campo %s", fieldName)
		return "text"
	case "radio":
		log.Printf("Mapeando tipo 'radio' a 'select' para campo %s", fieldName)
		return "select"
	default:
		return fieldType
	}
}

// TemplatesMap obtiene un mapa de plantillas de campos comunes
func TemplatesMap() map[string]models.CommonFieldTemplate {
	allTemplates := models.GetCommonFieldTemplates()
	m := make(map[string]models.CommonFieldTemplate, len(allTemplates))
	for _, t := range allTemplates {
		m[t.Name] = t
	}
	return m
}

