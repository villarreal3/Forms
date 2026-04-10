# Forms — Frontend

Aplicación **estática** (HTML, CSS, JavaScript) para:

- **Público:** catálogo de formularios e inscripción.
- **Administración:** dashboard, formularios, respuestas, exportación, configuración.

| Recurso | Ubicación |
|---------|-----------|
| **API backend** | [forms-api README](../../forms-api/README.md) |
| **Base de datos** | [forms-api/SQL/README.md](../../forms-api/SQL/README.md) |

---

## Convención de rutas

Cada pantalla vive en **una carpeta con `index.html`** (p. ej. `admin/login/index.html`, `form/index.html`).

---

## Ejecución local

Servir la carpeta `forms/` con cualquier servidor HTTP (la app usa `import()` de módulos ES; conviene no abrir el HTML como `file://`).

```bash
cd fontend
python3 -m http.server 3001
# Abrir http://127.0.0.1:3001/
```

---

## Configuración de APIs

Todo centralizado en **`services/config.js`**:

| Concepto | Detalle |
|----------|---------|
| **FORMS_API_BASE_URL** | Base hacia la API de formularios (`/api`). |
| **USERS_API_BASE_URL** | Base hacia users-api (login, usuarios). |
| **FORMS_ENDPOINTS** / **USERS_ENDPOINTS** | Rutas relativas agrupadas. |

**Entornos:**

- **Desarrollo** (localhost en la lista de hosts de dev): suele apuntar a `http://localhost:8080/api` (forms) y al puerto de users-api.
- **Producción:** URLs relativas a **proxies PHP** (`proxy-forms-api.php`, `proxy-users-api.php`) para mismo origen HTTPS y reenvío del header `Authorization`.

Ajustar hosts de “producción” vs “dev” en `isProductionHost()` dentro de `config.js` si cambia el despliegue.

---

## Páginas principales

### Público

| Ruta | Uso |
|------|-----|
| `index.html` | Listado público de formularios; enlaces e inscripción; modal **código QR** (URL canónica **`form/?id=…`** sin `index.html`; ver `.htaccess`). |
| `form/index.html` | Archivo físico del formulario público; en el navegador se usa **`/forms/form/?id=`**. |
| `hello/index.html` | Página de prueba / bienvenida (si aplica). |
| `test-connection/index.html` | Comprobación de conectividad a APIs. |

URLs sin extensión: ver **`.htaccess`** (Apache `mod_rewrite`).

### Administración

| Área | Entrada típica |
|------|----------------|
| Login | `admin/login/index.html` |
| Dashboard | `admin/dashboard/index.html` |
| Formularios | `admin/formularios/index.html` |
| Crear formulario | `admin/formularios/actions/create/index.html` |
| Ver / editar formulario | `admin/formularios/actions/view-form/index.html?id=` |
| Respuestas | `admin/respuestas/index.html` |
| Correos masivos | `admin/email/index.html` |
| Config / changelog / usuarios | `admin/config/index.html`, `admin/config/changelog/`, `admin/config/users/` |

---

## Proxies PHP

| Archivo | Función |
|---------|---------|
| `proxy-forms-api.php` | Reenvía peticiones al backend Go (públicas y admin con JWT). |
| `proxy-users-api.php` | Reenvía al servicio de usuarios. |
| `proxy-forms-images.php` | Imágenes de formularios (si está en uso). |

La variable **`$BACKEND_BASE`** en cada proxy debe apuntar a la URL interna correcta del API.

Apache: en **`.htaccess`** se expone `HTTP_AUTHORIZATION` a PHP para que el token llegue al backend.

---

## Estructura de carpetas (resumen)

```
forms/
├── services/ # config.js, llamadas API
├── utils/
├── components/
├── admin/
│ ├── shared/ # sidebar.css, nav-toggle, formularios, form_sections, customization
│ ├── login/index.html
│ ├── dashboard/index.html
│ ├── dashboard/index.js
│ ├── email/index.html
│ ├── formularios/
│ ├── respuestas/index.html
│ └── config/
├── form/index.html # formulario público dinámico
├── index.html
├── proxy-forms-api.php
├── proxy-users-api.php
└── .htaccess
```

---

## Flujos resumidos

1. **Inscripción pública:** el usuario abre `form/?id=` (Apache reescribe a `form/index.html`) → JSON a `POST /api/forms/submit`.
2. **Admin:** login en users-api → token en `localStorage` / envío en `Authorization` → rutas `/api/admin/...` vía proxy o directo.

Mantener **endpoints alineados** con [main.go del backend](../../forms-api/forms/main.go) al añadir rutas nuevas.

---

## Mantenimiento

- Cambios de URL o rutas: **`services/config.js`**.
- Changelog interno: **`admin/config/changelog/`** (`changelog.js`).
- Tras desplegar backend nuevo: invalidar caché del navegador en entornos con proxy agresivo.
