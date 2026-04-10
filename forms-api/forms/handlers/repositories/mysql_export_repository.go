package repositories

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"forms/database"
	"forms/schema"
)

type mysqlExportRepository struct{}

func NewExportRepository() ExportRepository {
	return &mysqlExportRepository{}
}

// colToExportKey maps response table column names to export map keys used by export service
func colToExportKey(col string) (key string, skip bool) {
	switch strings.ToLower(col) {
	case "id":
		return "submission_id", false
	case "user_id":
		return "", true
	case "submitted_at":
		return "submitted_at", false
	case "nombre":
		return "first_name", false
	case "apellido":
		return "last_name", false
	case "cedula":
		return "id_number", false
	case "email":
		return "email", false
	case "telefono":
		return "phone", false
	case "provincia":
		return "province", false
	case "genero":
		return "gender", false
	case "edad":
		return "_edad_raw", false // handled below
	case "negocio":
		return "business_name", false
	case "registro_empresarial":
		return "business_registration", false
	case "instagram":
		return "instagram", false
	case "facebook":
		return "facebook", false
	case "tiktok":
		return "tiktok", false
	case "twitter":
		return "twitter", false
	case "compartirinfo":
		return "data_consent", false
	case "newsletter":
		return "newsletter", false
	default:
		return col, false // extra columns as dynamic-like keys
	}
}

func normalizeExportMap(m map[string]interface{}) {
	if v, ok := m["_edad_raw"]; ok {
		delete(m, "_edad_raw")
		switch t := v.(type) {
		case int64:
			m["age"] = int(t)
		case int:
			m["age"] = t
		case []byte:
			if n, err := strconv.Atoi(string(t)); err == nil {
				m["age"] = n
			} else {
				m["age"] = 0
			}
		case string:
			if n, err := strconv.Atoi(t); err == nil {
				m["age"] = n
			} else {
				m["age"] = 0
			}
		default:
			m["age"] = 0
		}
	}
	if v, ok := m["data_consent"]; ok {
		switch t := v.(type) {
		case int64:
			m["data_consent"] = t != 0
		case int:
			m["data_consent"] = t != 0
		case []byte:
			m["data_consent"] = string(t) == "1"
		}
	}
	if v, ok := m["newsletter"]; ok {
		switch t := v.(type) {
		case int64:
			m["newsletter"] = t != 0
		case int:
			m["newsletter"] = t != 0
		case []byte:
			m["newsletter"] = string(t) == "1"
		}
	}
}

func (r *mysqlExportRepository) ExportSubmissions(formID *int64, startDate, endDate *string, province, searchTerm *string, limit, offset int) ([]map[string]interface{}, error) {
	if formID == nil || *formID == 0 {
		return nil, fmt.Errorf("formID requerido para exportar (table-per-form)")
	}
	effectiveLimit := limit
	if effectiveLimit <= 0 {
		effectiveLimit = 999999
	}
	table := schema.ResponseTableName(*formID)
	q := "SELECT * FROM `" + table + "` WHERE 1=1"
	args := []interface{}{}
	if startDate != nil && *startDate != "" {
		q += " AND DATE(submitted_at) >= ?"
		args = append(args, *startDate)
	}
	if endDate != nil && *endDate != "" {
		q += " AND DATE(submitted_at) <= ?"
		args = append(args, *endDate)
	}
	if province != nil && *province != "" {
		q += " AND (provincia = ? OR province = ?)"
		args = append(args, *province, *province)
	}
	if searchTerm != nil && *searchTerm != "" {
		q += " AND (nombre LIKE ? OR apellido LIKE ? OR email LIKE ? OR cedula LIKE ?)"
		p := "%" + *searchTerm + "%"
		args = append(args, p, p, p, p)
	}
	q += " ORDER BY submitted_at DESC LIMIT ? OFFSET ?"
	args = append(args, effectiveLimit, offset)

	rows, err := database.DB.Query(q, args...)
	if err != nil {
		log.Printf("ExportSubmissions query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var formName string
	_ = database.DB.QueryRow("SELECT form_name FROM forms WHERE id = ?", *formID).Scan(&formName)

	cols, _ := rows.Columns()
	var submissions []map[string]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			continue
		}
		sub := make(map[string]interface{})
		for i, c := range cols {
			v := vals[i]
			if b, ok := v.([]byte); ok {
				v = string(b)
			}
			key, skip := colToExportKey(c)
			if skip {
				continue
			}
			if key == "submission_id" {
				// ensure int64 for export service
				switch t := v.(type) {
				case int64:
					sub[key] = t
				case int:
					sub[key] = int64(t)
				default:
					sub[key] = v
				}
			} else {
				sub[key] = v
			}
		}
		sub["form_id"] = *formID
		sub["form_name"] = formName
		normalizeExportMap(sub)
		submissions = append(submissions, sub)
	}
	log.Printf("ExportSubmissions: %d filas desde %s", len(submissions), table)
	return submissions, nil
}

// GetDynamicFields: table-per-form — una fila ya contiene todas las columnas.
// Devolvemos vacío para no duplicar columnas ya mapeadas en ExportSubmissions.
func (r *mysqlExportRepository) GetDynamicFields(submissionID int64) (map[string]string, error) {
	return map[string]string{}, nil
}
