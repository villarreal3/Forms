// Importar funciones necesarias
import { apiRequest, API_CONFIG } from '../../services/index.js';
import { showNotification } from '../../components/index.js';

// Gestión de personalización de formularios

let currentFormId = null;
let currentFormName = '';

// Función para establecer el ID y nombre del formulario actual
function setCurrentForm(formId, formName = '') {
 currentFormId = formId;
 currentFormName = formName;
}

// Función para obtener el formId (del modal, del contenedor de la página, o de currentFormId)
function getCurrentFormId() {
 // Primero intentar obtener del modal (si existe)
 const formDetailsModal = document.getElementById('formDetailsModal');
 if (formDetailsModal) {
   const formId = formDetailsModal.getAttribute('data-form-id');
   if (formId) {
     return parseInt(formId);
   }
 }
 
 // Si no hay modal, intentar obtener del contenedor de la página (view-form.html)
 const formDetailsContent = document.getElementById('formDetailsContent');
 if (formDetailsContent) {
   const formId = formDetailsContent.getAttribute('data-form-id');
   if (formId) {
     return parseInt(formId);
   }
 }
 
 // Si no hay ninguno, intentar obtener de la URL (para view-form.html)
 const urlParams = new URLSearchParams(window.location.search);
 const formIdFromUrl = urlParams.get('id');
 if (formIdFromUrl) {
   const parsedId = parseInt(formIdFromUrl);
   if (!isNaN(parsedId)) {
     return parsedId;
   }
 }
 
 // Finalmente, usar currentFormId si está disponible
 return currentFormId;
}

// Abrir modal de personalización (mantenida para compatibilidad, pero ahora se usa desde el tab)
async function editFormCustomization(formId, formName) {
 currentFormId = formId;
 currentFormName = formName;
 
 // Si el modal de detalles está disponible, abrirlo en el tab de personalizar
 const formDetailsModal = document.getElementById('formDetailsModal');
 if (formDetailsModal) {
 // Cargar el formulario y abrir en el tab de personalizar
 if (typeof viewForm === 'function') {
 await viewForm(formId);
 if (typeof switchTab === 'function') {
 switchTab('personalize');
 }
 }
 }
}

// Cerrar modal (ahora solo limpia los datos, no cierra el modal)
function closeCustomizationModal() {
 // Ya no cerramos el modal, solo limpiamos los datos si es necesario
 const formImageInput = document.getElementById('formImageInput');
 if (formImageInput) formImageInput.value = '';
 const formImageInputMobile = document.getElementById('formImageInputMobile');
 if (formImageInputMobile) formImageInputMobile.value = '';
 
 const imagePreviewPc = document.getElementById('imagePreviewPc');
 if (imagePreviewPc) imagePreviewPc.style.display = 'none';
 const imagePreviewMobile = document.getElementById('imagePreviewMobile');
 if (imagePreviewMobile) imagePreviewMobile.style.display = 'none';
 
 const currentImagePc = document.getElementById('currentImagePc');
 if (currentImagePc) currentImagePc.style.display = 'none';
 const currentImageMobile = document.getElementById('currentImageMobile');
 if (currentImageMobile) currentImageMobile.style.display = 'none';
}

// Cargar personalización actual
async function loadCustomization(formId) {
 console.log('[loadCustomization] Iniciando carga de personalización para formId:', formId);
 
 try {
   const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_CUSTOMIZATION(formId));
   console.log('[loadCustomization] Respuesta de la API:', response);
   
   if (response.success && response.customization) {
     const custom = response.customization;
     console.log('[loadCustomization] Datos de personalización recibidos:', custom);
     
     // Función helper para actualizar un input de color de forma segura
     const updateColorInput = (pickerId, textId, value, defaultValue) => {
       const picker = document.getElementById(pickerId);
       const text = document.getElementById(textId);
       const finalValue = value || defaultValue;
       
       if (picker) {
         picker.value = finalValue.toUpperCase();
         console.log(`[loadCustomization] Actualizado ${pickerId}:`, finalValue);
       } else {
         console.warn(`[loadCustomization] Elemento ${pickerId} no encontrado`);
       }
       
       if (text) {
         text.value = finalValue.toUpperCase();
         console.log(`[loadCustomization] Actualizado ${textId}:`, finalValue);
       } else {
         console.warn(`[loadCustomization] Elemento ${textId} no encontrado`);
       }
     };
     
     // Cargar colores
     updateColorInput('primaryColor', 'primaryColorText', custom.primary_color, '#3B82F6');
     updateColorInput('secondaryColor', 'secondaryColorText', custom.secondary_color, '#1E40AF');
     updateColorInput('backgroundColor', 'backgroundColorText', custom.background_color, '#FFFFFF');
     updateColorInput('textColor', 'textColorText', custom.text_color, '#1F2937');
     updateColorInput('titleColor', 'titleColorText', custom.title_color, '#FFFFFF');
     
     // Cargar personalización de contenedores
     const formContainerColor = custom.form_container_color || '#FFFFFF';
     const formContainerOpacity = custom.form_container_opacity !== undefined ? custom.form_container_opacity : 1.0;
     const descriptionContainerColor = custom.description_container_color || '#FFFFFF';
     const descriptionContainerOpacity = custom.description_container_opacity !== undefined ? custom.description_container_opacity : 1.0;
     const formMetaBackgroundStartColor = custom.form_meta_background_start || custom.form_meta_background || '#3B82F6';
     const formMetaBackgroundEndColor = custom.form_meta_background_end || custom.secondary_color || '#1E40AF';
     const formMetaBackgroundOpacity = custom.form_meta_background_opacity !== undefined ? custom.form_meta_background_opacity : 1.0;
     const formMetaTextColor = custom.form_meta_text_color || '#FFFFFF';
     
     console.log('[loadCustomization] Valores de contenedores:', {
       formContainerColor,
       formContainerOpacity,
       descriptionContainerColor,
      descriptionContainerOpacity,
      formMetaBackgroundStartColor,
      formMetaBackgroundEndColor,
      formMetaBackgroundOpacity,
      formMetaTextColor
     });
     
     updateColorInput('formContainerColor', 'formContainerColorText', formContainerColor, '#FFFFFF');
     
     const formContainerOpacityEl = document.getElementById('formContainerOpacity');
     const formContainerOpacityTextEl = document.getElementById('formContainerOpacityText');
     if (formContainerOpacityEl && formContainerOpacityTextEl) {
       const opacityValue = Math.round(formContainerOpacity * 100);
       formContainerOpacityEl.value = opacityValue;
       formContainerOpacityTextEl.value = opacityValue + '%';
       console.log('[loadCustomization] Actualizado formContainerOpacity:', opacityValue + '%');
     }
     
     updateColorInput('descriptionContainerColor', 'descriptionContainerColorText', descriptionContainerColor, '#FFFFFF');
     
     const descriptionContainerOpacityEl = document.getElementById('descriptionContainerOpacity');
     const descriptionContainerOpacityTextEl = document.getElementById('descriptionContainerOpacityText');
     if (descriptionContainerOpacityEl && descriptionContainerOpacityTextEl) {
       const opacityValue = Math.round(descriptionContainerOpacity * 100);
       descriptionContainerOpacityEl.value = opacityValue;
       descriptionContainerOpacityTextEl.value = opacityValue + '%';
       console.log('[loadCustomization] Actualizado descriptionContainerOpacity:', opacityValue + '%');
     }

     updateColorInput('formMetaBackgroundStartColor', 'formMetaBackgroundStartColorText', formMetaBackgroundStartColor, '#3B82F6');
     updateColorInput('formMetaBackgroundEndColor', 'formMetaBackgroundEndColorText', formMetaBackgroundEndColor, '#1E40AF');
     updateColorInput('formMetaTextColor', 'formMetaTextColorText', formMetaTextColor, '#FFFFFF');

     const formMetaBackgroundOpacityEl = document.getElementById('formMetaBackgroundOpacity');
     const formMetaBackgroundOpacityTextEl = document.getElementById('formMetaBackgroundOpacityText');
     if (formMetaBackgroundOpacityEl && formMetaBackgroundOpacityTextEl) {
       const opacityValue = Math.round(formMetaBackgroundOpacity * 100);
       formMetaBackgroundOpacityEl.value = opacityValue;
       formMetaBackgroundOpacityTextEl.value = opacityValue + '%';
       console.log('[loadCustomization] Actualizado formMetaBackgroundOpacity:', opacityValue + '%');
     }
     
     // Cargar imagen PC si existe
     if (custom.logo_url) {
       console.log('[loadCustomization] Cargando imagen PC:', custom.logo_url);
       const imageUrl = getImageFullURL(custom.logo_url);
       console.log('[loadCustomization] URL completa imagen PC:', imageUrl);
       
       const currentImgPc = document.getElementById('currentImgPc');
       const currentImagePc = document.getElementById('currentImagePc');
       if (currentImgPc && currentImagePc) {
         currentImgPc.src = imageUrl;
         currentImgPc.onerror = () => {
           console.error('[loadCustomization] Error cargando imagen PC:', imageUrl);
           currentImagePc.style.display = 'none';
         };
         currentImgPc.onload = () => {
           console.log('[loadCustomization] Imagen PC cargada correctamente');
           currentImagePc.style.display = 'block';
         };
         // Forzar carga
         currentImgPc.src = imageUrl;
       } else {
         console.warn('[loadCustomization] Elementos de imagen PC no encontrados');
       }
       
       // También actualizar preview logo si existe
       const previewLogoImg = document.getElementById('previewLogoImg');
       const previewLogo = document.getElementById('previewLogo');
       if (previewLogoImg && previewLogo) {
         previewLogoImg.src = imageUrl;
         previewLogo.style.display = 'block';
       }
     } else {
       console.log('[loadCustomization] No hay imagen PC guardada');
       const currentImagePc = document.getElementById('currentImagePc');
       if (currentImagePc) {
         currentImagePc.style.display = 'none';
       }
     }
     
     // Cargar imagen móvil si existe
     if (custom.logo_url_mobile) {
       console.log('[loadCustomization] Cargando imagen móvil:', custom.logo_url_mobile);
       const imageUrlMobile = getImageFullURL(custom.logo_url_mobile);
       console.log('[loadCustomization] URL completa imagen móvil:', imageUrlMobile);
       
       const currentImgMobile = document.getElementById('currentImgMobile');
       const currentImageMobile = document.getElementById('currentImageMobile');
       if (currentImgMobile && currentImageMobile) {
         currentImgMobile.src = imageUrlMobile;
         currentImgMobile.onerror = () => {
           console.error('[loadCustomization] Error cargando imagen móvil:', imageUrlMobile);
           currentImageMobile.style.display = 'none';
         };
         currentImgMobile.onload = () => {
           console.log('[loadCustomization] Imagen móvil cargada correctamente');
           currentImageMobile.style.display = 'block';
         };
         // Forzar carga
         currentImgMobile.src = imageUrlMobile;
       } else {
         console.warn('[loadCustomization] Elementos de imagen móvil no encontrados');
       }
     } else {
       console.log('[loadCustomization] No hay imagen móvil guardada');
       const currentImageMobile = document.getElementById('currentImageMobile');
       if (currentImageMobile) {
         currentImageMobile.style.display = 'none';
       }
     }
     
     // Actualizar vista previa después de cargar todos los datos
     console.log('[loadCustomization] Actualizando vista previa...');
     updatePreview();
     console.log('[loadCustomization] Carga completada exitosamente');
   } else {
     console.error('[loadCustomization] Error: No se recibió personalización en la respuesta', response);
     if (!response.success) {
       showNotification('No se pudo cargar la personalización: ' + (response.message || 'Error desconocido'), 'error');
     }
   }
 } catch (error) {
   console.error('[loadCustomization] Error cargando personalización:', error);
   showNotification('Error cargando personalización: ' + error.message, 'error');
 }
}

// Configurar eventos de inputs de color
function setupColorInputs() {
 const colorInputs = [
 { picker: 'primaryColor', text: 'primaryColorText' },
 { picker: 'secondaryColor', text: 'secondaryColorText' },
 { picker: 'backgroundColor', text: 'backgroundColorText' },
 { picker: 'textColor', text: 'textColorText' },
 { picker: 'titleColor', text: 'titleColorText' },
 { picker: 'formContainerColor', text: 'formContainerColorText' },
 { picker: 'descriptionContainerColor', text: 'descriptionContainerColorText' },
 { picker: 'formMetaBackgroundStartColor', text: 'formMetaBackgroundStartColorText' },
 { picker: 'formMetaBackgroundEndColor', text: 'formMetaBackgroundEndColorText' },
 { picker: 'formMetaTextColor', text: 'formMetaTextColorText' }
 ];
 
 colorInputs.forEach(({ picker, text }) => {
 const pickerEl = document.getElementById(picker);
 const textEl = document.getElementById(text);
 
 if (!pickerEl || !textEl) return; // Saltar si no existen
 
 pickerEl.addEventListener('input', (e) => {
   textEl.value = e.target.value.toUpperCase();
   updatePreview();
 });
 
 textEl.addEventListener('input', (e) => {
   const value = e.target.value;
   if (/^#[0-9A-Fa-f]{6}$/i.test(value)) {
     pickerEl.value = value.toUpperCase();
     updatePreview();
   }
 });
 });
 
 // Configurar eventos de opacidad
 const opacityInputs = [
 { slider: 'formContainerOpacity', text: 'formContainerOpacityText' },
 { slider: 'descriptionContainerOpacity', text: 'descriptionContainerOpacityText' },
 { slider: 'formMetaBackgroundOpacity', text: 'formMetaBackgroundOpacityText' }
 ];
 
 opacityInputs.forEach(({ slider, text }) => {
 const sliderEl = document.getElementById(slider);
 const textEl = document.getElementById(text);
 
 if (!sliderEl || !textEl) return;
 
 sliderEl.addEventListener('input', (e) => {
   textEl.value = e.target.value + '%';
   updatePreview();
 });
 });
 
 // También actualizar preview cuando cambian los colores de contenedores
 const containerColorInputs = [
   { picker: 'formContainerColor', text: 'formContainerColorText' },
   { picker: 'descriptionContainerColor', text: 'descriptionContainerColorText' }
 ];
 
 containerColorInputs.forEach(({ picker, text }) => {
   const pickerEl = document.getElementById(picker);
   const textEl = document.getElementById(text);
   
   if (!pickerEl || !textEl) return;
   
   pickerEl.addEventListener('input', (e) => {
     textEl.value = e.target.value.toUpperCase();
     updatePreview();
   });
   
   textEl.addEventListener('input', (e) => {
     const value = e.target.value;
     if (/^#[0-9A-Fa-f]{6}$/i.test(value)) {
       pickerEl.value = value.toUpperCase();
       updatePreview();
     }
   });
 });
}

// Actualizar preview
function updatePreview() {
 const preview = document.getElementById('colorPreview');
 if (!preview) return;
 
 // Obtener valores de los inputs
 const primaryColor = document.getElementById('primaryColor')?.value || '#3B82F6';
 const secondaryColor = document.getElementById('secondaryColor')?.value || '#1E40AF';
 const backgroundColor = document.getElementById('backgroundColor')?.value || '#FFFFFF';
 const textColor = document.getElementById('textColor')?.value || '#1F2937';
 const titleColor = document.getElementById('titleColor')?.value || '#FFFFFF';
 
 // Obtener valores de contenedores
 const formContainerColor = document.getElementById('formContainerColor')?.value || '#FFFFFF';
 const formContainerOpacity = document.getElementById('formContainerOpacity')?.value || 100;
 const descriptionContainerColor = document.getElementById('descriptionContainerColor')?.value || '#FFFFFF';
 const descriptionContainerOpacity = document.getElementById('descriptionContainerOpacity')?.value || 100;
const formMetaBackgroundStartColor = document.getElementById('formMetaBackgroundStartColor')?.value || '#3B82F6';
const formMetaBackgroundEndColor = document.getElementById('formMetaBackgroundEndColor')?.value || '#1E40AF';
const formMetaBackgroundOpacity = document.getElementById('formMetaBackgroundOpacity')?.value || 100;
const formMetaTextColor = document.getElementById('formMetaTextColor')?.value || '#FFFFFF';
 
 // Aplicar colores al preview
 const previewDiv = preview.querySelector('div');
 if (previewDiv) {
   // Fondo del preview
   preview.style.backgroundColor = backgroundColor;
   
   // Contenedor interno (simula el contenedor del formulario)
   const containerOpacity = parseFloat(formContainerOpacity) / 100;
   previewDiv.style.backgroundColor = formContainerColor;
   previewDiv.style.opacity = containerOpacity;
   previewDiv.style.border = `2px solid ${primaryColor}`;
   
   // Título
   const title = previewDiv.querySelector('h3');
   if (title) {
     title.style.color = titleColor;
   }
   
   // Descripción (simula el contenedor de descripción)
   const description = previewDiv.querySelector('p');
   if (description) {
     const descOpacity = parseFloat(descriptionContainerOpacity) / 100;
    const metaOpacity = parseFloat(formMetaBackgroundOpacity) / 100;
    const hexStart = formMetaBackgroundStartColor.replace('#', '');
    const r1 = parseInt(hexStart.substring(0, 2), 16);
    const g1 = parseInt(hexStart.substring(2, 4), 16);
    const b1 = parseInt(hexStart.substring(4, 6), 16);
    const hexEnd = formMetaBackgroundEndColor.replace('#', '');
    const r2 = parseInt(hexEnd.substring(0, 2), 16);
    const g2 = parseInt(hexEnd.substring(2, 4), 16);
    const b2 = parseInt(hexEnd.substring(4, 6), 16);
    description.style.color = formMetaTextColor;
    description.style.background = `linear-gradient(to right, rgba(${r1}, ${g1}, ${b1}, ${metaOpacity}), rgba(${r2}, ${g2}, ${b2}, ${metaOpacity}))`;
     description.style.padding = 'var(--spacing-sm)';
     description.style.borderRadius = 'var(--border-radius-sm)';
    description.style.opacity = descOpacity;
   }
   
   // Botón
   const button = previewDiv.querySelector('button');
   if (button) {
     button.style.backgroundColor = primaryColor;
     button.style.color = '#FFFFFF';
     button.style.border = `2px solid ${secondaryColor}`;
   }
 }
}

// Preview de imagen antes de subir
function previewImage(input, type) {
 if (input.files && input.files[0]) {
 const reader = new FileReader();
 reader.onload = (e) => {
 if (type === 'mobile') {
 const previewImgMobile = document.getElementById('previewImgMobile');
 const imagePreviewMobile = document.getElementById('imagePreviewMobile');
 if (previewImgMobile && imagePreviewMobile) {
 previewImgMobile.src = e.target.result;
 imagePreviewMobile.style.display = 'block';
 }
 // Habilitar botón de subir
 const uploadMobileBtn = document.getElementById('uploadMobileBtn');
 if (uploadMobileBtn) uploadMobileBtn.disabled = false;
 } else {
 // type === 'desktop'
 const previewImgPc = document.getElementById('previewImgPc');
 const imagePreviewPc = document.getElementById('imagePreviewPc');
 if (previewImgPc && imagePreviewPc) {
 previewImgPc.src = e.target.result;
 imagePreviewPc.style.display = 'block';
 }
 // También actualizar preview logo si existe
 const previewLogoImg = document.getElementById('previewLogoImg');
 const previewLogo = document.getElementById('previewLogo');
 if (previewLogoImg && previewLogo) {
 previewLogoImg.src = e.target.result;
 previewLogo.style.display = 'block';
 }
 // Habilitar botón de subir
 const uploadPcBtn = document.getElementById('uploadPcBtn');
 if (uploadPcBtn) uploadPcBtn.disabled = false;
 }
 };
 reader.readAsDataURL(input.files[0]);
 }
}

// Subir imagen
async function uploadFormImage(type = 'desktop') {
 const inputId = type === 'mobile' ? 'formImageInputMobile' : 'formImageInput'; // type puede ser 'desktop' o 'mobile'
 const input = document.getElementById(inputId);
 if (!input || !input.files || !input.files[0]) {
 showNotification('Por favor selecciona una imagen', 'error');
 return;
 }
 
 const formId = getCurrentFormId();
 if (!formId) {
 showNotification('Error: ID de formulario no válido', 'error');
 return;
 }
 
 console.log(`[FRONTEND] Iniciando subida de imagen - Tipo: ${type}, FormID: ${formId}`);
 
 const formData = new FormData();
 formData.append('image', input.files[0]);
 
 try {
 const token = Auth.getToken();
 const baseURL = API_CONFIG.FORMS_API_BASE_URL;
 const url = `${baseURL}${API_CONFIG.ENDPOINTS.FORM_IMAGE_UPLOAD(formId)}?type=${type}`;
 
 console.log(`[FRONTEND] URL de petición: ${url}`);
 
 // NO incluir Content-Type - el navegador lo establece automáticamente con FormData
 const response = await fetch(url, {
 method: 'POST',
 headers: {
 'Authorization': `Bearer ${token}`
 // NO incluir 'Content-Type' aquí - el navegador lo establece automáticamente
 },
 body: formData
 });
 
 console.log(`[FRONTEND] Respuesta recibida - Status: ${response.status}, OK: ${response.ok}`);
 
 // Verificar si la respuesta es JSON antes de parsear
 let data;
 const contentType = response.headers.get('content-type');
 if (contentType && contentType.includes('application/json')) {
 data = await response.json();
 } else {
 const text = await response.text();
 console.error(`[FRONTEND] Error - Respuesta no es JSON. Status: ${response.status}, Text: ${text}`);
 throw new Error(`Error del servidor: ${response.status} ${response.statusText}. ${text}`);
 }
 
 if (!response.ok) {
 console.error(`[FRONTEND] Error - Respuesta no OK. Status: ${response.status}, Error: ${data.message || data.error || 'Error subiendo imagen'}`);
 throw new Error(data.message || data.error || 'Error subiendo imagen');
 }
 
 if (data.success) {
 const savedUrl = type === 'mobile' ? (data.logo_url_mobile || data.image_url) : (data.logo_url || data.image_url);
 console.log(`[FRONTEND] Resultado final - Success: ${data.success}, Tipo: ${type}, URL guardada: ${savedUrl || 'N/A'}`);
 showNotification('Imagen subida correctamente', 'success');
 
 // Actualizar preview según el tipo
 const imageUrl = getImageFullURL(data.image_url);
 if (type === 'mobile') {
 const currentImgMobile = document.getElementById('currentImgMobile');
 const currentImageMobile = document.getElementById('currentImageMobile');
 const imagePreviewMobile = document.getElementById('imagePreviewMobile');
 if (currentImgMobile && currentImageMobile) {
 currentImgMobile.src = imageUrl;
 currentImageMobile.style.display = 'block';
 }
 if (imagePreviewMobile) imagePreviewMobile.style.display = 'none';
 } else {
 const currentImgPc = document.getElementById('currentImgPc');
 const currentImagePc = document.getElementById('currentImagePc');
 const imagePreviewPc = document.getElementById('imagePreviewPc');
 if (currentImgPc && currentImagePc) {
 currentImgPc.src = imageUrl;
 currentImagePc.style.display = 'block';
 }
 if (imagePreviewPc) imagePreviewPc.style.display = 'none';
 // También actualizar preview logo si existe
 const previewLogoImg = document.getElementById('previewLogoImg');
 const previewLogo = document.getElementById('previewLogo');
 if (previewLogoImg && previewLogo) {
 previewLogoImg.src = imageUrl;
 previewLogo.style.display = 'block';
 }
 }
 input.value = '';
 
 // Deshabilitar botón de subir
 const btnId = type === 'mobile' ? 'uploadMobileBtn' : 'uploadPcBtn';
 const btn = document.getElementById(btnId);
 if (btn) btn.disabled = true;
 } else {
 console.error(`[FRONTEND] Error - data.success es false. Tipo: ${type}, Data:`, data);
 }
 } catch (error) {
 console.error(`[FRONTEND] Error subiendo imagen - Tipo: ${type}, FormID: ${formId}, Error:`, error);
 // Mejorar mensajes de error para problemas de conexión
 let errorMessage = error.message;
 if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
 errorMessage = 'No se pudo conectar al servidor. Verifica que el servicio esté corriendo en ' + API_CONFIG.FORMS_API_BASE_URL;
 }
 showNotification('Error subiendo imagen: ' + errorMessage, 'error');
 }
}

// Eliminar imagen
async function deleteFormImage(type = 'desktop') {
 const formId = getCurrentFormId();
 if (!formId) {
 showNotification('Error: ID de formulario no válido', 'error');
 return;
 }
 
 const typeLabel = type === 'mobile' ? 'móvil' : 'Desktop';
 if (!confirm(`¿Estás seguro de que deseas eliminar la imagen para ${typeLabel}?`)) {
 return;
 }
 
 try {
 const url = `${API_CONFIG.ENDPOINTS.FORM_IMAGE_DELETE(formId)}?type=${type}`;
 const response = await apiRequest(url, {
 method: 'DELETE'
 });
 
 if (response.success) {
 showNotification('Imagen eliminada correctamente', 'success');
 
 // Ocultar imagen actual según el tipo
 if (type === 'mobile') {
 const currentImageMobile = document.getElementById('currentImageMobile');
 if (currentImageMobile) {
 currentImageMobile.style.display = 'none';
 }
 const formImageInputMobile = document.getElementById('formImageInputMobile');
 if (formImageInputMobile) formImageInputMobile.value = '';
 } else {
 const currentImagePc = document.getElementById('currentImagePc');
 if (currentImagePc) {
 currentImagePc.style.display = 'none';
 }
 const formImageInput = document.getElementById('formImageInput');
 if (formImageInput) formImageInput.value = '';
 // También ocultar preview logo si existe
 const previewLogo = document.getElementById('previewLogo');
 if (previewLogo) previewLogo.style.display = 'none';
 }
 }
 } catch (error) {
 console.error('Error eliminando imagen:', error);
 showNotification('Error eliminando imagen: ' + error.message, 'error');
 }
}

// Obtener URL completa de imagen
function getImageFullURL(relativeUrl) {
 if (!relativeUrl) return '';
 if (relativeUrl.startsWith('http')) return relativeUrl;
 return `${API_CONFIG.FORMS_API_BASE_URL}${relativeUrl}`;
}

// Guardar personalización
const customizationForm = document.getElementById('customizationForm');
if (customizationForm) {
 customizationForm.addEventListener('submit', async (e) => {
 e.preventDefault();
 
 // Obtener el formId usando la función getCurrentFormId que busca en múltiples lugares
 const formId = getCurrentFormId();
 
 if (!formId || isNaN(formId)) {
   showNotification('Error: ID de formulario no válido', 'error');
   console.error('Error: No se pudo obtener el ID del formulario. formId:', formId);
   return;
 }
 
 const formIdNum = parseInt(formId);
 
 // Obtener valores de opacidad (convertir de porcentaje a decimal)
 const formContainerOpacityEl = document.getElementById('formContainerOpacity');
 const descriptionContainerOpacityEl = document.getElementById('descriptionContainerOpacity');
const formMetaBackgroundOpacityEl = document.getElementById('formMetaBackgroundOpacity');
 
 const formContainerOpacity = formContainerOpacityEl ? parseFloat(formContainerOpacityEl.value) / 100 : 1.0;
 const descriptionContainerOpacity = descriptionContainerOpacityEl ? parseFloat(descriptionContainerOpacityEl.value) / 100 : 1.0;
const formMetaBackgroundOpacity = formMetaBackgroundOpacityEl ? parseFloat(formMetaBackgroundOpacityEl.value) / 100 : 1.0;
 
 // Obtener valores de colores de contenedores
 const formContainerColorEl = document.getElementById('formContainerColor');
 const descriptionContainerColorEl = document.getElementById('descriptionContainerColor');
const formMetaBackgroundStartColorEl = document.getElementById('formMetaBackgroundStartColor');
const formMetaBackgroundEndColorEl = document.getElementById('formMetaBackgroundEndColor');
const formMetaTextColorEl = document.getElementById('formMetaTextColor');
 
 // Verificar que el elemento existe Y tiene un valor válido (no vacío)
 const formContainerColor = (formContainerColorEl?.value && formContainerColorEl.value.trim() !== '') 
 ? formContainerColorEl.value 
 : '#FFFFFF';
 
 const descriptionContainerColor = (descriptionContainerColorEl?.value && descriptionContainerColorEl.value.trim() !== '') 
 ? descriptionContainerColorEl.value 
 : '#FFFFFF';

const formMetaBackgroundStartColor = (formMetaBackgroundStartColorEl?.value && formMetaBackgroundStartColorEl.value.trim() !== '')
? formMetaBackgroundStartColorEl.value
: '#3B82F6';
const formMetaBackgroundEndColor = (formMetaBackgroundEndColorEl?.value && formMetaBackgroundEndColorEl.value.trim() !== '')
? formMetaBackgroundEndColorEl.value
: '#1E40AF';

const formMetaTextColor = (formMetaTextColorEl?.value && formMetaTextColorEl.value.trim() !== '')
? formMetaTextColorEl.value
: '#FFFFFF';
 
 console.log('Valores de contenedores a enviar:', {
 formContainerColor,
 formContainerOpacity,
 descriptionContainerColor,
 descriptionContainerOpacity,
 formMetaBackgroundStartColor,
 formMetaBackgroundEndColor,
 formMetaBackgroundOpacity,
 formMetaTextColor
 });
 
 const titleColorEl = document.getElementById('titleColor');
 const titleColor = (titleColorEl?.value && titleColorEl.value.trim() !== '') 
 ? titleColorEl.value 
 : '#FFFFFF';
 
 console.log('titleColor obtenido del input:', {
 element: titleColorEl,
 value: titleColorEl?.value,
 finalValue: titleColor
 });
 
 const customization = {
 primary_color: document.getElementById('primaryColor').value,
 secondary_color: document.getElementById('secondaryColor').value,
 background_color: document.getElementById('backgroundColor').value,
 text_color: document.getElementById('textColor').value,
 title_color: titleColor,
 font_family: 'Arial, sans-serif',
 button_style: 'rounded',
 form_container_color: formContainerColor,
 form_container_opacity: formContainerOpacity,
 description_container_color: descriptionContainerColor,
 description_container_opacity: descriptionContainerOpacity,
 form_meta_background: formMetaBackgroundStartColor,
 form_meta_background_start: formMetaBackgroundStartColor,
 form_meta_background_end: formMetaBackgroundEndColor,
 form_meta_background_opacity: formMetaBackgroundOpacity,
 form_meta_text_color: formMetaTextColor
 };
 
 console.log('Objeto customization completo:', customization);
 console.log('JSON a enviar:', JSON.stringify(customization, null, 2));
 
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.FORM_CUSTOMIZATION(formIdNum), {
 method: 'PUT',
 body: JSON.stringify(customization)
 });
 
 console.log('Respuesta del servidor:', response);
 
 if (response.success) {
   showNotification('Personalización guardada correctamente', 'success');
   
   // Recargar los datos desde la base de datos para asegurar que los inputs muestren los valores correctos
   await loadCustomization(formIdNum);
   
   // Actualizar la vista previa
   updatePreview();
 } else {
   showNotification('Error guardando personalización: ' + (response.message || 'Error desconocido'), 'error');
 }
} catch (error) {
 console.error('Error guardando personalización:', error);
 showNotification('Error guardando personalización: ' + error.message, 'error');
}
});
}

// Función para inicializar personalización
async function initializeCustomization(formId) {
  console.log('[initializeCustomization] Iniciando con formId:', formId);
  
  if (!formId || isNaN(formId)) {
    console.error('[initializeCustomization] Error: formId inválido:', formId);
    return;
  }
  
  setCurrentForm(formId);
  
  // Esperar un momento para asegurar que el DOM esté listo
  // Especialmente importante cuando se cambia de tab
  await new Promise(resolve => setTimeout(resolve, 100));
  
  // Verificar que los elementos existan antes de cargar
  const primaryColorEl = document.getElementById('primaryColor');
  if (!primaryColorEl) {
    console.warn('[initializeCustomization] Los elementos del formulario aún no están disponibles, esperando...');
    await new Promise(resolve => setTimeout(resolve, 200));
  }
  
  await loadCustomization(formId);
  
  // Configurar event listeners después de cargar los datos
  if (typeof setupColorInputs === 'function') {
    setupColorInputs();
  } else {
    console.warn('[initializeCustomization] setupColorInputs no está disponible');
  }
  
  console.log('[initializeCustomization] Inicialización completada');
}

// Exportar funciones a window para que estén disponibles en onclick del HTML
window.uploadFormImage = uploadFormImage;
window.deleteFormImage = deleteFormImage;
window.previewImage = previewImage;
window.setupColorInputs = setupColorInputs;

// Exportar funciones como ES6 modules
export { 
  initializeCustomization,
  previewImage,
  uploadFormImage,
  deleteFormImage,
  setupColorInputs,
  loadCustomization,
  setCurrentForm
};

