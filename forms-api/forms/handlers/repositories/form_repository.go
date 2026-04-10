package repositories

import (
	"forms/models"
)

// FormRepository define las operaciones de acceso a datos para formularios
type FormRepository interface {
	// Create crea un nuevo formulario
	Create(formName, description string, eventDate, openAt, expiresAt interface{}, isDraft bool, inscriptionLimit *int, conferencista, ubicacion *string) (int64, error)

	// GetForms obtiene todos los formularios
	GetForms() ([]models.Form, error)

	// GetByID obtiene un formulario por ID
	GetByID(formID int64) (*models.Form, error)

	// FormExists verifica si existe un formulario
	FormExists(formID int64) (bool, error)

	// GetFields obtiene los campos de un formulario
	GetFields(formID int64) ([]models.FormField, error)

	// GetFieldNames obtiene los nombres de los campos existentes
	GetFieldNames(formID int64) ([]string, error)

	// CreateField crea un nuevo campo
	CreateField(formID int64, fieldLabel, fieldName, fieldType, fieldOptions string, isRequired bool) (int64, error)

	// UpdateField actualiza un campo existente
	UpdateField(fieldID, formID int64, fieldLabel, fieldName, fieldType, fieldOptions string, isRequired bool) (bool, error)

	// DeleteField elimina un campo
	DeleteField(fieldID, formID int64) (bool, string, error)

	// ReorderField reordena un campo (mueve arriba o abajo)
	ReorderField(fieldID, formID int64, direction string) (bool, string, error)

	// AddDefaultFields agrega los campos por defecto a un formulario
	AddDefaultFields(formID int64) (int, error)

	// CreateSection crea una nueva sección
	CreateSection(formID int64, title, icon string) (int64, error)

	// GetFormSections obtiene todas las secciones de un formulario
	GetFormSections(formID int64) ([]models.FormSection, error)

	// UpdateSection actualiza una sección
	UpdateSection(sectionID, formID int64, title, icon string) (bool, error)

	// DeleteSection elimina una sección
	DeleteSection(sectionID, formID int64) (bool, string, error)

	// ReorderSection reordena una sección
	ReorderSection(sectionID, formID int64, direction string) (bool, string, error)

	// AssignFieldToSection asigna un campo a una sección
	AssignFieldToSection(fieldID, sectionID, formID int64) (bool, error)

	// GetFormName obtiene el nombre de un formulario
	GetFormName(formID int64) (string, error)

	// GetFormExpiresAt obtiene la fecha de expiración de un formulario
	GetFormExpiresAt(formID int64) (interface{}, bool, error)

	// CloseForm cierra un formulario
	CloseForm(formID, adminUserID int64, ipAddress, userAgent string) (bool, error)

	// SetDraftStatus setea el flag is_draft (borrador/activo público)
	SetDraftStatus(formID int64, isDraft bool) (bool, error)

	// OpenForm marca el formulario como abierto para envíos (is_closed=0, is_draft=0).
	OpenForm(formID int64) (bool, error)

	// UpdateFormSchema actualiza forms.schema; si no hay envíos, recrea form_{id}_responses
	UpdateFormSchema(formID int64, schemaJSON []byte) error

	UpdateFormMetadata(formID int64, conferencista, ubicacion *string) error
}
