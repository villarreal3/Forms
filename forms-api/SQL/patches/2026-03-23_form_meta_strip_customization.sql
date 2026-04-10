-- Patch: agregar persistencia para personalización de formMetaStrip
-- Fecha: 2026-03-23

ALTER TABLE form_customization
    ADD COLUMN form_meta_background VARCHAR(7) DEFAULT '#3B82F6' AFTER description_container_opacity,
    ADD COLUMN form_meta_background_opacity DECIMAL(3,2) DEFAULT 1.00 AFTER form_meta_background,
    ADD COLUMN form_meta_text_color VARCHAR(7) DEFAULT '#FFFFFF' AFTER form_meta_background_opacity;

UPDATE form_customization
SET
    form_meta_background = COALESCE(form_meta_background, '#3B82F6'),
    form_meta_background_opacity = COALESCE(form_meta_background_opacity, 1.00),
    form_meta_text_color = COALESCE(form_meta_text_color, '#FFFFFF');

DROP PROCEDURE IF EXISTS sp_create_or_update_customization;
DELIMITER $$
CREATE PROCEDURE sp_create_or_update_customization(
    IN p_form_id INT, IN p_primary_color VARCHAR(7), IN p_secondary_color VARCHAR(7), IN p_background_color VARCHAR(7), IN p_text_color VARCHAR(7),
    IN p_title_color VARCHAR(7), IN p_logo_url VARCHAR(500), IN p_logo_url_mobile VARCHAR(500), IN p_font_family VARCHAR(100), IN p_button_style VARCHAR(50),
    IN p_form_container_color VARCHAR(7), IN p_form_container_opacity DECIMAL(3,2), IN p_description_container_color VARCHAR(7), IN p_description_container_opacity DECIMAL(3,2),
    IN p_form_meta_background VARCHAR(7), IN p_form_meta_background_opacity DECIMAL(3,2), IN p_form_meta_text_color VARCHAR(7),
    OUT p_customization_id INT
)
BEGIN
    DECLARE v_existing_id INT;
    SELECT id INTO v_existing_id FROM form_customization WHERE form_id = p_form_id LIMIT 1;
    IF v_existing_id IS NULL THEN
        INSERT INTO form_customization (
            form_id, primary_color, secondary_color, background_color, text_color, title_color,
            logo_url, logo_url_mobile, font_family, button_style,
            form_container_color, form_container_opacity, description_container_color, description_container_opacity,
            form_meta_background, form_meta_background_opacity, form_meta_text_color
        )
        VALUES (
            p_form_id, p_primary_color, p_secondary_color, p_background_color, p_text_color, p_title_color,
            p_logo_url, p_logo_url_mobile, p_font_family, p_button_style,
            p_form_container_color, p_form_container_opacity, p_description_container_color, p_description_container_opacity,
            p_form_meta_background, p_form_meta_background_opacity, p_form_meta_text_color
        );
        SET p_customization_id = LAST_INSERT_ID();
    ELSE
        UPDATE form_customization
        SET
            primary_color = p_primary_color,
            secondary_color = p_secondary_color,
            background_color = p_background_color,
            text_color = p_text_color,
            title_color = p_title_color,
            logo_url = p_logo_url,
            logo_url_mobile = p_logo_url_mobile,
            font_family = p_font_family,
            button_style = p_button_style,
            form_container_color = p_form_container_color,
            form_container_opacity = p_form_container_opacity,
            description_container_color = p_description_container_color,
            description_container_opacity = p_description_container_opacity,
            form_meta_background = p_form_meta_background,
            form_meta_background_opacity = p_form_meta_background_opacity,
            form_meta_text_color = p_form_meta_text_color
        WHERE id = v_existing_id;
        SET p_customization_id = v_existing_id;
    END IF;
END$$
DELIMITER ;

DROP VIEW IF EXISTS vw_form_customization;
CREATE VIEW vw_form_customization AS
SELECT
    id, form_id, primary_color, secondary_color, background_color, text_color, title_color,
    logo_url, logo_url_mobile, font_family, button_style,
    form_container_color, form_container_opacity, description_container_color, description_container_opacity,
    form_meta_background, form_meta_background_opacity, form_meta_text_color
FROM form_customization;
