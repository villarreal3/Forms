package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"forms/database"
	"forms/models"

	"github.com/gorilla/mux"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
)

// GetFormCustomization obtiene la personalización de un formulario
func GetFormCustomization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "ID de formulario inválido")
		return
	}

	// Usar el repositorio que incluye logo_url_mobile
	customization, err := repositories.CustomizationRepo.GetFormCustomization(formID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Devolver valores por defecto
			defaultCustomization := models.FormCustomization{
				FormID:                      formID,
				PrimaryColor:                "#3B82F6",
				SecondaryColor:              "#1E40AF",
				BackgroundColor:             "#FFFFFF",
				TextColor:                   "#1F2937",
				TitleColor:                  "#FFFFFF",
				FontFamily:                  "Arial, sans-serif",
				ButtonStyle:                 "rounded",
				FormContainerColor:          "#FFFFFF",
				FormContainerOpacity:        1.0,
				DescriptionContainerColor:   "#FFFFFF",
				DescriptionContainerOpacity: 1.0,
				FormMetaBackground:          "#3B82F6",
				FormMetaBackgroundStart:     "#3B82F6",
				FormMetaBackgroundEnd:       "#1E40AF",
				FormMetaBackgroundOpacity:   1.0,
				FormMetaTextColor:           "#FFFFFF",
			}
			helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
				"success":       true,
				"customization": defaultCustomization,
				"is_default":    true,
			})
			return
		}
		log.Printf("Error obteniendo personalización: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo personalización: %v", err))
		return
	}

	// Valores por defecto si no están configurados
	if customization.TitleColor == "" {
		customization.TitleColor = "#FFFFFF"
	}
	if customization.FormContainerColor == "" {
		customization.FormContainerColor = "#FFFFFF"
	}
	if customization.DescriptionContainerColor == "" {
		customization.DescriptionContainerColor = "#FFFFFF"
	}
	if customization.FormMetaBackground == "" {
		customization.FormMetaBackground = "#3B82F6"
	}
	if customization.FormMetaBackgroundStart == "" {
		customization.FormMetaBackgroundStart = customization.FormMetaBackground
	}
	if customization.FormMetaBackgroundEnd == "" {
		customization.FormMetaBackgroundEnd = "#1E40AF"
	}
	if customization.FormMetaTextColor == "" {
		customization.FormMetaTextColor = "#FFFFFF"
	}

	log.Printf("[GetFormCustomization] Personalización obtenida para form_id=%d: logo_url=%s, logo_url_mobile=%s",
		formID, customization.LogoURL, customization.LogoURLMobile)

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":       true,
		"customization": customization,
		"is_default":    false,
	})
}

// GetPublicFormLogo devuelve solo las URLs de logo (escritorio y móvil) para el formulario público.
// Endpoint ligero pensado para cargar la imagen de forma independiente del resto de la personalización.
func GetPublicFormLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "ID de formulario inválido")
		return
	}

	customization, err := repositories.CustomizationRepo.GetFormCustomization(formID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
				"success":         true,
				"logo_url":        "",
				"logo_url_mobile": "",
			})
			return
		}
		log.Printf("[GetPublicFormLogo] Error obteniendo logos form_id=%d: %v", formID, err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error obteniendo logos: %v", err))
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":         true,
		"logo_url":        customization.LogoURL,
		"logo_url_mobile": customization.LogoURLMobile,
	})
}

// CreateOrUpdateCustomization crea o actualiza la personalización de un formulario
func CreateOrUpdateCustomization(w http.ResponseWriter, r *http.Request) {
	log.Printf(">>> CreateOrUpdateCustomization llamado - Método: %s, Path: %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	vars := mux.Vars(r)
	formID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("ERROR: ID de formulario inválido: %s", vars["id"])
		helpers.RespondError(w, http.StatusBadRequest, "ID de formulario inválido")
		return
	}

	var req models.CreateCustomizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR decodificando JSON: %v", err)
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	// Usar el form_id de la URL
	req.FormID = formID

	// Log inmediato después de decodificar JSON
	log.Printf("=== INICIO CreateOrUpdateCustomization ===")
	log.Printf("Método HTTP: %s, URL: %s, form_id=%d", r.Method, r.URL.Path, req.FormID)
	log.Printf("JSON decodificado RAW - form_container_color='%s' (len=%d), form_container_opacity=%v, description_container_color='%s' (len=%d), description_container_opacity=%v",
		req.FormContainerColor, len(req.FormContainerColor), req.FormContainerOpacity,
		req.DescriptionContainerColor, len(req.DescriptionContainerColor), req.DescriptionContainerOpacity)
	log.Printf("JSON decodificado RAW (meta strip) - form_meta_background='%s' (len=%d), form_meta_background_opacity=%v, form_meta_text_color='%s'",
		req.FormMetaBackground, len(req.FormMetaBackground), req.FormMetaBackgroundOpacity, req.FormMetaTextColor)
	log.Printf("JSON decodificado RAW (meta strip gradiente) - form_meta_background_start='%s', form_meta_background_end='%s'",
		req.FormMetaBackgroundStart, req.FormMetaBackgroundEnd)
	log.Printf("Todos los campos recibidos - PrimaryColor=%s, SecondaryColor=%s, BackgroundColor=%s, TextColor=%s, TitleColor=%s, LogoURL=%s, FontFamily=%s, ButtonStyle=%s",
		req.PrimaryColor, req.SecondaryColor, req.BackgroundColor, req.TextColor, req.TitleColor, req.LogoURL, req.FontFamily, req.ButtonStyle)

	// Log para debugging
	log.Printf("Recibida personalización para form_id=%d: form_container_color=%s, form_container_opacity=%v, description_container_color=%s, description_container_opacity=%v",
		req.FormID, req.FormContainerColor, req.FormContainerOpacity, req.DescriptionContainerColor, req.DescriptionContainerOpacity)

	// Validar que el formulario existe
	var exists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM forms WHERE id = ?", req.FormID).Scan(&exists)
	if err != nil || exists == 0 {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	// Validar colores (formato hex)
	if !isValidHexColor(req.PrimaryColor) || !isValidHexColor(req.SecondaryColor) ||
		!isValidHexColor(req.BackgroundColor) || !isValidHexColor(req.TextColor) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color inválido. Use formato hexadecimal (ej: #3B82F6)")
		return
	}

	// Validar color de título (opcional, usar valor por defecto si no se proporciona)
	if req.TitleColor != "" && !isValidHexColor(req.TitleColor) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color del título inválido")
		return
	}

	// Validar colores de contenedores (opcional, usar valores por defecto si no se proporcionan)
	if req.FormContainerColor != "" && !isValidHexColor(req.FormContainerColor) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color del contenedor del formulario inválido")
		return
	}
	if req.DescriptionContainerColor != "" && !isValidHexColor(req.DescriptionContainerColor) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color del contenedor de descripción inválido")
		return
	}
	if req.FormMetaBackground != "" && !isValidHexColor(req.FormMetaBackground) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color del fondo de la franja meta inválido")
		return
	}
	if req.FormMetaBackgroundStart != "" && !isValidHexColor(req.FormMetaBackgroundStart) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color inicial del fondo de la franja meta inválido")
		return
	}
	if req.FormMetaBackgroundEnd != "" && !isValidHexColor(req.FormMetaBackgroundEnd) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color final del fondo de la franja meta inválido")
		return
	}
	if req.FormMetaTextColor != "" && !isValidHexColor(req.FormMetaTextColor) {
		helpers.RespondError(w, http.StatusBadRequest, "Formato de color del texto de la franja meta inválido")
		return
	}

	// Validar opacidades (0.0 - 1.0)
	if req.FormContainerOpacity != nil {
		if *req.FormContainerOpacity < 0.0 || *req.FormContainerOpacity > 1.0 {
			helpers.RespondError(w, http.StatusBadRequest, "La opacidad del contenedor del formulario debe estar entre 0.0 y 1.0")
			return
		}
	}
	if req.DescriptionContainerOpacity != nil {
		if *req.DescriptionContainerOpacity < 0.0 || *req.DescriptionContainerOpacity > 1.0 {
			helpers.RespondError(w, http.StatusBadRequest, "La opacidad del contenedor de descripción debe estar entre 0.0 y 1.0")
			return
		}
	}
	if req.FormMetaBackgroundOpacity != nil {
		if *req.FormMetaBackgroundOpacity < 0.0 || *req.FormMetaBackgroundOpacity > 1.0 {
			helpers.RespondError(w, http.StatusBadRequest, "La opacidad del fondo de la franja meta debe estar entre 0.0 y 1.0")
			return
		}
	}

	// Valores por defecto si no se proporcionan
	log.Printf("ANTES de aplicar valores por defecto - TitleColor='%s', form_container_color='%s', form_container_opacity=%v, description_container_color='%s', description_container_opacity=%v",
		req.TitleColor, req.FormContainerColor, req.FormContainerOpacity, req.DescriptionContainerColor, req.DescriptionContainerOpacity)

	if req.TitleColor == "" {
		req.TitleColor = "#FFFFFF"
		log.Printf("TitleColor estaba vacío, aplicado valor por defecto: #FFFFFF")
	}
	if req.FormContainerColor == "" {
		req.FormContainerColor = "#FFFFFF"
		log.Printf("FormContainerColor estaba vacío, aplicado valor por defecto: #FFFFFF")
	}
	if req.DescriptionContainerColor == "" {
		req.DescriptionContainerColor = "#FFFFFF"
		log.Printf("DescriptionContainerColor estaba vacío, aplicado valor por defecto: #FFFFFF")
	}
	if req.FormMetaBackground == "" {
		req.FormMetaBackground = "#3B82F6"
		log.Printf("FormMetaBackground estaba vacío, aplicado valor por defecto: #3B82F6")
	}
	if req.FormMetaBackgroundStart == "" {
		req.FormMetaBackgroundStart = req.FormMetaBackground
		log.Printf("FormMetaBackgroundStart estaba vacío, aplicado valor por defecto desde form_meta_background")
	}
	if req.FormMetaBackgroundEnd == "" {
		req.FormMetaBackgroundEnd = "#1E40AF"
		log.Printf("FormMetaBackgroundEnd estaba vacío, aplicado valor por defecto: #1E40AF")
	}
	if req.FormMetaTextColor == "" {
		req.FormMetaTextColor = "#FFFFFF"
		log.Printf("FormMetaTextColor estaba vacío, aplicado valor por defecto: #FFFFFF")
	}
	// Si la opacidad no viene en el request, aplicar default 1.0.
	// Importante: ya no asumimos que `0` significa "no enviado", porque `0.0` es un valor válido.
	if req.FormContainerOpacity == nil {
		v := 1.0
		req.FormContainerOpacity = &v
		log.Printf("FormContainerOpacity no enviada, aplicado valor por defecto: 1.0")
	}
	if req.DescriptionContainerOpacity == nil {
		v := 1.0
		req.DescriptionContainerOpacity = &v
		log.Printf("DescriptionContainerOpacity no enviada, aplicado valor por defecto: 1.0")
	}
	if req.FormMetaBackgroundOpacity == nil {
		v := 1.0
		req.FormMetaBackgroundOpacity = &v
		log.Printf("FormMetaBackgroundOpacity no enviada, aplicado valor por defecto: 1.0")
	}

	log.Printf("DESPUÉS de aplicar valores por defecto - TitleColor='%s', form_container_color='%s', form_container_opacity=%v, description_container_color='%s', description_container_opacity=%v",
		req.TitleColor, req.FormContainerColor, req.FormContainerOpacity, req.DescriptionContainerColor, req.DescriptionContainerOpacity)

	// Verificar si ya existe personalización
	var existingID int64
	err = database.DB.QueryRow("SELECT id FROM form_customization WHERE form_id = ?", req.FormID).Scan(&existingID)

	log.Printf("Verificación de personalización existente - existingID=%d, error=%v", existingID, err)

	// Obtener los valores actuales de logo_url y logo_url_mobile para preservarlos si el request viene vacío
	var currentLogoURL, currentLogoURLMobile sql.NullString
	checkErr := err // Guardar el error de verificación de existencia
	if err == nil {
		// Solo intentar obtener los logos si existe la personalización
		logoErr := database.DB.QueryRow("SELECT logo_url, logo_url_mobile FROM form_customization WHERE form_id = ?", req.FormID).Scan(&currentLogoURL, &currentLogoURLMobile)
		if logoErr != nil && logoErr != sql.ErrNoRows {
			log.Printf("Error obteniendo logos actuales: %v", logoErr)
			// Continuar con valores vacíos si hay error, pero no afectar el flujo principal
		} else {
			log.Printf("Logos actuales obtenidos - logo_url: %v, logo_url_mobile: %v", currentLogoURL, currentLogoURLMobile)
		}
	}

	if checkErr != nil {
		if checkErr == sql.ErrNoRows {
			// Crear nueva personalización
			log.Printf("No existe personalización para form_id=%d, creando nueva...", req.FormID)
			result, err := database.DB.Exec(`
				INSERT INTO form_customization 
				(form_id, primary_color, secondary_color, background_color, text_color, title_color, logo_url, font_family, button_style,
				 form_container_color, form_container_opacity, description_container_color, description_container_opacity,
				 form_meta_background, form_meta_background_start, form_meta_background_end, form_meta_background_opacity, form_meta_text_color)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, req.FormID, req.PrimaryColor, req.SecondaryColor, req.BackgroundColor,
				req.TextColor, req.TitleColor, sql.NullString{String: req.LogoURL, Valid: req.LogoURL != ""},
				req.FontFamily, req.ButtonStyle,
				req.FormContainerColor, *req.FormContainerOpacity,
				req.DescriptionContainerColor, *req.DescriptionContainerOpacity,
				req.FormMetaBackground, req.FormMetaBackgroundStart, req.FormMetaBackgroundEnd, *req.FormMetaBackgroundOpacity, req.FormMetaTextColor)

			if err != nil {
				log.Printf("Error creando personalización: %v", err)
				helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error creando personalización: %v", err))
				return
			}

			customID, _ := result.LastInsertId()
			log.Printf("Personalización creada exitosamente con id=%d para form_id=%d", customID, req.FormID)
			helpers.RespondJSON(w, http.StatusCreated, map[string]interface{}{
				"success":          true,
				"message":          "Personalización creada correctamente",
				"customization_id": customID,
			})
			return
		} else {
			// Error diferente a "no rows"
			log.Printf("Error verificando personalización existente: %v", err)
			helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error verificando personalización: %v", err))
			return
		}
	}

	// Actualizar personalización existente
	log.Printf("=== EJECUTANDO UPDATE ===")
	log.Printf("Ejecutando UPDATE para form_id=%d con valores: form_container_color=%s, form_container_opacity=%v, description_container_color=%s, description_container_opacity=%v",
		req.FormID, req.FormContainerColor, req.FormContainerOpacity, req.DescriptionContainerColor, req.DescriptionContainerOpacity)
	log.Printf("Parámetros completos del UPDATE - PrimaryColor=%s, SecondaryColor=%s, BackgroundColor=%s, TextColor=%s, TitleColor=%s, LogoURL=%s, LogoURLMobile=%s, FontFamily=%s, ButtonStyle=%s",
		req.PrimaryColor, req.SecondaryColor, req.BackgroundColor, req.TextColor, req.TitleColor, req.LogoURL, req.LogoURLMobile, req.FontFamily, req.ButtonStyle)

	// Preparar valores para el UPDATE: usar los valores del request si vienen, sino usar los valores actuales
	logoURL := sql.NullString{String: req.LogoURL, Valid: req.LogoURL != ""}
	if !logoURL.Valid && currentLogoURL.Valid {
		logoURL = currentLogoURL
		log.Printf("Preservando logo_url existente: %s", currentLogoURL.String)
	}

	logoURLMobile := sql.NullString{String: req.LogoURLMobile, Valid: req.LogoURLMobile != ""}
	if !logoURLMobile.Valid && currentLogoURLMobile.Valid {
		logoURLMobile = currentLogoURLMobile
		log.Printf("Preservando logo_url_mobile existente: %s", currentLogoURLMobile.String)
	}

	result, err := database.DB.Exec(`
		UPDATE form_customization
		SET primary_color = ?, secondary_color = ?, background_color = ?, 
		    text_color = ?, title_color = ?, logo_url = ?, logo_url_mobile = ?, 
		    font_family = ?, button_style = ?,
		    form_container_color = ?, form_container_opacity = ?,
		    description_container_color = ?, description_container_opacity = ?,
		    form_meta_background = ?, form_meta_background_start = ?, form_meta_background_end = ?, form_meta_background_opacity = ?, form_meta_text_color = ?
		WHERE form_id = ?
	`, req.PrimaryColor, req.SecondaryColor, req.BackgroundColor,
		req.TextColor, req.TitleColor, logoURL, logoURLMobile,
		req.FontFamily, req.ButtonStyle,
		req.FormContainerColor, *req.FormContainerOpacity,
		req.DescriptionContainerColor, *req.DescriptionContainerOpacity,
		req.FormMetaBackground, req.FormMetaBackgroundStart, req.FormMetaBackgroundEnd, *req.FormMetaBackgroundOpacity, req.FormMetaTextColor,
		req.FormID)

	if err != nil {
		log.Printf("ERROR en UPDATE: %v", err)
		log.Printf("Error actualizando personalización: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error actualizando personalización: %v", err))
		return
	}

	// Verificar que se actualizó al menos una fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("ERROR obteniendo filas afectadas: %v", err)
	} else {
		log.Printf("=== RESULTADO UPDATE ===")
		log.Printf("UPDATE exitoso: %d fila(s) afectada(s) para form_id=%d", rowsAffected, req.FormID)
		log.Printf("Valores actualizados - TitleColor=%s, form_container_color=%s, form_container_opacity=%v, description_container_color=%s, description_container_opacity=%v",
			req.TitleColor, req.FormContainerColor, req.FormContainerOpacity, req.DescriptionContainerColor, req.DescriptionContainerOpacity)
	}

	if rowsAffected == 0 {
		log.Printf("ADVERTENCIA: UPDATE no afectó ninguna fila para form_id=%d. ¿Existe la personalización?", req.FormID)
		helpers.RespondError(w, http.StatusNotFound, "No se encontró personalización para actualizar")
		return
	}

	log.Printf("=== FIN CreateOrUpdateCustomization (UPDATE exitoso) ===")

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":       true,
		"message":       "Personalización actualizada correctamente",
		"rows_affected": rowsAffected,
	})
}

// isValidHexColor valida si un string es un color hexadecimal válido
func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
