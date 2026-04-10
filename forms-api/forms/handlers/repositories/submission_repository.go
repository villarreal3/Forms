package repositories

import (
	"forms/models"
)

// SubmissionRepository define las operaciones de acceso a datos para submissions
type SubmissionRepository interface {
	// GetSubmissions obtiene submissions con filtros
	GetSubmissions(formID *int64, startDate, endDate *string, province, searchTerm *string, limit, offset int) ([]map[string]interface{}, error)
	
	// CheckDuplicateSubmission verifica si ya existe una submission
	CheckDuplicateSubmission(idNumber, email string, formID int64) (bool, error)
	
	// GetAttendance obtiene el estado de asistencia
	GetAttendance(userID, formID int64) (bool, int64, error)
	
	// UpdateAttendance actualiza o crea el registro de asistencia
	UpdateAttendance(userID, formID int64, attended bool) (int64, error)
	
	// AutoSubmitFromPrevious crea una nueva inscripción en formID
	// reutilizando la última respuesta (por fecha) de otro formulario
	// con el mismo schema, si existe una inscripción previa con la misma cédula.
	AutoSubmitFromPrevious(targetFormID int64, idNumber string) (bool, error)

	// SubmitForm ejecuta el stored procedure sp_submit_form
	SubmitForm(req *models.SubmitFormRequest) (int, int64, error)
	
	// GetUserByCredentials obtiene un usuario desde vw_users_credentials por id_number y email
	GetUserByCredentials(idNumber, email string) (*models.User, error)
}

