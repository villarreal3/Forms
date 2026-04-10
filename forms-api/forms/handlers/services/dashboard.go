package services

import (
	"log"
	"net/http"
	"strconv"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
)

// GetDashboardStats obtiene las estadísticas del dashboard
func GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	stats, err := repositories.DashboardRepo.GetDashboardStats()
	if err != nil {
		log.Printf("Error obteniendo estadísticas: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo estadísticas")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"stats":   stats,
	})
}

// GetFormStats obtiene estadísticas por formulario
func GetFormStats(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	formStats, err := repositories.DashboardRepo.GetFormStats()
	if err != nil {
		log.Printf("Error obteniendo estadísticas de formularios: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo estadísticas de formularios")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"stats":   formStats,
	})
}

// GetRecentSubmissions obtiene las submissions recientes
func GetRecentSubmissions(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	submissions, err := repositories.DashboardRepo.GetRecentSubmissions(limit)
	if err != nil {
		log.Printf("Error obteniendo submissions recientes: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo submissions recientes")
		return
	}

	// Convertir a formato de respuesta
	type SubmissionResponse struct {
		ID          int64  `json:"id"`
		SubmittedAt string `json:"submitted_at"`
		FormName    string `json:"form_name"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		IDNumber    string `json:"id_number"`
	}

	var response []SubmissionResponse
	for _, sub := range submissions {
		response = append(response, SubmissionResponse{
			ID:          sub.ID,
			SubmittedAt: sub.SubmittedAt.Format("2006-01-02 15:04:05"),
			FormName:    sub.FormName,
			FirstName:   sub.FirstName,
			LastName:    sub.LastName,
			Email:       sub.Email,
			IDNumber:    sub.IDNumber,
		})
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"submissions": response,
	})
}

// GetResponseTimeline obtiene respuestas por fecha
func GetResponseTimeline(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	timeline, err := repositories.DashboardRepo.GetResponseTimeline(days)
	if err != nil {
		log.Printf("Error obteniendo timeline de respuestas: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo timeline de respuestas")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    timeline,
	})
}

// GetTopActiveForms obtiene formularios más activos
func GetTopActiveForms(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	forms, err := repositories.DashboardRepo.GetTopActiveForms(limit)
	if err != nil {
		log.Printf("Error obteniendo formularios más activos: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo formularios más activos")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"forms":   forms,
	})
}

// GetGeographicDistribution obtiene distribución por provincia
func GetGeographicDistribution(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	distributions, err := repositories.DashboardRepo.GetGeographicDistribution()
	if err != nil {
		log.Printf("Error obteniendo distribución geográfica: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo distribución geográfica")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"provinces": distributions,
	})
}

// GetUserMetrics obtiene métricas de usuarios
func GetUserMetrics(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	metrics, err := repositories.DashboardRepo.GetUserMetrics()
	if err != nil {
		log.Printf("Error obteniendo métricas de usuarios: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo métricas de usuarios")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"metrics": metrics,
	})
}

