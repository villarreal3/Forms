// Importar funciones necesarias
import { apiRequest, API_CONFIG } from '../../../../services/index.js';
import { showNotification } from '../../../../components/index.js';

// Función para inicializar la página de creación
export function initCreateFormPage() {
  // Establecer fecha mínima (ahora, hora local) para el input de expiración.
  const now = new Date();
  const localISO = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}T${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`;
  const expiresAtInput = document.getElementById('formExpiresAt');
  if (expiresAtInput) {
    expiresAtInput.min = localISO;
  }

  // Habilitar/deshabilitar límite de inscripciones
  const limitToggle = document.getElementById('formUseInscriptionLimit');
  const limitInput = document.getElementById('formInscriptionLimit');
  if (limitToggle && limitInput) {
    const sync = () => {
      if (limitToggle.checked) {
        limitInput.disabled = false;
      } else {
        limitInput.disabled = true;
        limitInput.value = '';
      }
    };
    limitToggle.addEventListener('change', sync);
    sync();
  }

  // Manejar envío del formulario
  const form = document.getElementById('createFormForm');
  if (form) {
    form.addEventListener('submit', handleFormSubmit);
  }
}

// Función para manejar el envío del formulario
async function handleFormSubmit(e) {
  e.preventDefault();

  const formName = document.getElementById('formName')?.value.trim();
  const description = document.getElementById('formDescription')?.value.trim();
  const expiresAt = document.getElementById('formExpiresAt')?.value;
  const eventDate = document.getElementById('formEventDate')?.value;
  const openAt = document.getElementById('formOpenAt')?.value;
  const template = document.getElementById('formTemplate')?.value;
  const limitEnabled = document.getElementById('formUseInscriptionLimit')?.checked;
  const limitRaw = document.getElementById('formInscriptionLimit')?.value;
  const conferencista = document.getElementById('formConferencista')?.value.trim() ?? '';
  const ubicacion = document.getElementById('formUbicacion')?.value.trim() ?? '';

  const formNameEl = document.getElementById('formName');
  const eventDateEl = document.getElementById('formEventDate');
  const openAtEl = document.getElementById('formOpenAt');
  const expiresAtEl = document.getElementById('formExpiresAt');
  const templateEl = document.getElementById('formTemplate');
  const limitInputEl = document.getElementById('formInscriptionLimit');

  function errorIdFor(el) {
    if (!el || !el.id) return null;
    return `${el.id}-error`;
  }

  function clearFieldError(el) {
    if (!el) return;
    const errorId = errorIdFor(el);
    el.removeAttribute('aria-invalid');
    if (errorId) el.removeAttribute('aria-describedby');
    if (errorId) {
      const existing = document.getElementById(errorId);
      if (existing) existing.remove();
    }
  }

  function setFieldError(el, message) {
    if (!el) return;
    clearFieldError(el);
    const errorId = errorIdFor(el);
    if (errorId) {
      el.setAttribute('aria-invalid', 'true');
      el.setAttribute('aria-describedby', errorId);

      const p = document.createElement('p');
      p.id = errorId;
      p.setAttribute('role', 'alert');
      p.style.cssText = 'color: var(--color-error); font-size: var(--font-size-xs); margin-top: var(--spacing-xs);';
      p.textContent = message;

      // Insertar mensaje justo después del input/select/textarea
      el.insertAdjacentElement('afterend', p);
    }
    try {
      el.focus();
    } catch (_) {
      // Ignorar si el elemento no puede hacer foco
    }
  }

  // Validaciones
  [formNameEl, eventDateEl, openAtEl, expiresAtEl, templateEl, limitInputEl].forEach(clearFieldError);

  if (!formName) {
    setFieldError(formNameEl, 'Ingresa el nombre del formulario.');
    return;
  }

  if (!eventDate) {
    setFieldError(eventDateEl, 'Selecciona la fecha del evento.');
    return;
  }

  if (!openAt) {
    setFieldError(openAtEl, 'Selecciona la fecha de apertura del formulario.');
    return;
  }

  if (!expiresAt) {
    setFieldError(expiresAtEl, 'Selecciona una fecha de expiración.');
    return;
  }

  if (!template) {
    setFieldError(templateEl, 'Selecciona una plantilla para el formulario.');
    return;
  }

  // Validación básica de orden: apertura no puede ser después de expiración.
  // No se valida "fecha futura" aquí porque el servidor puede estar en otra zona horaria;
  // esa validación la hace el backend contra su propio reloj.
  if (openAt > expiresAt) {
    setFieldError(openAtEl, 'La fecha de apertura debe ser anterior a la de expiración.');
    return;
  }

  // datetime-local devuelve hora local como "YYYY-MM-DDTHH:MM".
  // Se envía tal cual al backend, sin conversión a UTC.
  const expiresAtISO = expiresAt.replace('T', ' ') + ':00';
  const openAtISO = openAt.replace('T', ' ') + ':00';
  const eventDateISO = eventDate.replace('T', ' ') + ':00';

  const inscriptionLimit = limitEnabled
    ? (limitRaw && String(limitRaw).trim() !== '' ? parseInt(limitRaw, 10) : null)
    : null;

  if (limitEnabled && (inscriptionLimit === null || isNaN(inscriptionLimit))) {
    setFieldError(limitInputEl, 'Ingresa un número válido para el límite de inscripciones (ej. 50).');
    return;
  }

  // Determinar si usar plantilla por defecto
  const useDefaultTemplate = template === 'default';

  try {
    // Mostrar indicador de carga
    const submitButton = e.target.querySelector('button[type="submit"]');
    const originalText = submitButton?.textContent;
    if (submitButton) {
      submitButton.disabled = true;
      submitButton.textContent = 'Creando...';
    }

    const response = await apiRequest(API_CONFIG.ENDPOINTS.FORMS, {
      method: 'POST',
      body: JSON.stringify({
        form_name: formName,
        description: description || '',
        expires_at: expiresAtISO,
        event_date: eventDateISO,
        open_at: openAtISO,
        inscription_limit: inscriptionLimit,
        conferencista: conferencista || undefined,
        ubicacion: ubicacion || undefined,
        use_default_template: useDefaultTemplate
      })
    });

    if (response.success) {
      const fieldsAdded = response.fields_added ?? 0;
      // Table-per-form: backend puede devolver fields_added desde schema sin warning
      const hasWarning = response.warning && fieldsAdded === 0 && useDefaultTemplate;
      const message = response.message || (useDefaultTemplate
        ? `Formulario creado correctamente con ${fieldsAdded} campos en el esquema`
        : 'Formulario creado correctamente');
      showNotification(message, hasWarning ? 'warning' : 'success');
      if (hasWarning && response.warning) {
        showNotification(response.warning, 'error');
      }

      // Redirigir a la lista de formularios después de un breve delay
      setTimeout(() => {
        window.location.href = '../../';
      }, 1000);
    } else {
      throw new Error(response.message || 'Error al crear el formulario');
    }
  } catch (error) {
    console.error('Error creando formulario:', error);
    showNotification('Error creando formulario: ' + error.message, 'error');

    // Restaurar botón
    const submitButton = e.target.querySelector('button[type="submit"]');
    if (submitButton) {
      submitButton.disabled = false;
      submitButton.textContent = 'Crear Formulario';
    }
  }
}

// Inicializar cuando el DOM esté listo
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initCreateFormPage);
} else {
  initCreateFormPage();
}

// Exportar para uso global si es necesario
window.initCreateFormPage = initCreateFormPage;



