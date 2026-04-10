// Utilidades de estado
const STATUS_TEXT = {
  active: 'Activo',
  expired: 'Expirado',
  closed: 'Cerrado',
  pending: 'Pendiente',
  sent: 'Enviado',
  failed: 'Fallido'
};

export function getStatusText(status) {
  return STATUS_TEXT[status] || status;
}

export function getStatusColor(status) {
  // Función placeholder para futura implementación
  return '';
}

// Exportar a window para compatibilidad global
window.getStatusText = getStatusText;
window.getStatusColor = getStatusColor;
