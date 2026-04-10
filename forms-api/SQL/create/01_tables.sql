-- ============================================================================
-- TABLAS - FORMS (Table-per-Form)
-- ============================================================================
-- Esquema sin EAV: formulario define schema en forms.schema (JSON);
-- cada formulario tiene tabla form_{id}_responses creada por el backend.
-- Requiere 00_init.sql previo (USE app_db).
-- ============================================================================

-- ============================================================================
-- TABLA: USERS / PARTICIPANTS
-- ============================================================================
CREATE TABLE users (
 id INT AUTO_INCREMENT PRIMARY KEY,
 id_number VARCHAR(20) NOT NULL,
 email VARCHAR(150) NOT NULL,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 UNIQUE KEY unique_user (id_number, email),
 INDEX idx_users_email (email),
 INDEX idx_users_id_number (id_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLA: FORMS
-- ============================================================================
CREATE TABLE forms (
 id INT AUTO_INCREMENT PRIMARY KEY,
 form_name VARCHAR(255) NOT NULL,
 description TEXT,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 event_date DATETIME NOT NULL COMMENT 'Fecha del evento (informativa; usada en correos/UX)',
 expires_at DATETIME NOT NULL,
 open_at DATETIME NOT NULL COMMENT 'Fecha/hora desde la cual el formulario puede ser público',
 - status: 2=borrador (no público), 1=abierto (publicable si cumple ventana temporal), 0=cerrado (manual).
 status TINYINT NOT NULL DEFAULT 2 COMMENT 'Estado del formulario: 2=borrador, 1=abierto, 0=cerrado',
 inscription_limit INT NULL COMMENT 'Límite de inscripciones; NULL = sin límite',
 conferencista VARCHAR(512) NULL COMMENT 'Nombre del conferencista (opcional)',
 ubicacion VARCHAR(512) NULL COMMENT 'Lugar o dirección del evento (opcional)',
 closed_at DATETIME NULL COMMENT 'Fecha y hora de cierre manual',
 `schema` JSON NULL COMMENT 'Secciones y campos (label, name, type, options, required, max_length, display_order, section_id)',
 INDEX idx_forms_expires_at (expires_at),
 INDEX idx_forms_open_at (open_at),
 INDEX idx_forms_created_at (created_at),
 INDEX idx_forms_status (status),
 INDEX idx_forms_public_window (open_at, expires_at, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLA: FORM_SUBMISSIONS (relación user_id + form_id + response_id)
-- ============================================================================
CREATE TABLE form_submissions (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
 user_id INT NOT NULL,
 form_id INT NOT NULL,
 response_id BIGINT UNSIGNED NOT NULL COMMENT 'PK en form_{form_id}_responses',
 submitted_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
 UNIQUE KEY uk_user_form (user_id, form_id),
 INDEX idx_form_submitted (form_id, submitted_at),
 INDEX idx_user_form (user_id, form_id),
 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
 FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLA: ATTENDANCE (FOR EVENTS / ACTIVITIES)
-- ============================================================================
CREATE TABLE attendance (
 id INT AUTO_INCREMENT PRIMARY KEY,
 user_id INT NOT NULL,
 form_id INT NOT NULL,
 attended BOOLEAN DEFAULT 0,
 attendance_date DATETIME DEFAULT CURRENT_TIMESTAMP,
 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
 FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE,
 INDEX idx_attendance_user_id (user_id),
 INDEX idx_attendance_form_id (form_id),
 INDEX idx_attendance_date (attendance_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLA: ROLES (copia local para el frontend)
-- ============================================================================
CREATE TABLE roles (
 id INT AUTO_INCREMENT PRIMARY KEY,
 name VARCHAR(255) NOT NULL,
 type VARCHAR(50) NOT NULL,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 UNIQUE KEY uk_roles_name_type (name, type),
 INDEX idx_roles_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLA: FORM_CUSTOMIZATION (Personalización visual de formularios)
-- ============================================================================
CREATE TABLE form_customization (
 id INT AUTO_INCREMENT PRIMARY KEY,
 form_id INT NOT NULL,
 primary_color VARCHAR(7) DEFAULT '#3B82F6',
 secondary_color VARCHAR(7) DEFAULT '#1E40AF',
 background_color VARCHAR(7) DEFAULT '#FFFFFF',
 text_color VARCHAR(7) DEFAULT '#1F2937',
 title_color VARCHAR(7) DEFAULT '#FFFFFF',
 logo_url VARCHAR(500),
 logo_url_mobile VARCHAR(500),
 font_family VARCHAR(100) DEFAULT 'Arial, sans-serif',
 button_style VARCHAR(50) DEFAULT 'rounded',
 form_container_color VARCHAR(7) DEFAULT '#FFFFFF',
 form_container_opacity DECIMAL(3,2) DEFAULT 1.00,
 description_container_color VARCHAR(7) DEFAULT '#FFFFFF',
 description_container_opacity DECIMAL(3,2) DEFAULT 1.00,
 form_meta_background VARCHAR(7) DEFAULT '#3B82F6',
 form_meta_background_start VARCHAR(7) DEFAULT '#3B82F6',
 form_meta_background_end VARCHAR(7) DEFAULT '#1E40AF',
 form_meta_background_opacity DECIMAL(3,2) DEFAULT 1.00,
 form_meta_text_color VARCHAR(7) DEFAULT '#FFFFFF',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE,
 UNIQUE KEY unique_form_customization (form_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLA: FORM_AUDIT_LOG (Auditoría de cambios en formularios)
-- ============================================================================
CREATE TABLE form_audit_log (
 id BIGINT AUTO_INCREMENT PRIMARY KEY,
 action VARCHAR(50) NOT NULL COMMENT 'form_created, form_closed, form_schema_updated, form_customization_created, form_customization_updated',
 form_id INT NULL COMMENT 'ID del formulario afectado (si aplica)',
 entity_table VARCHAR(64) NOT NULL COMMENT 'forms, form_customization',
 entity_id INT NULL COMMENT 'PK del registro afectado en entity_table',
 details JSON NULL COMMENT 'Resumen de cambios',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 INDEX idx_audit_form_id (form_id),
 INDEX idx_audit_created_at (created_at),
 INDEX idx_audit_action (action),
 FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- NOTA: Tablas form_{form_id}_responses se crean dinámicamente con CREATE TABLE
-- desde el backend Go según forms.schema.
-- ============================================================================
