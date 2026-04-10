-- ============================================================================
-- TRIGGERS DE AUDITORÍA - FORMS (Table-per-Form)
-- ============================================================================
DELIMITER $$

DROP TRIGGER IF EXISTS trg_forms_after_insert$$
CREATE TRIGGER trg_forms_after_insert
AFTER INSERT ON forms FOR EACH ROW
BEGIN
 INSERT INTO form_audit_log (action, form_id, entity_table, entity_id, details)
 VALUES ('form_created', NEW.id, 'forms', NEW.id, NULL);
END$$

DROP TRIGGER IF EXISTS trg_forms_after_update$$
CREATE TRIGGER trg_forms_after_update
AFTER UPDATE ON forms FOR EACH ROW
BEGIN
 IF COALESCE(OLD.status, 0) <> COALESCE(NEW.status, 0) THEN
 IF COALESCE(NEW.status, 0) = 0 THEN
 INSERT INTO form_audit_log (action, form_id, entity_table, entity_id, details)
 VALUES ('form_closed', NEW.id, 'forms', NEW.id, NULL);
 ELSEIF COALESCE(NEW.status, 0) = 2 THEN
 INSERT INTO form_audit_log (action, form_id, entity_table, entity_id, details)
 VALUES ('form_set_draft', NEW.id, 'forms', NEW.id, NULL);
 ELSEIF COALESCE(NEW.status, 0) = 1 THEN
 INSERT INTO form_audit_log (action, form_id, entity_table, entity_id, details)
 VALUES ('form_published', NEW.id, 'forms', NEW.id, NULL);
 END IF;
 END IF;
 IF NOT (OLD.`schema` <=> NEW.`schema`) THEN
 INSERT INTO form_audit_log (action, form_id, entity_table, entity_id, details)
 VALUES ('form_schema_updated', NEW.id, 'forms', NEW.id, JSON_OBJECT('schema_changed', TRUE));
 END IF;
END$$

-- form_customization triggers (idénticos a antes)
DROP TRIGGER IF EXISTS trg_form_customization_after_insert$$
CREATE TRIGGER trg_form_customization_after_insert
AFTER INSERT ON form_customization FOR EACH ROW
BEGIN
 INSERT INTO form_audit_log (action, form_id, entity_table, entity_id, details)
 VALUES ('form_customization_created', NEW.form_id, 'form_customization', NEW.id,
 JSON_OBJECT('primary_color', NEW.primary_color, 'secondary_color', NEW.secondary_color));
END$$

DROP TRIGGER IF EXISTS trg_form_customization_after_update$$
CREATE TRIGGER trg_form_customization_after_update
AFTER UPDATE ON form_customization FOR EACH ROW
BEGIN
 DECLARE j JSON DEFAULT JSON_OBJECT();
 IF NOT (IFNULL(OLD.primary_color, '') <=> IFNULL(NEW.primary_color, '')) THEN
 SET j = JSON_MERGE_PRESERVE(j, JSON_OBJECT('old_primary_color', OLD.primary_color, 'new_primary_color', NEW.primary_color));
 END IF;
 INSERT INTO form_audit_log (action, form_id, entity_table, entity_id, details)
 VALUES ('form_customization_updated', NEW.form_id, 'form_customization', NEW.id, NULLIF(j, JSON_OBJECT()));
END$$

DELIMITER ;

