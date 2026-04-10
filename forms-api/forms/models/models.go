package models

import (
	"encoding/json"
	"time"
)

// SubmitFormRequest representa la petición para enviar un formulario
type SubmitFormRequest struct {
	FirstName            string        `json:"first_name"`
	LastName             string        `json:"last_name"`
	IDNumber             string        `json:"id_number"`
	Email                string        `json:"email"`
	Phone                string        `json:"phone,omitempty"`
	Province             string        `json:"province,omitempty"`
	Gender               string        `json:"gender,omitempty"`
	Age                  int           `json:"age,omitempty"`
	BusinessName         string        `json:"business_name,omitempty"`
	BusinessRegistration string        `json:"business_registration,omitempty"`
	Instagram            string        `json:"instagram,omitempty"`
	Facebook             string        `json:"facebook,omitempty"`
	TikTok               string        `json:"tiktok,omitempty"`
	Twitter              string        `json:"twitter,omitempty"`
	DataConsent          bool          `json:"data_consent"`
	Newsletter           bool          `json:"newsletter"`
	FormID               int           `json:"form_id"`
	Answers              []AnswerField `json:"answers"`
}

// AnswerField representa una respuesta a un campo del formulario
type AnswerField struct {
	FieldID    int    `json:"field_id"`
	FieldValue string `json:"field_value"`
}

// SubmitFormResponse representa la respuesta después de enviar el formulario
type SubmitFormResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  int    `json:"result"` // 1=existe ambos, 2=actualizado, 3=nuevo
	UserID  int64  `json:"user_id,omitempty"`
}

// GetUserRequest representa la petición para obtener un usuario
type GetUserRequest struct {
	IDNumber string `json:"id_number"`
	Email    string `json:"email"`
}

// User representa un usuario
type User struct {
	ID                   int64     `json:"id"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	IDNumber             string    `json:"id_number"`
	Email                string    `json:"email"`
	Phone                string    `json:"phone"`
	Province             string    `json:"province"`
	Gender               string    `json:"gender"`
	Age                  int       `json:"age"`
	BusinessName         string    `json:"business_name"`
	BusinessRegistration string    `json:"business_registration"`
	Instagram            string    `json:"instagram"`
	Facebook             string    `json:"facebook"`
	TikTok               string    `json:"tiktok"`
	Twitter              string    `json:"twitter"`
	DataConsent          bool      `json:"data_consent"`
	Newsletter           bool      `json:"newsletter"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Form representa un formulario
type Form struct {
	ID               int64           `json:"id"`
	FormName         string          `json:"form_name"`
	Description      string          `json:"description"`
	CreatedAt        time.Time       `json:"created_at"`
	EventDate        time.Time       `json:"event_date"`
	ExpiresAt        time.Time       `json:"expires_at"`
	OpenAt           time.Time       `json:"open_at"`
	// status: 2=borrador, 1=abierto, 0=cerrado
	Status           int             `json:"status"`
	IsDraft          bool            `json:"is_draft"`
	IsClosed         bool            `json:"is_closed"`
	InscriptionLimit *int            `json:"inscription_limit,omitempty"`
	Conferencista    *string         `json:"conferencista,omitempty"`
	Ubicacion        *string         `json:"ubicacion,omitempty"`
	ClosedAt         *time.Time      `json:"closed_at,omitempty"`
	Schema           json.RawMessage `json:"schema,omitempty"` // JSON esquema sections+fields (objeto JSON, no base64)
	Fields           []FormField     `json:"fields,omitempty"` // derivado del schema si se expande
}

// FormField representa un campo de formulario
type FormField struct {
	ID           int64  `json:"id"`
	FormID       int64  `json:"form_id"`
	FieldLabel   string `json:"field_label"`
	FieldName    string `json:"field_name"`
	FieldType    string `json:"field_type"`
	FieldOptions string `json:"field_options,omitempty"`
	IsRequired   bool   `json:"is_required"`
	DisplayOrder int    `json:"display_order"`
	SectionID    *int64 `json:"section_id,omitempty"` // NULL si no tiene sección
}

// FormSection representa una sección/subtítulo de un formulario
type FormSection struct {
	ID           int64  `json:"id"`
	FormID       int64  `json:"form_id"`
	SectionTitle string `json:"section_title"`
	SectionIcon  string `json:"section_icon"`
	DisplayOrder int    `json:"display_order"`
	CreatedAt    string `json:"created_at,omitempty"`
}

// CreateFormRequest representa la petición para crear un formulario
type CreateFormRequest struct {
	FormName           string `json:"form_name"`
	Description        string `json:"description,omitempty"`
	EventDate          string `json:"event_date"`                     // Formato: "2006-01-02T15:04:05" o "2006-01-02 15:04:05" - REQUERIDO
	OpenAt             string `json:"open_at"`                        // Formato: "2006-01-02T15:04:05" o "2006-01-02 15:04:05" - REQUERIDO
	ExpiresAt          string `json:"expires_at"`                     // Formato: "2006-01-02T15:04:05" o "2006-01-02 15:04:05" - REQUERIDO (debe ser futura)
	InscriptionLimit   *int   `json:"inscription_limit,omitempty"`    // NULL = sin límite
	Conferencista      string `json:"conferencista,omitempty"`
	Ubicacion          string `json:"ubicacion,omitempty"`
	UseDefaultTemplate bool   `json:"use_default_template,omitempty"` // Si es true, agrega los campos por defecto
}

// PatchFormMetadataRequest actualiza conferencista y/o ubicación (clave omitida = no cambiar; "" = NULL en BD).
type PatchFormMetadataRequest struct {
	Conferencista *string `json:"conferencista"`
	Ubicacion     *string `json:"ubicacion"`
}

// CreateFieldRequest representa la petición para crear un campo
type CreateFieldRequest struct {
	FieldLabel   string `json:"field_label"`
	FieldName    string `json:"field_name"`
	FieldType    string `json:"field_type"`
	FieldOptions string `json:"field_options,omitempty"`
	IsRequired   bool   `json:"is_required"`
}

// CreateSectionRequest representa la petición para crear una sección
type CreateSectionRequest struct {
	SectionTitle string `json:"section_title"`
	SectionIcon  string `json:"section_icon"`
}

// UpdateSectionRequest representa la petición para actualizar una sección
type UpdateSectionRequest struct {
	SectionTitle string `json:"section_title"`
	SectionIcon  string `json:"section_icon"`
}

// AssignFieldToSectionRequest representa la petición para asignar un campo a una sección
type AssignFieldToSectionRequest struct {
	SectionID *int64 `json:"section_id"` // NULL para quitar de sección
}

// ReorderSectionRequest representa la petición para reordenar una sección
type ReorderSectionRequest struct {
	Direction string `json:"direction"` // "up" o "down"
}

// FormsListResponse representa la lista de formularios
type FormsListResponse struct {
	Success bool   `json:"success"`
	Forms   []Form `json:"forms"`
}

// FormResponse representa un formulario con sus campos
type FormResponse struct {
	Success bool `json:"success"`
	Form    Form `json:"form"`
}

// Role representa un rol (tabla/vista roles para el frontend)
type Role struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
