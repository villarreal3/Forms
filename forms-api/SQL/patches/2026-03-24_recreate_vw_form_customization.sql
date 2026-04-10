-- Fix: Error 1054 columna desconocida al leer personalización.
-- Suele ocurrir si la vista apunta a una columna ya eliminada de form_customization.
-- Ejecutar en la BD afectada. Fecha: 2026-03-24

DROP VIEW IF EXISTS vw_form_customization;

-- Eliminar columna legacy (nombre construido en tiempo de ejecución)
SET @legacy_col := CONCAT('form_meta', '_border_', 'color');
SELECT COUNT(*) INTO @cnt FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'form_customization' AND COLUMN_NAME = @legacy_col;
SET @sql := IF(@cnt > 0, CONCAT('ALTER TABLE form_customization DROP COLUMN `', @legacy_col, '`'), 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE VIEW vw_form_customization AS
SELECT
    id, form_id, primary_color, secondary_color, background_color, text_color, title_color,
    logo_url, logo_url_mobile, font_family, button_style,
    form_container_color, form_container_opacity, description_container_color, description_container_opacity,
    form_meta_background, form_meta_background_start, form_meta_background_end,
    form_meta_background_opacity, form_meta_text_color
FROM form_customization;
