import { API_CONFIG, refreshAccessToken } from './config.js';
import { Auth } from './auth.js';

let refreshPromise = null;
async function doRefresh() {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      const data = await refreshAccessToken();
      Auth.setToken(data.accessToken);
      if (data.refreshToken) Auth.setRefreshToken(data.refreshToken);
      Auth.setLastActivity();
      return data.accessToken;
    })().finally(() => { refreshPromise = null; });
  }
  return refreshPromise;
}

// Función para determinar qué servicio usar según el endpoint
export function getBaseURL(endpoint) {
  // Endpoints específicos del servicio Login (Users API)
  const loginEndpoints = [
    '/auth/login',
    '/auth/refresh',
    '/auth/logout',
    '/auth/change-password',
    '/auth/request-password-reset',
    '/auth/reset-password',
    '/auth/verify',
    '/users/me',
    '/users', // /users y todas sus variantes
    '/roles',
    '/departments'
  ];
  
  // Verificar si el endpoint coincide exactamente o empieza con alguno de los endpoints de login
  for (const loginEndpoint of loginEndpoints) {
    // Coincidencia exacta
    if (endpoint === loginEndpoint) {
      return API_CONFIG.USERS_API_BASE_URL;
    }
    // O si empieza con el endpoint de login (para rutas con parámetros)
    if (endpoint.startsWith(loginEndpoint + '/') || endpoint.startsWith(loginEndpoint + '?')) {
      return API_CONFIG.USERS_API_BASE_URL;
    }
  }

  // Todos los demás endpoints van al servicio Forms
  return API_CONFIG.FORMS_API_BASE_URL;
}

// Función para hacer peticiones a la API
export async function apiRequest(endpoint, options = {}) {
  const token = Auth.getToken();
  const baseURL = getBaseURL(endpoint);
  const url = `${baseURL}${endpoint}`;
  
  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` })
    }
  };

  const finalOptions = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...(options.headers || {})
    }
  };

  try {
    let response = await fetch(url, finalOptions);
    let contentType = response.headers.get('content-type');
    let data = contentType && contentType.includes('application/json') ? await response.json() : await response.text();

    if (!response.ok && response.status === 401) {
      const loginEndpoint = API_CONFIG.ENDPOINTS?.LOGIN || API_CONFIG.USERS_ENDPOINTS?.LOGIN || '/auth/login';
      const refreshEndpoint = API_CONFIG.USERS_ENDPOINTS?.REFRESH || '/auth/refresh';
      if (endpoint !== loginEndpoint && endpoint !== refreshEndpoint && Auth.getRefreshToken()) {
        try {
          const newToken = await doRefresh();
          const retryHeaders = { ...finalOptions.headers, 'Authorization': `Bearer ${newToken}` };
          const retryResponse = await fetch(url, { ...finalOptions, headers: retryHeaders });
          const retryContentType = retryResponse.headers.get('content-type');
          const retryData = retryContentType && retryContentType.includes('application/json') ? await retryResponse.json() : await retryResponse.text();
          if (retryResponse.ok) return retryData;
          response = retryResponse;
          data = retryData;
        } catch (_) {
          /* refresh falló, continuar con el 401 original */
        }
      }
    }

    if (!response.ok) {
      const errMsg = typeof data === 'object' ? (data.message || data.error || '') : String(data);
      if (response.status === 401) {
        const loginEndpoint = API_CONFIG.ENDPOINTS?.LOGIN || API_CONFIG.USERS_ENDPOINTS?.LOGIN || '/auth/login';
        if (endpoint === loginEndpoint) throw new Error(errMsg || '');
        const isUsersApi = getBaseURL(endpoint) === API_CONFIG.USERS_API_BASE_URL;
        if (isUsersApi) {
          if (!Auth.getToken()) {
            Auth.logout();
            throw new Error('Sesión expirada. Por favor, inicia sesión nuevamente.');
          }
          if (errMsg.includes('expirado') || errMsg.includes('inválido') || errMsg.includes('invalid')) {
            Auth.logout();
            throw new Error('Sesión expirada. Por favor, inicia sesión nuevamente.');
          }
        }
        throw new Error(errMsg || 'Error de autenticación. Verifica la configuración del servidor.');
      }
      throw new Error(errMsg || 'Error en la petición');
    }

    return data;
  } catch (error) {
    console.error('Error en apiRequest:', error);
    // Mejorar mensajes de error para problemas de conexión
    if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
      throw new Error(`No se pudo conectar al servidor. Verifica que los servicios estén corriendo en ${baseURL}`);
    }
    throw error;
  }
}

// Función para exportar a Excel (maneja archivos binarios)
export async function apiRequestBlob(endpoint, options = {}) {
  const token = Auth.getToken();
  const baseURL = getBaseURL(endpoint);
  const url = `${baseURL}${endpoint}`;
  
  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` })
    }
  };

  const finalOptions = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...(options.headers || {})
    }
  };

  try {
    let response = await fetch(url, finalOptions);

    if (!response.ok && response.status === 401) {
      const refreshEndpoint = API_CONFIG.USERS_ENDPOINTS?.REFRESH || '/auth/refresh';
      if (endpoint !== refreshEndpoint && Auth.getRefreshToken()) {
        try {
          const newToken = await doRefresh();
          const retryHeaders = { ...finalOptions.headers, 'Authorization': `Bearer ${newToken}` };
          response = await fetch(url, { ...finalOptions, headers: retryHeaders });
        } catch (_) {
          /* refresh falló */
        }
      }
    }

    if (!response.ok) {
      if (response.status === 401) {
        const isUsersApi = getBaseURL(endpoint) === API_CONFIG.USERS_API_BASE_URL;
        if (isUsersApi) {
          Auth.logout();
          throw new Error('Sesión expirada. Por favor, inicia sesión nuevamente.');
        }
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || errorData.error || 'Error de autenticación.');
      }
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || errorData.error || 'Error en la petición');
    }

    return await response.blob();
  } catch (error) {
    console.error('Error en apiRequestBlob:', error);
    throw error;
  }
}

