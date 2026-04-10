// Script común para inicializar el nav toggle en todas las páginas
(function() {
  function initNavToggle() {
    const toggle = document.querySelector('.nav-toggle');
    const nav = document.querySelector('.navigation');
    
    if (!toggle || !nav) {
      // Reintentar después de un breve delay si el navbar aún no está listo
      setTimeout(initNavToggle, 100);
      return;
    }
  
    toggle.addEventListener('click', () => {
      const isOpen = nav.classList.toggle('active');
      toggle.classList.toggle('active', isOpen);
      toggle.setAttribute('aria-expanded', isOpen);
    });
  
    nav.addEventListener('click', e => {
      if (e.target === nav) {
        nav.classList.remove('active');
        toggle.classList.remove('active');
        toggle.setAttribute('aria-expanded', false);
      }
    });
  }
  
  // Inicializar cuando el DOM esté listo
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initNavToggle);
  } else {
    // Si el DOM ya está listo, esperar un poco para que el navbar se cree
    setTimeout(initNavToggle, 100);
  }
})();

