package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
	"forms/models"
)

// GetPublicForms: solo abiertos y vigentes (status=1, open_at <= NOW() <= expires_at).
func GetPublicForms(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	forms, err := repositories.PublicRepo.GetPublicForms()
	if err != nil {
		log.Printf("Error obteniendo formularios públicos: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo formularios: %v", err))
		return
	}

	helpers.RespondJSON(w, http.StatusOK, models.FormsListResponse{
		Success: true,
		Forms:   forms,
	})
}

// GetPublicForm: mismo criterio que el listado; adjunta campos si existen.
func GetPublicForm(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	form, err := repositories.PublicRepo.GetPublicForm(formID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado o no disponible")
			return
		}
		log.Printf("Error obteniendo formulario: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo formulario: %v", err))
		return
	}

	// Obtener campos del formulario
	fields, err := repositories.FormRepo.GetFields(formID)
	if err != nil {
		log.Printf("Error obteniendo campos: %v", err)
	} else {
		form.Fields = fields
	}

	helpers.RespondJSON(w, http.StatusOK, models.FormResponse{
		Success: true,
		Form:    *form,
	})
}

// GetPublicFormSections obtiene las secciones de un formulario público
func GetPublicFormSections(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	// Verificar que el formulario existe y está activo
	_, err := repositories.PublicRepo.GetPublicForm(formID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado o no disponible")
			return
		}
		log.Printf("Error verificando formulario: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error verificando formulario: %v", err))
		return
	}

	// Obtener secciones del formulario
	sections, err := repositories.FormRepo.GetFormSections(formID)
	if err != nil {
		log.Printf("Error obteniendo secciones: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo secciones: %v", err))
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"sections": sections,
	})
}
