// Función para crear el header de página
export function createPageHeader(title, subtitle = '', action = '') {
  return `
    <div>
      <div>
        <h1>${title}</h1>
        ${subtitle ? `<p>${subtitle}</p>` : ''}
      </div>
      ${action ? `<div>${action}</div>` : ''}
    </div>
  `;
}

