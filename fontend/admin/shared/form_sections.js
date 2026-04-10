// Importar funciones necesarias
import { apiRequest, API_CONFIG } from '../../services/index.js';
import { showNotification } from '../../components/index.js';
import { getHeroiconSVG, loadHeroiconSVGFromFile, loadHeroiconSVGFromWeb, normalizeIconName } from '../../components/icons.js';

// Funciones para gestionar secciones de formularios

let currentFormSections = [];
let currentFormIdForSections = null;

// Lista de iconos comunes usando nombres de FontAwesome (se guardan así, pero se muestran como Heroicons)
const COMMON_ICONS = [
 'fa-user', 'fa-envelope', 'fa-phone', 'fa-map-marker-alt', 'fa-briefcase',
 'fa-share-alt', 'fa-shield-alt', 'fa-info-circle', 'fa-calendar', 'fa-clock',
 'fa-file-alt', 'fa-image', 'fa-link', 'fa-globe', 'fa-building',
 'fa-users', 'fa-user-tie', 'fa-graduation-cap', 'fa-heart', 'fa-star',
 'fa-home', 'fa-car', 'fa-plane', 'fa-shopping-cart', 'fa-credit-card',
 'fa-lock', 'fa-unlock', 'fa-key', 'fa-bell', 'fa-cog',
 'fa-check-circle', 'fa-times-circle', 'fa-question-circle', 'fa-exclamation-circle',
 'fa-arrow-right', 'fa-arrow-left', 'fa-arrow-up', 'fa-arrow-down',
 'fa-search', 'fa-filter', 'fa-sort', 'fa-list', 'fa-th',
 'fa-edit', 'fa-trash', 'fa-save', 'fa-download', 'fa-upload'
];

// Cargar secciones de un formulario
async function loadFormSections(formId) {
 try {
 currentFormIdForSections = formId;
 const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_SECTIONS(formId));
 if (response.success && response.sections) {
 currentFormSections = response.sections;
 // Hacer disponible globalmente para otras funciones
 if (typeof window !== 'undefined') {
   window.currentFormSections = currentFormSections;
 }
 displayFormSections(currentFormSections);
 } else {
 currentFormSections = [];
 if (typeof window !== 'undefined') {
   window.currentFormSections = [];
 }
 displayFormSections([]);
 }
 } catch (error) {
 console.error('Error cargando secciones:', error);
 currentFormSections = [];
 if (typeof window !== 'undefined') {
   window.currentFormSections = [];
 }
 displayFormSections([]);
 }
}

// Table-per-Form: secciones vienen de forms.schema; edición vía PUT .../schema
const TABLE_PER_FORM_READONLY = true;

// En view-form.html, si existe currentFormSchema + saveFormSchema, Secciones y Campos editan el schema
function useSchemaEditing() {
  return !!(typeof window !== 'undefined' && window.currentFormSchema && typeof window.saveFormSchema === 'function');
}

// Mostrar secciones en el modal
function displayFormSections(sections) {
 const sectionsContainer = document.getElementById('formDetailsSections');
 if (!sectionsContainer) return;

 if (!sections || sections.length === 0) {
   if (!useSchemaEditing()) {
     sectionsContainer.innerHTML = `<div><p>No hay secciones en el esquema.</p></div>`;
     return;
   }
   sectionsContainer.innerHTML = `<div class="mb-2 p-3 border rounded"><p>No hay secciones. Usa el botón <strong>Añadir Sección</strong> para crear una.</p></div>`;
   return;
 }

 if (TABLE_PER_FORM_READONLY && !useSchemaEditing()) {
   sectionsContainer.innerHTML = `
     <div class="mb-3 p-3 rounded border border-amber-200 bg-amber-50 text-amber-900 text-sm" style="margin-bottom:12px">
       <strong>Table-per-Form:</strong> secciones y campos están en el esquema JSON.
       Cambios estructurales: <code>PUT .../admin/forms/{id}/schema</code> (si no hay respuestas, se recrea la tabla).
     </div>
   ` + (sections || []).map((section) => {
     const iconSVG = getHeroiconSVG(section.section_icon) || '';
     return `<div class="mb-2 p-3 border rounded" style="margin-bottom:8px">
 <div style="display:flex;align-items:center;gap:8px">
 <span style="width:24px;height:24px;display:inline-flex;align-items:center;justify-content:center">${iconSVG}</span>
 <div>
 <h4 style="margin:0">${section.section_title || 'Sin título'}</h4>
 <p style="margin:0;font-size:0.875rem;color:#666">Orden: ${section.display_order}</p>
 </div></div></div>`;
   }).join('');
   return;
 }

 const list = sections || [];
 sectionsContainer.innerHTML = list.map((section, index) => {
   const iconSVG = getHeroiconSVG(section.section_icon) || '';
   return `
 <div>
 <div>
 <div>
 <div style="display: inline-flex; align-items: center; justify-content: center; width: 24px; height: 24px;">
 ${iconSVG}
 </div>
 <div>
 <h4>${section.section_title}</h4>
 <p>Orden: ${section.display_order}</p>
 </div>
 </div>
 <div>
 <div>
 <button 
 type="button"
 onclick="moveSectionUp(${currentFormIdForSections}, ${section.id}, ${index})"
 title="Mover arriba"
 ${index === 0 ? 'disabled' : ''}
 >
 ${getHeroiconSVG('chevron-up')}
 </button>
 <button 
 type="button"
 onclick="moveSectionDown(${currentFormIdForSections}, ${section.id}, ${index})"
 title="Mover abajo"
 ${index === sections.length - 1 ? 'disabled' : ''}
 >
 ${getHeroiconSVG('chevron-down')}
 </button>
 </div>
 <button 
 onclick="editSection(${currentFormIdForSections}, ${section.id})"
 title="Editar sección"
 >
 ${getHeroiconSVG('pencil')}
 </button>
 <button 
 onclick="deleteSection(${currentFormIdForSections}, ${section.id}, '${section.section_title}')"
 title="Eliminar sección"
 >
 ${getHeroiconSVG('trash')}
 </button>
 </div>
 </div>
 </div>
 `;
 }).join('');
 
 // Cargar iconos async si no están en inline
 sections.forEach((section, index) => {
   const iconSVG = getHeroiconSVG(section.section_icon);
   if (!iconSVG && section.section_icon) {
     // Normalizar el nombre para cargar desde la web
     const normalizedName = normalizeIconName(section.section_icon);
     if (normalizedName) {
       loadHeroiconSVGFromWeb(normalizedName).then(svg => {
       if (svg) {
         const sectionElement = sectionsContainer.children[index];
         if (sectionElement) {
           const iconContainer = sectionElement.querySelector('[style*="display: inline-flex"]');
           if (iconContainer) {
             iconContainer.innerHTML = svg;
           }
         }
       }
     }).catch(err => console.warn('Error cargando icono:', err));
     }
   }
 });
}

// Abrir modal para crear sección
async function openCreateSectionModal() {
 if (!currentFormIdForSections) {
 showNotification('Error: No hay formulario seleccionado', 'error');
 return;
 }
 document.getElementById('sectionModalTitle').textContent = 'Crear Sección';
 document.getElementById('sectionForm').reset();
 document.getElementById('sectionFormId').value = '';
 document.getElementById('sectionFormFormId').value = currentFormIdForSections;
 document.getElementById('sectionIcon').value = 'fa-info-circle';
 document.getElementById('sectionModal').style.display = 'block';
 await loadIconGrid();
 await updateIconPreview();
}

// Abrir modal para editar sección (sectionId puede ser numérico desde API o string desde schema)
async function editSection(formId, sectionId) {
 const section = currentFormSections.find(s => String(s.id) === String(sectionId)) ||
   (window.currentFormSchema && (() => {
     const idx = typeof sectionId === 'number' ? sectionId - 1 : window.currentFormSchema.sections.findIndex(s => s.id === sectionId);
     return idx >= 0 ? window.currentFormSchema.sections[idx] : null;
   })());
 if (!section) {
 showNotification('Error: Sección no encontrada', 'error');
 return;
 }
 
 document.getElementById('sectionModalTitle').textContent = 'Editar Sección';
 document.getElementById('sectionFormId').value = section.id;
 document.getElementById('sectionFormFormId').value = formId;
 document.getElementById('sectionTitle').value = section.section_title || '';
 document.getElementById('sectionIcon').value = section.section_icon || 'fa-folder';
 document.getElementById('sectionModal').style.display = 'block';
 await loadIconGrid();
 await updateIconPreview();
}

// Cerrar modal de sección
function closeSectionModal() {
 document.getElementById('sectionModal').style.display = 'none';
}

// Crear sección (table-per-form: actualizar currentFormSchema y PUT schema)
async function createSection(formId, title, icon) {
 if (useSchemaEditing() && window.currentFormSchema && window.saveFormSchema) {
   const schema = window.currentFormSchema;
   const nextId = 's' + (schema.sections.length + 1);
   schema.sections.push({
     id: nextId,
     section_title: title,
     section_icon: icon || 'fa-folder',
     display_order: schema.sections.length
   });
   try {
     await window.saveFormSchema();
     showNotification('Sección creada correctamente', 'success');
     closeSectionModal();
   } catch (e) {
     schema.sections.pop();
     showNotification('Error creando sección: ' + e.message, 'error');
   }
   return;
 }
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_SECTION_CREATE(formId), {
 method: 'POST',
 body: JSON.stringify({
 section_title: title,
 section_icon: icon
 })
 });
 
 if (response.success) {
 showNotification('Sección creada correctamente', 'success');
 closeSectionModal();
 
 const isInSeparatePage = window.location.pathname.includes('/view-form/');
 if (isInSeparatePage && window.loadFormDetails && typeof window.loadFormDetails === 'function') {
   await window.loadFormDetails();
 } else {
   await loadFormSections(formId);
 }
 }
 } catch (error) {
 console.error('Error creando sección:', error);
 showNotification('Error creando sección: ' + error.message, 'error');
 }
}

// Actualizar sección (table-per-form: actualizar currentFormSchema y PUT schema)
async function updateSection(formId, sectionId, title, icon) {
 if (useSchemaEditing() && window.currentFormSchema && window.saveFormSchema) {
   const schema = window.currentFormSchema;
   const idx = typeof sectionId === 'number' ? sectionId - 1 : schema.sections.findIndex(s => String(s.id) === String(sectionId));
   if (idx < 0 || idx >= schema.sections.length) {
     showNotification('Sección no encontrada', 'error');
     return;
   }
   schema.sections[idx].section_title = title;
   schema.sections[idx].section_icon = icon || 'fa-folder';
   try {
     await window.saveFormSchema();
     showNotification('Sección actualizada correctamente', 'success');
     closeSectionModal();
   } catch (e) {
     showNotification('Error actualizando sección: ' + e.message, 'error');
   }
   return;
 }
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_SECTION_UPDATE(formId, sectionId), {
 method: 'PUT',
 body: JSON.stringify({
 section_title: title,
 section_icon: icon
 })
 });
 
 if (response.success) {
 showNotification('Sección actualizada correctamente', 'success');
 closeSectionModal();
 
 const isInSeparatePage = window.location.pathname.includes('/view-form/');
 if (isInSeparatePage && window.loadFormDetails && typeof window.loadFormDetails === 'function') {
   await window.loadFormDetails();
 } else {
   await loadFormSections(formId);
 }
 }
 } catch (error) {
 console.error('Error actualizando sección:', error);
 showNotification('Error actualizando sección: ' + error.message, 'error');
 }
}

// Eliminar sección (table-per-form: quitar del schema y quitar section_id de campos que la usen)
async function deleteSection(formId, sectionId, sectionTitle) {
 if (!confirm(`¿Estás seguro de que deseas eliminar la sección "${sectionTitle}"?\n\nLos campos de esta sección quedarán sin sección asignada.`)) {
 return;
 }
 
 if (useSchemaEditing() && window.currentFormSchema && window.saveFormSchema) {
   const schema = window.currentFormSchema;
   const sectionIdStr = typeof sectionId === 'number' ? schema.sections[sectionId - 1]?.id : sectionId;
   const idx = schema.sections.findIndex(s => String(s.id) === String(sectionIdStr || sectionId));
   if (idx < 0) {
     showNotification('Sección no encontrada', 'error');
     return;
   }
   schema.sections.splice(idx, 1);
   schema.fields.forEach(f => { if (f.section_id === sectionIdStr) f.section_id = ''; });
   try {
     await window.saveFormSchema();
     showNotification('Sección eliminada correctamente', 'success');
   } catch (e) {
     showNotification('Error eliminando sección: ' + e.message, 'error');
   }
   return;
 }
 
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_SECTION_DELETE(formId, sectionId), {
 method: 'DELETE'
 });
 
 if (response.success) {
 showNotification('Sección eliminada correctamente', 'success');
 
 const isInSeparatePage = window.location.pathname.includes('/view-form/');
 if (isInSeparatePage && window.loadFormDetails && typeof window.loadFormDetails === 'function') {
   await window.loadFormDetails();
 } else {
   await loadFormSections(formId);
 }
 }
 } catch (error) {
 console.error('Error eliminando sección:', error);
 showNotification('Error eliminando sección: ' + error.message, 'error');
 }
}

// Reordenar sección
async function moveSectionUp(formId, sectionId, currentIndex) {
 if (currentIndex === 0) {
 showNotification('La sección ya está en la primera posición.', 'info');
 return;
 }
 await reorderSection(formId, sectionId, 'up');
}

async function moveSectionDown(formId, sectionId, currentIndex) {
 const sections = currentFormSections;
 if (!sections || currentIndex === sections.length - 1) {
 showNotification('La sección ya está en la última posición.', 'info');
 return;
 }
 await reorderSection(formId, sectionId, 'down');
}

async function reorderSection(formId, sectionId, direction) {
 if (useSchemaEditing() && window.currentFormSchema && window.saveFormSchema) {
   const schema = window.currentFormSchema;
   const idx = typeof sectionId === 'number' ? sectionId - 1 : schema.sections.findIndex(s => String(s.id) === String(sectionId));
   if (idx < 0 || (direction === 'up' && idx === 0) || (direction === 'down' && idx >= schema.sections.length - 1)) {
     if (direction === 'up') showNotification('La sección ya está en la primera posición.', 'info');
     else showNotification('La sección ya está en la última posición.', 'info');
     return;
   }
   const swap = direction === 'up' ? idx - 1 : idx + 1;
   [schema.sections[idx], schema.sections[swap]] = [schema.sections[swap], schema.sections[idx]];
   schema.sections.forEach((s, i) => { s.display_order = i; });
   try {
     await window.saveFormSchema();
     showNotification('Sección reordenada correctamente', 'success');
   } catch (e) {
     [schema.sections[idx], schema.sections[swap]] = [schema.sections[swap], schema.sections[idx]];
     schema.sections.forEach((s, i) => { s.display_order = i; });
     showNotification('Error reordenando: ' + e.message, 'error');
   }
   return;
 }
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_SECTION_REORDER(formId, sectionId), {
 method: 'POST',
 body: JSON.stringify({ direction: direction })
 });
 
 if (response.success) {
 showNotification('Sección reordenada correctamente', 'success');
 await loadFormSections(formId);
 }
 } catch (error) {
 console.error('Error reordenando sección:', error);
 showNotification('Error reordenando sección: ' + error.message, 'error');
 }
}

// Asignar campo a sección (table-per-form: actualizar field.section_id en schema y PUT)
async function assignFieldToSection(formId, fieldId, sectionId) {
 if (useSchemaEditing() && window.currentFormSchema && window.saveFormSchema) {
   const schema = window.currentFormSchema;
   const fieldIdx = typeof fieldId === 'number' ? fieldId - 1 : schema.fields.findIndex(f => String(f.id) === String(fieldId));
   if (fieldIdx < 0 || fieldIdx >= schema.fields.length) {
     showNotification('Campo no encontrado', 'error');
     return;
   }
   const secIdx = (sectionId === 0 || sectionId === null || sectionId === '') ? -1 : (typeof sectionId === 'number' ? sectionId - 1 : schema.sections.findIndex(s => String(s.id) === String(sectionId)));
   const sectionIdStr = secIdx >= 0 && schema.sections[secIdx] ? schema.sections[secIdx].id : '';
   schema.fields[fieldIdx].section_id = sectionIdStr;
   try {
     await window.saveFormSchema();
     showNotification('Campo asignado a sección correctamente', 'success');
   } catch (e) {
     showNotification('Error: ' + e.message, 'error');
   }
   return;
 }
 try {
   const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_FIELD_ASSIGN_SECTION(formId, fieldId), {
     method: 'POST',
     body: JSON.stringify({
       section_id: sectionId === 0 ? null : sectionId
     })
   });
   
   if (response.success) {
     showNotification('Campo asignado a sección correctamente', 'success');
     const isInSeparatePage = window.location.pathname.includes('/view-form/');
     if (isInSeparatePage && window.loadFormDetails && typeof window.loadFormDetails === 'function') {
       await window.loadFormDetails();
     } else if (typeof viewForm === 'function') {
       viewForm(formId);
     } else {
       window.location.reload();
     }
   } else {
     showNotification('Error: ' + (response.message || 'No se pudo asignar el campo a la sección'), 'error');
   }
 } catch (error) {
   console.error('[assignFieldToSection] Error asignando campo a sección:', error);
   showNotification('Error asignando campo a sección: ' + error.message, 'error');
 }
}

// Selector de iconos
async function openIconSelector() {
 document.getElementById('iconSelectorModal').style.display = 'block';
 await loadIconGrid();
}

function closeIconSelector() {
 document.getElementById('iconSelectorModal').style.display = 'none';
}

async function loadIconGrid() {
 const iconGrid = document.getElementById('iconGrid');
 if (!iconGrid) return;
 
 // Mostrar loading
 iconGrid.innerHTML = '<p style="text-align: center; padding: 20px;">Cargando iconos...</p>';
 
 // Crear botones con iconos cargados desde la web
 const buttonsHTML = await Promise.all(COMMON_ICONS.map(async (iconName) => {
   // Normalizar el nombre para mostrar el icono de Heroicons
   const normalizedName = normalizeIconName(iconName);
   // Primero intentar icono inline (más rápido) - getHeroiconSVG normaliza automáticamente
   let iconSVG = getHeroiconSVG(iconName);
   // Si no está inline, cargar desde la web usando Iconify API con el nombre normalizado
   if (!iconSVG && normalizedName) {
     iconSVG = await loadHeroiconSVGFromWeb(normalizedName);
   }
   // Guardar el nombre original (fa-*) pero mostrar el heroicon normalizado
   return `
 <button 
 type="button"
 data-icon="${iconName.replace(/"/g, '&quot;')}"
 title="${iconName.replace(/"/g, '&quot;')}"
 style="display: inline-flex; align-items: center; justify-content: center; width: 48px; height: 48px; padding: 8px; border: 1px solid #e5e7eb; border-radius: 4px; background: white; cursor: pointer; transition: all 0.2s;"
 onmouseover="this.style.backgroundColor='#f3f4f6'; this.style.borderColor='#9ca3af';"
 onmouseout="this.style.backgroundColor='white'; this.style.borderColor='#e5e7eb';"
 >
 ${iconSVG || '<span style="font-size: 12px;">?</span>'}
 </button>
 `;
 }));
 
 iconGrid.innerHTML = buttonsHTML.join('');
 
 // Agregar event listeners a los botones después de insertar el HTML
 const buttons = iconGrid.querySelectorAll('button[data-icon]');
 buttons.forEach(button => {
   button.addEventListener('click', async () => {
     const iconName = button.getAttribute('data-icon');
     if (iconName && window.selectIcon) {
       await window.selectIcon(iconName);
     }
   });
 });
 
 // Aplicar estilo de grid si no existe
 if (!iconGrid.style.display) {
   iconGrid.style.display = 'grid';
   iconGrid.style.gridTemplateColumns = 'repeat(auto-fill, minmax(48px, 1fr))';
   iconGrid.style.gap = '8px';
   iconGrid.style.padding = '16px';
 }
}

function filterIcons() {
 const searchTerm = document.getElementById('iconSearch').value.toLowerCase();
 const buttons = document.getElementById('iconGrid').querySelectorAll('button');
 
 buttons.forEach(button => {
 const icon = button.getAttribute('title');
 if (icon.toLowerCase().includes(searchTerm)) {
 button.style.display = 'inline-flex';
 } else {
 button.style.display = 'none';
 }
 });
}

async function selectIcon(icon) {
 document.getElementById('sectionIcon').value = icon;
 await updateIconPreview();
 closeIconSelector();
}

async function updateIconPreview() {
 const iconInput = document.getElementById('sectionIcon');
 const preview = document.getElementById('iconPreview');
 if (!iconInput || !preview) return;
 
 const iconValue = iconInput.value.trim();
 if (iconValue) {
   // Primero intentar icono inline (más rápido)
   let iconSVG = getHeroiconSVG(iconValue);
   // Si no está inline, cargar desde la web usando Iconify API
   if (!iconSVG) {
     iconSVG = await loadHeroiconSVGFromWeb(iconValue);
   }
   // Si tampoco funciona desde web, intentar desde archivo local como fallback
   if (!iconSVG) {
     iconSVG = await loadHeroiconSVGFromFile(iconValue);
   }
   preview.innerHTML = iconSVG || '';
 } else {
   preview.innerHTML = '';
 }
}

// Manejar submit del formulario de sección
// Exportar funciones a window para que estén disponibles en onclick del HTML
window.openCreateSectionModal = openCreateSectionModal;
window.closeSectionModal = closeSectionModal;
window.openIconSelector = openIconSelector;
window.closeIconSelector = closeIconSelector;
window.selectIcon = selectIcon;

document.addEventListener('DOMContentLoaded', function() {
 const sectionForm = document.getElementById('sectionForm');
 if (sectionForm) {
 sectionForm.addEventListener('submit', async (e) => {
 e.preventDefault();
 
 const formId = parseInt(document.getElementById('sectionFormFormId').value);
 const sectionId = document.getElementById('sectionFormId').value;
 const title = document.getElementById('sectionTitle').value.trim();
 const icon = document.getElementById('sectionIcon').value.trim();
 
 if (!title || !icon) {
 showNotification('Por favor completa todos los campos', 'error');
 return;
 }
 
 if (sectionId) {
 const id = /^\d+$/.test(sectionId) ? parseInt(sectionId, 10) : sectionId;
 await updateSection(formId, id, title, icon);
 } else {
 await createSection(formId, title, icon);
 }
 });
 }
 
// Inicializar preview de icono
const iconInput = document.getElementById('sectionIcon');
if (iconInput) {
iconInput.addEventListener('input', () => {
 updateIconPreview().catch(err => console.error('Error actualizando preview de icono:', err));
});
}
});

// Exportar funciones para uso en otras páginas
export { 
  loadFormSections, 
  openCreateSectionModal, 
  closeSectionModal, 
  editSection, 
  deleteSection, 
  moveSectionUp, 
  moveSectionDown,
  openIconSelector,
  closeIconSelector,
  filterIcons,
  updateIconPreview,
  assignFieldToSection
};

// Exportar funciones a window para uso en onclick del HTML
window.assignFieldToSection = assignFieldToSection;
window.openCreateSectionModal = openCreateSectionModal;
window.closeSectionModal = closeSectionModal;
window.editSection = editSection;
window.deleteSection = deleteSection;
window.moveSectionUp = moveSectionUp;
window.moveSectionDown = moveSectionDown;
window.openIconSelector = openIconSelector;
window.closeIconSelector = closeIconSelector;
window.selectIcon = selectIcon;
window.filterIcons = filterIcons;
window.updateIconPreview = updateIconPreview;

