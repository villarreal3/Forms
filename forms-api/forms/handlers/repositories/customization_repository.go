package repositories

import (
	"forms/models"
)

// CustomizationRepository define las operaciones de acceso a datos para personalización
type CustomizationRepository interface {
	// GetFormCustomization obtiene la personalización de un formulario
	GetFormCustomization(formID int64) (*models.FormCustomization, error)

	// CreateOrUpdateCustomization crea o actualiza la personalización
	CreateOrUpdateCustomization(formID int64, customization *models.CreateCustomizationRequest) (int64, error)

	// GetCustomizationID obtiene el ID de personalización si existe
	GetCustomizationID(formID int64) (int64, error)

	// UpdateLogo actualiza el logo de un formulario (isMobile indica si es para móvil)
	UpdateLogo(formID int64, logoURL string, isMobile bool) (bool, error)

	// GetLogo obtiene el logo de un formulario
	GetLogo(formID int64) (string, bool, error)

	// DeleteLogo elimina el logo de un formulario (isMobile indica si es para móvil)
	DeleteLogo(formID int64, isMobile bool) (bool, error)
}
