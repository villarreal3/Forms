package repositories

// EmailRepository define las operaciones de acceso a datos para emails
type EmailRepository interface {
	// AddToQueue agrega un email a la cola
	AddToQueue(recipientEmail, recipientName, subject, emailFilePath string) (int64, error)
	
	// GetPendingEmails obtiene emails pendientes
	GetPendingEmails(limit int) ([]struct {
		ID            int64
		RecipientEmail string
		RecipientName  string
		Subject        string
		EmailFilePath  string
	}, error)
	
	// UpdateEmailStatus actualiza el estado de un email
	UpdateEmailStatus(queueID int64, status string, errorMessage *string) error
	
	// LogEmail registra un email en los logs
	LogEmail(recipientEmail, subject, status string, errorMessage *string) error
	
	// GetRecipients obtiene destinatarios según el tipo
	GetRecipients(recipientType string, formID *int64) ([]struct {
		Email string
		Name  string
	}, error)
	
	// GetEmailTemplate obtiene una plantilla de email
	GetEmailTemplate(templateName string) (*struct {
		ID        int64
		Subject   string
		BodyHTML  string
		BodyText  string
		IsActive  bool
	}, error)
}














