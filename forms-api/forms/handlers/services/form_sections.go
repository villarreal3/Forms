package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
	"forms/models"
)

// CreateSection crea una nueva sección
func CreateSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formIDStr := vars["form_id"]

	if formIDStr == "" {
		helpers.RespondError(w, http.StatusBadRequest, "form_id es requerido")
		return
	}

	formID, err := strconv.ParseInt(formIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "form_id inválido")
		return
	}

	var req models.CreateSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	if req.SectionTitle == "" || req.SectionIcon == "" {
		helpers.RespondError(w, http.StatusBadRequest, "section_title y section_icon son requeridos")
		return
	}

	sectionID, err := repositories.FormRepo.CreateSection(formID, req.SectionTitle, req.SectionIcon)
	if err != nil {
		log.Printf("Error creando sección: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error creando sección: %v", err))
		return
	}

	if sectionID == 0 {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"success":    true,
		"message":    "Sección creada correctamente",
		"section_id": sectionID,
	})
}

// GetFormSections obtiene todas las secciones de un formulario
func GetFormSections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formIDStr := vars["form_id"]

	if formIDStr == "" {
		helpers.RespondError(w, http.StatusBadRequest, "form_id es requerido")
		return
	}

	formID, err := strconv.ParseInt(formIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "form_id inválido")
		return
	}

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

// UpdateSection actualiza una sección
func UpdateSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formIDStr := vars["form_id"]
	sectionIDStr := vars["section_id"]

	if formIDStr == "" || sectionIDStr == "" {
		helpers.RespondError(w, http.StatusBadRequest, "form_id y section_id son requeridos")
		return
	}

	formID, err := strconv.ParseInt(formIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "form_id inválido")
		return
	}

	sectionID, err := strconv.ParseInt(sectionIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "section_id inválido")
		return
	}

	var req models.UpdateSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	if req.SectionTitle == "" || req.SectionIcon == "" {
		helpers.RespondError(w, http.StatusBadRequest, "section_title y section_icon son requeridos")
		return
	}

	success, err := repositories.FormRepo.UpdateSection(sectionID, formID, req.SectionTitle, req.SectionIcon)
	if err != nil {
		log.Printf("Error actualizando sección: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error actualizando sección: %v", err))
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusNotFound, "Sección no encontrada o no pertenece al formulario")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sección actualizada correctamente",
	})
}

// DeleteSection elimina una sección
func DeleteSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formIDStr := vars["form_id"]
	sectionIDStr := vars["section_id"]

	if formIDStr == "" || sectionIDStr == "" {
		helpers.RespondError(w, http.StatusBadRequest, "form_id y section_id son requeridos")
		return
	}

	formID, err := strconv.ParseInt(formIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "form_id inválido")
		return
	}

	sectionID, err := strconv.ParseInt(sectionIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "section_id inválido")
		return
	}

	success, errorMessage, err := repositories.FormRepo.DeleteSection(sectionID, formID)
	if err != nil {
		log.Printf("Error eliminando sección: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error eliminando sección: %v", err))
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusBadRequest, errorMessage)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sección eliminada correctamente",
	})
}

// ReorderSection reordena una sección
func ReorderSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formIDStr := vars["form_id"]
	sectionIDStr := vars["section_id"]

	if formIDStr == "" || sectionIDStr == "" {
		helpers.RespondError(w, http.StatusBadRequest, "form_id y section_id son requeridos")
		return
	}

	formID, err := strconv.ParseInt(formIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "form_id inválido")
		return
	}

	sectionID, err := strconv.ParseInt(sectionIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "section_id inválido")
		return
	}

	var req models.ReorderSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	if req.Direction != "up" && req.Direction != "down" {
		helpers.RespondError(w, http.StatusBadRequest, "direction debe ser 'up' o 'down'")
		return
	}

	success, errorMessage, err := repositories.FormRepo.ReorderSection(sectionID, formID, req.Direction)
	if err != nil {
		log.Printf("Error reordenando sección: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error reordenando sección: %v", err))
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusBadRequest, errorMessage)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sección reordenada correctamente",
	})
}

// AssignFieldToSection asigna un campo a una sección
func AssignFieldToSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formIDStr := vars["form_id"]
	fieldIDStr := vars["field_id"]

	if formIDStr == "" || fieldIDStr == "" {
		helpers.RespondError(w, http.StatusBadRequest, "form_id y field_id son requeridos")
		return
	}

	formID, err := strconv.ParseInt(formIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "form_id inválido")
		return
	}

	fieldID, err := strconv.ParseInt(fieldIDStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "field_id inválido")
		return
	}

	var req models.AssignFieldToSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	var sectionID int64
	if req.SectionID != nil {
		sectionID = *req.SectionID
	}

	success, err := repositories.FormRepo.AssignFieldToSection(fieldID, sectionID, formID)
	if err != nil {
		log.Printf("Error asignando campo a sección: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error asignando campo a sección: %v", err))
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusBadRequest, "Campo o sección no encontrados")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Campo asignado a sección correctamente",
	})
}
