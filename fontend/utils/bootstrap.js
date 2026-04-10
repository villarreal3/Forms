// Punto de entrada principal para inicializar la aplicación
// Este archivo debe cargarse como módulo ES6

// Importar servicios
import { API_CONFIG, Auth, apiRequest, apiRequestBlob, checkAuth, SessionTimeout, showLoading, hideLoading } from '../services/index.js';

// Importar componentes
import { createNavbar, showNotification } from '../components/index.js';

// Importar utils
import { formatDate, formatNumber } from './format.js';
import { downloadFile } from './file.js';
import { getStatusText, getStatusColor } from './status.js';

// Importar layout
import { initLayout } from '../layout/adminLayout.js';

// Exportar todo a window para compatibilidad con código inline en HTML
// Nota: API_CONFIG ya se exporta en config.js e index.js, pero lo hacemos aquí para asegurar disponibilidad
if (typeof window !== 'undefined') {
  window.API_CONFIG = API_CONFIG;
  
  // Verificar que ENDPOINTS esté disponible
  if (!window.API_CONFIG.ENDPOINTS) {
    console.error('[Bootstrap] ERROR: API_CONFIG.ENDPOINTS no está disponible');
    console.log('[Bootstrap] API_CONFIG estructura:', {
      keys: Object.keys(window.API_CONFIG || {}),
      hasEndpoints: 'ENDPOINTS' in (window.API_CONFIG || {}),
      hasFormsEndpoints: 'FORMS_ENDPOINTS' in (window.API_CONFIG || {}),
      configType: typeof window.API_CONFIG
    });
  }
}
window.Auth = Auth;
window.apiRequest = apiRequest;
window.apiRequestBlob = apiRequestBlob;
window.checkAuth = checkAuth;
window.SessionTimeout = SessionTimeout;
window.showLoading = showLoading;
window.hideLoading = hideLoading;
window.createNavbar = createNavbar;
window.showNotification = showNotification;
window.formatDate = formatDate;
window.formatNumber = formatNumber;
window.downloadFile = downloadFile;
window.getStatusText = getStatusText;
window.getStatusColor = getStatusColor;
window.initLayout = initLayout;

// Verificar que todo esté correctamente exportado
console.log('[Bootstrap] API_CONFIG exportado:', {
  hasEndpoints: !!window.API_CONFIG?.ENDPOINTS,
  hasFormsEndpoints: !!window.API_CONFIG?.FORMS_ENDPOINTS,
  hasApiRequest: !!window.apiRequest,
  hasShowNotification: !!window.showNotification
});

// Función para inicializar la página
export function initPage() {
  // Verificar autenticación
  if (!checkAuth()) {
    return;
  }

  // Inicializar layout (carga el sidebar automáticamente)
  initLayout();

  // Inicializar timeout de sesión
  if (SessionTimeout && typeof SessionTimeout.init === 'function') {
    SessionTimeout.init();
  }
}

// Auto-inicializar cuando el DOM esté listo
function startInit() {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      requestAnimationFrame(() => {
        initPage();
      });
    });
  } else {
    requestAnimationFrame(() => {
      initPage();
    });
  }
}

startInit();
