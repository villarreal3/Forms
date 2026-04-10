# Recursos compartidos del panel admin

Scripts y estilos reutilizados por varias pantallas (`dashboard`, `formularios`, `config`, `email`, `respuestas`, etc.):

| Archivo | Rol |
|---------|-----|
| `sidebar.css` | Estilos del menú lateral y botón **nav-toggle** (móvil). |
| `nav-toggle.js` | Abre/cierra la navegación en pantallas pequeñas. |
| `formularios.js` | Lógica de detalle/edición de formulario (tabs, modales de campos, etc.). |
| `form_sections.js` | Editor de secciones y campos del schema. |
| `customization.js` | Personalización visual (colores, logos, franja meta). |

Las páginas cargan estos archivos con rutas relativas (`../shared/…`, `../../shared/…`, etc.).
