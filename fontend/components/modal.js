// Función para crear un modal
export function createModal(id, title, content, footer = '') {
  return `
    <div id="${id}" style="display: none;">
      <div>
        <div onclick="document.getElementById('${id}').style.display = 'none'"></div>
        
        <div>
          <div>
            <h3>${title}</h3>
            <button onclick="document.getElementById('${id}').style.display = 'none'">
              
            </button>
          </div>
          <div>
            ${content}
          </div>
          ${footer ? `<div>${footer}</div>` : ''}
        </div>
      </div>
    </div>
  `;
}

