// Clase base abstracta para Cards - Permite polimorfismo y herencia
export class BaseCard {
  constructor(data) {
    if (this.constructor === BaseCard) {
      throw new Error('BaseCard es una clase abstracta y no puede ser instanciada directamente');
    }
    this.data = data;
  }

  // Método abstracto que debe ser implementado por las clases hijas
  render() {
    throw new Error('El método render() debe ser implementado por la clase hija');
  }

  // Método común para obtener el HTML base de una card
  getBaseCardHTML(content) {
    return `
      <div class="card">
        ${content}
      </div>
    `;
  }

  // Método común para crear badges de estado
  createStatusBadge(status, text) {
    return `<span class="status-badge ${status}">${text}</span>`;
  }

  // Método común para crear botones
  createButton(text, href, className, onClick = null) {
    const onClickAttr = onClick ? `onclick="${onClick}"` : '';
    const tag = href ? 'a' : 'button';
    const hrefAttr = href ? `href="${href}"` : '';
    
    return `<${tag} ${hrefAttr} ${onClickAttr} class="card-button ${className}">${text}</${tag}>`;
  }
}

// Clase específica para Cards de Formularios - Hereda de BaseCard
export class FormCard extends BaseCard {
  constructor(formData) {
    super(formData);
    this.form = formData;
  }

  // Implementación del método render usando polimorfismo
  render() {
    const status = this.getFormStatus();
    const statusText = this.getStatusText(status);
    const statusColor = this.getStatusColor(status);

    return this.getBaseCardHTML(`
      ${this.renderHeader(statusText, statusColor)}
      ${this.renderDescription()}
      ${this.renderMeta()}
      ${this.renderFooter()}
    `);
  }

  // Métodos privados para construir las secciones de la card
  renderHeader(statusText, statusColor) {
    return `
      <div class="card-header">
        <h3 class="card-title">${this.form.form_name}</h3>
        ${this.createStatusBadge(statusColor, statusText)}
      </div>
    `;
  }

  renderDescription() {
    if (!this.form.description) return '';
    return `<p class="card-description">${this.form.description}</p>`;
  }

  renderMeta() {
    const conf = this.form.conferencista ? String(this.form.conferencista).trim() : '';
    const ubi = this.form.ubicacion ? String(this.form.ubicacion).trim() : '';
    const extra = [
      conf ? `<div class="card-meta-item"><span>Conferencista:</span><span class="card-meta-value">${this.escapeHtml(conf)}</span></div>` : '',
      ubi ? `<div class="card-meta-item"><span>Ubicación:</span><span class="card-meta-value">${this.escapeHtml(ubi)}</span></div>` : ''
    ].join('');
    return `
      <div class="card-content">
        <div class="card-meta">
          <div class="card-meta-item">
            <span>Creado:</span>
            <span class="card-meta-value">${this.formatDate(this.form.created_at)}</span>
          </div>
          <div class="card-meta-item">
            <span>Apertura:</span>
            <span class="card-meta-value">${this.formatDate(this.form.open_at)}</span>
          </div>
          <div class="card-meta-item">
            <span>Evento:</span>
            <span class="card-meta-value">${this.formatDate(this.form.event_date)}</span>
          </div>
          <div class="card-meta-item">
            <span>Cierre:</span>
            <span class="card-meta-value">${this.formatDate(this.form.expires_at)}</span>
          </div>
          ${extra}
        </div>
      </div>
    `;
  }

  escapeHtml(s) {
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }

  renderFooter() {
    const status = this.getFormStatus();
    const selectorValue = this.form.is_closed ? 'closed'
      : this.form.is_draft ? 'draft' : 'open';

    return `
      <div class="card-footer">
        ${this.createButton(
          'Ver Detalles',
          `actions/view-form/?id=${this.form.id}`,
          'card-button-primary'
        )}
        <div class="status-dropdown" data-form-id="${this.form.id}">
          <button type="button" class="card-button card-button-status" aria-haspopup="true" aria-expanded="false">
            Cambiar estado ▾
          </button>
          <ul class="status-dropdown__menu" role="menu">
            <li role="menuitem" data-value="open" class="${selectorValue === 'open' ? 'active' : ''}">Abierto</li>
            <li role="menuitem" data-value="closed" class="${selectorValue === 'closed' ? 'active' : ''}">Cerrado</li>
            <li role="menuitem" data-value="draft" class="${selectorValue === 'draft' ? 'active' : ''}">Borrador</li>
          </ul>
        </div>
      </div>
    `;
  }

  // Métodos auxiliares
  getFormStatus() {
    const now = new Date();
    // Parse de DATETIME del backend sin corrimiento UTC
    const parseBackendDate = (dateValue) => {
      if (!dateValue) return null;
      if (dateValue instanceof Date) return dateValue;
      if (typeof dateValue !== 'string') return new Date(dateValue);

      let s = dateValue.trim();
      if (/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}/.test(s)) {
        s = s.replace(' ', 'T');
      }
      s = s.replace(/Z$/, '');
      s = s.replace(/([+-]\d{2}:\d{2})$/, '');
      return new Date(s);
    };

    const expiresAt = this.form.expires_at ? parseBackendDate(this.form.expires_at) : null;
    const openAt = this.form.open_at ? parseBackendDate(this.form.open_at) : null;

    if (this.form.is_closed) {
      return 'closed';
    }
    if (this.form.is_draft) {
      return 'draft';
    }
    if (openAt && openAt > now) {
      return 'future_open';
    } else if (expiresAt < now) {
      return 'expired';
    }
    return 'active';
  }

  getStatusText(status) {
    const statusMap = {
      active: 'Activo',
      draft: 'Borrador',
      future_open: 'Próxima apertura',
      closed: 'Cerrado',
      expired: 'Expirado'
    };
    return statusMap[status] || 'Desconocido';
  }

  getStatusColor(status) {
    return status;
  }

  formatDate(dateString) {
    // Usar la función formatDate del módulo utils si está disponible
    if (window.formatDate) {
      return window.formatDate(dateString);
    }
    // Fallback básico
    return new Date(dateString).toLocaleDateString('es-ES');
  }
}

// Clase Renderer para manejar la renderización de múltiples cards
export class CardRenderer {
  constructor(containerId, CardClass) {
    this.container = document.getElementById(containerId);
    this.CardClass = CardClass; // Permite polimorfismo - puede recibir cualquier clase que extienda BaseCard
  }

  renderCards(dataArray) {
    if (!this.container) {
      console.error(`Contenedor con ID "${this.containerId}" no encontrado`);
      return;
    }

    if (!dataArray || dataArray.length === 0) {
      this.renderEmpty();
      return;
    }

    const cardsHTML = dataArray.map(data => {
      const card = new this.CardClass(data);
      return card.render();
    }).join('');

    this.container.innerHTML = cardsHTML;
    this._bindStatusDropdowns();
  }

  _bindStatusDropdowns() {
    if (!this.container) return;

    this.container.addEventListener('click', (e) => {
      const toggleBtn = e.target.closest('.card-button-status');
      const menuItem = e.target.closest('.status-dropdown__menu li');

      if (toggleBtn) {
        e.stopPropagation();
        const dropdown = toggleBtn.closest('.status-dropdown');
        const wasOpen = dropdown.classList.contains('open');
        this.container.querySelectorAll('.status-dropdown.open').forEach(d => d.classList.remove('open'));
        if (!wasOpen) dropdown.classList.add('open');
        return;
      }

      if (menuItem) {
        e.stopPropagation();
        const dropdown = menuItem.closest('.status-dropdown');
        const formId = parseInt(dropdown.dataset.formId);
        const targetValue = menuItem.dataset.value;
        dropdown.classList.remove('open');
        this._changeFormStatus(formId, targetValue, dropdown);
        return;
      }
    });

    document.addEventListener('click', () => {
      if (!this.container) return;
      this.container.querySelectorAll('.status-dropdown.open').forEach(d => d.classList.remove('open'));
    });
  }

  async _changeFormStatus(formId, target, dropdown) {
    if (!window.API_CONFIG?.ENDPOINTS || !window.apiRequest) return;

    const EP = window.API_CONFIG.ENDPOINTS;
    const endpoint = target === 'closed' ? EP.FORM_CLOSE(formId)
      : target === 'draft' ? EP.FORM_DRAFT(formId)
      : EP.FORM_OPEN(formId);

    const btn = dropdown.querySelector('.card-button-status');
    btn.textContent = 'Actualizando...';
    btn.disabled = true;

    try {
      const res = await window.apiRequest(endpoint, { method: 'POST' });
      if (res?.success) {
        if (window.showNotification) window.showNotification(res.message || 'Estado actualizado.', 'success');
        if (window.formsManager) window.formsManager.loadForms();
      } else {
        if (window.showNotification) window.showNotification(res?.message || 'No se pudo cambiar el estado.', 'error');
        btn.textContent = 'Cambiar estado ▾';
        btn.disabled = false;
      }
    } catch (err) {
      if (window.showNotification) window.showNotification('Error: ' + err.message, 'error');
      btn.textContent = 'Cambiar estado ▾';
      btn.disabled = false;
    }
  }

  renderLoading() {
    if (!this.container) return;
    this.container.innerHTML = `
      <div class="forms-loading">
        <div class="forms-loading-spinner">
          <div class="forms-loading-spinner-inner"></div>
        </div>
        <p>Cargando formularios...</p>
      </div>
    `;
  }

  renderEmpty() {
    if (!this.container) return;
    this.container.innerHTML = `
      <div class="forms-empty">
        <p>No hay formularios. Crea uno nuevo para comenzar.</p>
      </div>
    `;
  }

  renderError(message) {
    if (!this.container) return;
    this.container.innerHTML = `
      <div class="forms-error">
        <p>Error al cargar formularios: ${message}</p>
      </div>
    `;
  }
}

// Clase principal para gestionar la lista de formularios
export class FormsManager {
  constructor() {
    this.forms = [];
    this.renderer = new CardRenderer('formsContainer', FormCard);
  }

  async loadForms() {
    try {
      this.renderer.renderLoading();
      
      // Verificar que API_CONFIG y ENDPOINTS estén disponibles
      if (!window.API_CONFIG) {
        throw new Error('API_CONFIG no está disponible. Verifica que los módulos se hayan cargado correctamente.');
      }
      if (!window.API_CONFIG.ENDPOINTS) {
        throw new Error('API_CONFIG.ENDPOINTS no está disponible. Verifica la configuración.');
      }
      if (!window.apiRequest) {
        throw new Error('apiRequest no está disponible. Verifica que los módulos se hayan cargado correctamente.');
      }
      
      const endpoint = window.API_CONFIG.ENDPOINTS.FORMS;
      console.log('[FormsManager] Cargando formularios desde:', endpoint);
      const response = await window.apiRequest(endpoint);
      
      if (response.success && response.forms) {
        this.forms = response.forms;
        this.renderer.renderCards(this.forms);
      } else {
        this.renderer.renderEmpty();
      }
    } catch (error) {
      console.error('Error cargando formularios:', error);
      if (window.showNotification) {
        window.showNotification('Error cargando formularios: ' + error.message, 'error');
      } else {
        console.error('showNotification no está disponible');
      }
      this.renderer.renderError(error.message);
    }
  }

  getForms() {
    return this.forms;
  }
}

// Inicialización cuando el DOM esté listo
let formsManager = null;

export function initFormsPage() {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      formsManager = new FormsManager();
      formsManager.loadForms();
      // Hacer disponible globalmente para recargar después de crear formularios
      window.formsManager = formsManager;
    });
  } else {
    formsManager = new FormsManager();
    formsManager.loadForms();
    // Hacer disponible globalmente para recargar después de crear formularios
    window.formsManager = formsManager;
  }
}

// Función global para recargar formularios (usada después de crear/cerrar formularios)
export function reloadForms() {
  if (window.formsManager) {
    window.formsManager.loadForms();
  } else if (window.loadForms) {
    // Fallback al sistema antiguo
    window.loadForms();
  }
}

// Exportar para uso global si es necesario
window.FormsManager = FormsManager;
window.FormCard = FormCard;
window.BaseCard = BaseCard;
window.CardRenderer = CardRenderer;
window.reloadForms = reloadForms;

