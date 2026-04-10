// Función helper para obtener SVG de Heroicons
// Soporta nombres de Heroicons directamente y mapeo desde FontAwesome

// Mapeo de nombres FontAwesome a nombres Heroicons
const FONTAWESOME_TO_HEROICONS = {
  'fa-user': 'user',
  'fa-envelope': 'envelope',
  'fa-phone': 'phone',
  'fa-map-marker-alt': 'map-pin',
  'fa-briefcase': 'briefcase',
  'fa-share-alt': 'share',
  'fa-shield-alt': 'shield-check',
  'fa-info-circle': 'information-circle',
  'fa-calendar': 'calendar',
  'fa-clock': 'clock',
  'fa-file-alt': 'document-text',
  'fa-image': 'photo',
  'fa-link': 'link',
  'fa-globe': 'globe-alt',
  'fa-building': 'building-office',
  'fa-users': 'users',
  'fa-user-tie': 'user',
  'fa-graduation-cap': 'academic-cap',
  'fa-heart': 'heart',
  'fa-star': 'star',
  'fa-home': 'home',
  'fa-car': 'truck',
  'fa-plane': 'paper-airplane',
  'fa-shopping-cart': 'shopping-cart',
  'fa-credit-card': 'credit-card',
  'fa-lock': 'lock-closed',
  'fa-unlock': 'lock-open',
  'fa-key': 'key',
  'fa-bell': 'bell',
  'fa-cog': 'cog-6-tooth',
  'fa-check-circle': 'check-circle',
  'fa-times-circle': 'x-circle',
  'fa-question-circle': 'question-mark-circle',
  'fa-exclamation-circle': 'exclamation-circle',
  'fa-arrow-right': 'arrow-right',
  'fa-arrow-left': 'arrow-left',
  'fa-arrow-up': 'arrow-up',
  'fa-arrow-down': 'arrow-down',
  'fa-search': 'magnifying-glass',
  'fa-filter': 'funnel',
  'fa-sort': 'bars-arrow-down',
  'fa-list': 'list-bullet',
  'fa-th': 'squares-2x2',
  'fa-edit': 'pencil',
  'fa-trash': 'trash',
  'fa-save': 'document-arrow-down',
  'fa-download': 'arrow-down-tray',
  'fa-upload': 'arrow-up-tray'
};

// Función para normalizar el nombre del icono (mapear FontAwesome a Heroicons o usar directamente)
function normalizeIconName(iconName) {
  if (!iconName) return null;
  
  // Si empieza con 'fa-', buscar en el mapeo
  if (iconName.startsWith('fa-')) {
    return FONTAWESOME_TO_HEROICONS[iconName] || iconName.replace('fa-', '');
  }
  
  // Si ya es un nombre de heroicon, usar directamente
  return iconName;
}

// Función para cargar SVG desde archivo (async) - fallback local
async function loadHeroiconSVGFromFile(iconName) {
  try {
    const normalizedName = normalizeIconName(iconName);
    if (!normalizedName) return '';
    
    const response = await fetch(`/src/icons/heroicons-24-outline/${normalizedName}.svg`);
    if (!response.ok) return '';
    
    let svgContent = await response.text();
    // Reemplazar stroke="#0F172A" por stroke="currentColor" para que herede el color
    svgContent = svgContent.replace(/stroke="#[0-9A-Fa-f]{6}"/g, 'stroke="currentColor"');
    // Asegurar que tenga width="24" height="24" o ajustar según necesidad
    svgContent = svgContent.replace(/width="\d+"/, 'width="24"').replace(/height="\d+"/, 'height="24"');
    return svgContent;
  } catch (error) {
    console.warn(`Error cargando icono ${iconName}:`, error);
    return '';
  }
}

// Función para cargar SVG desde la web usando Iconify API
async function loadHeroiconSVGFromWeb(iconName) {
  try {
    const normalizedName = normalizeIconName(iconName);
    if (!normalizedName) return '';
    
    const response = await fetch(`https://api.iconify.design/heroicons/${normalizedName}.svg?color=currentColor`);
    if (!response.ok) return '';
    
    let svgContent = await response.text();
    // Asegurar que use currentColor y tamaño apropiado
    svgContent = svgContent.replace(/stroke="#[0-9A-Fa-f]{6}"/g, 'stroke="currentColor"').replace(/fill="#[0-9A-Fa-f]{6}"/g, 'fill="currentColor"');
    svgContent = svgContent.replace(/width="\d+"/, 'width="24"').replace(/height="\d+"/, 'height="24"');
    return svgContent;
  } catch (error) {
    console.warn(`Error cargando icono ${iconName} desde web:`, error);
    return '';
  }
}

// Función principal para obtener SVG de Heroicons
// Puede ser usada de forma síncrona (retorna SVG inline si está disponible) 
// o asíncrona (carga desde archivo)
export function getHeroiconSVG(iconName) {
  if (!iconName) return '';
  
  const normalizedName = normalizeIconName(iconName);
  if (!normalizedName) return '';
  
  // Algunos iconos comunes inline para evitar fetch (opcional, para performance)
  const inlineIcons = {
    'user': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M15.75 6C15.75 8.07107 14.071 9.75 12 9.75C9.9289 9.75 8.24996 8.07107 8.24996 6C8.24996 3.92893 9.9289 2.25 12 2.25C14.071 2.25 15.75 3.92893 15.75 6Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><path d="M4.5011 20.1182C4.5714 16.0369 7.90184 12.75 12 12.75C16.0982 12.75 19.4287 16.0371 19.4988 20.1185C17.216 21.166 14.6764 21.75 12.0003 21.75C9.32396 21.75 6.78406 21.1659 4.5011 20.1182Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'envelope': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M21.75 6.75V17.25C21.75 18.4926 20.7426 19.5 19.5 19.5H4.5C3.25736 19.5 2.25 18.4926 2.25 17.25V6.75M21.75 6.75C21.75 5.50736 20.7426 4.5 19.5 4.5H4.5C3.25736 4.5 2.25 5.50736 2.25 6.75M21.75 6.75V6.99271C21.75 7.77405 21.3447 8.49945 20.6792 8.90894L13.1792 13.5243C12.4561 13.9694 11.5439 13.9694 10.8208 13.5243L3.32078 8.90894C2.65535 8.49945 2.25 7.77405 2.25 6.99271V6.75" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'information-circle': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M11.25 11.25L11.2915 11.2293C11.8646 10.9427 12.5099 11.4603 12.3545 12.082L11.6455 14.918C11.4901 15.5397 12.1354 16.0573 12.7085 15.7707L12.75 15.75M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12ZM12 8.25H12.0075V8.2575H12V8.25Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'home': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M2.25 12L11.2045 3.04549C11.6438 2.60615 12.3562 2.60615 12.7955 3.04549L21.75 12M4.5 9.75V19.875C4.5 20.4963 5.00368 21 5.625 21H9.75V16.125C9.75 15.5037 10.2537 15 10.875 15H13.125C13.7463 15 14.25 15.5037 14.25 16.125V21H18.375C18.9963 21 19.5 20.4963 19.5 19.875V9.75M8.25 21H16.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'document-text': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M19.5 14.25V11.625C19.5 9.76104 17.989 8.25 16.125 8.25H14.625C14.0037 8.25 13.5 7.74632 13.5 7.125V5.625C13.5 3.76104 11.989 2.25 10.125 2.25H8.25M8.25 15H15.75M8.25 18H12M10.5 2.25H5.625C5.00368 2.25 4.5 2.75368 4.5 3.375V20.625C4.5 21.2463 5.00368 21.75 5.625 21.75H18.375C18.9963 21.75 19.5 21.2463 19.5 20.625V11.25C19.5 6.27944 15.4706 2.25 10.5 2.25Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'cog-6-tooth': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M9.594 3.94C9.684 3.39768 10.1533 3 10.7033 3H13.2972C13.8472 3 14.3165 3.39768 14.4069 3.94L14.6204 5.22119C14.6828 5.59523 14.9327 5.9068 15.2645 6.09045C15.3387 6.13151 15.412 6.17393 15.4844 6.21766C15.8095 6.41393 16.2048 6.47495 16.5604 6.34175L17.7772 5.88587C18.2922 5.69293 18.8712 5.9006 19.1462 6.37687L20.4432 8.6233C20.7181 9.09957 20.6085 9.70482 20.1839 10.0544L19.1795 10.8812C18.887 11.122 18.742 11.4938 18.7491 11.8726C18.7498 11.915 18.7502 11.9575 18.7502 12.0001C18.7502 12.0427 18.7498 12.0852 18.7491 12.1275C18.742 12.5064 18.887 12.8782 19.1795 13.119L20.1839 13.9458C20.6085 14.2953 20.7181 14.9006 20.4432 15.3769L19.1462 17.6233C18.8712 18.0996 18.2922 18.3072 17.7772 18.1143L16.5604 17.6584C16.2048 17.5252 15.8095 17.5862 15.4844 17.7825C15.412 17.8263 15.3387 17.8687 15.2645 17.9097C14.9327 18.0934 14.6828 18.4049 14.6204 18.779L14.4069 20.06C14.3165 20.6025 13.8472 21 13.2972 21H10.7033C10.1533 21 9.68397 20.6025 9.59356 20.06L9.38005 18.779C9.31771 18.4049 9.06774 18.0934 8.73597 17.9097C8.66179 17.8687 8.58847 17.8263 8.51604 17.7825C8.19101 17.5863 7.79568 17.5252 7.44011 17.6584L6.22325 18.1143C5.70826 18.3072 5.12926 18.0996 4.85429 17.6233L3.55731 15.3769C3.28234 14.9006 3.39199 14.2954 3.81657 13.9458L4.82092 13.119C5.11343 12.8782 5.25843 12.5064 5.25141 12.1276C5.25063 12.0852 5.25023 12.0427 5.25023 12.0001C5.25023 11.9575 5.25063 11.915 5.25141 11.8726C5.25843 11.4938 5.11343 11.122 4.82092 10.8812L3.81657 10.0544C3.39199 9.70484 3.28234 9.09958 3.55731 8.62332L4.85429 6.37688C5.12926 5.90061 5.70825 5.69295 6.22325 5.88588L7.4401 6.34176C7.79566 6.47496 8.19099 6.41394 8.51603 6.21767C8.58846 6.17393 8.66179 6.13151 8.73597 6.09045C9.06774 5.9068 9.31771 5.59523 9.38005 5.22119L9.59356 3.94014Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><path d="M15 12C15 13.6569 13.6569 15 12 15C10.3431 15 9 13.6569 9 12C9 10.3432 10.3431 9 12 9C13.6569 9 15 10.3432 15 12Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'chevron-up': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M4.5 15.75L12 8.25L19.5 15.75" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'chevron-down': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M19.5 8.25L12 15.75L4.5 8.25" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'pencil': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M16.8617 4.48667L18.5492 2.79917C19.2814 2.06694 20.4686 2.06694 21.2008 2.79917C21.9331 3.53141 21.9331 4.71859 21.2008 5.45083L6.83218 19.8195C6.30351 20.3481 5.65144 20.7368 4.93489 20.9502L2.25 21.75L3.04978 19.0651C3.26323 18.3486 3.65185 17.6965 4.18052 17.1678L16.8617 4.48667ZM16.8617 4.48667L19.5 7.12499" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    'trash': `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M14.7404 9L14.3942 18M9.60577 18L9.25962 9M19.2276 5.79057C19.5696 5.84221 19.9104 5.89747 20.25 5.95629M19.2276 5.79057L18.1598 19.6726C18.0696 20.8448 17.0921 21.75 15.9164 21.75H8.08357C6.90786 21.75 5.93037 20.8448 5.8402 19.6726L4.77235 5.79057M19.2276 5.79057C18.0812 5.61744 16.9215 5.48485 15.75 5.39432M3.75 5.95629C4.08957 5.89747 4.43037 5.84221 4.77235 5.79057M4.77235 5.79057C5.91878 5.61744 7.07849 5.48485 8.25 5.39432M15.75 5.39432V4.47819C15.75 3.29882 14.8393 2.31423 13.6606 2.27652C13.1092 2.25889 12.5556 2.25 12 2.25C11.4444 2.25 10.8908 2.25889 10.3394 2.27652C9.16065 2.31423 8.25 3.29882 8.25 4.47819V5.39432M15.75 5.39432C14.5126 5.2987 13.262 5.25 12 5.25C10.738 5.25 9.48744 5.2987 8.25 5.39432" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>`
  };
  
  // Si está en iconos inline, retornarlo
  if (inlineIcons[normalizedName]) {
    return inlineIcons[normalizedName];
  }
  
  // Si no está inline, retornar string vacío (el código deberá usar loadHeroiconSVGFromFile)
  return '';
}

// Exportar función async para cargar desde archivo y desde web
export { loadHeroiconSVGFromFile, loadHeroiconSVGFromWeb, normalizeIconName, FONTAWESOME_TO_HEROICONS };
