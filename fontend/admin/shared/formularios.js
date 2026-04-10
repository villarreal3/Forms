/**
 * Vista detalle / edición de formulario (view-form.html).
 * Exporta: displayFormDetailsInPage, switchTab, toggleFormInfoAccordion
 * y funciones globales para modales (openAddFieldModal, editField, etc.).
 */

import { loadFormSections } from './form_sections.js';
import { initializeCustomization } from './customization.js';
import { showNotification } from '../../components/index.js';

function escapeHtml(s) {
  if (s == null) return '';
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function formatDateSafe(v) {
  if (!v) return '—';
  if (typeof window.formatDate === 'function') return window.formatDate(v);
  try {
    return new Date(v).toLocaleString('es-PA');
  } catch {
    return String(v);
  }
}

function statusBadgeHtml(form) {
  if (form.is_closed) {
    return '<span class="status-badge closed">Cerrado</span>';
  }
  if (form.is_draft) {
    return '<span class="status-badge draft">Borrador</span>';
  }
  const now = new Date();
  const parseBackendDate = (dateValue) => {
    if (!dateValue) return null;
    if (typeof dateValue !== 'string') return new Date(dateValue);
    let s = dateValue.trim();
    if (/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}/.test(s)) s = s.replace(' ', 'T');
    s = s.replace(/Z$/, '').replace(/([+-]\d{2}:\d{2})$/, '');
    return new Date(s);
  };
  const expiresAt = form.expires_at ? parseBackendDate(form.expires_at) : null;
  const openAt = form.open_at ? parseBackendDate(form.open_at) : null;
  if (openAt && openAt > now) {
    return '<span class="status-badge future_open">Próxima apertura</span>';
  }
  if (expiresAt && expiresAt < now) {
    return '<span class="status-badge expired">Expirado</span>';
  }
  return '<span class="status-badge active">Activo</span>';
}

function setDetailPair(suffix, value) {
  const a = document.getElementById(`formDetails${suffix}`);
  const b = document.getElementById(`formDetails${suffix}Mobile`);
  const text = value == null || value === '' ? '—' : String(value);
  if (a) a.textContent = text;
  if (b) b.textContent = text;
}

/**
 * Lista campos en #formDetailsFields (usa window.currentFormSchema si existe; si no, form.schema / form.fields).
 */
export function renderFormFieldsList(form) {
  const container = document.getElementById('formDetailsFields');
  if (!container) return;

  let fields = [];
  if (window.currentFormSchema && Array.isArray(window.currentFormSchema.fields)) {
    fields = window.currentFormSchema.fields;
  } else if (form.fields && form.fields.length) {
    fields = form.fields;
  } else if (form.schema) {
    try {
      const sch = typeof form.schema === 'string' ? JSON.parse(form.schema) : form.schema;
      if (sch && Array.isArray(sch.fields)) fields = sch.fields;
    } catch {
      /* ignore */
    }
  }

  if (!fields.length) {
    container.innerHTML =
      '<p style="padding:1rem;color:var(--color-text-secondary)">No hay campos. Usa <strong>Añadir Campo</strong> o el tab <strong>Esquema (JSON)</strong>.</p>';
    return;
  }

  const formId = form.id;
  container.innerHTML = fields
    .map((f, idx) => {
      const i = idx + 1;
      const label = f.field_label || f.fieldLabel || 'Campo';
      const name = f.field_name || f.fieldName || '';
      const type = f.field_type || f.fieldType || 'text';
      const req = f.is_required ? ' · Requerido' : '';
      return `
      <div class="field-card-view">
        <div>
          <h4>${escapeHtml(label)}</h4>
          <p><code>${escapeHtml(name)}</code> · ${escapeHtml(type)}${req}</p>
        </div>
        <div class="field-card-view__actions">
          <button type="button" class="btn-edit-field" onclick="window.editField(${formId}, ${i})">Editar</button>
          <button type="button" class="btn-delete-field" onclick="window.deleteField(${formId}, ${i})">Eliminar</button>
        </div>
      </div>`;
    })
    .join('');
}

/**
 * Carga cabecera lateral, secciones, lista de campos y personalización.
 */
export async function displayFormDetailsInPage(form) {
  if (!form || !form.id) return;

  const titleEl = document.getElementById('formTitle');
  if (titleEl) titleEl.textContent = form.form_name || 'Formulario';

  setDetailPair('Name', form.form_name || '—');

  const statusHtml = statusBadgeHtml(form);
  const st = document.getElementById('formDetailsStatus');
  const stM = document.getElementById('formDetailsStatusMobile');
  if (st) st.innerHTML = statusHtml;
  if (stM) stM.innerHTML = statusHtml;

  setDetailPair('Created', formatDateSafe(form.created_at));
  setDetailPair('Expires', formatDateSafe(form.expires_at));
  setDetailPair('EventDate', formatDateSafe(form.event_date));
  setDetailPair('OpenAt', formatDateSafe(form.open_at));

  const lim =
    form.inscription_limit === undefined || form.inscription_limit === null
      ? 'Sin límite'
      : String(form.inscription_limit);
  setDetailPair('InscriptionLimit', lim);

  const desc = (form.description && String(form.description).trim()) || '—';
  setDetailPair('Description', desc);

  const conf =
    form.conferencista != null && String(form.conferencista).trim() !== ''
      ? String(form.conferencista).trim()
      : '—';
  const ubi =
    form.ubicacion != null && String(form.ubicacion).trim() !== ''
      ? String(form.ubicacion).trim()
      : '—';
  setDetailPair('Conferencista', conf);
  setDetailPair('Ubicacion', ubi);

  await loadFormSections(form.id);
  renderFormFieldsList(form);
  await initializeCustomization(form.id);
}

export function switchTab(name) {
  const map = {
    sections: { tab: 'tabSections', content: 'tabContentSections' },
    fields: { tab: 'tabFields', content: 'tabContentFields' },
    personalize: { tab: 'tabPersonalize', content: 'tabContentPersonalize' },
    schema: { tab: 'tabSchema', content: 'tabContentSchema' }
  };
  const m = map[name];
  if (!m) return;

  document.querySelectorAll('.tabs-list .tab-button').forEach((b) => b.classList.remove('active'));
  document.querySelectorAll('.tabs-container .tab-content').forEach((c) => c.classList.remove('active'));

  const tb = document.getElementById(m.tab);
  const tc = document.getElementById(m.content);
  if (tb) tb.classList.add('active');
  if (tc) tc.classList.add('active');

  const addSec = document.getElementById('addSectionButton');
  const addFld = document.getElementById('addFieldButton');
  if (addSec) addSec.style.display = name === 'sections' ? '' : 'none';
  if (addFld) addFld.style.display = name === 'fields' ? '' : 'none';
}

export function toggleFormInfoAccordion() {
  const content = document.getElementById('formInfoAccordionContent');
  const btn = document.querySelector('.form-info-accordion .accordion-button span:last-child');
  if (!content) return;
  const isOpen = content.style.display === 'block' || content.classList.contains('accordion-open');
  if (isOpen) {
    content.style.display = 'none';
    content.classList.remove('accordion-open');
    if (btn) btn.textContent = '▼';
  } else {
    content.style.display = 'block';
    content.classList.add('accordion-open');
    if (btn) btn.textContent = '▲';
  }
}

function showModalEl(id) {
  const el = document.getElementById(id);
  if (el) {
    el.classList.add('active');
    el.style.display = 'flex';
  }
}

function hideModalEl(id) {
  const el = document.getElementById(id);
  if (el) {
    el.classList.remove('active');
    el.style.display = 'none';
  }
}

window.handleFieldTypeChange = function handleFieldTypeChange() {
  const sel = document.getElementById('fieldType');
  const optCont = document.getElementById('fieldOptionsContainer');
  if (!sel || !optCont) return;
  const show = sel.value === 'select';
  optCont.style.display = show ? 'block' : 'none';
  const inp = document.getElementById('fieldOptions');
  if (inp) inp.required = show;
};

window.handleEditFieldTypeChange = function handleEditFieldTypeChange() {
  const sel = document.getElementById('editFieldType');
  const optCont = document.getElementById('editFieldOptionsContainer');
  if (!sel || !optCont) return;
  const show = sel.value === 'select';
  optCont.style.display = show ? 'block' : 'none';
  const inp = document.getElementById('editFieldOptions');
  if (inp) inp.required = show;
};

window.openAddFieldModal = function openAddFieldModal() {
  const formEl = document.getElementById('addFieldForm');
  if (formEl) formEl.reset();
  showModalEl('addFieldModal');
  window.handleFieldTypeChange();
};

window.closeAddFieldModal = function closeAddFieldModal() {
  hideModalEl('addFieldModal');
};

window.closeEditFieldModal = function closeEditFieldModal() {
  hideModalEl('editFieldModal');
};

window.editField = function editField(formId, fieldIndex) {
  const schema = window.currentFormSchema;
  if (!schema || !Array.isArray(schema.fields)) {
    showNotification('Esquema no cargado. Recarga la página.', 'error');
    return;
  }
  const idx = fieldIndex - 1;
  if (idx < 0 || idx >= schema.fields.length) {
    showNotification('Campo no encontrado.', 'error');
    return;
  }
  const f = schema.fields[idx];
  const form = document.getElementById('editFieldForm');
  if (form) {
    form.setAttribute('data-form-id', String(formId));
    form.setAttribute('data-field-id', String(fieldIndex));
  }

  const setVal = (id, v) => {
    const el = document.getElementById(id);
    if (el) el.value = v != null ? v : '';
  };

  setVal('editFieldLabel', f.field_label || f.fieldLabel || '');
  setVal('editFieldName', f.field_name || f.fieldName || '');
  setVal('editFieldType', f.field_type || f.fieldType || 'text');
  const opts = f.field_options != null ? String(f.field_options) : '';
  setVal('editFieldOptions', opts === 'null' ? '' : opts);
  const req = document.getElementById('editFieldRequired');
  if (req) req.checked = !!f.is_required;
  const ml = f.max_length != null ? f.max_length : f.maxLength;
  setVal('editFieldMaxLength', ml != null && ml !== '' ? String(ml) : '');

  showModalEl('editFieldModal');
  window.handleEditFieldTypeChange();
};

window.deleteField = async function deleteField(formId, fieldIndex) {
  if (!confirm('¿Eliminar este campo del esquema?')) return;
  const schema = window.currentFormSchema;
  if (!schema || !Array.isArray(schema.fields)) {
    showNotification('Esquema no disponible.', 'error');
    return;
  }
  const idx = fieldIndex - 1;
  if (idx < 0 || idx >= schema.fields.length) return;
  schema.fields.splice(idx, 1);
  if (typeof window.saveFormSchema === 'function') {
    try {
      await window.saveFormSchema();
    } catch (e) {
      showNotification(e.message || 'Error al guardar', 'error');
    }
  } else {
    showNotification('saveFormSchema no está disponible.', 'error');
  }
};

window.switchTab = switchTab;
window.toggleFormInfoAccordion = toggleFormInfoAccordion;
window.displayFormDetailsInPage = displayFormDetailsInPage;
window.renderFormFieldsList = renderFormFieldsList;
