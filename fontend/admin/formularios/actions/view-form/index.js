// Importar funciones necesarias
import { apiRequest, API_CONFIG } from '../../../../services/index.js';
import { showNotification } from '../../../../components/index.js';
import { displayFormDetailsInPage, switchTab, toggleFormInfoAccordion } from '../../../shared/formularios.js';
import { initializeCustomization } from '../../../shared/customization.js';

// Establecer página actual para el layout
window.page = 'formularios';

// Funciones helper para manejar modales con clases CSS
function showModal(modalId) {
  const modal = document.getElementById(modalId);
  if (modal) {
    modal.classList.add('active');
  }
}

function hideModal(modalId) {
  const modal = document.getElementById(modalId);
  if (modal) {
    modal.classList.remove('active');
  }
}

// Exportar funciones globales para modales
window.showModal = showModal;
window.hideModal = hideModal;

// Variable global para el ID del formulario
let currentFormId = null;

// Obtener el ID del formulario desde la URL
const urlParams = new URLSearchParams(window.location.search);
const formId = urlParams.get('id');

if (formId) {
  currentFormId = parseInt(formId);
  
  // Validar que el ID sea un número válido
  if (isNaN(currentFormId)) {
    console.error('Error: ID de formulario inválido en URL:', formId);
    const formDetailsContent = document.getElementById('formDetailsContent');
    if (formDetailsContent) {
      formDetailsContent.innerHTML = `
        <div>
          <p>Error: ID de formulario inválido en la URL</p>
          <a href="../../">Volver a Formularios</a>
        </div>
      `;
    }
  } else {
    console.log('ID de formulario desde URL:', currentFormId);
    
    // Crear un elemento oculto para almacenar el ID del formulario (similar al modal)
    const formDetailsContainer = document.getElementById('formDetailsContent');
    if (formDetailsContainer) {
      formDetailsContainer.setAttribute('data-form-id', currentFormId);
    }

  // Cargar detalles del formulario directamente
  window.loadFormDetails = async function loadFormDetails() {
    try {
      // Validar que el ID del formulario sea válido
      if (!currentFormId || isNaN(currentFormId)) {
        throw new Error('ID de formulario inválido');
      }

      console.log('Cargando formulario con ID:', currentFormId);
      console.log('URL completa:', API_CONFIG.ENDPOINTS.FORM(currentFormId));
      
      const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM(currentFormId));
      console.log('Respuesta de la API:', response);
      
      if (response.success && response.form) {
        const form = response.form;
        console.log('Formulario recibido de la API:', form.id, form.form_name);
        
        // Validar que el formulario cargado coincida con el ID solicitado
        if (parseInt(form.id) !== parseInt(currentFormId)) {
          console.error('Error: El formulario cargado no coincide con el ID solicitado', {
            requested: currentFormId,
            received: form.id,
            requestedType: typeof currentFormId,
            receivedType: typeof form.id
          });
          throw new Error(`Error: Se cargó el formulario ${form.id} en lugar del ${currentFormId}`);
        }

        console.log('Formulario cargado correctamente:', form.id, form.form_name);
        
        // Limpiar contenido anterior antes de cargar nuevo
        const formDetailsContent = document.getElementById('formDetailsContent');
        if (formDetailsContent) {
          // Limpiar solo los contenedores dinámicos, no la estructura HTML
          const sectionsContainer = document.getElementById('formDetailsSections');
          const fieldsContainer = document.getElementById('formDetailsFields');
          if (sectionsContainer) sectionsContainer.innerHTML = '';
          if (fieldsContainer) fieldsContainer.innerHTML = '';
        }

        // Schema en memoria ANTES de pintar secciones/campos (table-per-form)
        setCurrentFormSchema(form);

        // Usar displayFormDetailsInPage para páginas HTML separadas
        await displayFormDetailsInPage(form);

        // Rellenar metadatos extra (borrador/fechas/límite) y botones de visibilidad
        const formatMaybeDateTime = (v) => {
          if (!v) return '';
          if (window.formatDate) return window.formatDate(v);
          try {
            return new Date(v).toLocaleString('es-ES');
          } catch {
            return String(v);
          }
        };

        const setText = (id, text) => {
          const el = document.getElementById(id);
          if (el) el.textContent = text ?? '';
        };

        const limitText = (form.inscription_limit === undefined || form.inscription_limit === null)
          ? 'Sin límite'
          : String(form.inscription_limit);

        setText('formDetailsEventDate', formatMaybeDateTime(form.event_date));
        setText('formDetailsOpenAt', formatMaybeDateTime(form.open_at));
        setText('formDetailsInscriptionLimit', limitText);
        setText('formDetailsEventDateMobile', formatMaybeDateTime(form.event_date));
        setText('formDetailsOpenAtMobile', formatMaybeDateTime(form.open_at));
        setText('formDetailsInscriptionLimitMobile', limitText);

        const publishBtn = document.getElementById('formDetailsPublishBtn');
        const draftBtn = document.getElementById('formDetailsDraftBtn');
        const publishBtnMobile = document.getElementById('formDetailsPublishBtnMobile');
        const draftBtnMobile = document.getElementById('formDetailsDraftBtnMobile');

        const applyButtonVisibility = () => {
          const showPublish = !!form.is_draft;
          const showDraft = !showPublish;

          if (publishBtn) publishBtn.style.display = showPublish ? 'inline-block' : 'none';
          if (draftBtn) draftBtn.style.display = showDraft ? 'inline-block' : 'none';
          if (publishBtnMobile) publishBtnMobile.style.display = showPublish ? 'inline-block' : 'none';
          if (draftBtnMobile) draftBtnMobile.style.display = showDraft ? 'inline-block' : 'none';
        };

        applyButtonVisibility();

        const setDraftStatus = async (shouldBeDraft) => {
          try {
            const endpoint = shouldBeDraft
              ? API_CONFIG.ENDPOINTS.FORM_DRAFT(currentFormId)
              : API_CONFIG.ENDPOINTS.FORM_PUBLISH(currentFormId);

            const btns = [publishBtn, draftBtn, publishBtnMobile, draftBtnMobile].filter(Boolean);
            btns.forEach(b => b.disabled = true);

            const res = await apiRequest(endpoint, { method: 'POST' });
            if (res.success) {
              showNotification(res.message || 'Actualizado correctamente', 'success');
              await window.loadFormDetails();
              btns.forEach(b => b.disabled = false);
            } else {
              showNotification(res.message || res.error || 'Error al actualizar visibilidad', 'error');
              btns.forEach(b => b.disabled = false);
            }
          } catch (err) {
            console.error('Error setDraftStatus:', err);
            showNotification((err && err.message) ? err.message : 'Error al actualizar visibilidad', 'error');
            const btns = [publishBtn, draftBtn, publishBtnMobile, draftBtnMobile].filter(Boolean);
            btns.forEach(b => b.disabled = false);
          }
        };

        if (publishBtn) publishBtn.onclick = () => setDraftStatus(false);
        if (draftBtn) draftBtn.onclick = () => setDraftStatus(true);
        if (publishBtnMobile) publishBtnMobile.onclick = () => setDraftStatus(false);
        if (draftBtnMobile) draftBtnMobile.onclick = () => setDraftStatus(true);

        initSchemaEditor(form);
        // setCurrentFormSchema ya se llamó antes de displayFormDetailsInPage
        // Para que editField/deleteField en formularios.js tengan el form en view-form (no está en forms[])
        window.currentForm = form;
        // Inicializar tabs: mostrar tab de Secciones por defecto
        switchTab('sections');
      } else {
        throw new Error('No se pudo cargar el formulario');
      }
    } catch (error) {
      console.error('Error cargando formulario:', error);
      showNotification('Error cargando formulario: ' + error.message, 'error');
      const contentContainer = document.getElementById('formDetailsContent');
      if (contentContainer) {
        contentContainer.innerHTML = `
          <div>
            <p>Error: ${error.message}</p>
            <a href="../../">Volver a Formularios</a>
          </div>
        `;
      }
    }
  }

// Mantener copia del schema para editar desde Secciones y Campos (table-per-form)
function setCurrentFormSchema(form) {
  if (!form || !currentFormId) return;
  try {
    let schema = form.schema;
    if (schema == null) {
      window.currentFormSchema = { sections: [], fields: [] };
      return;
    }
    if (typeof schema === 'string') schema = JSON.parse(schema);
    window.currentFormSchema = {
      sections: Array.isArray(schema.sections) ? schema.sections.map(s => ({
        id: s.id || ('s' + (schema.sections.indexOf(s) + 1)),
        section_title: s.section_title || s.sectionTitle || '',
        section_icon: s.section_icon || s.sectionIcon || 'fa-folder',
        display_order: s.display_order != null ? s.display_order : schema.sections.indexOf(s)
      })) : [],
      fields: Array.isArray(schema.fields) ? schema.fields.map((f, i) => ({
        id: f.id || ('f' + (i + 1)),
        field_label: f.field_label || f.fieldLabel || '',
        field_name: f.field_name || f.fieldName || '',
        field_type: f.field_type || f.fieldType || 'text',
        field_options: f.field_options != null ? f.field_options : (f.fieldOptions != null ? f.fieldOptions : null),
        is_required: !!f.is_required,
        max_length: f.max_length != null ? f.max_length : (f.maxLength != null ? f.maxLength : null),
        display_order: f.display_order != null ? f.display_order : (f.displayOrder != null ? f.displayOrder : i),
        section_id: f.section_id != null ? f.section_id : (f.sectionId != null ? f.sectionId : '')
      })) : []
    };
  } catch (e) {
    console.warn('setCurrentFormSchema:', e);
    window.currentFormSchema = { sections: [], fields: [] };
  }
}

// Guardar schema en backend y recargar formulario (para Secciones y Campos)
window.saveFormSchema = async function saveFormSchema() {
  if (!currentFormId || !window.currentFormSchema) return;
  try {
    const endpoint = API_CONFIG.FORMS_ENDPOINTS.FORM_SCHEMA(currentFormId);
    const res = await apiRequest(endpoint, {
      method: 'PUT',
      body: JSON.stringify({ schema: window.currentFormSchema })
    });
    if (res.success !== false) {
      showNotification(res.message || 'Cambios guardados', 'success');
      await window.loadFormDetails();
    } else {
      throw new Error(res.message || 'Error al guardar');
    }
  } catch (e) {
    showNotification('Error: ' + e.message, 'error');
    throw e;
  }
};

// Editor de esquema JSON (table-per-form) — PUT /admin/forms/{id}/schema
function initSchemaEditor(form) {
  const ta = document.getElementById('schemaJsonEditor');
  const saveBtn = document.getElementById('schemaSaveBtn');
  const formatBtn = document.getElementById('schemaFormatBtn');
  const msg = document.getElementById('schemaSaveMsg');
  if (!ta || !form) return;
  try {
    const schema = form.schema;
    if (schema == null) {
      ta.value = '';
      return;
    }
    if (typeof schema === 'string') {
      try {
        ta.value = JSON.stringify(JSON.parse(schema), null, 2);
      } catch {
        ta.value = schema;
      }
    } else {
      ta.value = JSON.stringify(schema, null, 2);
    }
  } catch (e) {
    ta.value = String(form.schema || '');
  }
  if (formatBtn) {
    formatBtn.onclick = () => {
      try {
        ta.value = JSON.stringify(JSON.parse(ta.value), null, 2);
        showNotification('JSON formateado', 'success');
      } catch (e) {
        showNotification('JSON inválido: ' + e.message, 'error');
      }
    };
  }
  if (saveBtn) {
    saveBtn.onclick = async () => {
      if (!currentFormId) return;
      let parsed;
      try {
        parsed = JSON.parse(ta.value);
      } catch (e) {
        showNotification('JSON inválido: ' + e.message, 'error');
        return;
      }
      try {
        const endpoint = API_CONFIG.FORMS_ENDPOINTS.FORM_SCHEMA(currentFormId);
        const res = await apiRequest(endpoint, {
          method: 'PUT',
          body: JSON.stringify({ schema: parsed })
        });
        if (res.success) {
          showNotification(res.message || 'Esquema guardado', 'success');
          if (msg) {
            msg.textContent = res.message || 'OK';
            msg.classList.remove('hidden');
          }
          await window.loadFormDetails();
        } else {
          throw new Error(res.message || 'Error al guardar');
        }
      } catch (e) {
        showNotification('Error: ' + e.message, 'error');
      }
    };
  }
}

  function setupModalEventListeners() {
    // Manejar creación de campo (table-per-form: actualizar schema y PUT)
    const addFieldForm = document.getElementById('addFieldForm');
    if (addFieldForm) {
      addFieldForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        if (!currentFormId) {
          showNotification('Error: No se pudo obtener el ID del formulario', 'error');
          return;
        }

        const fieldLabel = document.getElementById('fieldLabel')?.value?.trim();
        const fieldName = document.getElementById('fieldName')?.value?.trim().replace(/\s+/g, '_') || 'field';
        const fieldType = document.getElementById('fieldType')?.value || 'text';
        const fieldOptions = document.getElementById('fieldOptions')?.value?.trim() || '';
        const fieldRequired = document.getElementById('fieldRequired')?.checked;
        const maxLengthVal = document.getElementById('fieldMaxLength')?.value?.trim();
        const maxLength = maxLengthVal ? parseInt(maxLengthVal, 10) : null;
        const sectionIdSelect = document.getElementById('addFieldSectionId');
        const sectionId = sectionIdSelect ? sectionIdSelect.value : '';

        if (window.currentFormSchema && window.saveFormSchema) {
          const schema = window.currentFormSchema;
          const nextId = 'f' + (schema.fields.length + 1);
          schema.fields.push({
            id: nextId,
            field_label: fieldLabel,
            field_name: fieldName,
            field_type: fieldType,
            field_options: fieldOptions || null,
            is_required: !!fieldRequired,
            max_length: (maxLength != null && !isNaN(maxLength)) ? maxLength : null,
            display_order: schema.fields.length,
            section_id: sectionId || ''
          });
          try {
            await window.saveFormSchema();
            if (window.closeAddFieldModal) window.closeAddFieldModal();
          } catch (err) {
            schema.fields.pop();
          }
          return;
        }

        const payload = { field_label: fieldLabel, field_name: fieldName, field_type: fieldType, field_options: fieldOptions || '', is_required: fieldRequired };
        if (maxLength != null && !isNaN(maxLength)) payload.max_length = maxLength;
        try {
          const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_FIELD_CREATE(currentFormId), { method: 'POST', body: JSON.stringify(payload) });
          if (response.success) {
            showNotification('Campo añadido correctamente', 'success');
            if (window.closeAddFieldModal) window.closeAddFieldModal();
            await loadFormDetails();
          }
        } catch (error) {
          showNotification('Error creando campo: ' + error.message, 'error');
        }
      });
    }

    // Manejar edición de campo (table-per-form: actualizar schema y PUT)
    const editFieldForm = document.getElementById('editFieldForm');
    if (editFieldForm) {
      editFieldForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formId = parseInt(editFieldForm.getAttribute('data-form-id'));
        const fieldId = parseInt(editFieldForm.getAttribute('data-field-id'));

        if (!formId || !fieldId) {
          showNotification('Error: No se pudo obtener el ID del formulario o campo', 'error');
          return;
        }

        const fieldLabel = document.getElementById('editFieldLabel')?.value?.trim();
        const fieldName = document.getElementById('editFieldName')?.value?.trim().replace(/\s+/g, '_') || 'field';
        const fieldType = document.getElementById('editFieldType')?.value || 'text';
        const fieldOptions = document.getElementById('editFieldOptions')?.value?.trim() || '';
        const fieldRequired = document.getElementById('editFieldRequired')?.checked;
        const maxLengthVal = document.getElementById('editFieldMaxLength')?.value?.trim();
        const maxLength = maxLengthVal ? parseInt(maxLengthVal, 10) : null;

        if (window.currentFormSchema && window.saveFormSchema) {
          const schema = window.currentFormSchema;
          const idx = fieldId - 1;
          if (idx >= 0 && idx < schema.fields.length) {
            const f = schema.fields[idx];
            f.field_label = fieldLabel;
            f.field_name = fieldName;
            f.field_type = fieldType;
            f.field_options = fieldOptions || null;
            f.is_required = !!fieldRequired;
            f.max_length = (maxLength != null && !isNaN(maxLength)) ? maxLength : null;
            try {
              await window.saveFormSchema();
              if (window.closeEditFieldModal) window.closeEditFieldModal();
            } catch (err) {
              showNotification('Error al guardar: ' + err.message, 'error');
            }
          }
          return;
        }

        const payload = { field_label: fieldLabel, field_name: fieldName, field_type: fieldType, field_options: fieldOptions || '', is_required: fieldRequired };
        if (maxLength != null && !isNaN(maxLength)) payload.max_length = maxLength;
        try {
          const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_FIELD_UPDATE(formId, fieldId), { method: 'PUT', body: JSON.stringify(payload) });
          if (response.success) {
            showNotification('Campo actualizado correctamente', 'success');
            if (window.closeEditFieldModal) window.closeEditFieldModal();
            await loadFormDetails();
          }
        } catch (error) {
          showNotification('Error actualizando campo: ' + error.message, 'error');
        }
      });
    }
  }

    // Inicializar cuando el DOM esté listo
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => {
        loadFormDetails();
        setupModalEventListeners();
      });
    } else {
      loadFormDetails();
      setupModalEventListeners();
    }
  }
} else {
  document.getElementById('formDetailsContent').innerHTML = `
    <div>
      <p>Error: No se proporcionó un ID de formulario</p>
      <a href="../../">Volver a Formularios</a>
    </div>
  `;
}

