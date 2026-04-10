/**
 * Configuración centralizada de APIs
 *
 * Este archivo centraliza todas las conexiones a las APIs de Forms y Users,
 * tanto para llamadas públicas como autenticadas.
 *
 * Modo desarrollo: el frontend debe abrirse en localhost o 127.0.0.1.
 * Servicios necesarios:
 * - Forms API: http://localhost:8080 (ej: cd forms-api/forms && go run .)
 * - Users API: http://localhost:3000 (ej: cd users-api && npm run start:dev)
 */

// ================== CONFIGURACIÓN DE URLs BASE ==================

// En producción se usa el proxy PHP; en desarrollo (localhost o IP de dev) se llama directo a las APIs.
function isProductionHost() {
 if (typeof window === 'undefined' || !window.location) return false;
 const h = window.location.hostname || '';
 const devHosts = ['localhost', '127.0.0.1', '172.28.125.76'];
 return !devHosts.includes(h);
}

const PROJECT_BASE_PATH = 'forms';

// URLs de desarrollo (localhost).
const DEV_USERS_API_BASE = 'http://localhost:3000/api/v1';
const DEV_FORMS_API_BASE = 'http://localhost:8080/api';

function getApiBaseUrl(proxyPathSuffix, devBaseUrl) {
 if (isProductionHost()) return `/${PROJECT_BASE_PATH}/${proxyPathSuffix}`;
 return devBaseUrl;
}

const _usersApiBaseUrl = getApiBaseUrl('proxy-users-api.php', DEV_USERS_API_BASE);
const _formsApiBaseUrl = getApiBaseUrl('proxy-forms-api.php', DEV_FORMS_API_BASE);

const API_CONFIG = {
 // Users API (producción: proxy; desarrollo: URL directa)
 USERS_API_BASE_URL: _usersApiBaseUrl,

 // Forms API (producción: proxy; desarrollo: URL directa)
 FORMS_API_BASE_URL: _formsApiBaseUrl,
 
 // Endpoints de Users API
 USERS_ENDPOINTS: {
 // Autenticación
 LOGIN: '/auth/login',
 REFRESH: '/auth/refresh',
 LOGOUT: '/auth/logout',
 CHANGE_PASSWORD: '/auth/change-password',
 REQUEST_PASSWORD_RESET: '/auth/request-password-reset',
 RESET_PASSWORD: '/auth/reset-password',
 VERIFY: '/auth/verify',
 
 // Usuarios
 ME: '/users/me',
 USERS: '/users',
 USER_BY_ID: (id) => `/users/${id}`,
 USERS_ADMIN: '/users',
 USER_SOFT_DELETE: (id) => `/users/${id}/soft`,
 USER_RESTORE: (id) => `/users/${id}/restore`,
 // Roles
 ROLES: '/roles',
 ROLE_BY_ID: (id) => `/roles/${id}`,
 // Departamentos
 DEPARTMENTS: '/departments',
 DEPARTMENT_BY_ID: (id) => `/departments/${id}`,
 },
 
 // Endpoints de Forms API
 FORMS_ENDPOINTS: {
 // Dashboard
 DASHBOARD_STATS: '/admin/dashboard/stats',
 DASHBOARD_FORMS: '/admin/dashboard/forms',
 DASHBOARD_SUBMISSIONS: '/admin/dashboard/submissions',
 DASHBOARD_RESPONSES_TIMELINE: '/admin/dashboard/responses-timeline',
 DASHBOARD_TOP_FORMS: '/admin/dashboard/top-active-forms',
 DASHBOARD_GEOGRAPHIC: '/admin/dashboard/geographic-distribution',
 DASHBOARD_USER_METRICS: '/admin/dashboard/user-metrics',
 
 // Formularios (Admin)
 FORMS: '/admin/forms',
 FORM_BY_ID: (id) => `/admin/forms/${id}`,
 FORM_SCHEMA: (id) => `/admin/forms/${id}/schema`,
 FORM: (id) => `/admin/forms/${id}`, // Alias para compatibilidad
 FORM_CLOSE: (id) => `/admin/forms/${id}/close`,
 FORM_OPEN: (id) => `/admin/forms/${id}/open`,
 FORM_PUBLISH: (id) => `/admin/forms/${id}/publish`,
 FORM_DRAFT: (id) => `/admin/forms/${id}/draft`,
 FORM_IMAGE_UPLOAD: (id) => `/admin/forms/${id}/image`,
 FORM_IMAGE_DELETE: (id) => `/admin/forms/${id}/image`,
 
 // Campos de formulario (Admin)
 FORM_FIELDS: (id) => `/admin/forms/${id}/fields`,
 FORM_FIELDS_COMMON: (id) => `/admin/forms/${id}/fields/common`,
 FORM_FIELD_CREATE: (id) => `/admin/forms/${id}/fields`,
 FORM_FIELD_UPDATE: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}`,
 FORM_FIELD_DELETE: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}`,
 FORM_FIELD_REORDER: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}/reorder`,
 FORM_FIELD_ASSIGN_SECTION: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}/assign-section`,
 
 // Secciones de formulario (Admin)
 FORM_SECTIONS: (formId) => `/admin/forms/${formId}/sections`,
 FORM_SECTION_CREATE: (formId) => `/admin/forms/${formId}/sections`,
 FORM_SECTION_UPDATE: (formId, sectionId) => `/admin/forms/${formId}/sections/${sectionId}`,
 FORM_SECTION_DELETE: (formId, sectionId) => `/admin/forms/${formId}/sections/${sectionId}`,
 FORM_SECTION_REORDER: (formId, sectionId) => `/admin/forms/${formId}/sections/${sectionId}/reorder`,
 // Alias para compatibilidad
 PUBLIC_FORM_SECTIONS: (formId) => `/public/forms/${formId}/sections`,
 
 // Personalización de formularios (Admin)
 FORM_CUSTOMIZATION: (id) => `/admin/customization/forms/${id}`,
 
 // Templates y campos
 FIELD_TEMPLATES: '/admin/fields/templates',

 // Roles (tabla local en forms-api para select en config/usuarios)
 ROLES: '/admin/roles',
 
 // Envíos/Respuestas (Admin)
 SUBMISSIONS: '/admin/submissions',
 EXPORT: '/admin/export/submissions',
 
 // Asistencia
 ATTENDANCE_UPDATE: '/admin/attendance/update',

 // Auto-inscripción pública basada en cédula + teléfono y schema compatible
 FORM_AUTO_SUBMIT: '/forms/auto-submit',
 
 // Email
 EMAIL_BULK: '/admin/email/bulk',
 
 // Formularios públicos (table-per-form; base = FORMS_API_BASE_URL, ej. /api)
 // GET /public/forms, GET /public/forms/{id}, GET /public/forms/{id}/sections, GET /public/forms/{id}/customization, POST /forms/submit
 PUBLIC_FORMS: '/public/forms',
 PUBLIC_FORM_BY_ID: (id) => `/public/forms/${id}`,
 PUBLIC_FORM_SECTIONS: (formId) => `/public/forms/${formId}/sections`,
 PUBLIC_FORM_CUSTOMIZATION: (formId) => `/public/forms/${formId}/customization`,
 /** Solo URLs de imagen de cabecera (GET /public/forms/{id}/logo) */
 PUBLIC_FORM_LOGO: (formId) => `/public/forms/${formId}/logo`,
 FORM_SUBMIT: '/forms/submit',
 
 // Usuarios (público)
 USERS_GET: '/users/get',
 
 // Imágenes
 IMAGE: (formId, filename) => `/api/images/${formId}/${filename}`,
 },
 
 // Objeto ENDPOINTS unificado para compatibilidad con código antiguo
 // Combina endpoints de Users y Forms
 ENDPOINTS: {
 // Endpoints del servicio Login (Users API)
 LOGIN: '/auth/login',
 ME: '/users/me',
 USERS_ADMIN: '/users',
 // Endpoints del servicio Forms
 DASHBOARD_STATS: '/admin/dashboard/stats',
 DASHBOARD_FORMS: '/admin/dashboard/forms',
 DASHBOARD_SUBMISSIONS: '/admin/dashboard/submissions',
 DASHBOARD_RESPONSES_TIMELINE: '/admin/dashboard/responses-timeline',
 DASHBOARD_TOP_FORMS: '/admin/dashboard/top-active-forms',
 DASHBOARD_GEOGRAPHIC: '/admin/dashboard/geographic-distribution',
 DASHBOARD_USER_METRICS: '/admin/dashboard/user-metrics',
 FORMS: '/admin/forms',
 FORM: (id) => `/admin/forms/${id}`,
 FORM_SCHEMA: (id) => `/admin/forms/${id}/schema`,
 FORM_CLOSE: (id) => `/admin/forms/${id}/close`,
 FORM_OPEN: (id) => `/admin/forms/${id}/open`,
 FORM_PUBLISH: (id) => `/admin/forms/${id}/publish`,
 FORM_DRAFT: (id) => `/admin/forms/${id}/draft`,
 FORM_FIELD_CREATE: (id) => `/admin/forms/${id}/fields`,
 FORM_FIELD_UPDATE: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}`,
 FORM_FIELD_DELETE: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}`,
 FORM_FIELD_REORDER: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}/reorder`,
 FORM_FIELD_ASSIGN_SECTION: (formId, fieldId) => `/admin/forms/${formId}/fields/${fieldId}/assign-section`,
 FORM_FIELDS: (id) => `/admin/forms/${id}/fields`,
 FORM_FIELDS_COMMON: (id) => `/admin/forms/${id}/fields/common`,
 // Secciones
 FORM_SECTIONS: (formId) => `/admin/forms/${formId}/sections`,
 FORM_SECTION_CREATE: (formId) => `/admin/forms/${formId}/sections`,
 FORM_SECTION_UPDATE: (formId, sectionId) => `/admin/forms/${formId}/sections/${sectionId}`,
 FORM_SECTION_DELETE: (formId, sectionId) => `/admin/forms/${formId}/sections/${sectionId}`,
 FORM_SECTION_REORDER: (formId, sectionId) => `/admin/forms/${formId}/sections/${sectionId}/reorder`,
 // Públicos
 PUBLIC_FORM_SECTIONS: (formId) => `/public/forms/${formId}/sections`,
 FIELD_TEMPLATES: '/admin/fields/templates',
 SUBMISSIONS: '/admin/submissions',
 ATTENDANCE_UPDATE: '/admin/attendance/update',
 EXPORT: '/admin/export/submissions',
 EMAIL_BULK: '/admin/email/bulk',
 // Personalización de formularios
 FORM_CUSTOMIZATION: (id) => `/admin/customization/forms/${id}`,
 FORM_IMAGE_UPLOAD: (id) => `/admin/forms/${id}/image`,
 FORM_IMAGE_DELETE: (id) => `/admin/forms/${id}/image`,
 // Endpoints públicos del servicio Forms
 PUBLIC_FORMS: '/public/forms',
 PUBLIC_FORM: (id) => `/public/forms/${id}`,
 FORM_SUBMIT: '/forms/submit',
 USERS_GET: '/users/get',
 FORM_AUTO_SUBMIT: '/forms/auto-submit',
 // Servir imágenes
 IMAGE: (formId, filename) => `/api/images/${formId}/${filename}`
 }
};

// ================== FUNCIONES PARA LLAMADAS PÚBLICAS ==================
// Estas funciones NO requieren autenticación

/**
 * Realiza una petición pública (sin autenticación) a la Forms API
 * @param {string} endpoint - Endpoint relativo (ej: '/public/forms')
 * @param {object} options - Opciones de fetch (method, body, headers, etc.)
 * @returns {Promise<object>} - Respuesta JSON
 */
export async function publicRequest(endpoint, options = {}) {
 const url = `${API_CONFIG.FORMS_API_BASE_URL}${endpoint}`;
 
 console.log('[publicRequest] URL completa:', url);
 console.log('[publicRequest] Método:', options.method || 'GET');
 console.log('[publicRequest] Endpoint:', endpoint);
 console.log('[publicRequest] Base URL:', API_CONFIG.FORMS_API_BASE_URL);
 
 const defaultOptions = {
 headers: {
 'Content-Type': 'application/json',
 },
 ...options,
 };

 try {
 console.log('[publicRequest] Realizando petición fetch...');
 const response = await fetch(url, defaultOptions);
 console.log('[publicRequest] Respuesta recibida:', {
 status: response.status,
 statusText: response.statusText,
 ok: response.ok,
 headers: Object.fromEntries(response.headers.entries())
 });
 
 if (!response.ok) {
 const errorData = await response.json().catch(() => ({ 
 message: `Error ${response.status}: ${response.statusText}` 
 }));
 console.error('[publicRequest] Error en respuesta:', errorData);
 throw new Error(errorData.message || errorData.error || 'Error en la petición');
 }

 const contentType = response.headers.get('content-type');
 console.log('[publicRequest] Content-Type:', contentType);
 
 if (contentType && contentType.includes('application/json')) {
 const jsonData = await response.json();
 console.log('[publicRequest] Datos JSON recibidos:', jsonData);
 return jsonData;
 }
 
 const textData = await response.text();
 console.log('[publicRequest] Datos de texto recibidos:', textData);
 return textData;
 } catch (error) {
 console.error('[publicRequest] Error capturado:', error);
 console.error('[publicRequest] Tipo de error:', error.constructor.name);
 console.error('[publicRequest] Mensaje:', error.message);
 if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
 const base = API_CONFIG.FORMS_API_BASE_URL;
 throw new Error(
 `No se pudo conectar a la API de formularios. En desarrollo, inicia el Forms API (puerto 8080). URL esperada: ${base}`
 );
 }
 throw error;
 }
}

/**
 * Realiza una petición pública a la Users API
 * @param {string} endpoint - Endpoint relativo (ej: '/auth/login')
 * @param {object} options - Opciones de fetch
 * @returns {Promise<object>} - Respuesta JSON
 */
export async function publicUsersRequest(endpoint, options = {}) {
 const url = `${API_CONFIG.USERS_API_BASE_URL}${endpoint}`;
 
 const defaultOptions = {
 headers: {
 'Content-Type': 'application/json',
 },
 ...options,
 };

 try {
 const response = await fetch(url, defaultOptions);
 
 if (!response.ok) {
 const errorData = await response.json().catch(() => ({ 
 message: `Error ${response.status}: ${response.statusText}` 
 }));
 throw new Error(errorData.message || errorData.error || 'Error en la petición');
 }

 return await response.json();
 } catch (error) {
 console.error('Error en publicUsersRequest:', error);
 if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
 throw new Error(`No se pudo conectar al servidor. Verifica que el servicio esté corriendo en ${API_CONFIG.USERS_API_BASE_URL}`);
 }
 throw error;
 }
}

// ================== FUNCIONES PARA LLAMADAS AUTENTICADAS ==================
// Estas funciones SÍ requieren autenticación (token JWT)

/**
 * Obtiene el token de autenticación del localStorage
 * @returns {string|null} - Token JWT o null si no existe
 */
function getAuthToken() {
 return localStorage.getItem('auth_token');
}

/**
 * Determina qué URL base usar según el endpoint
 * @param {string} endpoint - Endpoint relativo
 * @returns {string} - URL base (Users API o Forms API)
 */
function getBaseURLForEndpoint(endpoint) {
 // Endpoints que van a Users API
 const usersEndpoints = [
 '/auth/',
 '/users',
 '/roles',
 '/departments',
 ];
 
 for (const usersEndpoint of usersEndpoints) {
 if (endpoint.startsWith(usersEndpoint)) {
 return API_CONFIG.USERS_API_BASE_URL;
 }
 }
 
 // Todos los demás van a Forms API
 return API_CONFIG.FORMS_API_BASE_URL;
}

/**
 * Realiza una petición autenticada a cualquier API
 * @param {string} endpoint - Endpoint relativo
 * @param {object} options - Opciones de fetch
 * @returns {Promise<object>} - Respuesta JSON
 */
export async function authenticatedRequest(endpoint, options = {}) {
 const token = getAuthToken();
 const baseURL = getBaseURLForEndpoint(endpoint);
 const url = `${baseURL}${endpoint}`;
 
 const defaultOptions = {
 headers: {
 'Content-Type': 'application/json',
 ...(token && { 'Authorization': `Bearer ${token}` })
 },
 ...options,
 headers: {
 'Content-Type': 'application/json',
 ...(token && { 'Authorization': `Bearer ${token}` }),
 ...(options.headers || {})
 }
 };

 try {
 const response = await fetch(url, defaultOptions);
 
 let data;
 const contentType = response.headers.get('content-type');
 if (contentType && contentType.includes('application/json')) {
 data = await response.json();
 } else {
 const text = await response.text();
 throw new Error(`Error del servidor: ${response.status} ${response.statusText}. ${text}`);
 }

 if (!response.ok) {
 if (response.status === 401) {
 const errorMsg = data.message || data.error || '';
 const isUsersApi = getBaseURLForEndpoint(endpoint) === API_CONFIG.USERS_API_BASE_URL;
 if (isUsersApi && (errorMsg.includes('expirado') || errorMsg.includes('inválido') || errorMsg.includes('invalid'))) {
 localStorage.removeItem('auth_token');
 localStorage.removeItem('auth_user');
 const loginPath = window.location.pathname.includes('/admin/') 
 ? '../login/' 
 : 'admin/login/';
 window.location.href = loginPath;
 throw new Error('Sesión expirada. Por favor, inicia sesión nuevamente.');
 }
 throw new Error(errorMsg || 'Error de autenticación. Verifica la configuración del servidor.');
 }
 throw new Error(data.message || data.error || 'Error en la petición');
 }

 return data;
 } catch (error) {
 console.error('Error en authenticatedRequest:', error);
 if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
 throw new Error(`No se pudo conectar al servidor. Verifica que los servicios estén corriendo.`);
 }
 throw error;
 }
}

/**
 * Realiza una petición autenticada que devuelve un Blob (para descargas)
 * @param {string} endpoint - Endpoint relativo
 * @param {object} options - Opciones de fetch
 * @returns {Promise<Blob>} - Blob de la respuesta
 */
export async function authenticatedRequestBlob(endpoint, options = {}) {
 const token = getAuthToken();
 const baseURL = getBaseURLForEndpoint(endpoint);
 const url = `${baseURL}${endpoint}`;
 
 const defaultOptions = {
 headers: {
 'Content-Type': 'application/json',
 ...(token && { 'Authorization': `Bearer ${token}` })
 },
 ...options,
 headers: {
 'Content-Type': 'application/json',
 ...(token && { 'Authorization': `Bearer ${token}` }),
 ...(options.headers || {})
 }
 };

 try {
 const response = await fetch(url, defaultOptions);

 if (!response.ok) {
 if (response.status === 401) {
 const isUsersApi = getBaseURLForEndpoint(endpoint) === API_CONFIG.USERS_API_BASE_URL;
 if (isUsersApi) {
 localStorage.removeItem('auth_token');
 localStorage.removeItem('auth_user');
 const loginPath = window.location.pathname.includes('/admin/') 
 ? '../login/' 
 : 'admin/login/';
 window.location.href = loginPath;
 throw new Error('Sesión expirada. Por favor, inicia sesión nuevamente.');
 }
 const errorData = await response.json().catch(() => ({}));
 throw new Error(errorData.message || errorData.error || 'Error de autenticación.');
 }
 const errorData = await response.json();
 throw new Error(errorData.message || errorData.error || 'Error en la petición');
 }

 return await response.blob();
 } catch (error) {
 console.error('Error en authenticatedRequestBlob:', error);
 throw error;
 }
}

// ================== FUNCIONES HELPER ESPECÍFICAS ==================

/**
 * Obtiene la lista de formularios públicos
 * @returns {Promise<object>} - Lista de formularios
 */
export async function getPublicForms() {
 const endpoint = API_CONFIG.FORMS_ENDPOINTS.PUBLIC_FORMS;
 const fullUrl = `${API_CONFIG.FORMS_API_BASE_URL}${endpoint}`;
 console.log('[getPublicForms] Llamando a:', fullUrl);
 return publicRequest(endpoint);
}

/**
 * Obtiene un formulario público por ID
 * @param {number|string} formId - ID del formulario
 * @returns {Promise<object>} - Datos del formulario
 */
export async function getPublicForm(formId) {
 const endpoint = API_CONFIG.FORMS_ENDPOINTS.PUBLIC_FORM_BY_ID(formId);
 const fullUrl = `${API_CONFIG.FORMS_API_BASE_URL}${endpoint}`;
 console.log('[getPublicForm] Llamando a:', fullUrl);
 return publicRequest(endpoint);
}

/**
 * Obtiene las secciones de un formulario público
 * @param {number|string} formId - ID del formulario
 * @returns {Promise<object>} - Secciones del formulario
 */
export async function getPublicFormSections(formId) {
 const endpoint = API_CONFIG.FORMS_ENDPOINTS.PUBLIC_FORM_SECTIONS(formId);
 const fullUrl = `${API_CONFIG.FORMS_API_BASE_URL}${endpoint}`;
 console.log('[getPublicFormSections] Llamando a:', fullUrl);
 return publicRequest(endpoint);
}

/**
 * Obtiene la personalización de un formulario público
 * @param {number|string} formId - ID del formulario
 * @returns {Promise<object>} - Personalización del formulario
 */
export async function getPublicFormCustomization(formId) {
 const endpoint = API_CONFIG.FORMS_ENDPOINTS.PUBLIC_FORM_CUSTOMIZATION(formId);
 const fullUrl = `${API_CONFIG.FORMS_API_BASE_URL}${endpoint}`;
 console.log('[getPublicFormCustomization] Llamando a:', fullUrl);
 return publicRequest(endpoint);
}

/**
 * Obtiene solo las URLs de logo del formulario público (endpoint dedicado, ligero).
 * @param {number|string} formId - ID del formulario
 * @returns {Promise<{ success?: boolean, logo_url?: string, logo_url_mobile?: string }>}
 */
export async function getPublicFormLogo(formId) {
 const endpoint = API_CONFIG.FORMS_ENDPOINTS.PUBLIC_FORM_LOGO(formId);
 const fullUrl = `${API_CONFIG.FORMS_API_BASE_URL}${endpoint}`;
 console.log('[getPublicFormLogo] Llamando a:', fullUrl);
 return publicRequest(endpoint);
}

/**
 * Envía un formulario (público, sin autenticación)
 * @param {object} payload - Datos del formulario a enviar
 * @returns {Promise<object>} - Respuesta del servidor
 */
export async function submitForm(payload) {
 const endpoint = API_CONFIG.FORMS_ENDPOINTS.FORM_SUBMIT;
 const fullUrl = `${API_CONFIG.FORMS_API_BASE_URL}${endpoint}`;
 console.log('[submitForm] Llamando a:', fullUrl, 'con payload:', payload);
 return publicRequest(endpoint, {
 method: 'POST',
 body: JSON.stringify(payload)
 });
}

/**
 * Realiza login
 * @param {string} email - Email del usuario
 * @param {string} password - Contraseña
 * @returns {Promise<object>} - Respuesta con token y usuario
 */
export async function login(email, password) {
 return publicUsersRequest(API_CONFIG.USERS_ENDPOINTS.LOGIN, {
 method: 'POST',
 body: JSON.stringify({ email, password })
 });
}

/**
 * Renueva el access token usando el refresh token.
 * @returns {Promise<{accessToken: string, refreshToken?: string, user?: object}>} Nuevos tokens
 */
export async function refreshAccessToken() {
 const refreshToken = typeof localStorage !== 'undefined' ? localStorage.getItem('auth_refresh_token') : null;
 if (!refreshToken) throw new Error('No hay refresh token');
 const url = `${API_CONFIG.USERS_API_BASE_URL}${API_CONFIG.USERS_ENDPOINTS.REFRESH}`;
 const response = await fetch(url, {
 method: 'POST',
 headers: { 'Content-Type': 'application/json' },
 body: JSON.stringify({ refreshToken }),
 });
 const data = await response.json().catch(() => ({}));
 if (!response.ok) throw new Error(data.message || data.error || 'Error al renovar sesión');
 return data;
}

/**
 * Obtiene información del usuario actual (autenticado)
 * @returns {Promise<object>} - Datos del usuario
 */
export async function getCurrentUser() {
 return authenticatedRequest(API_CONFIG.USERS_ENDPOINTS.ME);
}

/**
 * Obtiene la lista de formularios (admin, autenticado)
 * @returns {Promise<object>} - Lista de formularios
 */
export async function getForms() {
 return authenticatedRequest(API_CONFIG.FORMS_ENDPOINTS.FORMS);
}

/**
 * Obtiene un formulario por ID (admin, autenticado)
 * @param {number|string} formId - ID del formulario
 * @returns {Promise<object>} - Datos del formulario
 */
export async function getForm(formId) {
 return authenticatedRequest(API_CONFIG.FORMS_ENDPOINTS.FORM_BY_ID(formId));
}

/**
 * Obtiene estadísticas del dashboard (admin, autenticado)
 * @returns {Promise<object>} - Estadísticas
 */
export async function getDashboardStats() {
 return authenticatedRequest(API_CONFIG.FORMS_ENDPOINTS.DASHBOARD_STATS);
}

// ================== EXPORTAR CONFIGURACIÓN ==================

export { API_CONFIG };

// Exportar también a window para compatibilidad con código inline en HTML
if (typeof window !== 'undefined') {
 window.API_CONFIG = API_CONFIG;
 window.publicRequest = publicRequest;
 window.publicUsersRequest = publicUsersRequest;
 window.authenticatedRequest = authenticatedRequest;
 window.authenticatedRequestBlob = authenticatedRequestBlob;
 window.getPublicForms = getPublicForms;
 window.getPublicForm = getPublicForm;
 window.getPublicFormSections = getPublicFormSections;
 window.getPublicFormCustomization = getPublicFormCustomization;
 window.getPublicFormLogo = getPublicFormLogo;
 window.submitForm = submitForm;
 window.login = login;
 window.getCurrentUser = getCurrentUser;
 window.getForms = getForms;
 window.getForm = getForm;
 window.getDashboardStats = getDashboardStats;
 
 // Log para debugging (solo en desarrollo)
 if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1' || window.location.hostname.includes('172.28')) {
 console.log('[config.js] Funciones exportadas a window:', {
 hasAPI_CONFIG: !!window.API_CONFIG,
 hasEndpoints: !!window.API_CONFIG?.ENDPOINTS,
 hasGetPublicForms: typeof window.getPublicForms === 'function',
 hasGetPublicForm: typeof window.getPublicForm === 'function',
 FORMS_API_BASE_URL: window.API_CONFIG?.FORMS_API_BASE_URL,
 USERS_API_BASE_URL: window.API_CONFIG?.USERS_API_BASE_URL
 });
 }
}

