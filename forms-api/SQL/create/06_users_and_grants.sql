-- ============================================================================
-- USUARIOS Y PERMISOS - FORMS
-- ============================================================================
-- Crea admin_user (privilegios completos) y public_user (solo procedures de
-- formularios públicos). Ejecutar como root después de crear tablas y procedures.
--
-- IMPORTANTE: Sustituir las contraseñas antes de ejecutar (o usar variables de
-- entorno al invocar). No versionar contraseñas reales.
--
-- Base de datos: si usas FORMS_PROD en lugar de app_db, sustituir "app_db"
-- en todos los GRANT.
-- ============================================================================

USE app_db;

DROP USER IF EXISTS 'admin_user'@'%';
CREATE USER 'admin_user'@'%' IDENTIFIED BY 'REPLACE_ADMIN_PASSWORD';

DROP USER IF EXISTS 'public_user'@'%';
CREATE USER 'public_user'@'%' IDENTIFIED BY 'REPLACE_PUBLIC_PASSWORD';

GRANT ALL PRIVILEGES ON app_db.* TO 'admin_user'@'%';

-- public_user: lectura vía vistas + tabla `forms` (listado público ya no usa vw_public_forms).
GRANT SELECT ON app_db.vw_users_credentials TO 'public_user'@'%';
GRANT SELECT ON app_db.vw_forms TO 'public_user'@'%';
GRANT SELECT ON app_db.vw_form_customization TO 'public_user'@'%';
GRANT SELECT ON app_db.forms TO 'public_user'@'%';

FLUSH PRIVILEGES;
