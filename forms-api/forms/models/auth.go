package models

import "time"

// DashboardStats representa las estadísticas del dashboard
type DashboardStats struct {
	TotalForms       int `json:"total_forms"`
	ActiveForms      int `json:"active_forms"`
	ExpiredForms     int `json:"expired_forms"`
	TotalUsers       int `json:"total_users"`
	TotalSubmissions int `json:"total_submissions"`
	SubmissionsToday int `json:"submissions_today"`
	SubmissionsWeek  int `json:"submissions_week"`
	SubmissionsMonth int `json:"submissions_month"`
}

// FormStats representa estadísticas de un formulario
type FormStats struct {
	ID              int64     `json:"id"`
	FormName        string    `json:"form_name"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	SubmissionCount int       `json:"submission_count"`
	UniqueUsers     int       `json:"unique_users"`
	LastSubmission  time.Time `json:"last_submission,omitempty"`
	Status          string    `json:"status"` // active, expired
}

// FormCustomization representa la personalización visual de un formulario
type FormCustomization struct {
	ID                          int64   `json:"id"`
	FormID                      int64   `json:"form_id"`
	PrimaryColor                string  `json:"primary_color"`
	SecondaryColor              string  `json:"secondary_color"`
	BackgroundColor             string  `json:"background_color"`
	TextColor                   string  `json:"text_color"`
	TitleColor                  string  `json:"title_color"`
	LogoURL                     string  `json:"logo_url,omitempty"`
	LogoURLMobile               string  `json:"logo_url_mobile,omitempty"`
	FontFamily                  string  `json:"font_family"`
	ButtonStyle                 string  `json:"button_style"`
	FormContainerColor          string  `json:"form_container_color"`
	FormContainerOpacity        float64 `json:"form_container_opacity"`
	DescriptionContainerColor   string  `json:"description_container_color"`
	DescriptionContainerOpacity float64 `json:"description_container_opacity"`
	FormMetaBackground          string  `json:"form_meta_background"`
	FormMetaBackgroundStart     string  `json:"form_meta_background_start"`
	FormMetaBackgroundEnd       string  `json:"form_meta_background_end"`
	FormMetaBackgroundOpacity   float64 `json:"form_meta_background_opacity"`
	FormMetaTextColor           string  `json:"form_meta_text_color"`
}

// CreateCustomizationRequest representa la petición para crear/actualizar personalización
type CreateCustomizationRequest struct {
	FormID                      int64   `json:"form_id"`
	PrimaryColor                string  `json:"primary_color"`
	SecondaryColor              string  `json:"secondary_color"`
	BackgroundColor             string  `json:"background_color"`
	TextColor                   string  `json:"text_color"`
	TitleColor                  string  `json:"title_color"`
	LogoURL                     string  `json:"logo_url,omitempty"`
	LogoURLMobile               string  `json:"logo_url_mobile,omitempty"`
	FontFamily                  string  `json:"font_family"`
	ButtonStyle                 string  `json:"button_style"`
	FormContainerColor          string  `json:"form_container_color"`
	FormContainerOpacity        *float64 `json:"form_container_opacity,omitempty"`
	DescriptionContainerColor   string  `json:"description_container_color"`
	DescriptionContainerOpacity *float64 `json:"description_container_opacity,omitempty"`
	FormMetaBackground          string   `json:"form_meta_background,omitempty"`
	FormMetaBackgroundStart     string   `json:"form_meta_background_start,omitempty"`
	FormMetaBackgroundEnd       string   `json:"form_meta_background_end,omitempty"`
	FormMetaBackgroundOpacity   *float64 `json:"form_meta_background_opacity,omitempty"`
	FormMetaTextColor           string   `json:"form_meta_text_color,omitempty"`
}

// EmailTemplate representa una plantilla de correo
type EmailTemplate struct {
	ID           int64     `json:"id"`
	TemplateName string    `json:"template_name"`
	Subject      string    `json:"subject"`
	BodyHTML     string    `json:"body_html"`
	BodyText     string    `json:"body_text"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BulkEmailRequest representa la petición para envío masivo
type BulkEmailRequest struct {
	Recipients []string `json:"recipients"` // Lista de emails o "all_users", "form_submitters:FORM_ID"
	Subject    string   `json:"subject"`
	BodyHTML   string   `json:"body_html"`
	BodyText   string   `json:"body_text,omitempty"`
	TemplateID int64    `json:"template_id,omitempty"`
}

// SubmissionExport representa una respuesta exportable
type SubmissionExport struct {
	SubmissionID  int64             `json:"submission_id"`
	FormName      string            `json:"form_name"`
	SubmittedAt   time.Time         `json:"submitted_at"`
	FirstName     string            `json:"first_name"`
	LastName      string            `json:"last_name"`
	IDNumber      string            `json:"id_number"`
	Email         string            `json:"email"`
	Phone         string            `json:"phone"`
	Province      string            `json:"province"`
	Gender        string            `json:"gender"`
	Age           int               `json:"age"`
	BusinessName  string            `json:"business_name"`
	DynamicFields map[string]string `json:"dynamic_fields"`
}

// ResponseTimeline representa datos de respuestas por fecha
type ResponseTimeline struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// TopForm representa un formulario en el ranking
type TopForm struct {
	FormName        string `json:"form_name"`
	SubmissionCount int    `json:"submission_count"`
}

// ProvinceDistribution representa distribución por provincia
type ProvinceDistribution struct {
	Province string `json:"province"`
	Count    int    `json:"count"`
}

// UserMetrics representa métricas de usuarios
type UserMetrics struct {
	NewUsersMonth  int `json:"new_users_month"`
	ReturningUsers int `json:"returning_users"`
	InactiveUsers  int `json:"inactive_users"`
}

// FilterRequest representa los filtros de búsqueda
type FilterRequest struct {
	FormID     *int64  `json:"form_id,omitempty"`
	StartDate  *string `json:"start_date,omitempty"`
	EndDate    *string `json:"end_date,omitempty"`
	Province   *string `json:"province,omitempty"`
	SearchTerm *string `json:"search_term,omitempty"`
	Limit      int     `json:"limit,omitempty"`
	Offset     int     `json:"offset,omitempty"`
}
