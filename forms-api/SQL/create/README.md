# Creación de base de datos - Forms

Scripts SQL divididos por tipo de objeto. El orden de ejecución es fijo.

## Arquitectura table-per-form

- **forms.schema** (JSON): secciones y campos; ya no hay `form_fields` / `form_sections`.
- **form_submissions**: relación `user_id` + `form_id` + `response_id` (una fila por usuario y formulario).
- **form_{id}_responses**: una tabla física por formulario, creada al crear el formulario (columnas según el schema).

## Archivos

| Archivo | Contenido |
|---------|-----------|
| `00_init.sql` | DROP/CREATE database `app_db`, USE, time_zone |
| `01_tables.sql` | Tablas: users, forms (+schema JSON), form_submissions, attendance, roles, form_customization, form_audit_log |
| `02_procedures.sql` | SPs mínimos (sin EAV: sin sp_submit_form ni sp_get_form_fields, etc.) |
| `03_triggers.sql` | Auditoría en forms (incl. cambios de schema), form_customization |
| `04_views.sql` | v_dashboard_stats, v_form_stats (basadas en form_submissions) |
| `06_users_and_grants.sql` | admin_user / public_user |
| `05_verify.sql` | Verificación |

## Ejecución

Desde `SQL/create`:

```bash
mysql -u usuario -p < create_db.sql
```

**Importante:** `create_db.sql` usa `SOURCE` con rutas relativas; ejecutar desde el directorio de los `.sql` o usar rutas absolutas. Para entornos restaurados, recrear tablas `form_*_responses` desde `forms.schema` (script de migración o lógica en Go).

### public_user

Solo EXECUTE sobre SPs públicos restantes. El submit se hace con SQL directo desde la API (usuario con privilegios DML sobre `form_submissions` y cada `form_{id}_responses`).
