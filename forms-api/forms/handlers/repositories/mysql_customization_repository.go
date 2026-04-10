package repositories

import (
	"database/sql"
	"log"

	"forms/database"
	"forms/models"
)

type mysqlCustomizationRepository struct{}

func NewCustomizationRepository() CustomizationRepository {
	return &mysqlCustomizationRepository{}
}

func (r *mysqlCustomizationRepository) GetFormCustomization(formID int64) (*models.FormCustomization, error) {
	row := database.DB.QueryRow(
		"SELECT id, form_id, primary_color, secondary_color, background_color, text_color, title_color, logo_url, logo_url_mobile, font_family, button_style, form_container_color, form_container_opacity, description_container_color, description_container_opacity, form_meta_background, form_meta_background_start, form_meta_background_end, form_meta_background_opacity, form_meta_text_color FROM vw_form_customization WHERE form_id = ?",
		formID,
	)

	var customization models.FormCustomization
	var logoURL, logoURLMobile sql.NullString
	var formContainerOpacity, descriptionContainerOpacity, formMetaBackgroundOpacity sql.NullFloat64

	err := row.Scan(
		&customization.ID, &customization.FormID, &customization.PrimaryColor,
		&customization.SecondaryColor, &customization.BackgroundColor, &customization.TextColor,
		&customization.TitleColor, &logoURL, &logoURLMobile, &customization.FontFamily, &customization.ButtonStyle,
		&customization.FormContainerColor, &formContainerOpacity,
		&customization.DescriptionContainerColor, &descriptionContainerOpacity,
		&customization.FormMetaBackground, &customization.FormMetaBackgroundStart, &customization.FormMetaBackgroundEnd, &formMetaBackgroundOpacity, &customization.FormMetaTextColor,
	)
	if err != nil {
		return nil, err
	}

	if logoURL.Valid {
		customization.LogoURL = logoURL.String
	}
	if logoURLMobile.Valid {
		customization.LogoURLMobile = logoURLMobile.String
	}
	if formContainerOpacity.Valid {
		customization.FormContainerOpacity = formContainerOpacity.Float64
	} else {
		customization.FormContainerOpacity = 1.0
	}
	if descriptionContainerOpacity.Valid {
		customization.DescriptionContainerOpacity = descriptionContainerOpacity.Float64
	} else {
		customization.DescriptionContainerOpacity = 1.0
	}
	if formMetaBackgroundOpacity.Valid {
		customization.FormMetaBackgroundOpacity = formMetaBackgroundOpacity.Float64
	} else {
		customization.FormMetaBackgroundOpacity = 1.0
	}

	return &customization, nil
}

func (r *mysqlCustomizationRepository) CreateOrUpdateCustomization(formID int64, customization *models.CreateCustomizationRequest) (int64, error) {
	formOpacity := 1.0
	if customization.FormContainerOpacity != nil {
		formOpacity = *customization.FormContainerOpacity
	}
	descOpacity := 1.0
	if customization.DescriptionContainerOpacity != nil {
		descOpacity = *customization.DescriptionContainerOpacity
	}

	metaBgOpacity := 1.0
	if customization.FormMetaBackgroundOpacity != nil {
		metaBgOpacity = *customization.FormMetaBackgroundOpacity
	}

	_, err := database.DB.Exec("CALL sp_create_or_update_customization(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, @p_customization_id)",
		formID, customization.PrimaryColor, customization.SecondaryColor, customization.BackgroundColor,
		customization.TextColor, customization.TitleColor, customization.LogoURL, customization.LogoURLMobile, customization.FontFamily,
		customization.ButtonStyle, customization.FormContainerColor, formOpacity,
		customization.DescriptionContainerColor, descOpacity,
		customization.FormMetaBackground, customization.FormMetaBackgroundStart, customization.FormMetaBackgroundEnd, metaBgOpacity, customization.FormMetaTextColor)
	if err != nil {
		return 0, err
	}
	var customizationID int64
	err = database.DB.QueryRow("SELECT @p_customization_id").Scan(&customizationID)
	return customizationID, err
}

func (r *mysqlCustomizationRepository) GetCustomizationID(formID int64) (int64, error) {
	_, err := database.DB.Exec("CALL sp_get_customization_id(?, @p_customization_id)", formID)
	if err != nil {
		return 0, err
	}
	var customizationID int64
	err = database.DB.QueryRow("SELECT @p_customization_id").Scan(&customizationID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return customizationID, err
}

func (r *mysqlCustomizationRepository) UpdateLogo(formID int64, logoURL string, isMobile bool) (bool, error) {
	// Convertir booleano a entero explícitamente para MySQL (0 = false, 1 = true)
	isMobileInt := 0
	if isMobile {
		isMobileInt = 1
	}
	
	log.Printf("[Repository] UpdateLogo INICIO - formID=%d, logoURL='%s', isMobile=%v, isMobileInt=%d", formID, logoURL, isMobile, isMobileInt)
	log.Printf("[Repository] UpdateLogo ANTES DE SP - Ejecutando CALL sp_update_customization_logo(%d, '%s', %d, @p_success)", formID, logoURL, isMobileInt)
	
	_, err := database.DB.Exec("CALL sp_update_customization_logo(?, ?, ?, @p_success)", formID, logoURL, isMobileInt)
	
	log.Printf("[Repository] UpdateLogo DESPUÉS DE SP - error=%v", err)
	
	if err != nil {
		log.Printf("[Repository] UpdateLogo ERROR - Error ejecutando stored procedure: %v", err)
		return false, err
	}
	
	var success bool
	err = database.DB.QueryRow("SELECT @p_success").Scan(&success)
	
	log.Printf("[Repository] UpdateLogo RESULTADO - success=%v, error=%v", success, err)
	
	return success, err
}

func (r *mysqlCustomizationRepository) GetLogo(formID int64) (string, bool, error) {
	_, err := database.DB.Exec("CALL sp_get_customization_logo(?, @p_logo_url, @p_exists)", formID)
	if err != nil {
		return "", false, err
	}
	var logoURL string
	var exists bool
	err = database.DB.QueryRow("SELECT @p_logo_url, @p_exists").Scan(&logoURL, &exists)
	return logoURL, exists, err
}

func (r *mysqlCustomizationRepository) DeleteLogo(formID int64, isMobile bool) (bool, error) {
	_, err := database.DB.Exec("CALL sp_delete_customization_logo(?, ?, @p_success)", formID, isMobile)
	if err != nil {
		return false, err
	}
	var success bool
	err = database.DB.QueryRow("SELECT @p_success").Scan(&success)
	return success, err
}
