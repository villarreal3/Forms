// Función para mostrar notificaciones
// @param {string} message - Texto a mostrar
// @param {'info'|'success'|'error'|'warning'} type - Tipo de notificación
// @param {number} [durationMs=3000] - Tiempo visible en ms antes de ocultar
export function showNotification(message, type = 'info', durationMs = 3000) {
  const notification = document.createElement('div');
  notification.style.cssText = 'position: fixed; top: 1rem; right: 1rem; padding: 0.75rem 1.5rem; border-radius: 0.5rem; box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1); z-index: 50; transform: translateX(100%); transition: transform 0.3s;';
  notification.setAttribute('role', 'status');
  notification.setAttribute('aria-live', 'polite');
  if (type === 'info') notification.style.backgroundColor = '#3b82f6';
  else if (type === 'success') notification.style.backgroundColor = '#10b981';
  else if (type === 'error') notification.style.backgroundColor = '#ef4444';
  else if (type === 'warning') notification.style.backgroundColor = '#f59e0b';
  notification.style.color = '#ffffff';
  notification.innerHTML = `
    <div>
      <span>${message}</span>
      <button
        type="button"
        aria-label="Cerrar notificación"
        title="Cerrar"
        onclick="this.parentElement.parentElement.remove()"
        style="margin-left: 0.75rem; background: transparent; border: none; color: #fff; cursor: pointer; font-size: 1.2rem; line-height: 1; padding: 0;"
      >×</button>
    </div>
  `;
  
  document.body.appendChild(notification);

  setTimeout(() => {
    notification.style.transform = 'translateX(0)';
  }, 10);

  setTimeout(() => {
    notification.style.transform = 'translateX(100%)';
    setTimeout(() => {
      notification.remove();
    }, 300);
  }, durationMs);
}

