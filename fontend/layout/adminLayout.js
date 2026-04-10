// Inicialización del layout administrativo
import { createNavbar } from '../components/sidebar.js';

export function initLayout() {
  // Obtener la página actual desde window.page o del pathname
  const currentPage = window.page || getCurrentPageFromPath();
  
  // Crear e inicializar el navbar
  createNavbar({ 
    currentPage,
    containerSelector: '#navbar'
  });
}

function getCurrentPageFromPath() {
  const path = window.location.pathname;
  if (path.includes('dashboard')) return 'dashboard';
  if (path.includes('formularios')) return 'formularios';
  if (path.includes('respuestas')) return 'respuestas';
  if (path.includes('correos')) return 'correos';
  if (path.includes('config')) return 'config';
  if (path.includes('changelog')) return 'changelog';
  if (path.includes('admin') && path.includes('config')) return 'admin';
  return 'dashboard';
}

// Exportar también a window para compatibilidad con código inline en HTML
window.initLayout = initLayout;
