// Función para crear un botón
export function createButton(text, type = 'primary', icon = '', onclick = '', size = 'md') {
  const baseClasses = 'inline-flex items-center justify-center font-semibold rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2';
  
  const typeClasses = {
    primary: 'bg-gradient-to-r from-blue-700 to-blue-900 text-white hover:shadow-lg transform hover:-translate-y-0.5 focus:ring-blue-500',
    secondary: 'bg-gray-100 text-gray-700 hover:bg-gray-200 focus:ring-gray-500',
    success: 'bg-green-500 text-white hover:bg-green-600 focus:ring-green-500',
    danger: 'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500',
    outline: 'border-2 border-gray-300 text-gray-700 hover:bg-gray-50 focus:ring-gray-500'
  };
  
  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base'
  };

  const iconHtml = icon ? `` : '';
  const onclickAttr = onclick ? `onclick="${onclick}"` : '';

  return `
    <button type="button" ${onclickAttr}>
      ${iconHtml}${text}
    </button>
  `;
}

