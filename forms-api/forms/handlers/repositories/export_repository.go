package repositories

// ExportRepository define las operaciones de acceso a datos para exportación
type ExportRepository interface {
	// ExportSubmissions exporta submissions con filtros
	ExportSubmissions(formID *int64, startDate, endDate *string, province, searchTerm *string, limit, offset int) ([]map[string]interface{}, error)
	
	// GetDynamicFields obtiene los campos dinámicos de una submission
	GetDynamicFields(submissionID int64) (map[string]string, error)
}














