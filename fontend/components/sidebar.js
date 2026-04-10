import { Auth } from '../services/auth.js';
import { showNotification } from './notification.js';

// components/navbar.js
export function createNavbar({
  currentPage = '',
  containerSelector = '#navbar'
} = {}) {
  const container = document.querySelector(containerSelector);
  if (!container) return;

  // Validación de autenticación
  const user = Auth.getUser();
  if (!user) return;

  // Detectar ruta para construir paths correctos
  const pathInfo = detectPath();
  const { pathPrefix } = pathInfo;

  const nav = document.createElement('nav');
  nav.className = 'navigation';

  // Construir items de navegación (solo rutas de carpeta, sin index.html)
  const items = [
    { page: 'dashboard', text: 'Dashboard', href: `${pathPrefix}dashboard/` },
    { page: 'formularios', text: 'Formularios', href: `${pathPrefix}formularios/` },
    { page: 'respuestas', text: 'Respuestas', href: `${pathPrefix}respuestas/` },
    { page: 'encuestas', text: 'Encuestas', href: `${pathPrefix}encuestas/`, disabled: true },
    { page: 'config', text: 'Configuración', href: `${pathPrefix}config/` }
  ];

  // Agregar item de correos masivos solo para admin
  if (user.role === 'admin') {
    items.splice(items.length - 1, 0, {
      page: 'correos',
      text: 'Correos Masivos',
      href: `${pathPrefix}email/`,
      disabled: false
    });
  }

  // Construir HTML de la lista
  const listItems = items.map(item => {
    const activeClass = item.page === currentPage ? 'active' : '';
    const disabledAttr = item.disabled ? 'data-disabled="true"' : '';
    const ariaDisabledAttr = item.disabled ? 'aria-disabled="true" tabindex="-1"' : '';
    const href = item.disabled ? '#' : item.href;

    return `<li><a href="${href}" data-page="${item.page}" class="${activeClass}" ${disabledAttr} ${ariaDisabledAttr}>${item.text}</a></li>`;
  }).join('');

  nav.innerHTML = `
    <ul class="navigation-list">
      ${listItems}
    </ul>
  `;

  // Marcar página activa
  if (currentPage) {
    const activeLink = nav.querySelector(`[data-page="${currentPage}"]`);
    if (activeLink) activeLink.classList.add('active');
  }

  // Bindear eventos para items deshabilitados
  nav.querySelectorAll('[data-disabled="true"]').forEach(el => {
    el.addEventListener('click', e => {
      e.preventDefault();
      showNotification(
        'Esta sección estará disponible en una próxima actualización.',
        'info'
      );
    });
  });

  container.appendChild(nav);
}

// Función helper para detectar la ruta actual
function detectPath() {
  const path = window.location.pathname;
  const isInConfig = path.includes('/admin/config/');
  const isInFormularios = path.includes('/admin/formularios/');
  const isInAdmin = path.includes('/admin/') && !isInConfig && !isInFormularios;

  // Calcular la profundidad desde /admin/
  // Ejemplo: /admin/formularios/actions/view-form/ -> profundidad (segmentos bajo /admin/)
  // Necesitamos ../../ para volver a /admin/
  let pathPrefix = '';
  
  if (path.includes('/admin/')) {
    // Encontrar la posición de /admin/ en el path
    const adminIndex = path.indexOf('/admin/');
    if (adminIndex !== -1) {
      // Obtener la parte del path después de /admin/
      const afterAdmin = path.substring(adminIndex + 7); // 7 = longitud de '/admin/'
      
      // Contar cuántos niveles de profundidad hay (contando directorios, no el archivo)
      // Dividir por '/' y filtrar elementos vacíos y el nombre del archivo
      const parts = afterAdmin.split('/').filter(part => part && !part.includes('.html'));
      const depth = parts.length;
      
      // Construir el pathPrefix con '../' repetido según la profundidad
      pathPrefix = depth > 0 ? '../'.repeat(depth) : '';
    }
  } else {
    // Si no estamos en /admin/, necesitamos 'admin/'
    pathPrefix = 'admin/';
  }

  return {
    isInConfig,
    isInAdmin,
    isInFormularios,
    pathPrefix,
    path
  };
}

