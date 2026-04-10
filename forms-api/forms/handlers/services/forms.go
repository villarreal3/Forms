package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
	"forms/models"
)

// CreateForm crea un nuevo formulario
func CreateForm(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	var req models.CreateFormRequest
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}

	if req.FormName == "" {
		helpers.RespondError(w, http.StatusBadRequest, "El nombre del formulario es requerido")
		return
	}
	if req.EventDate == "" {
		helpers.RespondError(w, http.StatusBadRequest, "La fecha del evento es requerida")
		return
	}
	if req.OpenAt == "" {
		helpers.RespondError(w, http.StatusBadRequest, "La fecha de apertura es requerida")
		return
	}
	if req.ExpiresAt == "" {
		helpers.RespondError(w, http.StatusBadRequest, "La fecha de expiración es requerida")
		return
	}

	eventDate, err := helpers.ParseDateTime(req.EventDate)
	if err != nil {
		log.Printf("Error parseando fecha de evento (%s): %v", req.EventDate, err)
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	openAt, err := helpers.ParseDateTime(req.OpenAt)
	if err != nil {
		log.Printf("Error parseando fecha de apertura (%s): %v", req.OpenAt, err)
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	expiresAt, err := helpers.ParseExpiry(req.ExpiresAt)
	if err != nil {
		log.Printf("Error parseando fecha de expiración (%s): %v", req.ExpiresAt, err)
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if openAt.After(expiresAt) {
		helpers.RespondError(w, http.StatusBadRequest, "La fecha de apertura debe ser anterior o igual a la fecha de expiración")
		return
	}

	description := req.Description
	if description == "" {
		description = ""
	}

	var confPtr, ubiPtr *string
	if s := strings.TrimSpace(req.Conferencista); s != "" {
		confPtr = &s
	}
	if s := strings.TrimSpace(req.Ubicacion); s != "" {
		ubiPtr = &s
	}

	// Por diseño: al crear un formulario entra en modo borrador (no público)
	formID, err := repositories.FormRepo.Create(req.FormName, description, eventDate, openAt, expiresAt, true, req.InscriptionLimit, confPtr, ubiPtr)
	if err != nil {
		log.Printf("Error creando formulario: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error creando formulario: %v", err))
		return
	}

	// Table-per-Form: la plantilla ya está en forms.schema al crear
	if req.UseDefaultTemplate {
		fields, _ := repositories.FormRepo.GetFields(formID)
		fieldsAdded := len(fields)
		helpers.RespondJSON(w, http.StatusCreated, map[string]interface{}{
			"success":      true,
			"message":      "Formulario creado con esquema por defecto (table-per-form)",
			"form_id":      formID,
			"fields_added": fieldsAdded,
		})
	} else {
		helpers.RespondJSON(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Formulario creado correctamente",
			"form_id": formID,
		})
	}
}

// GetForms obtiene la lista de todos los formularios
func GetForms(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	forms, err := repositories.FormRepo.GetForms()
	if err != nil {
		log.Printf("Error obteniendo formularios: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo formularios")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, models.FormsListResponse{
		Success: true,
		Forms:   forms,
	})
}

// GetForm obtiene un formulario por ID con sus campos
func GetForm(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	form, err := repositories.FormRepo.GetByID(formID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
			return
		}
		log.Printf("Error obteniendo formulario: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo formulario")
		return
	}

	// Campos del formulario
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

// UpdateFormSchema PUT body: { "schema": { "sections": [...], "fields": [...] } } o directamente { "sections", "fields" }
func UpdateFormSchema(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPut) {
		return
	}
	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil || len(raw) == 0 {
		helpers.RespondError(w, http.StatusBadRequest, "cuerpo requerido")
		return
	}
	var wrap struct {
		Schema json.RawMessage `json:"schema"`
	}
	schemaBytes := raw
	if json.Unmarshal(raw, &wrap) == nil && len(wrap.Schema) > 0 {
		schemaBytes = wrap.Schema
	}
	var check struct {
		Sections []interface{} `json:"sections"`
		Fields   []interface{} `json:"fields"`
	}
	if json.Unmarshal(schemaBytes, &check) != nil || (len(check.Sections) == 0 && len(check.Fields) == 0) {
		helpers.RespondError(w, http.StatusBadRequest, "schema debe incluir sections y/o fields")
		return
	}
	exists, err := repositories.FormRepo.FormExists(formID)
	if err != nil || !exists {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}
	if err := repositories.FormRepo.UpdateFormSchema(formID, schemaBytes); err != nil {
		log.Printf("UpdateFormSchema: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Esquema actualizado; si ya había respuestas, la tabla de respuestas no se recreó",
	})
}

// CreateField crea un nuevo campo en un formulario
func CreateField(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	var req models.CreateFieldRequest
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}

	if req.FieldLabel == "" || req.FieldName == "" || req.FieldType == "" {
		helpers.RespondError(w, http.StatusBadRequest, "field_label, field_name y field_type son requeridos")
		return
	}

	if !helpers.ValidFieldTypes[req.FieldType] {
		helpers.RespondError(w, http.StatusBadRequest,
			"Tipo de campo inválido. Tipos válidos: text, number, email, tel, select, radio, checkbox, textarea, date")
		return
	}

	exists, err := repositories.FormRepo.FormExists(formID)
	if err != nil || !exists {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	dbFieldType := helpers.MapDBFieldType(req.FieldType, req.FieldName)
	fieldOptions := ""
	if req.FieldOptions != "" {
		fieldOptions = req.FieldOptions
	}

	fieldID, err := repositories.FormRepo.CreateField(formID, req.FieldLabel, req.FieldName, dbFieldType, fieldOptions, req.IsRequired)
	if err != nil {
		log.Printf("Error creando campo: %v", err)
		if strings.Contains(err.Error(), "not supported") {
			helpers.RespondError(w, http.StatusNotImplemented, "Table-per-Form: edite forms.schema (PATCH) en lugar de crear campos sueltos")
			return
		}
		helpers.RespondError(w, http.StatusInternalServerError, "Error creando campo")
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"success":  true,
		"message":  "Campo creado correctamente",
		"field_id": fieldID,
	})
}

// GetCommonFieldTemplates obtiene las plantillas de campos comunes
func GetCommonFieldTemplates(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	templates := models.GetCommonFieldTemplates()
	categories := map[string][]models.CommonFieldTemplate{
		"personal": models.GetFieldsByCategory("personal"),
		"contact":  models.GetFieldsByCategory("contact"),
		"business": models.GetFieldsByCategory("business"),
		"social":   models.GetFieldsByCategory("social"),
		"consent":  models.GetFieldsByCategory("consent"),
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"templates":  templates,
		"categories": categories,
	})
}

// AddCommonFields agrega campos comunes a un formulario
func AddCommonFields(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	var req struct {
		FieldNames []string `json:"field_names"`
	}
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}
	if len(req.FieldNames) == 0 {
		helpers.RespondError(w, http.StatusBadRequest, "Debe especificar al menos un campo")
		return
	}

	exists, err := repositories.FormRepo.FormExists(formID)
	if err != nil || !exists {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	tmplMap := helpers.TemplatesMap()

	// Campos existentes
	existing := make(map[string]bool)
	fieldNames, err := repositories.FormRepo.GetFieldNames(formID)
	if err == nil {
		for _, name := range fieldNames {
			existing[name] = true
		}
	}

	var addedIDs []int64
	var skipped []string

	for _, name := range req.FieldNames {
		if existing[name] {
			skipped = append(skipped, name)
			continue
		}
		tmpl, ok := tmplMap[name]
		if !ok {
			skipped = append(skipped, name+" (no encontrado)")
			continue
		}

		dbType := helpers.MapDBFieldType(tmpl.Type, tmpl.Name)
		fieldOptions := ""
		if tmpl.Options != "" {
			fieldOptions = tmpl.Options
		}

		fieldID, err := repositories.FormRepo.CreateField(formID, tmpl.Label, tmpl.Name, dbType, fieldOptions, tmpl.IsRequired)
		if err != nil {
			if strings.Contains(err.Error(), "not supported") {
				helpers.RespondError(w, http.StatusNotImplemented, "Table-per-Form: no se pueden añadir campos con este endpoint; use forms.schema")
				return
			}
			log.Printf("Error insertando campo %s: %v", name, err)
			skipped = append(skipped, name+" (error)")
			continue
		}

		if fieldID > 0 {
			addedIDs = append(addedIDs, fieldID)
		}
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":        true,
		"message":        fmt.Sprintf("Agregados %d campo(s) correctamente", len(addedIDs)),
		"added_fields":   addedIDs,
		"skipped_fields": skipped,
	})
}

// CloseForm cierra un formulario usando el stored procedure
func CloseForm(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	adminUserID, _ := r.Context().Value("user_id").(int64)
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	userAgent := r.Header.Get("User-Agent")

	success, err := repositories.FormRepo.CloseForm(formID, adminUserID, ipAddress, userAgent)
	if err != nil {
		log.Printf("Error ejecutando sp_close_form: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error cerrando formulario")
		return
	}

	if !success {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Formulario cerrado correctamente",
	})
}

// PublishForm marca el formulario como público (is_draft=0).
func PublishForm(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	ok, err := repositories.FormRepo.SetDraftStatus(formID, false)
	if err != nil {
		log.Printf("Error seteando formulario público: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error publicando formulario")
		return
	}
	if !ok {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Formulario publicado correctamente",
	})
}

// SetFormDraft marca el formulario como borrador (is_draft=1).
func SetFormDraft(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	ok, err := repositories.FormRepo.SetDraftStatus(formID, true)
	if err != nil {
		log.Printf("Error seteando formulario en borrador: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error moviendo formulario a borrador")
		return
	}
	if !ok {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Formulario movido a borrador correctamente",
	})
}

// PatchFormMetadata PATCH body: { "conferencista": "...", "ubicacion": "..." } — omitir clave = no cambiar; "" = NULL en BD.
func PatchFormMetadata(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPatch) {
		return
	}
	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}
	var req models.PatchFormMetadataRequest
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}
	if req.Conferencista == nil && req.Ubicacion == nil {
		helpers.RespondError(w, http.StatusBadRequest, "Indique al menos conferencista o ubicacion")
		return
	}
	exists, err := repositories.FormRepo.FormExists(formID)
	if err != nil || !exists {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}
	var cPtr, uPtr *string
	if req.Conferencista != nil {
		s := strings.TrimSpace(*req.Conferencista)
		cPtr = &s
	}
	if req.Ubicacion != nil {
		s := strings.TrimSpace(*req.Ubicacion)
		uPtr = &s
	}
	if err := repositories.FormRepo.UpdateFormMetadata(formID, cPtr, uPtr); err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
			return
		}
		log.Printf("PatchFormMetadata: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Metadatos actualizados",
	})
}

// OpenForm marca el formulario como abierto para envíos (is_closed=0, is_draft=0).
func OpenForm(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	formID, ok := helpers.ParseFormIDFromVars(w, r)
	if !ok {
		return
	}

	ok, err := repositories.FormRepo.OpenForm(formID)
	if err != nil {
		log.Printf("Error moviendo formulario a abierto: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error moviendo formulario a abierto")
		return
	}
	if !ok {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Formulario abierto correctamente",
	})
}
