-- Patch: extender formMetaStrip para gradiente (color inicio/fin)
-- Fecha: 2026-03-23
--
-- IMPORTANTE: primero DROP VIEW. Si la vista antigua depende de una columna a quitar,
-- hay que eliminar la vista antes del ALTER.

DROP VIEW IF EXISTS vw_form_customization;

ALTER TABLE form_customization
    ADD COLUMN IF NOT EXISTS form_meta_background VARCHAR(7) DEFAULT '#3B82F6' AFTER description_container_opacity,
    ADD COLUMN IF NOT EXISTS form_meta_background_opacity DECIMAL(3,2) DEFAULT 1.00 AFTER form_meta_background,
    ADD COLUMN IF NOT EXISTS form_meta_text_color VARCHAR(7) DEFAULT '#FFFFFF' AFTER form_meta_background_opacity,
    ADD COLUMN IF NOT EXISTS form_meta_background_start VARCHAR(7) DEFAULT '#3B82F6' AFTER form_meta_background,
    ADD COLUMN IF NOT EXISTS form_meta_background_end VARCHAR(7) DEFAULT '#1E40AF' AFTER form_meta_background_start;

SET @legacy_col := CONCAT('form_meta', '_border_', 'color');
SELECT COUNT(*) INTO @cnt FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'form_customization' AND COLUMN_NAME = @legacy_col;
SET @sql := IF(@cnt > 0, CONCAT('ALTER TABLE form_customization DROP COLUMN `', @legacy_col, '`'), 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE form_customization
SET
    form_meta_background = COALESCE(form_meta_background, '#3B82F6'),
    form_meta_background_start = COALESCE(form_meta_background_start, form_meta_background, '#3B82F6'),
    form_meta_background_end = COALESCE(form_meta_background_end, '#1E40AF'),
    form_meta_background_opacity = COALESCE(form_meta_background_opacity, 1.00),
    form_meta_text_color = COALESCE(form_meta_text_color, '#FFFFFF')
WHERE id > 0;

DROP PROCEDURE IF EXISTS sp_create_or_update_customization;
DELIMITER $$
CREATE PROCEDURE sp_create_or_update_customization(
    IN p_form_id INT, IN p_primary_color VARCHAR(7), IN p_secondary_color VARCHAR(7), IN p_background_color VARCHAR(7), IN p_text_color VARCHAR(7),
    IN p_title_color VARCHAR(7), IN p_logo_url VARCHAR(500), IN p_logo_url_mobile VARCHAR(500), IN p_font_family VARCHAR(100), IN p_button_style VARCHAR(50),
    IN p_form_container_color VARCHAR(7), IN p_form_container_opacity DECIMAL(3,2), IN p_description_container_color VARCHAR(7), IN p_description_container_opacity DECIMAL(3,2),
    IN p_form_meta_background VARCHAR(7), IN p_form_meta_background_start VARCHAR(7), IN p_form_meta_background_end VARCHAR(7),
    IN p_form_meta_background_opacity DECIMAL(3,2), IN p_form_meta_text_color VARCHAR(7),
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
            form_meta_background, form_meta_background_start, form_meta_background_end, form_meta_background_opacity, form_meta_text_color
        )
        VALUES (
            p_form_id, p_primary_color, p_secondary_color, p_background_color, p_text_color, p_title_color,
            p_logo_url, p_logo_url_mobile, p_font_family, p_button_style,
            p_form_container_color, p_form_container_opacity, p_description_container_color, p_description_container_opacity,
            p_form_meta_background, p_form_meta_background_start, p_form_meta_background_end, p_form_meta_background_opacity, p_form_meta_text_color
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
            form_meta_background_start = p_form_meta_background_start,
            form_meta_background_end = p_form_meta_background_end,
            form_meta_background_opacity = p_form_meta_background_opacity,
            form_meta_text_color = p_form_meta_text_color
        WHERE id = v_existing_id;
        SET p_customization_id = v_existing_id;
    END IF;
END$$
DELIMITER ;

CREATE VIEW vw_form_customization AS
SELECT
    id, form_id, primary_color, secondary_color, background_color, text_color, title_color,
    logo_url, logo_url_mobile, font_family, button_style,
    form_container_color, form_container_opacity, description_container_color, description_container_opacity,
    form_meta_background, form_meta_background_start, form_meta_background_end,
    form_meta_background_opacity, form_meta_text_color
FROM form_customization;
