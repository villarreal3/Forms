-- ============================================================================
-- SCRIPT PRINCIPAL - FORMS v0.1
-- ============================================================================
-- Ejecuta en orden: init → tablas → procedures → triggers → views → verify.
-- Ejecutar desde este directorio para que SOURCE encuentre los .sql:
--
-- cd forms-api/SQL/create
-- mysql -u usuario -p < create_db.sql
--
-- O desde el cliente MySQL (después de conectar):
-- SOURCE /ruta/completa/a/SQL/create/create_db.sql;
-- ============================================================================
-- FECHA: 2025-11-11
-- VERSIÓN: 0.1
-- ============================================================================

SOURCE 00_init.sql;
SOURCE 01_tables.sql;
SOURCE 02_procedures.sql;
SOURCE 03_triggers.sql;
SOURCE 04_views.sql;
SOURCE 06_users_and_grants.sql;
SOURCE 05_verify.sql;
