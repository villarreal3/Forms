package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"forms/database"
	"forms/models"
	"forms/schema"
)

// #region agent log
func debugLog(location, message, hypothesisId string, data map[string]interface{}) {
	const logPath = "/home/davillarreal/forms-automate/.cursor/debug-784680.log"
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	payload := map[string]interface{}{"sessionId": "784680", "location": location, "message": message, "hypothesisId": hypothesisId, "data": data, "timestamp": time.Now().UnixMilli()}
	b, _ := json.Marshal(payload)
	f.Write(append(b, '\n'))
	f.Close()
}

// #endregion

type mysqlSubmissionRepository struct{}

func NewSubmissionRepository() SubmissionRepository {
	return &mysqlSubmissionRepository{}
}

func (r *mysqlSubmissionRepository) GetSubmissions(formID *int64, startDate, endDate *string, province, searchTerm *string, limit, offset int) ([]map[string]interface{}, error) {
	if formID == nil || *formID == 0 {
		return []map[string]interface{}{}, nil
	}
	// Table-per-Form:
	// - form_submissions: quién está inscrito y cuándo (submitted_at)
	// - form_{id}_responses: respuestas del formulario (nombre, apellido, cedula, email, telefono, provincia, ...)
	// - users: email "canónico" (fallback)
	// - attendance: si asistió o no
	respTable := schema.ResponseTableName(*formID)
	q := fmt.Sprintf(`
		SELECT
			fs.id                    AS id,
			fs.user_id               AS user_id,
			f.form_name              AS form_name,
			fs.submitted_at          AS submitted_at,
			r.nombre                 AS first_name,
			r.apellido               AS last_name,
			r.cedula                 AS id_number,
			COALESCE(r.email, u.email)    AS email,
			r.telefono               AS phone,
			r.provincia              AS province,
			COALESCE(a.attended, 0)  AS attended
		FROM form_submissions fs
		JOIN forms f ON f.id = fs.form_id
		JOIN %s r ON r.id = fs.response_id
		LEFT JOIN users u ON u.id = fs.user_id
		LEFT JOIN attendance a
		       ON a.user_id = fs.user_id AND a.form_id = fs.form_id
		WHERE fs.form_id = ?
	`, "`"+respTable+"`")

	args := []interface{}{*formID}
	if startDate != nil && *startDate != "" {
		q += " AND DATE(fs.submitted_at) >= ?"
		args = append(args, *startDate)
	}
	if endDate != nil && *endDate != "" {
		q += " AND DATE(fs.submitted_at) <= ?"
		args = append(args, *endDate)
	}
	if province != nil && *province != "" {
		q += " AND r.provincia = ?"
		args = append(args, *province)
	}
	if searchTerm != nil && *searchTerm != "" {
		q += " AND (r.nombre LIKE ? OR r.apellido LIKE ? OR r.email LIKE ? OR r.cedula LIKE ? OR r.telefono LIKE ?)"
		p := "%" + *searchTerm + "%"
		args = append(args, p, p, p, p, p)
	}
	q += " ORDER BY fs.submitted_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	var out []map[string]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			continue
		}
		m := make(map[string]interface{})
		for i, c := range cols {
			v := vals[i]
			if b, ok := v.([]byte); ok {
				m[c] = string(b)
			} else {
				m[c] = v
			}
		}
		out = append(out, m)
	}
	return out, nil
}

func (r *mysqlSubmissionRepository) CheckDuplicateSubmission(idNumber, email string, formID int64) (bool, error) {
	var uid int64
	err := database.DB.QueryRow(
		"SELECT id FROM users WHERE id_number = ? AND email = ? LIMIT 1",
		idNumber, email,
	).Scan(&uid)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var n int
	err = database.DB.QueryRow(
		"SELECT COUNT(*) FROM form_submissions WHERE user_id = ? AND form_id = ?",
		uid, formID,
	).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *mysqlSubmissionRepository) GetAttendance(userID, formID int64) (bool, int64, error) {
	_, err := database.DB.Exec("CALL sp_get_attendance(?, ?, @p_attended, @p_attendance_id)", userID, formID)
	if err != nil {
		return false, 0, err
	}
	var attended bool
	var attendanceID int64
	err = database.DB.QueryRow("SELECT @p_attended, @p_attendance_id").Scan(&attended, &attendanceID)
	return attended, attendanceID, err
}

func (r *mysqlSubmissionRepository) UpdateAttendance(userID, formID int64, attended bool) (int64, error) {
	_, err := database.DB.Exec("CALL sp_update_attendance(?, ?, ?, @p_attendance_id)", userID, formID, attended)
	if err != nil {
		return 0, err
	}
	var attendanceID int64
	err = database.DB.QueryRow("SELECT @p_attendance_id").Scan(&attendanceID)
	return attendanceID, err
}

// AutoSubmitFromPrevious crea una nueva inscripción en targetFormID
// reutilizando la última respuesta (por fecha de submitted_at) de otro
// formulario con el mismo schema, siempre que exista una inscripción
// previa cuya cédula (id_number) coincida.
func (r *mysqlSubmissionRepository) AutoSubmitFromPrevious(targetFormID int64, idNumber string) (bool, error) {
	if targetFormID == 0 || strings.TrimSpace(idNumber) == "" {
		return false, nil
	}

	// 1) Formularios con el mismo schema (excluyendo el actual)
	rows, err := database.DB.Query(
		"SELECT id FROM forms WHERE id <> ? AND `schema` = (SELECT `schema` FROM forms WHERE id = ?)",
		targetFormID,
		targetFormID,
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var formIDs []int64
	for rows.Next() {
		var fid int64
		if err := rows.Scan(&fid); err != nil {
			continue
		}
		formIDs = append(formIDs, fid)
	}
	// #region agent log
	debugLog("mysql_submission_repository.go:AutoSubmitFromPrevious:after_same_schema", "same-schema form IDs", "H2", map[string]interface{}{"targetFormID": targetFormID, "idNumber": idNumber, "formIDsCount": len(formIDs), "formIDs": formIDs})
	// #endregion
	if len(formIDs) == 0 {
		// No hay otros formularios compatibles
		return false, nil
	}

	// 2) Buscar la última inscripción (submitted_at) entre todos esos formularios
	//    donde cédula y teléfono coincidan.
	var parts []string
	var args []interface{}
	for _, fid := range formIDs {
		respTable := schema.ResponseTableName(fid)
		part := fmt.Sprintf(`
			SELECT fs.user_id, fs.form_id, fs.response_id, fs.submitted_at
			FROM form_submissions fs
			JOIN %s r ON r.id = fs.response_id
			WHERE fs.form_id = ? AND r.cedula = ?`,
			"`"+respTable+"`",
		)
		parts = append(parts, part)
		args = append(args, fid, idNumber)
	}
	if len(parts) == 0 {
		return false, nil
	}

	unionSQL := strings.Join(parts, " UNION ALL ")
	q := fmt.Sprintf(`
		SELECT t.user_id, t.form_id, t.response_id
		FROM (
			%s
		) AS t
		ORDER BY t.submitted_at DESC
		LIMIT 1
	`, unionSQL)

	var userID, sourceFormID, sourceResponseID int64
	err = database.DB.QueryRow(q, args...).Scan(&userID, &sourceFormID, &sourceResponseID)
	// #region agent log
	debugLog("mysql_submission_repository.go:AutoSubmitFromPrevious:after_union", "UNION query result", "H3", map[string]interface{}{"err": fmt.Sprintf("%v", err), "isNoRows": err == sql.ErrNoRows, "userID": userID, "sourceFormID": sourceFormID, "sourceResponseID": sourceResponseID})
	// #endregion
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// 3) Obtener columnas dinámicas del schema (mismas en todos los formularios compatibles)
	var schemaJSON []byte
	if err = database.DB.QueryRow("SELECT `schema` FROM forms WHERE id = ?", targetFormID).Scan(&schemaJSON); err != nil {
		return false, err
	}

	cols, err := responseColumnsFromSchema(schemaJSON)
	if err != nil {
		return false, err
	}
	if len(cols) == 0 {
		return false, nil
	}

	// 4) Copiar la fila de la tabla de respuestas del formulario origen
	//    a la tabla de respuestas del formulario destino.
	sourceTable := schema.ResponseTableName(sourceFormID)
	targetTable := schema.ResponseTableName(targetFormID)
	colList := strings.Join(cols, ",")

	copySQL := fmt.Sprintf(
		"INSERT INTO `%s` (user_id, %s) SELECT user_id, %s FROM `%s` WHERE id = ?",
		targetTable,
		colList,
		colList,
		sourceTable,
	)

	// #region agent log
	debugLog("mysql_submission_repository.go:AutoSubmitFromPrevious:before_copy", "INSERT...SELECT params", "H4", map[string]interface{}{"sourceTable": sourceTable, "targetTable": targetTable, "sourceResponseID": sourceResponseID})
	// #endregion
	res, err := database.DB.Exec(copySQL, sourceResponseID)
	// #region agent log
	afterCopyData := map[string]interface{}{"err": fmt.Sprintf("%v", err)}
	if err == nil && res != nil {
		afterCopyData["newResponseID"], _ = res.LastInsertId()
	}
	debugLog("mysql_submission_repository.go:AutoSubmitFromPrevious:after_copy", "INSERT...SELECT result", "H4", afterCopyData)
	// #endregion
	if err != nil {
		return false, err
	}
	newResponseID, _ := res.LastInsertId()
	if newResponseID == 0 {
		return false, nil
	}

	// 5) Crear la fila en form_submissions para el formulario destino.
	_, err = database.DB.Exec(
		"INSERT INTO form_submissions (user_id, form_id, response_id) VALUES (?, ?, ?)",
		userID,
		targetFormID,
		newResponseID,
	)
	// #region agent log
	debugLog("mysql_submission_repository.go:AutoSubmitFromPrevious:after_submission_insert", "form_submissions INSERT result", "H5", map[string]interface{}{"err": fmt.Sprintf("%v", err), "userID": userID, "targetFormID": targetFormID, "newResponseID": newResponseID})
	// #endregion
	if err != nil {
		return false, err
	}

	return true, nil
}

// mapSubmitToColumns builds column->value from SubmitFormRequest using schema field names
func mapSubmitToColumns(req *models.SubmitFormRequest, schemaJSON []byte) (cols []string, vals []interface{}, err error) {
	var s schema.FormSchemaJSON
	if err = json.Unmarshal(schemaJSON, &s); err != nil {
		return nil, nil, err
	}
	// fixed mapping from request struct to common column names
	fixed := map[string]string{
		"nombre": req.FirstName, "apellido": req.LastName, "cedula": req.IDNumber,
		"email": req.Email, "telefono": req.Phone, "provincia": req.Province, "genero": req.Gender,
		"edad": "", "negocio": req.BusinessName, "registro_empresarial": req.BusinessRegistration,
		"instagram": req.Instagram, "facebook": req.Facebook, "tiktok": req.TikTok, "twitter": req.Twitter,
	}
	if req.Age > 0 {
		fixed["edad"] = strconv.Itoa(req.Age)
	}
	if req.DataConsent {
		fixed["compartirinfo"] = "1"
	} else {
		fixed["compartirinfo"] = "0"
	}
	if req.Newsletter {
		fixed["newsletter"] = "1"
	} else {
		fixed["newsletter"] = "0"
	}
	// answers by field_id - schema order gives synthetic ids
	answerByID := make(map[int]string)
	for _, a := range req.Answers {
		answerByID[a.FieldID] = a.FieldValue
	}
	for i, f := range s.Fields {
		col := strings.ToLower(strings.TrimSpace(f.FieldName))
		if col == "" {
			col = f.ID
		}
		var v string
		if fv, ok := fixed[col]; ok {
			v = fv
		} else if fv, ok := answerByID[i+1]; ok {
			v = fv
		}
		cols = append(cols, "`"+col+"`")
		if f.FieldType == "checkbox" {
			if v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "on") {
				vals = append(vals, 1)
			} else if v == "" {
				vals = append(vals, 0)
			} else {
				vals = append(vals, 0)
			}
		} else if f.FieldType == "number" && v != "" {
			if n, e := strconv.ParseInt(v, 10, 64); e == nil {
				vals = append(vals, n)
			} else {
				vals = append(vals, v)
			}
		} else {
			vals = append(vals, v)
		}
	}
	return cols, vals, nil
}

// responseColumnsFromSchema devuelve la lista de columnas dinámicas para la tabla
// de respuestas, en el mismo orden que el schema JSON.
func responseColumnsFromSchema(schemaJSON []byte) ([]string, error) {
	var s schema.FormSchemaJSON
	if err := json.Unmarshal(schemaJSON, &s); err != nil {
		return nil, err
	}
	var cols []string
	for _, f := range s.Fields {
		col := strings.ToLower(strings.TrimSpace(f.FieldName))
		if col == "" {
			col = f.ID
		}
		if col == "" {
			continue
		}
		cols = append(cols, "`"+col+"`")
	}
	return cols, nil
}

func (r *mysqlSubmissionRepository) SubmitForm(req *models.SubmitFormRequest) (int, int64, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	statusExists := hasStatusColumn()

	var formStatus int
	var formClosed bool
	var formDraft bool
	var openAt, expiresAt time.Time
	var inscriptionLimit sql.NullInt64
	var submissionCount int

	if statusExists {
		err = tx.QueryRow(`
			SELECT
				COALESCE(f.status, 0) AS status,
				f.open_at,
				f.expires_at,
				f.inscription_limit,
				(
					SELECT COUNT(*)
					FROM form_submissions fs
					WHERE fs.form_id = f.id
				) AS submission_count
			FROM forms f
			WHERE f.id = ?
		`, req.FormID).Scan(&formStatus, &openAt, &expiresAt, &inscriptionLimit, &submissionCount)
	} else {
		err = tx.QueryRow(`
			SELECT
				COALESCE(f.is_closed, 0) AS is_closed,
				COALESCE(f.is_draft, 0) AS is_draft,
				f.open_at,
				f.expires_at,
				f.inscription_limit,
				(
					SELECT COUNT(*)
					FROM form_submissions fs
					WHERE fs.form_id = f.id
				) AS submission_count
			FROM forms f
			WHERE f.id = ?
		`, req.FormID).Scan(&formClosed, &formDraft, &openAt, &expiresAt, &inscriptionLimit, &submissionCount)
	}
	if err != nil {
		return 0, 0, err
	}

	now := time.Now()
	if statusExists && formStatus == 0 || (!statusExists && formClosed) {
		return 0, 0, &FormUnavailableError{
			Reason:     FormUnavailableClosed,
			Message:    "Este formulario ha sido cerrado y ya no está disponible para envío",
			StatusCode: 400,
		}
	}
	if statusExists && formStatus == 2 || (!statusExists && formDraft) {
		return 0, 0, &FormUnavailableError{
			Reason:     FormUnavailableDraft,
			Message:    "Este formulario aún está en modo borrador y no está disponible para envío",
			StatusCode: 400,
		}
	}
	if now.Before(openAt) {
		return 0, 0, &FormUnavailableError{
			Reason:     FormUnavailableNotOpenYet,
			Message:    "Este formulario aún no está abierto. Intenta más adelante",
			StatusCode: 400,
		}
	}
	if !expiresAt.After(now) {
		return 0, 0, &FormUnavailableError{
			Reason:     FormUnavailableExpired,
			Message:    "Este formulario ha expirado y ya no está disponible para envío",
			StatusCode: 400,
		}
	}
	if inscriptionLimit.Valid && submissionCount >= int(inscriptionLimit.Int64) {
		return 0, 0, &FormUnavailableError{
			Reason:     FormUnavailableLimitReached,
			Message:    "Se alcanzó el límite de inscripciones para este formulario",
			StatusCode: 409,
		}
	}

	// resolve user
	var userID int64
	err = tx.QueryRow("SELECT id FROM users WHERE id_number = ? AND email = ? LIMIT 1", req.IDNumber, req.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow("SELECT id FROM users WHERE id_number = ? LIMIT 1", req.IDNumber).Scan(&userID)
		if err == sql.ErrNoRows {
			res, e := tx.Exec("INSERT INTO users (id_number, email) VALUES (?, ?)", req.IDNumber, req.Email)
			if e != nil {
				return 0, 0, e
			}
			userID, _ = res.LastInsertId()
		} else {
			tx.Exec("UPDATE users SET email = ? WHERE id = ?", req.Email, userID)
		}
	}

	var dup int
	_ = tx.QueryRow("SELECT COUNT(*) FROM form_submissions WHERE user_id = ? AND form_id = ?", userID, req.FormID).Scan(&dup)
	if dup > 0 {
		return 4, userID, nil
	}

	var schemaJSON []byte
	if err = tx.QueryRow("SELECT `schema` FROM forms WHERE id = ?", req.FormID).Scan(&schemaJSON); err != nil {
		return 0, 0, err
	}
	cols, vals, err := mapSubmitToColumns(req, schemaJSON)
	if err != nil {
		return 0, 0, err
	}
	table := schema.ResponseTableName(int64(req.FormID))
	placeholders := strings.Repeat("?,", len(cols))
	placeholders = strings.TrimSuffix(placeholders, ",")
	q := fmt.Sprintf("INSERT INTO `%s` (user_id, %s) VALUES (?,%s)", table, strings.Join(cols, ","), placeholders)
	args := []interface{}{userID}
	for _, v := range vals {
		args = append(args, v)
	}
	res, err := tx.Exec(q, args...)
	if err != nil {
		log.Printf("SubmitForm INSERT error: %v", err)
		return 0, 0, err
	}
	responseID, _ := res.LastInsertId()
	_, err = tx.Exec("INSERT INTO form_submissions (user_id, form_id, response_id) VALUES (?, ?, ?)", userID, req.FormID, responseID)
	if err != nil {
		return 0, 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}
	return 3, userID, nil
}

func (r *mysqlSubmissionRepository) GetUserByCredentials(idNumber, email string) (*models.User, error) {
	row := database.DB.QueryRow(
		"SELECT id, id_number, email, created_at FROM vw_users_credentials WHERE id_number = ? AND email = ? LIMIT 1",
		idNumber, email,
	)
	var user models.User
	err := row.Scan(&user.ID, &user.IDNumber, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
