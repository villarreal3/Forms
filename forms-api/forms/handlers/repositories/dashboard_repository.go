package repositories

import (
	"forms/models"
	"time"
)

// DashboardRepository define las operaciones de acceso a datos para dashboard
type DashboardRepository interface {
	// GetDashboardStats obtiene las estadísticas del dashboard
	GetDashboardStats() (*models.DashboardStats, error)
	
	// GetFormStats obtiene estadísticas por formulario
	GetFormStats() ([]models.FormStats, error)
	
	// GetRecentSubmissions obtiene las submissions recientes
	GetRecentSubmissions(limit int) ([]struct {
		ID          int64
		SubmittedAt time.Time
		FormName    string
		FirstName   string
		LastName    string
		Email       string
		IDNumber    string
	}, error)
	
	// GetResponseTimeline obtiene respuestas por fecha
	GetResponseTimeline(days int) ([]models.ResponseTimeline, error)
	
	// GetTopActiveForms obtiene formularios más activos
	GetTopActiveForms(limit int) ([]models.TopForm, error)
	
	// GetGeographicDistribution obtiene distribución por provincia
	GetGeographicDistribution() ([]models.ProvinceDistribution, error)
	
	// GetUserMetrics obtiene métricas de usuarios
	GetUserMetrics() (*models.UserMetrics, error)
}














