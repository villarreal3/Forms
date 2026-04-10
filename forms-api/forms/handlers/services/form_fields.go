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
)

// UpdateField actualiza un campo de formulario
func UpdateField(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Extraer form_id y field_id de la URL usando mux.Vars
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

	var req struct {
		FieldLabel   string `json:"field_label"`
		FieldName    string `json:"field_name"`
		FieldType    string `json:"field_type"`
		FieldOptions string `json:"field_options"`
		IsRequired   bool   `json:"is_required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	// Validar campos requeridos
	if req.FieldLabel == "" || req.FieldName == "" || req.FieldType == "" {
		helpers.RespondError(w, http.StatusBadRequest, "field_label, field_name y field_type son requeridos")
		return
	}

	success, err := repositories.FormRepo.UpdateField(fieldID, formID, req.FieldLabel, req.FieldName, req.FieldType, req.FieldOptions, req.IsRequired)
	if err != nil {
		log.Printf("Error actualizando campo: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error actualizando campo: %v", err))
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusNotFound, "Campo no encontrado o no pertenece al formulario")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Campo actualizado correctamente",
	})
}

// DeleteField elimina un campo de formulario
func DeleteField(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Extraer form_id y field_id de la URL usando mux.Vars
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

	success, errorMessage, err := repositories.FormRepo.DeleteField(fieldID, formID)
	if err != nil {
		log.Printf("Error eliminando campo: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error eliminando campo: %v", err))
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusBadRequest, errorMessage)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Campo eliminado correctamente",
	})
}

// ReorderField reordena un campo de formulario (mueve arriba o abajo)
func ReorderField(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Extraer form_id y field_id de la URL usando mux.Vars
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

	var req struct {
		Direction string `json:"direction"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	// Validar dirección
	if req.Direction != "up" && req.Direction != "down" {
		helpers.RespondError(w, http.StatusBadRequest, "direction debe ser 'up' o 'down'")
		return
	}

	success, errorMessage, err := repositories.FormRepo.ReorderField(fieldID, formID, req.Direction)
	if err != nil {
		log.Printf("Error reordenando campo: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error reordenando campo: %v", err))
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusBadRequest, errorMessage)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Campo reordenado correctamente",
	})
}
