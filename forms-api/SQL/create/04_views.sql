-- ============================================================================
-- VISTAS - FORMS (Table-per-Form)
-- ============================================================================
-- total_submissions desde form_submissions. Conteo por formulario vía LEFT JOIN.
-- Nota: submission_count es número de envíos globales por form_id (form_submissions).
-- Vistas que reemplazan procedures de solo lectura: usuarios, formularios, públicos, personalización.
-- ============================================================================

-- =========================================================
-- VIEW: vw_users_credentials (antes: sp_get_user_by_credentials)
-- =========================================================
DROP VIEW IF EXISTS vw_users_credentials;
CREATE VIEW vw_users_credentials AS
SELECT id, id_number, email, created_at FROM users;

-- =========================================================
-- VIEW: vw_forms (antes: sp_get_forms + sp_get_form_by_id)
-- =========================================================
DROP VIEW IF EXISTS vw_forms;
CREATE VIEW vw_forms AS
SELECT
 id,
 form_name,
 description,
 created_at,
 event_date,
 expires_at,
 open_at,
 COALESCE(status, 2) AS status,
 CASE WHEN COALESCE(status, 2) = 2 THEN 1 ELSE 0 END AS is_draft,
 CASE WHEN COALESCE(status, 0) = 0 THEN 1 ELSE 0 END AS is_closed,
 closed_at,
 inscription_limit,
 conferencista,
 ubicacion,
 `schema`
FROM forms;

-- =========================================================
-- VIEW: vw_form_customization (antes: sp_get_form_customization)
-- =========================================================
DROP VIEW IF EXISTS vw_form_customization;
CREATE VIEW vw_form_customization AS
SELECT
 id, form_id, primary_color, secondary_color, background_color, text_color, title_color,
 logo_url, logo_url_mobile, font_family, button_style,
 form_container_color, form_container_opacity, description_container_color, description_container_opacity,
 form_meta_background, form_meta_background_start, form_meta_background_end,
 form_meta_background_opacity, form_meta_text_color
FROM form_customization;

-- =========================================================
-- vw_public_forms (eliminada): el API público lee directamente `forms`
-- (GET /api/public/forms, GET /api/public/forms/:id) con reglas en Go:
-- COALESCE(status,2) IN (1, 2) — borrador y publicado; excluye cerrado (0).
-- Cupos, expires_at y open_at se validan al enviar, no al listar.
-- =========================================================
DROP VIEW IF EXISTS vw_public_forms;

-- =========================================================
-- Dashboard y stats (existentes)
-- =========================================================
DROP VIEW IF EXISTS v_dashboard_stats;
CREATE VIEW v_dashboard_stats AS
SELECT 
 (SELECT COUNT(*) FROM forms) AS total_forms,
 (
 SELECT COUNT(*)
 FROM forms f
 WHERE
 f.status = 1
 AND f.expires_at > NOW()
 AND (
 f.inscription_limit IS NULL
 OR f.inscription_limit <= 0
 OR (
 SELECT COUNT(*)
 FROM form_submissions fs
 WHERE fs.form_id = f.id
 ) < f.inscription_limit
 )
 ) AS active_forms,
 (
 SELECT COUNT(*)
 FROM forms f
 WHERE
 NOT (
 f.status = 1
 AND f.expires_at > NOW()
 AND (
 f.inscription_limit IS NULL
 OR f.inscription_limit <= 0
 OR (
 SELECT COUNT(*)
 FROM form_submissions fs
 WHERE fs.form_id = f.id
 ) < f.inscription_limit
 )
 )
 ) AS expired_forms,
 (SELECT COUNT(*) FROM users) AS total_users,
 (SELECT COUNT(*) FROM form_submissions) AS total_submissions,
 (SELECT COUNT(*) FROM form_submissions WHERE DATE(submitted_at) = CURDATE()) AS submissions_today,
 (SELECT COUNT(*) FROM form_submissions WHERE DATE(submitted_at) >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)) AS submissions_week,
 (SELECT COUNT(*) FROM form_submissions WHERE DATE(submitted_at) >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)) AS submissions_month;

DROP VIEW IF EXISTS v_form_stats;
CREATE VIEW v_form_stats AS
SELECT 
 f.id, f.form_name, f.created_at, f.expires_at,
 CASE WHEN COALESCE(f.status, 0) = 0 THEN 1 ELSE 0 END AS is_closed,
 COUNT(DISTINCT fs.id) AS submission_count,
 COUNT(DISTINCT fs.user_id) AS unique_users,
 MAX(fs.submitted_at) AS last_submission,
 CASE
 WHEN
 f.status = 1
 AND f.expires_at > NOW()
 AND (
 f.inscription_limit IS NULL
 OR f.inscription_limit <= 0
 OR COUNT(DISTINCT fs.user_id) < f.inscription_limit
 )
 THEN 'active'
 ELSE 'expired'
 END AS status
FROM forms f
LEFT JOIN form_submissions fs ON f.id = fs.form_id
GROUP BY f.id, f.form_name, f.created_at, f.expires_at, f.status, f.open_at, f.inscription_limit;

DROP VIEW IF EXISTS v_roles;
CREATE VIEW v_roles AS
SELECT id, name, type, created_at FROM roles;
