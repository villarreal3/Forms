package repositories

import (
	"forms/models"
)

// PublicRepository define las operaciones de acceso a datos para endpoints públicos
type PublicRepository interface {
	// GetPublicForms: solo formularios abiertos y vigentes (status=1, open_at <= NOW() <= expires_at).
	GetPublicForms() ([]models.Form, error)

	// GetPublicForm: mismo criterio de visibilidad que el listado.
	GetPublicForm(formID int64) (*models.Form, error)
}














