-- ============================================================================
-- STORED PROCEDURES - FORMS (Table-per-Form)
-- ============================================================================
-- Sin EAV: no sp_submit_form, no form_fields/form_sections SPs.
-- Envíos y listados se hacen en Go con SQL directo sobre form_{id}_responses.
-- ============================================================================
DELIMITER $$

-- PROCEDURE: sp_close_form
DROP PROCEDURE IF EXISTS sp_close_form$$
CREATE PROCEDURE sp_close_form(
 IN p_form_id INT,
 IN p_admin_user_id INT,
 IN p_ip_address VARCHAR(45),
 IN p_user_agent TEXT,
 OUT p_result TINYINT)
BEGIN
 DECLARE v_exists INT;
 SET p_result = 0;
 SELECT COUNT(*) INTO v_exists FROM forms WHERE id = p_form_id;
 IF v_exists > 0 THEN
 UPDATE forms SET status = 0, closed_at = NOW() WHERE id = p_form_id;
 SET p_result = 1;
 END IF;
END$$

-- PROCEDURE: sp_upsert_role
DROP PROCEDURE IF EXISTS sp_upsert_role$$
CREATE PROCEDURE sp_upsert_role(IN p_name VARCHAR(255), IN p_type VARCHAR(50), OUT p_role_id INT)
BEGIN
 INSERT INTO roles (name, type) VALUES (p_name, p_type)
 ON DUPLICATE KEY UPDATE name = VALUES(name), type = VALUES(type);
 SELECT id INTO p_role_id FROM roles WHERE name = p_name AND type = p_type LIMIT 1;
END$$

-- PROCEDURE: sp_create_form (incluye schema JSON)
DROP PROCEDURE IF EXISTS sp_create_form$$
CREATE PROCEDURE sp_create_form(
 IN p_form_name VARCHAR(255),
 IN p_description TEXT,
 IN p_event_date DATETIME,
 IN p_expires_at DATETIME,
 IN p_open_at DATETIME,
 IN p_is_draft BOOLEAN,
 IN p_inscription_limit INT,
 IN p_conferencista VARCHAR(512),
 IN p_ubicacion VARCHAR(512),
 IN p_schema JSON,
 OUT p_form_id INT
)
BEGIN
 INSERT INTO forms (form_name, description, event_date, expires_at, open_at, status, inscription_limit, conferencista, ubicacion, `schema`)
 VALUES (
 p_form_name,
 p_description,
 p_event_date,
 p_expires_at,
 p_open_at,
 CASE WHEN COALESCE(p_is_draft, 0) = 1 THEN 2 ELSE 1 END,
 p_inscription_limit,
 NULLIF(TRIM(p_conferencista), ''),
 NULLIF(TRIM(p_ubicacion), ''),
 p_schema
 );
 SET p_form_id = LAST_INSERT_ID();
END$$

-- PROCEDURE: sp_form_exists
DROP PROCEDURE IF EXISTS sp_form_exists$$
CREATE PROCEDURE sp_form_exists(IN p_form_id INT, OUT p_exists BOOLEAN)
BEGIN
 DECLARE v_count INT;
 SELECT COUNT(*) INTO v_count FROM forms WHERE id = p_form_id;
 SET p_exists = (v_count > 0);
END$$

-- PROCEDURE: sp_get_form_name
DROP PROCEDURE IF EXISTS sp_get_form_name$$
CREATE PROCEDURE sp_get_form_name(IN p_form_id INT, OUT p_form_name VARCHAR(255))
BEGIN
 SELECT form_name INTO p_form_name FROM forms WHERE id = p_form_id LIMIT 1;
END$$

-- PROCEDURE: sp_get_form_expires_at
DROP PROCEDURE IF EXISTS sp_get_form_expires_at$$
CREATE PROCEDURE sp_get_form_expires_at(IN p_form_id INT, OUT p_expires_at DATETIME, OUT p_exists BOOLEAN)
BEGIN
 SELECT expires_at INTO p_expires_at FROM forms WHERE id = p_form_id LIMIT 1;
 IF p_expires_at IS NULL THEN SET p_exists = FALSE; ELSE SET p_exists = TRUE; END IF;
END$$

-- PROCEDURE: sp_get_attendance
DROP PROCEDURE IF EXISTS sp_get_attendance$$
CREATE PROCEDURE sp_get_attendance(IN p_user_id INT, IN p_form_id INT, OUT p_attended BOOLEAN, OUT p_attendance_id INT)
BEGIN
 DECLARE v_id INT;
 DECLARE v_attended BOOLEAN;
 SELECT id, attended INTO v_id, v_attended FROM attendance WHERE user_id = p_user_id AND form_id = p_form_id LIMIT 1;
 IF v_id IS NULL THEN SET p_attended = FALSE; SET p_attendance_id = 0;
 ELSE SET p_attended = v_attended; SET p_attendance_id = v_id; END IF;
END$$

-- PROCEDURE: sp_update_attendance
DROP PROCEDURE IF EXISTS sp_update_attendance$$
CREATE PROCEDURE sp_update_attendance(IN p_user_id INT, IN p_form_id INT, IN p_attended BOOLEAN, OUT p_attendance_id INT)
BEGIN
 DECLARE v_existing_id INT;
 SELECT id INTO v_existing_id FROM attendance WHERE user_id = p_user_id AND form_id = p_form_id LIMIT 1;
 IF v_existing_id IS NULL THEN
 INSERT INTO attendance (user_id, form_id, attended, attendance_date) VALUES (p_user_id, p_form_id, p_attended, NOW());
 SET p_attendance_id = LAST_INSERT_ID();
 ELSE
 UPDATE attendance SET attended = p_attended, attendance_date = NOW() WHERE id = v_existing_id;
 SET p_attendance_id = v_existing_id;
 END IF;
END$$

-- CUSTOMIZATION
DROP PROCEDURE IF EXISTS sp_create_or_update_customization$$
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
 INSERT INTO form_customization (form_id, primary_color, secondary_color, background_color, text_color, title_color, logo_url, logo_url_mobile, font_family, button_style, form_container_color, form_container_opacity, description_container_color, description_container_opacity, form_meta_background, form_meta_background_start, form_meta_background_end, form_meta_background_opacity, form_meta_text_color)
 VALUES (p_form_id, p_primary_color, p_secondary_color, p_background_color, p_text_color, p_title_color, p_logo_url, p_logo_url_mobile, p_font_family, p_button_style, p_form_container_color, p_form_container_opacity, p_description_container_color, p_description_container_opacity, p_form_meta_background, p_form_meta_background_start, p_form_meta_background_end, p_form_meta_background_opacity, p_form_meta_text_color);
 SET p_customization_id = LAST_INSERT_ID();
 ELSE
 UPDATE form_customization SET primary_color = p_primary_color, secondary_color = p_secondary_color, background_color = p_background_color, text_color = p_text_color, title_color = p_title_color,
 logo_url = p_logo_url, logo_url_mobile = p_logo_url_mobile, font_family = p_font_family, button_style = p_button_style,
 form_container_color = p_form_container_color, form_container_opacity = p_form_container_opacity, description_container_color = p_description_container_color, description_container_opacity = p_description_container_opacity,
 form_meta_background = p_form_meta_background, form_meta_background_start = p_form_meta_background_start, form_meta_background_end = p_form_meta_background_end, form_meta_background_opacity = p_form_meta_background_opacity, form_meta_text_color = p_form_meta_text_color
 WHERE id = v_existing_id;
 SET p_customization_id = v_existing_id;
 END IF;
END$$

DROP PROCEDURE IF EXISTS sp_get_customization_id$$
CREATE PROCEDURE sp_get_customization_id(IN p_form_id INT, OUT p_customization_id INT)
BEGIN
 SELECT id INTO p_customization_id FROM form_customization WHERE form_id = p_form_id LIMIT 1;
END$$

DROP PROCEDURE IF EXISTS sp_update_customization_logo$$
CREATE PROCEDURE sp_update_customization_logo(IN p_form_id INT, IN p_logo_url VARCHAR(500), IN p_is_mobile BOOLEAN, OUT p_success BOOLEAN)
BEGIN
 DECLARE v_existing_id INT;
 SELECT id INTO v_existing_id FROM form_customization WHERE form_id = p_form_id LIMIT 1;
 IF v_existing_id IS NULL THEN
 IF p_is_mobile THEN
 INSERT INTO form_customization (form_id, primary_color, secondary_color, background_color, text_color, logo_url_mobile, font_family, button_style)
 VALUES (p_form_id, '#3B82F6', '#1E40AF', '#FFFFFF', '#1F2937', p_logo_url, 'Arial, sans-serif', 'rounded');
 ELSE
 INSERT INTO form_customization (form_id, primary_color, secondary_color, background_color, text_color, logo_url, font_family, button_style)
 VALUES (p_form_id, '#3B82F6', '#1E40AF', '#FFFFFF', '#1F2937', p_logo_url, 'Arial, sans-serif', 'rounded');
 END IF;
 SET p_success = TRUE;
 ELSE
 IF p_is_mobile THEN UPDATE form_customization SET logo_url_mobile = p_logo_url WHERE id = v_existing_id;
 ELSE UPDATE form_customization SET logo_url = p_logo_url WHERE id = v_existing_id; END IF;
 SET p_success = TRUE;
 END IF;
END$$

DROP PROCEDURE IF EXISTS sp_get_customization_logo$$
CREATE PROCEDURE sp_get_customization_logo(IN p_form_id INT, OUT p_logo_url VARCHAR(500), OUT p_exists BOOLEAN)
BEGIN
 SELECT logo_url INTO p_logo_url FROM form_customization WHERE form_id = p_form_id LIMIT 1;
 IF p_logo_url IS NULL THEN SET p_exists = FALSE; ELSE SET p_exists = TRUE; END IF;
END$$

DROP PROCEDURE IF EXISTS sp_delete_customization_logo$$
CREATE PROCEDURE sp_delete_customization_logo(IN p_form_id INT, IN p_is_mobile BOOLEAN, OUT p_success BOOLEAN)
BEGIN
 IF p_is_mobile THEN UPDATE form_customization SET logo_url_mobile = NULL WHERE form_id = p_form_id;
 ELSE UPDATE form_customization SET logo_url = NULL WHERE form_id = p_form_id; END IF;
 SET p_success = (ROW_COUNT() > 0);
END$$

DELIMITER ;
