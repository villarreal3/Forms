-- ============================================================================
-- INICIALIZACIÓN - FORMS
-- ============================================================================
-- Crea la base de datos, la selecciona y configura zona horaria.
-- Ejecutar primero. El resto de scripts asumen USE app_db.
-- ============================================================================

-- Eliminar base de datos existente (si existe)
DROP DATABASE IF EXISTS app_db;

-- Crear base de datos nueva
CREATE DATABASE app_db
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

-- Usar la base de datos creada
USE app_db;

-- Configurar zona horaria de Panamá
SET time_zone = '-05:00';
