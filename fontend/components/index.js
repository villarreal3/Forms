// Barrel export - Re-exportar todos los componentes
export { getHeroiconSVG } from './icons.js';
export { createNavbar } from './sidebar.js';
export { createPageHeader } from './header.js';
export { createCard } from './card.js';
export { createButton } from './button.js';
export { createStatCard } from './stat-card.js';
export { createTable } from './table.js';
export { createModal } from './modal.js';
export { showNotification } from './notification.js';

// Exportar también a window para compatibilidad con código inline en HTML
import { getHeroiconSVG } from './icons.js';
import { createNavbar } from './sidebar.js';
import { createPageHeader } from './header.js';
import { createCard } from './card.js';
import { createButton } from './button.js';
import { createStatCard } from './stat-card.js';
import { createTable } from './table.js';
import { createModal } from './modal.js';
import { showNotification } from './notification.js';

window.getHeroiconSVG = getHeroiconSVG;
window.createNavbar = createNavbar;
window.createPageHeader = createPageHeader;
window.createCard = createCard;
window.createButton = createButton;
window.createStatCard = createStatCard;
window.createTable = createTable;
window.createModal = createModal;
window.showNotification = showNotification;
