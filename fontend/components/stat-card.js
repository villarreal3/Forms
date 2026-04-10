// Función para crear una stat card
export function createStatCard(title, value, icon, color = 'blue', change = '') {
  const colorClasses = {
    blue: 'from-blue-700 to-blue-900',
    green: 'from-green-600 to-green-700',
    yellow: 'from-yellow-500 to-yellow-600',
    red: 'from-red-500 to-red-600',
    purple: 'from-blue-700 to-blue-900'
  };

  return `
    <div>
      <h3>${title}</h3>
      <div>
        
      </div>
      <div>
        <p>${value}</p>
        ${change ? `<span>${change}</span>` : ''}
      </div>
    </div>
  `;
}

