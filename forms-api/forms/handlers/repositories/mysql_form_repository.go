package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"forms/database"
	"forms/models"
	"forms/schema"
)

type mysqlFormRepository struct{}

var (
	hasStatusOnce sync.Once
	hasStatusCol  bool
)

func hasStatusColumn() bool {
	hasStatusOnce.Do(func() {
		err := database.DB.QueryRow(`
			SELECT COUNT(*) > 0
			FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE()
			  AND TABLE_NAME = 'forms'
			  AND COLUMN_NAME = 'status'
		`).Scan(&hasStatusCol)
		if err != nil {
			log.Printf("Error comprobando columna status: %v", err)
			hasStatusCol = false
		}
	})
	return hasStatusCol
}

func NewFormRepository() FormRepository {
	return &mysqlFormRepository{}
}

func (r *mysqlFormRepository) Create(formName, description string, eventDate, openAt, expiresAt interface{}, isDraft bool, inscriptionLimit *int, conferencista, ubicacion *string) (int64, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	schemaJSON, err := schema.DefaultSchemaJSON()
	if err != nil {
		return 0, err
	}

	var conf interface{}
	var ubi interface{}
	if conferencista != nil {
		conf = *conferencista
	}
	if ubicacion != nil {
		ubi = *ubicacion
	}

	_, err = tx.Exec(
		"CALL sp_create_form(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, @p_form_id)",
		formName, description, eventDate, expiresAt, openAt, isDraft, inscriptionLimit, conf, ubi, schemaJSON,
	)
	if err != nil {
		return 0, err
	}
	var formID int64
	if err = tx.QueryRow("SELECT @p_form_id").Scan(&formID); err != nil {
		return 0, err
	}

	ddl, err := schema.BuildCreateResponseTableDDL(formID, schemaJSON)
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ddl); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return formID, nil
}

func (r *mysqlFormRepository) GetForms() ([]models.Form, error) {
	rows, err := database.DB.Query("SELECT id, form_name, description, created_at, event_date, expires_at, open_at, status, is_draft, is_closed, closed_at, inscription_limit, conferencista, ubicacion, `schema` FROM vw_forms ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var forms []models.Form
	for rows.Next() {
		var form models.Form
		var desc sql.NullString
		var eventDate time.Time
		var expiresAt time.Time
		var openAt time.Time
		var status int
		var isDraft bool
		var closedAt sql.NullTime
		var isClosed bool
		var inscriptionLimit sql.NullInt64
		var conferencista sql.NullString
		var ubicacion sql.NullString
		var schemaBytes sql.NullString
		if err := rows.Scan(&form.ID, &form.FormName, &desc, &form.CreatedAt, &eventDate, &expiresAt, &openAt, &status, &isDraft, &isClosed, &closedAt, &inscriptionLimit, &conferencista, &ubicacion, &schemaBytes); err != nil {
			log.Printf("Error escaneando formulario: %v", err)
			continue
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
		if closedAt.Valid {
			form.ClosedAt = &closedAt.Time
		} else {
			form.ClosedAt = nil
		}
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
		if schemaBytes.Valid {
			form.Schema = json.RawMessage(schemaBytes.String)
		}
		forms = append(forms, form)
	}
	return forms, nil
}

func (r *mysqlFormRepository) GetByID(formID int64) (*models.Form, error) {
	row := database.DB.QueryRow("SELECT id, form_name, description, created_at, event_date, expires_at, open_at, status, is_draft, is_closed, closed_at, inscription_limit, conferencista, ubicacion, `schema` FROM vw_forms WHERE id = ?", formID)
	var form models.Form
	var desc sql.NullString
	var eventDate time.Time
	var expiresAt time.Time
	var openAt time.Time
	var status int
	var isDraft bool
	var closedAt sql.NullTime
	var isClosed bool
	var inscriptionLimit sql.NullInt64
	var conferencista sql.NullString
	var ubicacion sql.NullString
	var schemaBytes sql.NullString
	err := row.Scan(&form.ID, &form.FormName, &desc, &form.CreatedAt, &eventDate, &expiresAt, &openAt, &status, &isDraft, &isClosed, &closedAt, &inscriptionLimit, &conferencista, &ubicacion, &schemaBytes)
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
	if closedAt.Valid {
		form.ClosedAt = &closedAt.Time
	} else {
		form.ClosedAt = nil
	}
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
	if schemaBytes.Valid {
		form.Schema = json.RawMessage(schemaBytes.String)
	}
	return &form, nil
}

func (r *mysqlFormRepository) FormExists(formID int64) (bool, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	_, err = tx.Exec("CALL sp_form_exists(?, @p_exists)", formID)
	if err != nil {
		return false, err
	}
	var exists bool
	if err = tx.QueryRow("SELECT @p_exists").Scan(&exists); err != nil {
		return false, err
	}
	_ = tx.Commit()
	return exists, nil
}

func (r *mysqlFormRepository) GetFields(formID int64) ([]models.FormField, error) {
	var schemaJSON []byte
	row := database.DB.QueryRow("SELECT `schema` FROM forms WHERE id = ?", formID)
	if err := row.Scan(&schemaJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	if len(schemaJSON) == 0 {
		return []models.FormField{}, nil
	}
	return schema.ParseToFormFields(formID, schemaJSON)
}

func (r *mysqlFormRepository) GetFieldNames(formID int64) ([]string, error) {
	fields, err := r.GetFields(formID)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.FieldName
	}
	return names, nil
}

func (r *mysqlFormRepository) CreateField(formID int64, fieldLabel, fieldName, fieldType, fieldOptions string, isRequired bool) (int64, error) {
	return 0, fmt.Errorf("CreateField not supported in Table-per-Form; update forms.schema instead")
}
func (r *mysqlFormRepository) UpdateField(fieldID, formID int64, fieldLabel, fieldName, fieldType, fieldOptions string, isRequired bool) (bool, error) {
	return false, fmt.Errorf("UpdateField not supported")
}
func (r *mysqlFormRepository) DeleteField(fieldID, formID int64) (bool, string, error) {
	return false, "not supported", fmt.Errorf("DeleteField not supported")
}
func (r *mysqlFormRepository) ReorderField(fieldID, formID int64, direction string) (bool, string, error) {
	return false, "not supported", fmt.Errorf("ReorderField not supported")
}
func (r *mysqlFormRepository) AddDefaultFields(formID int64) (int, error) {
	return 0, fmt.Errorf("AddDefaultFields not supported; schema set at create")
}
func (r *mysqlFormRepository) CreateSection(formID int64, title, icon string) (int64, error) {
	return 0, fmt.Errorf("CreateSection not supported")
}
func (r *mysqlFormRepository) GetFormSections(formID int64) ([]models.FormSection, error) {
	f, err := r.GetByID(formID)
	if err != nil || len(f.Schema) == 0 {
		return []models.FormSection{}, err
	}
	var s schema.FormSchemaJSON
	if err := json.Unmarshal(f.Schema, &s); err != nil {
		return nil, err
	}
	out := make([]models.FormSection, 0, len(s.Sections))
	for i, sec := range s.Sections {
		out = append(out, models.FormSection{
			ID:           int64(i + 1), // sintético: UI legacy usa id numérico; edición va vía schema
			FormID:       formID,
			SectionTitle: sec.SectionTitle,
			SectionIcon:  sec.SectionIcon,
			DisplayOrder: sec.DisplayOrder,
		})
	}
	return out, nil
}
func (r *mysqlFormRepository) UpdateSection(sectionID, formID int64, title, icon string) (bool, error) {
	return false, fmt.Errorf("UpdateSection not supported")
}
func (r *mysqlFormRepository) DeleteSection(sectionID, formID int64) (bool, string, error) {
	return false, "not supported", fmt.Errorf("DeleteSection not supported")
}
func (r *mysqlFormRepository) ReorderSection(sectionID, formID int64, direction string) (bool, string, error) {
	return false, "not supported", fmt.Errorf("ReorderSection not supported")
}
func (r *mysqlFormRepository) AssignFieldToSection(fieldID, sectionID, formID int64) (bool, error) {
	return false, fmt.Errorf("AssignFieldToSection not supported")
}

func (r *mysqlFormRepository) GetFormName(formID int64) (string, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	_, err = tx.Exec("CALL sp_get_form_name(?, @p_form_name)", formID)
	if err != nil {
		return "", err
	}
	var name sql.NullString
	if err = tx.QueryRow("SELECT @p_form_name").Scan(&name); err != nil {
		return "", err
	}
	_ = tx.Commit()
	if !name.Valid {
		return "", sql.ErrNoRows
	}
	return name.String, nil
}

func (r *mysqlFormRepository) GetFormExpiresAt(formID int64) (interface{}, bool, error) {
	// Consulta directa para evitar variables de sesión (@p_expires_at), que el driver devuelve como []uint8
	// y no se pueden escanear en *time.Time.
	var expiresAt sql.NullTime
	err := database.DB.QueryRow("SELECT expires_at FROM forms WHERE id = ?", formID).Scan(&expiresAt)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if !expiresAt.Valid {
		return nil, true, nil
	}
	return expiresAt.Time, true, nil
}

func (r *mysqlFormRepository) CloseForm(formID, adminUserID int64, ipAddress, userAgent string) (bool, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	_, err = tx.Exec("CALL sp_close_form(?, ?, ?, ?, @p_result)", formID, adminUserID, ipAddress, userAgent)
	if err != nil {
		// Compatibilidad: si el procedimiento falla (p.ej. DB vieja con is_closed),
		// intentamos un UPDATE directo según el esquema existente.
		log.Printf("Error ejecutando sp_close_form (fallback a UPDATE): %v", err)
		var res sql.Result
		var lastErr error

		// Intentar schema nuevo (con closed_at).
		res, err = tx.Exec("UPDATE forms SET status = 0, closed_at = NOW() WHERE id = ?", formID)
		if err == nil {
			lastErr = nil
		} else {
			lastErr = err
			// Intentar schema nuevo (sin closed_at).
			res, err = tx.Exec("UPDATE forms SET status = 0 WHERE id = ?", formID)
			if err != nil {
				lastErr = err
				// Intentar esquema viejo (con closed_at).
				res, err = tx.Exec("UPDATE forms SET is_closed = 1, closed_at = NOW() WHERE id = ?", formID)
				if err != nil {
					lastErr = err
					// Intentar esquema viejo (sin closed_at).
					res, err = tx.Exec("UPDATE forms SET is_closed = 1 WHERE id = ?", formID)
				}
			}
		}

		if err != nil {
			return false, lastErr
		}
		ra, _ := res.RowsAffected()
		if err = tx.Commit(); err != nil {
			return false, err
		}
		return ra > 0, nil
	}

	var result int
	if err = tx.QueryRow("SELECT @p_result").Scan(&result); err != nil {
		return false, err
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return result == 1, nil
}

func (r *mysqlFormRepository) SetDraftStatus(formID int64, isDraft bool) (bool, error) {
	// Compatibilidad robusta:
	// 1) Intentar schema nuevo (status: 2=borrador, 1=abierto)
	// 2) Si falla, intentar schema viejo (is_draft/is_closed)
	newStatus := 1
	if isDraft {
		newStatus = 2
	}

	// status (con closed_at)
	var lastErr error
	if res, err := database.DB.Exec(
		"UPDATE forms SET status = ?, closed_at = NULL WHERE id = ?",
		newStatus, formID,
	); err == nil {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	} else {
		lastErr = err
	}

	// Fallback a esquema viejo.
	newIsDraft := 0
	if isDraft {
		newIsDraft = 1
	}

	// status (sin closed_at)
	if res, err := database.DB.Exec(
		"UPDATE forms SET status = ? WHERE id = ?",
		newStatus, formID,
	); err == nil {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	} else {
		lastErr = err
	}

	// is_draft/is_closed (con closed_at)
	res, err := database.DB.Exec(
		"UPDATE forms SET is_draft = ?, is_closed = 0, closed_at = NULL WHERE id = ?",
		newIsDraft, formID,
	)
	if err != nil {
		lastErr = err
	} else {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	}

	// is_draft/is_closed (sin closed_at)
	if res, err := database.DB.Exec(
		"UPDATE forms SET is_draft = ?, is_closed = 0 WHERE id = ?",
		newIsDraft, formID,
	); err == nil {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	} else {
		lastErr = err
	}

	return false, lastErr
}

func (r *mysqlFormRepository) OpenForm(formID int64) (bool, error) {
	// Compatibilidad robusta:
	// 1) Intentar schema nuevo (status)
	// 2) Si falla (columna inexistente), intentar schema viejo (is_draft/is_closed)
	// status (con closed_at)
	if res, err := database.DB.Exec(
		"UPDATE forms SET status = 1, closed_at = NULL WHERE id = ?",
		formID,
	); err == nil {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	}

	// status (sin closed_at)
	if res, err := database.DB.Exec(
		"UPDATE forms SET status = 1 WHERE id = ?",
		formID,
	); err == nil {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	}

	// is_draft/is_closed (con closed_at)
	if res, err := database.DB.Exec(
		"UPDATE forms SET is_draft = 0, is_closed = 0, closed_at = NULL WHERE id = ?",
		formID,
	); err == nil {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	}

	// is_draft/is_closed (sin closed_at)
	if res, err := database.DB.Exec(
		"UPDATE forms SET is_draft = 0, is_closed = 0 WHERE id = ?",
		formID,
	); err == nil {
		ra, _ := res.RowsAffected()
		return ra > 0, nil
	}

	// Si llegamos acá, devolvemos el error final intentando el primer UPDATE.
	_, err := database.DB.Exec(
		"UPDATE forms SET status = 1, closed_at = NULL WHERE id = ?",
		formID,
	)
	return false, err
}

func (r *mysqlFormRepository) UpdateFormSchema(formID int64, schemaJSON []byte) error {
	if len(schemaJSON) == 0 {
		return fmt.Errorf("schema vacío")
	}
	var n int
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM form_submissions WHERE form_id = ?", formID).Scan(&n); err != nil {
		return err
	}
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec("UPDATE forms SET `schema` = ? WHERE id = ?", schemaJSON, formID); err != nil {
		return err
	}
	if n == 0 {
		table := schema.ResponseTableName(formID)
		if _, err = tx.Exec("DROP TABLE IF EXISTS `" + table + "`"); err != nil {
			return err
		}
		ddl, err := schema.BuildCreateResponseTableDDL(formID, schemaJSON)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ddl); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *mysqlFormRepository) UpdateFormMetadata(formID int64, conferencista, ubicacion *string) error {
	if conferencista == nil && ubicacion == nil {
		return fmt.Errorf("no hay campos para actualizar")
	}
	var sets []string
	var args []interface{}
	if conferencista != nil {
		sets = append(sets, "conferencista = ?")
		if *conferencista == "" {
			args = append(args, nil)
		} else {
			args = append(args, *conferencista)
		}
	}
	if ubicacion != nil {
		sets = append(sets, "ubicacion = ?")
		if *ubicacion == "" {
			args = append(args, nil)
		} else {
			args = append(args, *ubicacion)
		}
	}
	args = append(args, formID)
	q := "UPDATE forms SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	res, err := database.DB.Exec(q, args...)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
