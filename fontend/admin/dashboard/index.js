// Importar funciones necesarias
import { apiRequest, API_CONFIG } from '../services/index.js';
import { formatNumber, formatDate } from '../utils/format.js';
import { showNotification } from '../components/index.js';

// Cargar gráfico de respuestas en el tiempo
async function loadResponseTimeChart(days = 30) {
 try {
 const response = await apiRequest(`${API_CONFIG.ENDPOINTS.DASHBOARD_RESPONSES_TIMELINE}?days=${days}`);
 const container = document.getElementById('chartDataContainer');
 
 if (!response.success || !response.data || response.data.length === 0) {
 container.innerHTML = '<p>No hay datos disponibles</p>';
 return;
 }
 
 // Limitar a los primeros 10 resultados
 const limitedData = response.data.slice(0, 10);
 
 // Crear estructura HTML tipo tarjeta
 if (limitedData.length === 0) {
 container.innerHTML = '<p>No hay datos disponibles</p>';
 return;
 }
 
 const html = '<div class="response-time-cards">' +
 limitedData.map(item => `
 <div class="response-time-card">
 <div class="response-time-date">${item.date}</div>
 <div class="response-time-count">${formatNumber(item.count)}</div>
 </div>
 `).join('') +
 '</div>';
 
 container.innerHTML = html;
 } catch (error) {
 console.error('Error cargando gráfico de respuestas:', error);
 document.getElementById('chartDataContainer').innerHTML = '<p>Error cargando datos</p>';
 }
}

// Cargar formularios más activos
async function loadTopActiveForms() {
 try {
 const response = await apiRequest(`${API_CONFIG.ENDPOINTS.DASHBOARD_TOP_FORMS}?limit=10`);
 const container = document.getElementById('formsRankingContainer');
 
 if (!response.success || !response.forms || response.forms.length === 0) {
 container.innerHTML = '<p>No hay formularios disponibles</p>';
 return;
 }
 
 // Limitar a los primeros 10 resultados
 const limitedForms = response.forms.slice(0, 10);
 
 // Debug: ver qué campos tiene el primer formulario
 if (limitedForms.length > 0) {
   console.log('Primer formulario recibido:', limitedForms[0]);
 }
 
 // Función helper para obtener la fecha de creación
 const getCreatedDate = (form) => {
   // Intentar diferentes campos posibles para la fecha
   const dateField = form.created_at || form.createdAt || form.date_created || form.dateCreated || form.created || form.created_date;
   if (dateField) {
     return formatDate(dateField);
   }
   return null;
 };
 
 // Crear estructura HTML tipo tarjeta
 const html = '<div class="forms-ranking-cards">' +
 limitedForms.map((form, index) => {
   const createdDate = getCreatedDate(form);
   const dateDisplay = createdDate ? `Creado: ${createdDate}` : '';
   
   return `
 <div class="form-ranking-card">
   <div class="form-ranking-header">
     <span class="form-ranking-number">#${index + 1}</span>
     <div class="form-ranking-info">
       <h3 class="form-ranking-name">${form.form_name || 'Sin nombre'}</h3>
       ${dateDisplay ? `<p class="form-ranking-date">${dateDisplay}</p>` : ''}
     </div>
   </div>
   <div class="form-ranking-count">
     <span class="form-ranking-count-value">${formatNumber(form.submission_count || 0)}</span>
     <span class="form-ranking-count-label">respuestas</span>
   </div>
 </div>
 `;
 }).join('') +
 '</div>';
 
 container.innerHTML = html;
 } catch (error) {
 console.error('Error cargando formularios activos:', error);
 document.getElementById('formsRankingContainer').innerHTML = '<p>Error cargando datos</p>';
 }
}

// Cargar distribución geográfica
async function loadGeographicDistribution() {
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.DASHBOARD_GEOGRAPHIC);
 const container = document.getElementById('provincesContainer');
 
 if (!response.success || !response.provinces || response.provinces.length === 0) {
 container.innerHTML = '<p>No hay datos de provincias disponibles</p>';
 return;
 }
 
 // Crear estructura HTML tipo tarjeta (mostrar todos los resultados)
 const html = '<div class="province-cards">' +
 response.provinces.map(item => `
 <div class="province-card">
 <div class="province-card-content">
 <div class="province-name">${item.province || 'No especificada'}</div>
 <div class="province-count">${formatNumber(item.count)}</div>
 </div>
 </div>
 `).join('') +
 '</div>';
 
 container.innerHTML = html;
 } catch (error) {
 console.error('Error cargando distribución geográfica:', error);
 document.getElementById('provincesContainer').innerHTML = '<p>Error cargando datos</p>';
 }
}

// Cargar métricas de usuarios
async function loadUserMetrics() {
 try {
 const response = await apiRequest(API_CONFIG.ENDPOINTS.DASHBOARD_USER_METRICS);
 
 if (!response.success || !response.metrics) {
 document.getElementById('newUsersThisMonth').textContent = '-';
 document.getElementById('returningUsers').textContent = '-';
 document.getElementById('inactiveUsers').textContent = '-';
 return;
 }
 
 const metrics = response.metrics;
 document.getElementById('newUsersThisMonth').textContent = formatNumber(metrics.new_users_month || 0);
 document.getElementById('returningUsers').textContent = formatNumber(metrics.returning_users || 0);
 document.getElementById('inactiveUsers').textContent = formatNumber(metrics.inactive_users || 0);
 } catch (error) {
 console.error('Error cargando métricas de usuarios:', error);
 document.getElementById('newUsersThisMonth').textContent = '-';
 document.getElementById('returningUsers').textContent = '-';
 document.getElementById('inactiveUsers').textContent = '-';
 }
}

// Función principal para cargar el dashboard
async function loadDashboard() {
 try {
 // Cargar gráfico de respuestas (por defecto 30 días)
 await loadResponseTimeChart(30);
 
 // Cargar formularios más activos
 await loadTopActiveForms();
 
 // Cargar distribución geográfica
 await loadGeographicDistribution();
 
 // Cargar métricas de usuarios
 await loadUserMetrics();
 } catch (error) {
 console.error('Error cargando dashboard:', error);
 showNotification('Error cargando el dashboard: ' + error.message, 'error');
 }
}

// Event listener para cambiar el período del gráfico
document.addEventListener('DOMContentLoaded', function() {
 const timePeriodSelector = document.getElementById('timePeriodSelector');
 if (timePeriodSelector) {
 timePeriodSelector.addEventListener('change', function() {
 const days = parseInt(this.value, 10);
 loadResponseTimeChart(days);
 });
 }
});

// Cargar dashboard cuando la página esté lista
if (document.readyState === 'loading') {
 document.addEventListener('DOMContentLoaded', loadDashboard);
} else {
 loadDashboard();
}
