import { Auth } from './auth.js';

// Sistema de timeout de sesión automático
export const SessionTimeout = {
  timeoutMs: 15 * 60 * 1000, // 15 minutos en milisegundos
  warningMs: 1 * 60 * 1000, // 1 minuto antes de expirar para mostrar advertencia
  timerId: null,
  warningTimerId: null,
  isActive: false,
  eventThrottle: 1000, // Actualizar máximo una vez por segundo

  // Inicializar el sistema de timeout
  init: function() {
    if (!Auth.isAuthenticated()) {
      return;
    }

    // Verificar si la sesión ya expiró
    if (Auth.isSessionExpired()) {
      Auth.logout();
      return;
    }

    // Establecer timestamp inicial si no existe
    if (!Auth.getLastActivity()) {
      Auth.setLastActivity();
    }

    // Limpiar timers anteriores si existen
    this.stop();

    // Iniciar el timer
    this.start();

    // Configurar listeners de eventos de usuario
    this.setupEventListeners();

    this.isActive = true;
  },

  // Iniciar el timer de timeout
  start: function() {
    const self = this;
    
    // Calcular tiempo restante
    const lastActivity = Auth.getLastActivity();
    const now = Date.now();
    const elapsed = now - lastActivity;
    const remaining = Math.max(0, this.timeoutMs - elapsed);

    // Timer para mostrar advertencia
    const warningTime = Math.max(0, remaining - this.warningMs);
    if (warningTime > 0) {
      this.warningTimerId = setTimeout(() => {
        self.showWarning();
      }, warningTime);
    }

    // Timer principal para cerrar sesión
    this.timerId = setTimeout(() => {
      self.handleTimeout();
    }, remaining);
  },

  // Detener todos los timers
  stop: function() {
    if (this.timerId) {
      clearTimeout(this.timerId);
      this.timerId = null;
    }
    if (this.warningTimerId) {
      clearTimeout(this.warningTimerId);
      this.warningTimerId = null;
    }
    this.isActive = false;
  },

  // Reiniciar el timer (llamado cuando hay actividad del usuario)
  reset: function() {
    if (!this.isActive || !Auth.isAuthenticated()) {
      return;
    }

    // Actualizar timestamp de última actividad
    Auth.setLastActivity();

    // Reiniciar timers
    this.stop();
    this.start();
  },

  // Configurar listeners de eventos de usuario
  setupEventListeners: function() {
    const self = this;
    let lastUpdate = 0;

    // Función throttled para actualizar actividad
    const updateActivity = () => {
      const now = Date.now();
      if (now - lastUpdate >= self.eventThrottle) {
        lastUpdate = now;
        self.reset();
      }
    };

    // Eventos a monitorear
    const events = ['mousemove', 'keydown', 'click', 'scroll', 'touchstart', 'mousedown'];
    
    events.forEach(eventType => {
      document.addEventListener(eventType, updateActivity, { passive: true });
    });

    // Guardar referencia a los listeners para poder removerlos después
    this.eventListeners = events.map(eventType => ({
      type: eventType,
      handler: updateActivity
    }));
  },

  // Remover listeners de eventos
  removeEventListeners: function() {
    if (this.eventListeners) {
      this.eventListeners.forEach(({ type, handler }) => {
        document.removeEventListener(type, handler);
      });
      this.eventListeners = null;
    }
  },

  // Mostrar advertencia antes de cerrar sesión
  showWarning: function() {
    if (typeof showNotification !== 'undefined') {
      showNotification(
        'Tu sesión expirará en 1 minuto por inactividad. Mueve el mouse o presiona una tecla para continuar.',
        'warning'
      );
    }
  },

  // Manejar el timeout y cerrar sesión
  handleTimeout: function() {
    this.stop();
    this.removeEventListeners();
    
    if (typeof window.showNotification !== 'undefined') {
      window.showNotification(
        'Tu sesión ha expirado por inactividad. Serás redirigido al login.',
        'error'
      );
    }

    // Esperar un momento para que el usuario vea el mensaje
    setTimeout(() => {
      Auth.logout();
    }, 2000);
  },

  // Destruir el sistema de timeout
  destroy: function() {
    this.stop();
    this.removeEventListeners();
  }
};

