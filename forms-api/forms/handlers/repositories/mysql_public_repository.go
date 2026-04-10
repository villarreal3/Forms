package repositories

import (
	"database/sql"
	"encoding/json"
	"time"

	"forms/database"
	"forms/models"
)

// Listado y detalle público: consulta directa sobre `forms` (sin vw_public_forms).
// Solo formularios abiertos y vigentes en ventana pública:
// status=1, open_at <= NOW() y expires_at >= NOW().

const sqlPublicFormsList = `
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
	inscription_limit,
	conferencista,
	ubicacion,
	` + "`schema`" + `
FROM forms
WHERE COALESCE(status, 2) = 1
  AND open_at <= NOW()
  AND expires_at >= NOW()
ORDER BY created_at DESC
`

const sqlPublicFormByID = `
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
	inscription_limit,
	conferencista,
	ubicacion,
	` + "`schema`" + `
FROM forms
WHERE id = ?
  AND COALESCE(status, 2) = 1
  AND open_at <= NOW()
  AND expires_at >= NOW()
`

type mysqlPublicRepository struct{}

func NewPublicRepository() PublicRepository {
	return &mysqlPublicRepository{}
}

func scanFormRow(
	scanner interface {
		Scan(dest ...any) error
	},
) (*models.Form, error) {
	var form models.Form
	var desc sql.NullString
	var eventDate time.Time
	var expiresAt time.Time
	var openAt time.Time
	var status int
	var isDraft bool
	var isClosed bool
	var inscriptionLimit sql.NullInt64
	var conferencista sql.NullString
	var ubicacion sql.NullString
	var schemaJSON []byte

	err := scanner.Scan(
		&form.ID, &form.FormName, &desc, &form.CreatedAt,
		&eventDate, &expiresAt, &openAt, &status, &isDraft, &isClosed,
		&inscriptionLimit, &conferencista, &ubicacion, &schemaJSON,
	)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		form.Description = desc.String
	}
	form.EventDate = eventDate
	form.ExpiresAt = expiresAt
	form.OpenAt = openAt
	form.Status = status
	form.IsDraft = isDraft
	form.IsClosed = isClosed
	if inscriptionLimit.Valid {
		v := int(inscriptionLimit.Int64)
		form.InscriptionLimit = &v
	} else {
		form.InscriptionLimit = nil
	}
	if conferencista.Valid {
		s := conferencista.String
		form.Conferencista = &s
	}
	if ubicacion.Valid {
		s := ubicacion.String
		form.Ubicacion = &s
	}
	form.Schema = json.RawMessage(schemaJSON)
	return &form, nil
}

func (r *mysqlPublicRepository) GetPublicForms() ([]models.Form, error) {
	rows, err := database.DB.Query(sqlPublicFormsList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []models.Form
	for rows.Next() {
		form, err := scanFormRow(rows)
		if err != nil {
			return nil, err
		}
		forms = append(forms, *form)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *mysqlPublicRepository) GetPublicForm(formID int64) (*models.Form, error) {
	row := database.DB.QueryRow(sqlPublicFormByID, formID)
	form, err := scanFormRow(row)
	if err != nil {
		return nil, err
	}
	return form, nil
}
