// Función para crear una tabla
export function createTable(headers, rows, actions = []) {
  return `
    <div>
      <table>
        <thead>
          <tr>
            ${headers.map(header => `<th>${header}</th>`).join('')}
            ${actions.length > 0 ? `<th>Acciones</th>` : ''}
          </tr>
        </thead>
        <tbody>
          ${rows.map(row => `
            <tr>
              ${row.map(cell => `<td>${cell}</td>`).join('')}
              ${actions.length > 0 ? `<td>${actions.map(action => action).join('')}</td>` : ''}
            </tr>
          `).join('')}
        </tbody>
      </table>
    </div>
  `;
}

