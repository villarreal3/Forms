-- Parche: el API público ya no usa vw_public_forms (consulta directa a `forms`).
-- En bases ya desplegadas: eliminar la vista y dar SELECT sobre `forms` a public_user.
--
--   mysql -u... -p app_db < SQL/patches/2026-03-20_vw_public_forms_listing.sql

USE app_db;

DROP VIEW IF EXISTS vw_public_forms;

GRANT SELECT ON app_db.forms TO 'public_user'@'%';
FLUSH PRIVILEGES;
