// Función para mostrar loading
export function showLoading(element) {
  if (typeof element === 'string') {
    element = document.querySelector(element);
  }
  if (element) {
    element.innerHTML = '<div></div>';
  }
}

// Función para ocultar loading
export function hideLoading(element) {
  if (typeof element === 'string') {
    element = document.querySelector(element);
  }
  // El contenido se reemplazará por los datos
}

