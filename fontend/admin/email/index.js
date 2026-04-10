// Importar funciones necesarias
import { apiRequest, API_CONFIG } from '../services/index.js';
import { showNotification } from '../components/index.js';

async function loadForms() {
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.FORMS);
 if (response.success && response.forms) {
 const select = document.getElementById('emailRecipients');
 response.forms.forEach(form => {
 const option = document.createElement('option');
 option.value = `form_submitters:${form.id}`;
 option.textContent = `Usuarios que enviaron: ${form.form_name}`;
 select.appendChild(option);
 });
 }
 } catch (error) {
 console.error('Error cargando formularios:', error);
 }
}

document.getElementById('emailForm').addEventListener('submit', async (e) => {
 e.preventDefault();
 
 const recipientsSelect = document.getElementById('emailRecipients');
 const customRecipients = document.getElementById('emailCustomRecipients').value;
 const subject = document.getElementById('emailSubject').value;
 const bodyHTML = document.getElementById('emailBody').value;
 const bodyText = document.getElementById('emailBodyText').value;
 
 // Obtener destinatarios seleccionados
 const recipients = Array.from(recipientsSelect.selectedOptions).map(opt => opt.value);
 
 // Agregar destinatarios personalizados
 if (customRecipients) {
 const customEmails = customRecipients.split(/[,\n]/)
 .map(email => email.trim())
 .filter(email => email && email.includes('@'));
 recipients.push(...customEmails);
 }
 
 if (recipients.length === 0) {
 showNotification('Debes seleccionar al menos un destinatario', 'error');
 return;
 }
 
 if (!subject || !bodyHTML) {
 showNotification('El asunto y el cuerpo del mensaje son requeridos', 'error');
 return;
 }
 
 const submitBtn = e.target.querySelector('button[type="submit"]');
 const submitText = submitBtn ? submitBtn.querySelector('span') : null;
 const originalText = submitText ? submitText.textContent : '';
 if (submitBtn) submitBtn.disabled = true;
 if (submitText) submitText.textContent = 'Enviando...';
 
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.EMAIL_BULK, {
 method: 'POST',
 body: JSON.stringify({
 recipients: recipients,
 subject: subject,
 body_html: bodyHTML,
 body_text: bodyText || undefined
 })
 });
 
 if (response.success) {
 showNotification(`Correos agregados a la cola: ${response.queued_count || 0}`, 'success');
 document.getElementById('emailForm').reset();
 }
 } catch (error) {
 console.error('Error enviando correos:', error);
 showNotification('Error enviando correos: ' + error.message, 'error');
 } finally {
 if (submitBtn) submitBtn.disabled = false;
 if (submitText) submitText.textContent = originalText;
 }
});

// Cargar formularios cuando la página esté lista
if (document.readyState === 'loading') {
 document.addEventListener('DOMContentLoaded', loadForms);
} else {
 loadForms();
}
