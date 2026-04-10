-- ============================================================================
-- VERIFICACIÓN FINAL - FORMS
-- ============================================================================
-- Muestra tablas, procedures, triggers y vistas creados. Opcional.
-- ============================================================================

SELECT '✅ Base de datos creada exitosamente' AS status;
SELECT '📊 Tablas creadas:' AS info;
SHOW TABLES;
SELECT '🔧 Stored Procedures creados:' AS info;
SHOW PROCEDURE STATUS WHERE Db = 'app_db';
SELECT '⚙️ Funciones creadas:' AS info;
SHOW FUNCTION STATUS WHERE Db = 'app_db';
SELECT '🧾 Triggers creados:' AS info;
SHOW TRIGGERS FROM app_db;
SELECT '👁️ Vistas creadas:' AS info;
SHOW FULL TABLES WHERE Table_type = 'VIEW';

-- ============================================================================
-- FIN DEL SCRIPT
-- ============================================================================
