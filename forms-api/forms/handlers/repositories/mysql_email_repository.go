package repositories

import (
	"database/sql"
	"log"

	"forms/database"
)

type mysqlEmailRepository struct{}

func NewEmailRepository() EmailRepository {
	return &mysqlEmailRepository{}
}

func (r *mysqlEmailRepository) AddToQueue(recipientEmail, recipientName, subject, emailFilePath string) (int64, error) {
	_, err := database.DB.Exec("CALL sp_add_email_to_queue(?, ?, ?, ?, @p_queue_id)", recipientEmail, recipientName, subject, emailFilePath)
	if err != nil {
		return 0, err
	}
	var queueID int64
	err = database.DB.QueryRow("SELECT @p_queue_id").Scan(&queueID)
	return queueID, err
}

func (r *mysqlEmailRepository) GetPendingEmails(limit int) ([]struct {
	ID            int64
	RecipientEmail string
	RecipientName  string
	Subject        string
	EmailFilePath  string
}, error) {
	rows, err := database.DB.Query("CALL sp_get_pending_emails(?)", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []struct {
		ID            int64
		RecipientEmail string
		RecipientName  string
		Subject        string
		EmailFilePath  string
	}

	for rows.Next() {
		var email struct {
			ID            int64
			RecipientEmail string
			RecipientName  string
			Subject        string
			EmailFilePath  string
		}
		if err := rows.Scan(&email.ID, &email.RecipientEmail, &email.RecipientName, &email.Subject, &email.EmailFilePath); err != nil {
			log.Printf("Error escaneando email: %v", err)
			continue
		}
		emails = append(emails, email)
	}

	for rows.NextResultSet() {
	}

	return emails, nil
}

func (r *mysqlEmailRepository) UpdateEmailStatus(queueID int64, status string, errorMessage *string) error {
	_, err := database.DB.Exec(
		"CALL sp_update_email_status(?, ?, ?)",
		queueID, status, errorMessage,
	)
	return err
}

func (r *mysqlEmailRepository) LogEmail(recipientEmail, subject, status string, errorMessage *string) error {
	_, err := database.DB.Exec(
		"CALL sp_log_email(?, ?, ?, ?)",
		recipientEmail, subject, status, errorMessage,
	)
	return err
}

func (r *mysqlEmailRepository) GetRecipients(recipientType string, formID *int64) ([]struct {
	Email string
	Name  string
}, error) {
	var rows *sql.Rows
	var err error

	if formID != nil {
		rows, err = database.DB.Query("CALL sp_get_email_recipients(?, ?)", recipientType, *formID)
	} else {
		rows, err = database.DB.Query("CALL sp_get_email_recipients(?, NULL)", recipientType)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipients []struct {
		Email string
		Name  string
	}

	for rows.Next() {
		var recipient struct {
			Email string
			Name  string
		}
		if err := rows.Scan(&recipient.Email, &recipient.Name); err != nil {
			log.Printf("Error escaneando recipient: %v", err)
			continue
		}
		recipients = append(recipients, recipient)
	}

	for rows.NextResultSet() {
	}

	return recipients, nil
}

func (r *mysqlEmailRepository) GetEmailTemplate(templateName string) (*struct {
	ID        int64
	Subject   string
	BodyHTML  string
	BodyText  string
	IsActive  bool
}, error) {
	row := database.DB.QueryRow("CALL sp_get_email_template(?)", templateName)

	var template struct {
		ID        int64
		Subject   string
		BodyHTML  string
		BodyText  sql.NullString
		IsActive  bool
	}

	err := row.Scan(&template.ID, &template.Subject, &template.BodyHTML, &template.BodyText, &template.IsActive)
	if err != nil {
		return nil, err
	}

	result := &struct {
		ID        int64
		Subject   string
		BodyHTML  string
		BodyText  string
		IsActive  bool
	}{
		ID:       template.ID,
		Subject:  template.Subject,
		BodyHTML: template.BodyHTML,
		IsActive: template.IsActive,
	}

	if template.BodyText.Valid {
		result.BodyText = template.BodyText.String
	}

	return result, nil
}

