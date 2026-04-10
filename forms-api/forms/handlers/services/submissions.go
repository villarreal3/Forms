package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
	"forms/models"
)

// GetSubmissions obtiene las respuestas con filtros
func GetSubmissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Obtener parámetros de query
	query := r.URL.Query()
	var filter models.FilterRequest

	if formIDStr := query.Get("form_id"); formIDStr != "" {
		if formID, err := strconv.ParseInt(formIDStr, 10, 64); err == nil {
			filter.FormID = &formID
		}
	}

	if startDate := query.Get("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}

	if endDate := query.Get("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	if province := query.Get("province"); province != "" {
		filter.Province = &province
	}

	if searchTerm := query.Get("search_term"); searchTerm != "" {
		filter.SearchTerm = &searchTerm
	}

	limit := 100
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	filter.Limit = limit

	offset := 0
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o > 0 {
			offset = o
		}
	}
	filter.Offset = offset

	// Usar repository para obtener submissions
	submissionMaps, err := repositories.SubmissionRepo.GetSubmissions(
		filter.FormID,
		filter.StartDate,
		filter.EndDate,
		filter.Province,
		filter.SearchTerm,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		log.Printf("Error obteniendo submissions: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo respuestas: %v", err))
		return
	}

	type Submission struct {
		ID           int64     `json:"id"`
		UserID       int64     `json:"user_id"`
		FormName     string    `json:"form_name"`
		SubmittedAt  time.Time `json:"submitted_at"`
		FirstName    string    `json:"first_name"`
		LastName     string    `json:"last_name"`
		IDNumber     string    `json:"id_number"`
		Email        string    `json:"email"`
		Phone        string    `json:"phone,omitempty"`
		Province     string    `json:"province,omitempty"`
		Gender       string    `json:"gender,omitempty"`
		Age          int       `json:"age,omitempty"`
		BusinessName string    `json:"business_name,omitempty"`
		Attended     bool      `json:"attended"`
	}

	var submissions []Submission
	for _, subMap := range submissionMaps {
		// Helpers para evitar panic cuando alguna clave viene nil o con tipo inesperado
		getString := func(key string) string {
			v, ok := subMap[key]
			if !ok || v == nil {
				return ""
			}
			if s, ok := v.(string); ok {
				return s
			}
			return fmt.Sprint(v)
		}

		getInt64 := func(key string) int64 {
			v, ok := subMap[key]
			if !ok || v == nil {
				return 0
			}
			switch t := v.(type) {
			case int64:
				return t
			case int:
				return int64(t)
			case uint64:
				return int64(t)
			default:
				return 0
			}
		}

		sub := Submission{
			ID:        getInt64("id"),
			UserID:    getInt64("user_id"),
			FormName:  getString("form_name"),
			FirstName: getString("first_name"),
			LastName:  getString("last_name"),
			IDNumber:  getString("id_number"),
			Email:     getString("email"),
		}

		if submittedAt, ok := subMap["submitted_at"].(time.Time); ok {
			sub.SubmittedAt = submittedAt
		}

		if phone := getString("phone"); phone != "" {
			sub.Phone = phone
		}
		if province := getString("province"); province != "" {
			sub.Province = province
		}
		if gender := getString("gender"); gender != "" {
			sub.Gender = gender
		}
		if v, ok := subMap["age"].(int); ok {
			sub.Age = v
		}
		if businessName := getString("business_name"); businessName != "" {
			sub.BusinessName = businessName
		}
		// attended puede venir como bool, int (0/1) o string ("0"/"1")
		if raw, ok := subMap["attended"]; ok && raw != nil {
			switch v := raw.(type) {
			case bool:
				sub.Attended = v
			case int:
				sub.Attended = v != 0
			case int64:
				sub.Attended = v != 0
			case uint8:
				sub.Attended = v != 0
			case string:
				l := strings.ToLower(strings.TrimSpace(v))
				sub.Attended = l == "1" || l == "true" || l == "t" || l == "sí" || l == "si"
			}
		}

		submissions = append(submissions, sub)
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"submissions": submissions,
		"count":       len(submissions),
	})
}

// AutoSubmitFromPrevious crea una nueva inscripción en el formulario indicado
// reutilizando la última inscripción previa (por fecha) encontrada en otros
// formularios con el mismo schema, identificando por cédula (id_number) únicamente.
func AutoSubmitFromPrevious(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	var req struct {
		FormID   int64  `json:"form_id"`
		IDNumber string `json:"id_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	if req.FormID <= 0 {
		helpers.RespondError(w, http.StatusBadRequest, "form_id es requerido")
		return
	}
	req.IDNumber = strings.TrimSpace(req.IDNumber)
	if req.IDNumber == "" {
		helpers.RespondError(w, http.StatusBadRequest, "id_number es requerido")
		return
	}

	// #region agent log
	if f, err := os.OpenFile("/home/davillarreal/forms-automate/.cursor/debug-784680.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		b, _ := json.Marshal(map[string]interface{}{"sessionId": "784680", "location": "submissions.go:AutoSubmitFromPrevious", "message": "request received", "hypothesisId": "H1", "data": map[string]interface{}{"form_id": req.FormID, "id_number": req.IDNumber}, "timestamp": time.Now().UnixMilli()})
		f.Write(append(b, '\n'))
		f.Close()
	}
	// #endregion

	created, err := repositories.SubmissionRepo.AutoSubmitFromPrevious(req.FormID, req.IDNumber)
	if err != nil {
		log.Printf("Error en AutoSubmitFromPrevious: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error al intentar crear inscripción automática")
		return
	}

	if !created {
		helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"success": false,
			"found":   false,
			"message": "No se encontró una inscripción previa compatible para esta cédula.",
		})
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"found":   true,
		"message": "Inscripción creada automáticamente a partir del último formulario compatible.",
	})
}

// UpdateAttendance actualiza o crea el registro de asistencia
func UpdateAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	var req struct {
		UserID   int64 `json:"user_id"`
		FormID   int64 `json:"form_id"`
		Attended bool  `json:"attended"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	if req.UserID <= 0 || req.FormID <= 0 {
		helpers.RespondError(w, http.StatusBadRequest, "user_id y form_id son requeridos")
		return
	}

	// Usar repository para actualizar asistencia
	attendanceID, err := repositories.SubmissionRepo.UpdateAttendance(req.UserID, req.FormID, req.Attended)
	if err != nil {
		log.Printf("Error actualizando asistencia: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error actualizando registro de asistencia")
		return
	}

	message := "Asistencia actualizada correctamente"
	if attendanceID > 0 {
		// Verificar si es nuevo o actualizado
		_, existingID, _ := repositories.SubmissionRepo.GetAttendance(req.UserID, req.FormID)
		if existingID == 0 {
			message = "Asistencia registrada correctamente"
		}
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": message,
		"id":      attendanceID,
	})
}

