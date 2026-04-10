package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"forms/database"
	"forms/models"

	"github.com/jordan-wright/email"

	"forms/handlers/helpers"
)

// emailConfig contiene la configuración de correo
var emailConfig = struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
	EmailsDir    string
}{
	SMTPHost:     "smtp.gmail.com", // Cambiar según el proveedor
	SMTPPort:     "587",
	SMTPUser:     "", // Configurar en .env
	SMTPPassword: "", // Configurar en .env
	FromEmail:    "noreply@example.com",
	FromName:     "Forms",
	EmailsDir:    "src/emails", // Directorio donde se guardan los correos
}

// Recipient representa un destinatario de correo
type Recipient struct {
	Email string
	Name  string
}

// InitEmailConfig inicializa la configuración de correos (llamar desde main.go)
func InitEmailConfig(emailsDir string) {
	if emailsDir != "" {
		emailConfig.EmailsDir = emailsDir
	}

	// Crear directorio si no existe
	if err := os.MkdirAll(emailConfig.EmailsDir, 0755); err != nil {
		log.Printf("Error creando directorio de emails: %v", err)
	} else {
		log.Printf("Directorio de emails configurado: %s", emailConfig.EmailsDir)
	}
}

// SendBulkEmail envía correos masivos
func SendBulkEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	var req models.BulkEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error decodificando JSON: %v", err))
		return
	}

	if req.Subject == "" || req.BodyHTML == "" {
		helpers.RespondError(w, http.StatusBadRequest, "Subject y body_html son requeridos")
		return
	}

	// Obtener lista de destinatarios
	recipients, err := getRecipients(req.Recipients)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Error obteniendo destinatarios: %v", err))
		return
	}

	if len(recipients) == 0 {
		helpers.RespondError(w, http.StatusBadRequest, "No se encontraron destinatarios")
		return
	}

	// Guardar el correo como archivo base (template con variables)
	timestamp := time.Now().Format("20060102_150405")
	emailFileName := fmt.Sprintf("email_template_%s.html", timestamp)
	emailFilePath := filepath.Join(emailConfig.EmailsDir, emailFileName)

	// Crear contenido del correo template (con variables sin reemplazar)
	emailContent := req.BodyHTML
	if req.BodyText != "" {
		// Si hay texto plano, agregarlo como comentario en el HTML
		emailContent = fmt.Sprintf("<!-- Plain text template: %s -->\n%s", req.BodyText, emailContent)
	}

	// Guardar archivo template
	if err := ioutil.WriteFile(emailFilePath, []byte(emailContent), 0644); err != nil {
		log.Printf("Error guardando archivo de correo: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error guardando correo")
		return
	}

	// Agregar a la cola de correos (todos comparten el mismo template)
	// Las variables se reemplazan al enviar
	successCount := 0
	for _, recipient := range recipients {
		_, err := database.DB.Exec(`
			INSERT INTO email_queue (recipient_email, recipient_name, subject, email_file_path, status)
			VALUES (?, ?, ?, ?, 'pending')
		`, recipient.Email, recipient.Name, req.Subject, emailFilePath)

		if err == nil {
			successCount++
		}
	}

	// Procesar cola en segundo plano
	go processEmailQueue()

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":      true,
		"message":      fmt.Sprintf("Se agregaron %d correos a la cola", successCount),
		"queued_count": successCount,
		"email_file":   emailFilePath,
	})
}

// getRecipients obtiene la lista de destinatarios según los criterios especificados
func getRecipients(recipientList []string) ([]Recipient, error) {
	var recipients []Recipient

	for _, recipient := range recipientList {
		if recipient == "all_users" {
			// Obtener todos los usuarios
			rows, err := database.DB.Query("SELECT email, CONCAT(first_name, ' ', last_name) as name FROM users WHERE email IS NOT NULL AND email != ''")
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var r Recipient
				if err := rows.Scan(&r.Email, &r.Name); err == nil {
					recipients = append(recipients, r)
				}
			}
		} else if strings.HasPrefix(recipient, "form_submitters:") {
			// Obtener usuarios que enviaron un formulario específico
			formID := strings.TrimPrefix(recipient, "form_submitters:")
			rows, err := database.DB.Query(`
				SELECT DISTINCT u.email, CONCAT(u.first_name, ' ', u.last_name) as name
				FROM users u
				INNER JOIN form_answers fa ON u.id = fa.user_id
				WHERE fa.form_id = ? AND u.email IS NOT NULL AND u.email != ''
			`, formID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var r Recipient
				if err := rows.Scan(&r.Email, &r.Name); err == nil {
					recipients = append(recipients, r)
				}
			}
		} else {
			// Email directo
			recipients = append(recipients, Recipient{Email: recipient, Name: ""})
		}
	}

	return recipients, nil
}

// processEmailQueue procesa la cola de correos
func processEmailQueue() {
	rows, err := database.DB.Query(`
		SELECT id, recipient_email, recipient_name, subject, email_file_path
		FROM email_queue
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT 50
	`)

	if err != nil {
		log.Printf("Error obteniendo cola de correos: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var emailID int64
		var recipientEmail, recipientName, subject, emailFilePath sql.NullString

		if err := rows.Scan(&emailID, &recipientEmail, &recipientName, &subject, &emailFilePath); err != nil {
			continue
		}

		if !recipientEmail.Valid || !subject.Valid || !emailFilePath.Valid {
			continue
		}

		// Leer el contenido del correo desde el archivo
		bodyHTML, bodyText, err := readEmailFromFile(emailFilePath.String)
		if err != nil {
			log.Printf("Error leyendo archivo de correo %s: %v", emailFilePath.String, err)
			database.DB.Exec("UPDATE email_queue SET status = 'failed', error_message = ? WHERE id = ?",
				fmt.Sprintf("Error leyendo archivo: %v", err), emailID)
			continue
		}

		// Reemplazar variables en el template si es necesario
		// Variables disponibles: {FIRST_NAME}, {LAST_NAME}, {EMAIL}, {RECIPIENT_NAME}
		bodyHTML = strings.ReplaceAll(bodyHTML, "{RECIPIENT_NAME}", recipientName.String)
		bodyHTML = strings.ReplaceAll(bodyHTML, "{EMAIL}", recipientEmail.String)
		if bodyText != "" {
			bodyText = strings.ReplaceAll(bodyText, "{RECIPIENT_NAME}", recipientName.String)
			bodyText = strings.ReplaceAll(bodyText, "{EMAIL}", recipientEmail.String)
		}

		err = sendEmail(recipientEmail.String, recipientName.String, subject.String, bodyHTML, bodyText)

		if err != nil {
			log.Printf("Error enviando correo a %s: %v", recipientEmail.String, err)
			database.DB.Exec("UPDATE email_queue SET status = 'failed', error_message = ? WHERE id = ?", err.Error(), emailID)
		} else {
			database.DB.Exec("UPDATE email_queue SET status = 'sent', sent_at = ? WHERE id = ?", time.Now(), emailID)
			// Registrar en logs
			database.DB.Exec("INSERT INTO email_logs (recipient_email, subject, status) VALUES (?, ?, 'sent')",
				recipientEmail.String, subject.String)
		}

		time.Sleep(100 * time.Millisecond) // Pequeña pausa entre envíos
	}
}

// readEmailFromFile lee el contenido del correo desde un archivo
func readEmailFromFile(filePath string) (bodyHTML, bodyText string, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("error leyendo archivo %s: %v", filePath, err)
	}

	bodyHTML = string(content)

	// Intentar extraer texto plano del comentario si existe
	// Formato: <!-- Plain text: contenido texto --> o <!-- Plain text template: contenido -->
	if strings.Contains(bodyHTML, "<!-- Plain text") {
		// Buscar el comentario de texto plano
		startIdx := strings.Index(bodyHTML, "<!-- Plain text")
		if startIdx != -1 {
			endIdx := strings.Index(bodyHTML[startIdx:], "-->")
			if endIdx != -1 {
				comment := bodyHTML[startIdx : startIdx+endIdx+3]
				// Extraer el texto plano
				if strings.Contains(comment, "Plain text template:") {
					bodyText = strings.TrimSpace(strings.TrimPrefix(
						strings.TrimSuffix(comment, "-->"),
						"<!-- Plain text template:"))
				} else if strings.Contains(comment, "Plain text:") {
					bodyText = strings.TrimSpace(strings.TrimPrefix(
						strings.TrimSuffix(comment, "-->"),
						"<!-- Plain text:"))
				}
				// Remover el comentario del HTML
				bodyHTML = strings.TrimSpace(bodyHTML[startIdx+endIdx+3:])
			}
		}
	}

	return bodyHTML, bodyText, nil
}

// sendEmail envía un email individual
func sendEmail(to, name, subject, bodyHTML, bodyText string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", emailConfig.FromName, emailConfig.FromEmail)
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(bodyHTML)
	if bodyText != "" {
		e.Text = []byte(bodyText)
	}

	// En producción, configurar credenciales SMTP reales
	if emailConfig.SMTPUser == "" || emailConfig.SMTPPassword == "" {
		log.Printf("⚠️ Configuración SMTP no disponible. Correo no enviado a %s", to)
		return nil // Simular éxito para desarrollo
	}

	return e.Send(fmt.Sprintf("%s:%s", emailConfig.SMTPHost, emailConfig.SMTPPort),
		smtp.PlainAuth("", emailConfig.SMTPUser, emailConfig.SMTPPassword, emailConfig.SMTPHost))
}

// SendConfirmationEmail envía correo de confirmación al usuario
func SendConfirmationEmail(userEmail, userName, formName, idNumber string) error {
	// Obtener plantilla
	var template models.EmailTemplate
	err := database.DB.QueryRow(`
		SELECT subject, body_html, body_text
		FROM email_templates
		WHERE template_name = 'form_submission_confirmation' AND is_active = 1
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&template.Subject, &template.BodyHTML, &template.BodyText)

	if err != nil {
		if err == sql.ErrNoRows {
			// Usar plantilla por defecto
			template.Subject = "Confirmación de Inscripción"
			template.BodyHTML = fmt.Sprintf(`
				<html><body>
					<h2>¡Gracias por tu inscripción!</h2>
					<p>Estimado/a %s,</p>
					<p>Hemos recibido tu inscripción correctamente.</p>
					<p><strong>Detalles de tu inscripción:</strong></p>
					<ul>
						<li>Formulario: %s</li>
						<li>Fecha: %s</li>
						<li>Cédula: %s</li>
						<li>Email: %s</li>
					</ul>
					<p>Nos pondremos en contacto contigo pronto.</p>
					<p>Saludos,<br>Equipo de formularios</p>
				</body></html>
			`, userName, formName, time.Now().Format("2006-01-02 15:04:05"), idNumber, userEmail)
			template.BodyText = fmt.Sprintf("Gracias por tu inscripción!\n\nEstimado/a %s,\n\nHemos recibido tu inscripción correctamente.\n\nFormulario: %s\nFecha: %s\nCédula: %s\nEmail: %s\n\nSaludos,\nEquipo de formularios",
				userName, formName, time.Now().Format("2006-01-02 15:04:05"), idNumber, userEmail)
		} else {
			return err
		}
	}

	// Reemplazar variables en la plantilla
	bodyHTML := strings.ReplaceAll(template.BodyHTML, "{FIRST_NAME}", userName)
	bodyHTML = strings.ReplaceAll(bodyHTML, "{LAST_NAME}", "")
	bodyHTML = strings.ReplaceAll(bodyHTML, "{FORM_NAME}", formName)
	bodyHTML = strings.ReplaceAll(bodyHTML, "{SUBMISSION_DATE}", time.Now().Format("2006-01-02 15:04:05"))
	bodyHTML = strings.ReplaceAll(bodyHTML, "{ID_NUMBER}", idNumber)
	bodyHTML = strings.ReplaceAll(bodyHTML, "{EMAIL}", userEmail)

	bodyText := strings.ReplaceAll(template.BodyText, "{FIRST_NAME}", userName)
	bodyText = strings.ReplaceAll(bodyText, "{LAST_NAME}", "")
	bodyText = strings.ReplaceAll(bodyText, "{FORM_NAME}", formName)
	bodyText = strings.ReplaceAll(bodyText, "{SUBMISSION_DATE}", time.Now().Format("2006-01-02 15:04:05"))
	bodyText = strings.ReplaceAll(bodyText, "{ID_NUMBER}", idNumber)
	bodyText = strings.ReplaceAll(bodyText, "{EMAIL}", userEmail)

	// Guardar correo como archivo
	timestamp := time.Now().Format("20060102_150405")
	emailFileName := fmt.Sprintf("confirmation_%s_%s.html", idNumber, timestamp)
	emailFilePath := filepath.Join(emailConfig.EmailsDir, emailFileName)

	// Crear contenido del correo
	emailContent := bodyHTML
	if bodyText != "" {
		emailContent = fmt.Sprintf("<!-- Plain text: %s -->\n%s", bodyText, emailContent)
	}

	// Guardar archivo
	if err := ioutil.WriteFile(emailFilePath, []byte(emailContent), 0644); err != nil {
		log.Printf("Error guardando archivo de correo de confirmación: %v", err)
		return err
	}

	// Agregar a la cola (solo la ruta del archivo)
	_, err = database.DB.Exec(`
		INSERT INTO email_queue (recipient_email, recipient_name, subject, email_file_path, status)
		VALUES (?, ?, ?, ?, 'pending')
	`, userEmail, userName, template.Subject, emailFilePath)

	if err != nil {
		return err
	}

	// Procesar cola en segundo plano
	go processEmailQueue()

	return nil
}

