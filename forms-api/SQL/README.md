# Base de datos — Forms

Documentación del esquema **table-per-form** y de los scripts SQL del microservicio de formularios.

| Documento relacionado | Descripción |
|----------------------|-------------|
| [create/README.md](create/README.md) | Orden de ejecución de scripts y detalle de cada archivo |
| [README del backend](../README.md) | API Go, Docker y endpoints |
| [README del frontend](../../../fontend/README.md) | Panel y formularios públicos |

---

## Modelo de datos (resumen)

- **Campos del formulario** viven en `forms.schema` (JSON: secciones + campos). No hay tablas `form_fields` / `form_sections` en el esquema actual.
- **Respuestas**: por cada formulario existe una tabla física `form_{id}_responses` (columnas alineadas al schema). El backend la crea al crear el formulario.
- **Inscripción**: `form_submissions` enlaza `user_id` + `form_id` + `response_id` (máximo una inscripción por usuario y formulario).

```
users ──┬── form_submissions ──► form_{N}_responses
 │ │
forms ──┘ └── (FK form_id)
```

---

## Tablas principales

| Tabla | Uso |
|-------|-----|
| **users** | Participantes: `id_number`, `email`, `created_at`. Único `(id_number, email)`. |
| **forms** | Formularios: nombre, descripción, fechas, cierre, **`schema` JSON**. |
| **form_submissions** | Quién se inscribió en qué formulario y en qué fila de respuestas (`response_id`). |
| **form_{id}_responses** | Datos de respuesta por formulario (tabla dinámica por `id`). |
| **form_customization** | Colores, logos y estilos del formulario público. |
| **attendance** | Asistencia por `user_id` + `form_id`. |
| **roles** | Roles locales para selects del admin (sincronización con users-api). |
| **form_audit_log** | Auditoría de cambios en formularios y personalización. |

Tablas de correo (si están en tu despliegue): plantillas, cola, logs — ver `create/02_procedures.sql` y scripts históricos según versión.

---

## Vistas

Lecturas de admin y utilidades suelen ir contra **vistas**. El **listado público** (`GET /api/public/forms`) consulta la tabla **`forms`** en el backend (sin `vw_public_forms`): incluye borrador (`status=2`) y publicado (`status=1`); excluye cerrado (`status=0`). Cupos y fechas se validan al enviar.

| Vista | Propósito |
|-------|-----------|
| **vw_forms** | Listado y detalle admin de formularios. |
| **vw_users_credentials** | Login público por cédula + email. |
| **vw_form_customization** | Personalización por `form_id`. |
| **v_roles** | Vista sobre `roles` (id, name, type, created_at). |
| **v_dashboard_stats** / **v_form_stats** | Métricas del dashboard (basadas en `form_submissions`). |

Definición: [create/04_views.sql](create/04_views.sql).

---

## Procedimientos almacenados (parcial)

El envío del formulario y gran parte de la lectura se hacen con **SQL parametrizado desde Go**. Siguen existiendo SPs para creación de formularios, personalización, correos, asistencia, etc. Lista actual: [create/02_procedures.sql](create/02_procedures.sql).

---

## Instalación

1. Revisar [create/README.md](create/README.md) para el orden (`00_init` → `06_users_and_grants`).
2. Desde el directorio `SQL/create`:

```bash
mysql -u root -p < create_db.sql
```

`create_db.sql` usa `SOURCE` con rutas relativas: ejecutar desde `SQL/create` o ajustar rutas.

**Base por defecto:** `app_db` (ajustar en `00_init.sql` si usas otro nombre).

---

## Consultas útiles

**Formularios públicos activos** (equivalente lógico a la vista):

```sql
SELECT * FROM forms
WHERE (is_closed = 0 OR is_closed IS NULL) AND expires_at > NOW();
```

**Inscripciones de un formulario** (ej. id = 1):

```sql
SELECT fs.id, fs.user_id, fs.submitted_at, u.id_number, u.email
FROM form_submissions fs
JOIN users u ON u.id = fs.user_id
WHERE fs.form_id = 1
ORDER BY fs.submitted_at DESC;
```

**Estadísticas agregadas:**

```sql
SELECT * FROM v_dashboard_stats;
SELECT * FROM v_form_stats;
```

---

## Backup y mantenimiento

```bash
mysqldump -u root -p app_db > backup_$(date +%Y%m%d).sql
mysql -u root -p app_db < backup_20250101.sql
```

Tras muchas altas/bajas de formularios, `OPTIMIZE TABLE` sobre tablas grandes según necesidad.

---

## Notas

- **Zona horaria**: Panamá (UTC−5) en scripts de creación habituales.
- **Charset**: `utf8mb4` / `utf8mb4_unicode_ci`.
- **Seguridad**: la API usa consultas parametrizadas; no concatenar input del usuario en SQL.

**Versión documentada:** table-per-form + vistas (2025). 
**Motor:** MySQL 8.0+.
