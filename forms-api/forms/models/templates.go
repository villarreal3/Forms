package models

// CommonFieldTemplate representa una plantilla de campo común
type CommonFieldTemplate struct {
	Label      string `json:"label"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Options    string `json:"options,omitempty"`
	IsRequired bool   `json:"is_required"`
	Category   string `json:"category"`
}

// GetCommonFieldTemplates retorna las plantillas de campos comunes
func GetCommonFieldTemplates() []CommonFieldTemplate {
	return []CommonFieldTemplate{
		// Información Personal Básica
		{
			Label:      "Nombre",
			Name:       "nombre",
			Type:       "text",
			IsRequired: true,
			Category:   "personal",
		},
		{
			Label:      "Apellido",
			Name:       "apellido",
			Type:       "text",
			IsRequired: true,
			Category:   "personal",
		},
		{
			Label:      "Cédula/Pasaporte",
			Name:       "cedula",
			Type:       "text",
			IsRequired: true,
			Category:   "personal",
		},
		{
			Label:      "Género",
			Name:       "genero",
			Type:       "select",
			Options:    "Femenino,Masculino,Otro,Prefiero no decir",
			IsRequired: true,
			Category:   "personal",
		},
		{
			Label:      "Edad",
			Name:       "edad",
			Type:       "text",
			IsRequired: true,
			Category:   "personal",
		},
		{
			Label:      "Provincia",
			Name:       "provincia",
			Type:       "select",
			Options:    "Bocas del Toro,Chiriquí,Coclé,Colón,Darién,Emberá,Guna Yala,Herrera,Los Santos,Ngäbe-Buglé,Panamá,Panamá Oeste,Veraguas",
			IsRequired: true,
			Category:   "personal",
		},
		
		// Información de Contacto
		{
			Label:      "Correo Electrónico",
			Name:       "email",
			Type:       "email",
			IsRequired: true,
			Category:   "contact",
		},
		{
			Label:      "Teléfono/Celular",
			Name:       "telefono",
			Type:       "tel",
			IsRequired: true,
			Category:   "contact",
		},
		
		// Información del Emprendimiento
		{
			Label:      "Nombre del Emprendimiento",
			Name:       "negocio",
			Type:       "text",
			IsRequired: true,
			Category:   "business",
		},
		{
			Label:      "¿Tienes Registro Empresarial?",
			Name:       "registro_empresarial",
			Type:       "radio",
			Options:    "si,no,en_proceso",
			IsRequired: true,
			Category:   "business",
		},
		
		// Redes Sociales - Toggle + Campo
		{
			Label:      "Instagram",
			Name:       "instagram",
			Type:       "text",
			IsRequired: false,
			Category:   "social",
		},
		{
			Label:      "Facebook",
			Name:       "facebook",
			Type:       "text",
			IsRequired: false,
			Category:   "social",
		},
		{
			Label:      "TikTok",
			Name:       "tiktok",
			Type:       "text",
			IsRequired: false,
			Category:   "social",
		},
		{
			Label:      "X (Twitter)",
			Name:       "twitter",
			Type:       "text",
			IsRequired: false,
			Category:   "social",
		},
		
		// Consentimientos
		{
			Label:      "Autorizo el uso de mis datos",
			Name:       "compartirinfo",
			Type:       "checkbox",
			IsRequired: true,
			Category:   "consent",
		},
		{
			Label:      "Deseo recibir información sobre futuros eventos",
			Name:       "newsletter",
			Type:       "checkbox",
			IsRequired: false,
			Category:   "consent",
		},
	}
}

// GetFieldsByCategory retorna campos filtrados por categoría
func GetFieldsByCategory(category string) []CommonFieldTemplate {
	allTemplates := GetCommonFieldTemplates()
	var filtered []CommonFieldTemplate
	
	for _, template := range allTemplates {
		if template.Category == category {
			filtered = append(filtered, template)
		}
	}
	
	return filtered
}

// GetDefaultFields retorna la lista de nombres de campos que se deben agregar por defecto a cada formulario
func GetDefaultFields() []string {
	return []string{
		"nombre",
		"apellido",
		"cedula",
		"genero",
		"edad",
		"provincia",
		"email",
		"telefono",
		"negocio",
		"registro_empresarial",
		"instagram",
		"facebook",
		"tiktok",
		"twitter",
		"compartirinfo",
		"newsletter",
	}
}

