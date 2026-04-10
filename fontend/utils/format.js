// Funciones de formateo
export function formatNumber(num) {
  return new Intl.NumberFormat('es-PA').format(num);
}

function parseBackendDate(date) {
  if (!date) return null;
  if (date instanceof Date) return date;

  if (typeof date === 'string') {
    let s = date.trim();

    // Convertir "YYYY-MM-DD HH:MM:SS" a "YYYY-MM-DDTHH:MM:SS"
    if (/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}/.test(s)) {
      s = s.replace(' ', 'T');
    }

    // Si viene con sufijo UTC ('Z') o con offset (+/-HH:MM), lo quitamos para tratarlo como hora local.
    // Esto evita el corrimiento cuando el backend marca DATETIME como UTC.
    s = s.replace(/Z$/, '');
    s = s.replace(/([+-]\d{2}:\d{2})$/, '');

    const d = new Date(s);
    return d;
  }

  const d = new Date(date);
  return d;
}

export function formatDate(date) {
  if (!date) return '-';
  const dateObj = parseBackendDate(date);
  if (!dateObj || isNaN(dateObj.getTime())) return '-';
  return dateObj.toLocaleDateString('es-PA', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
}

// Exportar a window para compatibilidad global
window.formatNumber = formatNumber;
window.formatDate = formatDate;
