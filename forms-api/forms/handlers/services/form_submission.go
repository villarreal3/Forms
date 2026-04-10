package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"forms/database"
	"forms/models"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
)

// SubmitForm maneja el envío de un formulario usando sp_submit_form
func SubmitForm(w http.ResponseWriter, r *http.Request) {
	log.Printf("[SubmitForm] POST /api/forms/submit recibido")
	if r.Method != http.MethodPost {
		log.Printf("[SubmitForm] Error: método no permitido %s", r.Method)
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	var req models.SubmitFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[SubmitForm] Error decodificando JSON: %v", err)
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}
	log.Printf("[SubmitForm] form_id=%d, email=%s, id_number=%s", req.FormID, req.Email, req.IDNumber)

	// Validar campos requeridos
	if req.FirstName == "" || req.LastName == "" || req.IDNumber == "" || req.Email == "" || req.FormID == 0 {
		msg := "Campos requeridos faltantes: first_name, last_name, id_number, email, form_id"
		log.Printf("[SubmitForm] Error validación: %s", msg)
		helpers.RespondError(w, http.StatusBadRequest, msg)
		return
	}

	// Verificar que el formulario exista (la disponibilidad real se valida en el repository)
	exists, err := repositories.FormRepo.FormExists(int64(req.FormID))
	if err != nil {
		log.Printf("[SubmitForm] Error verificando formulario (form_id=%d): %v", req.FormID, err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error verificando formulario: %v", err))
		return
	}
	if !exists {
		log.Printf("[SubmitForm] Formulario no encontrado: form_id=%d", req.FormID)
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	// Validar ANTES de procesar si la persona ya está inscrita en este formulario
	// Esto es una validación adicional de seguridad
	duplicate, err := repositories.SubmissionRepo.CheckDuplicateSubmission(req.IDNumber, req.Email, int64(req.FormID))
	if err != nil {
		log.Printf("[SubmitForm] Error verificando inscripción existente: %v", err)
		// Continuar con el proceso, el stored procedure también validará
	} else if duplicate {
		log.Printf("[SubmitForm] Inscripción duplicada: form_id=%d, email=%s", req.FormID, req.Email)
		helpers.RespondError(w, http.StatusConflict, "Esta persona ya está inscrita en este formulario")
		return
	}

	// Ejecutar el stored procedure usando el repository
	result, userID, err := repositories.SubmissionRepo.SubmitForm(&req)
	if err != nil {
		log.Printf("[SubmitForm] Error en SubmitForm (repository): %v", err)
		var formUnavailable *repositories.FormUnavailableError
		if errors.As(err, &formUnavailable) && formUnavailable != nil {
			helpers.RespondError(w, formUnavailable.StatusCode, formUnavailable.Message)
			return
		}
		if err == sql.ErrNoRows {
			helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
			return
		}

		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error ejecutando SubmitForm: %v", err))
		return
	}

	// Verificar si la persona ya está inscrita en este formulario (result == 4)
	if result == 4 {
		log.Printf("[SubmitForm] Inscripción duplicada (result=4): form_id=%d", req.FormID)
		helpers.RespondError(w, http.StatusConflict, "Esta persona ya está inscrita en este formulario")
		return
	}
	log.Printf("[SubmitForm] Éxito: form_id=%d, result=%d, user_id=%d", req.FormID, result, userID)

	var messages = map[int]string{
		1: "Usuario existente encontrado",
		2: "Usuario actualizado",
		3: "Usuario nuevo creado",
	}

	// Obtener nombre del formulario para el correo usando el repository
	formName, err := repositories.FormRepo.GetFormName(int64(req.FormID))
	if err != nil {
		log.Printf("Error obteniendo nombre del formulario: %v", err)
		formName = "Formulario" // Valor por defecto si falla
	}

	// Enviar correo de confirmación en segundo plano
	go func() {
		userName := req.FirstName + " " + req.LastName
		if err := SendConfirmationEmail(req.Email, userName, formName, req.IDNumber); err != nil {
			log.Printf("Error enviando correo de confirmación: %v", err)
		}
	}()

	response := models.SubmitFormResponse{
		Success: true,
		Message: messages[result],
		Result:  result,
		UserID:  userID,
	}

	helpers.RespondJSON(w, http.StatusOK, response)
}

// GetUserByCredentials maneja la obtención de usuario por id_number y email (vw_users_credentials)
func GetUserByCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	var req models.GetUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	// Validar campos requeridos
	if req.IDNumber == "" || req.Email == "" {
		helpers.RespondError(w, http.StatusBadRequest, "Campos requeridos: id_number, email")
		return
	}

	// Obtener usuario usando el repository
	user, err := repositories.SubmissionRepo.GetUserByCredentials(req.IDNumber, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondError(w, http.StatusNotFound, "Usuario no encontrado")
			return
		}
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo usuario: %v", err))
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

// HealthCheck verifica el estado del servidor y la conexión a la base de datos
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := database.DB.Ping(); err != nil {
		helpers.RespondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  "unhealthy",
			"message": "No se pudo conectar a la base de datos",
			"error":   err.Error(),
		})
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"message": "Servidor y base de datos funcionando correctamente",
	})
}
