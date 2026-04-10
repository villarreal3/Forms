# Almacenamiento local (API de formularios)

Rutas por defecto (relativas al directorio de trabajo del proceso, p. ej. `/app` en Docker):

| Directorio | Variable   | Uso |
|------------|------------|-----|
| `emails/`  | `EMAILS_DIR` | Borradores / registro de correos masivos |
| `img/`     | `IMAGES_DIR` | Imágenes de formularios (logos, etc.), por `form_id` |

En **docker-compose** se monta `./src` → `/app/src` para persistir datos fuera del contenedor.
