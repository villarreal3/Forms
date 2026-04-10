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

// Utilidades para el manejo de tokens
export const Auth = {
  getToken: () => localStorage.getItem('auth_token'),
  setToken: (token) => localStorage.setItem('auth_token', token),
  removeToken: () => localStorage.removeItem('auth_token'),
  getRefreshToken: () => localStorage.getItem('auth_refresh_token'),
  setRefreshToken: (token) => localStorage.setItem('auth_refresh_token', token),
  removeRefreshToken: () => localStorage.removeItem('auth_refresh_token'),
  getUser: () => {
    const userStr = localStorage.getItem('auth_user');
    return userStr ? JSON.parse(userStr) : null;
  },
  setUser: (user) => localStorage.setItem('auth_user', JSON.stringify(user)),
  removeUser: () => localStorage.removeItem('auth_user'),
  isAuthenticated: () => !!Auth.getToken(),
  // Gestión de timestamp de sesión
  getLastActivity: () => {
    const timestamp = localStorage.getItem('auth_last_activity');
    return timestamp ? parseInt(timestamp, 10) : null;
  },
  setLastActivity: () => {
    localStorage.setItem('auth_last_activity', Date.now().toString());
  },
  removeLastActivity: () => {
    localStorage.removeItem('auth_last_activity');
  },
  // Verificar si la sesión ha expirado
  isSessionExpired: (timeoutMinutes = 15) => {
    const lastActivity = Auth.getLastActivity();
    if (!lastActivity) return true;
    const now = Date.now();
    const timeoutMs = timeoutMinutes * 60 * 1000;
    return (now - lastActivity) > timeoutMs;
  },
  logout: () => {
    Auth.removeToken();
    Auth.removeRefreshToken();
    Auth.removeUser();
    Auth.removeLastActivity();
    window.location.href = getLoginPath();
  }
};

