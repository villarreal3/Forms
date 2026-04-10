import { Auth } from './auth.js';

// Función helper para calcular la ruta correcta al login
function getLoginPath() {
  const path = window.location.pathname;
  
  // Si estamos en /admin/ o subdirectorios de /admin/
  if (path.includes('/admin/')) {
    // Encontrar la posición de /admin/ en el path
    const adminIndex = path.indexOf('/admin/');
    if (adminIndex !== -1) {
      // Obtener la parte del path después de /admin/
      const afterAdmin = path.substring(adminIndex + 7); // 7 = longitud de '/admin/'
      
      // Contar cuántos niveles de profundidad hay (contando directorios, no el archivo)
      const parts = afterAdmin.split('/').filter(part => part && !part.includes('.html'));
      const depth = parts.length;
      
      // Construir la ruta con '../' repetido según la profundidad
      return depth > 0 ? '../'.repeat(depth) + 'login/' : 'login/';
    }
  }
  
  // Si no estamos en /admin/, necesitamos 'admin/login/'
  return 'admin/login/';
}

// Función para verificar autenticación en páginas protegidas
export function checkAuth() {
  // No verificar autenticación en la página de login
  if (window.location.pathname.includes('/admin/login')) {
    return true;
  }
  
  if (!Auth.isAuthenticated()) {
    window.location.href = getLoginPath();
    return false;
  }
  // Verificar si la sesión ha expirado por timeout
  if (Auth.isSessionExpired()) {
    Auth.logout();
    return false;
  }
  // Actualizar timestamp de última actividad
  Auth.setLastActivity();
  return true;
}

