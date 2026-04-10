// Barrel export - Re-exportar todos los servicios
// Configuración unificada desde config.js
export { API_CONFIG } from './config.js';
export { Auth } from './auth.js';
export { getBaseURL, apiRequest, apiRequestBlob } from './api-client.js';
export { checkAuth } from './auth-check.js';
export { SessionTimeout } from './session-timeout.js';
export { showLoading, hideLoading } from './utils.js';

// Exportar también con alias para compatibilidad
export { 
  API_CONFIG as CONFIG,
  publicRequest,
  publicUsersRequest,
  authenticatedRequest,
  authenticatedRequestBlob,
  getPublicForms,
  getPublicForm,
  getPublicFormSections,
  getPublicFormCustomization,
  submitForm,
  login,
  getCurrentUser,
  getForms,
  getForm,
  getDashboardStats
} from './config.js';

// Exportar también a window para compatibilidad con código inline en HTML
// Nota: config.js ya exporta a window, pero lo hacemos aquí también para asegurar disponibilidad
import { API_CONFIG } from './config.js';
import { Auth } from './auth.js';
import { getBaseURL, apiRequest, apiRequestBlob } from './api-client.js';
import { checkAuth } from './auth-check.js';
import { SessionTimeout } from './session-timeout.js';
import { showLoading, hideLoading } from './utils.js';
import {
  API_CONFIG as CONFIG,
  publicRequest,
  publicUsersRequest,
  authenticatedRequest,
  authenticatedRequestBlob,
  getPublicForms,
  getPublicForm,
  getPublicFormSections,
  getPublicFormCustomization,
  submitForm,
  login,
  getCurrentUser,
  getForms,
  getForm,
  getDashboardStats
} from './config.js';

// Asegurar que API_CONFIG esté disponible en window con todas sus propiedades
// (config.js ya lo exporta, pero esto asegura que esté disponible cuando se carga index.js)
if (typeof window !== 'undefined') {
  window.API_CONFIG = API_CONFIG;
  window.CONFIG = CONFIG;
  
  // Verificar que ENDPOINTS esté disponible
  if (!window.API_CONFIG.ENDPOINTS) {
    console.warn('[index.js] API_CONFIG.ENDPOINTS no está disponible, verificando estructura...');
    console.log('[index.js] API_CONFIG keys:', Object.keys(window.API_CONFIG || {}));
  }
}
window.Auth = Auth;
window.getBaseURL = getBaseURL;
window.apiRequest = apiRequest;
window.apiRequestBlob = apiRequestBlob;
window.checkAuth = checkAuth;
window.SessionTimeout = SessionTimeout;
window.showLoading = showLoading;
window.hideLoading = hideLoading;
// Nuevas funciones centralizadas
window.publicRequest = publicRequest;
window.publicUsersRequest = publicUsersRequest;
window.authenticatedRequest = authenticatedRequest;
window.authenticatedRequestBlob = authenticatedRequestBlob;
window.getPublicForms = getPublicForms;
window.getPublicForm = getPublicForm;
window.getPublicFormSections = getPublicFormSections;
window.getPublicFormCustomization = getPublicFormCustomization;
window.submitForm = submitForm;
window.login = login;
window.getCurrentUser = getCurrentUser;
window.getForms = getForms;
window.getForm = getForm;
window.getDashboardStats = getDashboardStats;
