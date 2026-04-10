package schema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"forms/models"
)

type FormSchemaJSON struct {
	Sections []Section `json:"sections"`
	Fields   []Field   `json:"fields"`
}

type Section struct {
	ID           string `json:"id"`
	SectionTitle string `json:"section_title"`
	SectionIcon  string `json:"section_icon"`
	DisplayOrder int    `json:"display_order"`
}

type Field struct {
	ID           string  `json:"id"`
	FieldLabel   string  `json:"field_label"`
	FieldName    string  `json:"field_name"`
	FieldType    string  `json:"field_type"`
	FieldOptions *string `json:"field_options,omitempty"`
	IsRequired   bool    `json:"is_required"`
	MaxLength    *int    `json:"max_length,omitempty"`
	DisplayOrder int     `json:"display_order"`
	SectionID    string  `json:"section_id,omitempty"`
}

var colNameRe = regexp.MustCompile(`^[a-z0-9_]{1,64}$`)

func ResponseTableName(formID int64) string {
	return fmt.Sprintf("form_%d_responses", formID)
}

func DefaultSchemaJSON() ([]byte, error) {
	s := FormSchemaJSON{
		Sections: []Section{
			{"s1", "Información Personal", "fa-user", 0},
			{"s2", "Información de Contacto", "fa-envelope", 1},
			{"s3", "Información del Emprendimiento", "fa-briefcase", 2},
			{"s4", "Redes Sociales", "fa-share-alt", 3},
			{"s5", "Consentimiento y Autorización", "fa-check-circle", 4},
		},
		Fields: []Field{
			{"f1", "Nombre", "nombre", "text", nil, true, intPtr(100), 0, "s1"},
			{"f2", "Apellido", "apellido", "text", nil, true, intPtr(100), 1, "s1"},
			{"f3", "Cédula/Pasaporte", "cedula", "text", nil, true, intPtr(100), 2, "s1"},
			{"f4", "Género", "genero", "select", strPtr("Femenino,Masculino,Otro,Prefiero no decir"), true, nil, 3, "s1"},
			{"f5", "Edad", "edad", "text", nil, true, intPtr(100), 4, "s1"},
			{"f6", "Provincia", "provincia", "select", strPtr("Bocas del Toro,Chiriquí,Coclé,Colón,Darién,Emberá,Guna Yala,Herrera,Los Santos,Ngäbe-Buglé,Panamá,Panamá Oeste,Veraguas"), true, nil, 5, "s1"},
			{"f7", "Correo Electrónico", "email", "email", nil, true, intPtr(320), 0, "s2"},
			{"f8", "Teléfono/Celular", "telefono", "text", nil, true, intPtr(20), 1, "s2"},
			{"f9", "Nombre del Emprendimiento", "negocio", "text", nil, true, intPtr(255), 0, "s3"},
			{"f10", "¿Tienes Registro Empresarial?", "registro_empresarial", "select", strPtr("Sí,No,En proceso"), true, nil, 1, "s3"},
			{"f11", "Instagram", "instagram", "text", nil, false, intPtr(255), 0, "s4"},
			{"f12", "Facebook", "facebook", "text", nil, false, intPtr(255), 1, "s4"},
			{"f13", "TikTok", "tiktok", "text", nil, false, intPtr(255), 2, "s4"},
			{"f14", "X (Twitter)", "twitter", "text", nil, false, intPtr(255), 3, "s4"},
			{"f15", "Autorizo el uso de mis datos", "compartirinfo", "checkbox", nil, true, nil, 0, "s5"},
			{"f16", "Deseo recibir información sobre futuros eventos", "newsletter", "checkbox", nil, false, nil, 1, "s5"},
		},
	}
	return json.Marshal(s)
}

func intPtr(n int) *int      { return &n }
func strPtr(s string) *string { return &s }

func BuildCreateResponseTableDDL(formID int64, schemaJSON []byte) (string, error) {
	var s FormSchemaJSON
	if err := json.Unmarshal(schemaJSON, &s); err != nil {
		return "", err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TABLE IF NOT EXISTS %s (\n", ResponseTableName(formID))
	b.WriteString("\tid BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,\n")
	b.WriteString("\tuser_id INT NOT NULL,\n")
	b.WriteString("\tsubmitted_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),\n")
	const collation = "utf8mb4_unicode_ci"
	for _, f := range s.Fields {
		col := strings.ToLower(strings.TrimSpace(f.FieldName))
		if col == "" {
			col = f.ID
		}
		if !colNameRe.MatchString(col) {
			return "", fmt.Errorf("invalid column name: %s", col)
		}
		b.WriteString("\t`")
		b.WriteString(col)
		b.WriteString("` ")
		b.WriteString(fieldTypeToSQL(f, collation))
		b.WriteString(" DEFAULT NULL,\n")
	}
	b.WriteString("\tINDEX idx_user_id (user_id),\n")
	b.WriteString("\tINDEX idx_submitted_at (submitted_at),\n")
	b.WriteString("\tFOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE\n")
	b.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;")
	return b.String(), nil
}

func fieldTypeToSQL(f Field, collation string) string {
	ml := 255
	if f.MaxLength != nil && *f.MaxLength > 0 {
		ml = *f.MaxLength
		if f.FieldType == "email" && ml > 320 {
			ml = 320
		}
	}
	collate := ""
	if collation != "" {
		collate = " COLLATE " + collation
	}
	switch f.FieldType {
	case "number":
		return "BIGINT"
	case "date":
		return "DATE"
	case "checkbox":
		return "TINYINT(1)"
	case "textarea":
		if f.MaxLength != nil && *f.MaxLength > 4000 {
			return "TEXT" + collate
		}
		return fmt.Sprintf("VARCHAR(%d)%s", ml, collate)
	default:
		return fmt.Sprintf("VARCHAR(%d)%s", ml, collate)
	}
}

func ParseToFormFields(formID int64, schemaJSON []byte) ([]models.FormField, error) {
	var s FormSchemaJSON
	if err := json.Unmarshal(schemaJSON, &s); err != nil {
		return nil, err
	}
	out := make([]models.FormField, 0, len(s.Fields))
	// Map schema section id (string) -> índice 1..n para alinear con GetFormSections (IDs sintéticos)
	sectionIndex := make(map[string]int)
	for i, sec := range s.Sections {
		sectionIndex[sec.ID] = i + 1
	}
	for i, f := range s.Fields {
		var opts string
		if f.FieldOptions != nil {
			opts = *f.FieldOptions
		}
		var sectionID *int64
		if f.SectionID != "" {
			if idx, ok := sectionIndex[f.SectionID]; ok {
				v := int64(idx)
				sectionID = &v
			}
		}
		out = append(out, models.FormField{
			ID:           int64(i + 1),
			FormID:       formID,
			FieldLabel:   f.FieldLabel,
			FieldName:    f.FieldName,
			FieldType:    f.FieldType,
			FieldOptions: opts,
			IsRequired:   f.IsRequired,
			DisplayOrder: f.DisplayOrder,
			SectionID:    sectionID,
		})
	}
	return out, nil
}
