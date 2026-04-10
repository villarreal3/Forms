package handlers

import (
	"forms/handlers/repositories"
	"forms/handlers/services"
)

// Inicializar repositories al cargar el package
func init() {
	repositories.Init()
}

// Re-exportar funciones de services
var (
	// Forms
	CreateForm              = services.CreateForm
	GetForms                = services.GetForms
	GetForm                 = services.GetForm
	PatchFormMetadata       = services.PatchFormMetadata
	UpdateFormSchema        = services.UpdateFormSchema
	CreateField             = services.CreateField
	UpdateField             = services.UpdateField
	DeleteField             = services.DeleteField
	ReorderField            = services.ReorderField
	AssignFieldToSection    = services.AssignFieldToSection
	GetCommonFieldTemplates = services.GetCommonFieldTemplates
	// Sections
	CreateSection   = services.CreateSection
	GetFormSections = services.GetFormSections
	UpdateSection   = services.UpdateSection
	DeleteSection   = services.DeleteSection
	ReorderSection  = services.ReorderSection
	AddCommonFields = services.AddCommonFields
	CloseForm       = services.CloseForm
	PublishForm     = services.PublishForm
	SetFormDraft    = services.SetFormDraft
	OpenForm        = services.OpenForm

	// Dashboard
	GetDashboardStats         = services.GetDashboardStats
	GetFormStats              = services.GetFormStats
	GetRecentSubmissions      = services.GetRecentSubmissions
	GetResponseTimeline       = services.GetResponseTimeline
	GetTopActiveForms         = services.GetTopActiveForms
	GetGeographicDistribution = services.GetGeographicDistribution
	GetUserMetrics            = services.GetUserMetrics

	// Submissions
	GetSubmissions           = services.GetSubmissions
	UpdateAttendance         = services.UpdateAttendance
	AutoSubmitFromPrevious   = services.AutoSubmitFromPrevious

	// Export
	ExportSubmissions = services.ExportSubmissions

	// Email
	InitEmailConfig       = services.InitEmailConfig
	SendBulkEmail         = services.SendBulkEmail
	SendConfirmationEmail = services.SendConfirmationEmail

	// Upload
	InitImageConfig = services.InitImageConfig
	UploadFormImage = services.UploadFormImage
	ServeImage      = services.ServeImage
	DeleteFormImage = services.DeleteFormImage

	// Customization
	GetFormCustomization        = services.GetFormCustomization
	GetPublicFormLogo             = services.GetPublicFormLogo
	CreateOrUpdateCustomization = services.CreateOrUpdateCustomization

	// Public
	GetPublicForms        = services.GetPublicForms
	GetPublicForm         = services.GetPublicForm
	GetPublicFormSections = services.GetPublicFormSections

	// Form Submission
	SubmitForm           = services.SubmitForm
	GetUserByCredentials = services.GetUserByCredentials
	HealthCheck          = services.HealthCheck
)
