-- conferencista y ubicación opcionales en forms + vista y sp_create_form.

ALTER TABLE forms
    ADD COLUMN conferencista VARCHAR(512) NULL COMMENT 'Nombre del conferencista (opcional)' AFTER inscription_limit,
    ADD COLUMN ubicacion VARCHAR(512) NULL COMMENT 'Lugar o dirección del evento (opcional)' AFTER conferencista;

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

DELIMITER $$

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

DELIMITER ;
