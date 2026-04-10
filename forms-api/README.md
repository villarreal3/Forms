# Microservicio de formularios

API REST en **Go** para formularios dinámicos, inscripciones, personalización, dashboard y exportación. Los datos siguen el modelo **table-per-form** en MySQL.

| Recurso | Ubicación |
|---------|-----------|
| **Esquema y SQL** | [SQL/README.md](SQL/README.md) · [SQL/create/README.md](SQL/create/README.md) |
| **Código fuente** | Carpeta `forms/` |
| **Frontend** | [fontend/README.md](../../fontend/README.md) |

---

## Qué hace el servicio

- CRUD de formularios con **schema JSON** (secciones y campos).
- Tablas dinámicas **`form_{id}_responses`** + **`form_submissions`** por inscripción.
- Endpoints públicos (listado, formulario, envío, auto-inscripción por cédula).
- Panel admin: dashboard, respuestas, exportación, personalización, imágenes, correos masivos (según rol).
- JWT compartido con **users-api** (no emite tokens; solo valida).

---

## Stack

| Componente | Tecnología |
|------------|------------|
| Lenguaje | Go 1.21+ |
| HTTP | Gorilla Mux |
| Base de datos | MySQL 8 (`database/sql`) |
| Auth | JWT (mismo secreto que users-api) |
| Contenedores | Docker / Podman Compose |

---

## Requisitos

- Docker o Podman + Compose
- Red compartida con **users-api** (nombre de red Docker: `forms-users-shared`; ver `docker-compose.dev.yml`)
- MySQL (incluido en `docker-compose.yml` del proyecto)

---

## Puesta en marcha

```bash
cd forms-api
# Crear red si no existe (p. ej. si users-api aún no la creó)
docker network create forms-users-shared
# o: podman network create forms-users-shared

docker compose up -d
# o: podman-compose up -d
```

API: `http://localhost:8080` (o el puerto de `SERVER_PORT`).

```bash
curl http://localhost:8080/api/health
docker compose logs -f forms
```

### Documentación OpenAPI (go-swagger)

- **Swagger UI vía Docker Compose** (carpeta `forms/docs` con nginx): `http://localhost:8081/` — el servicio **`api-docs`** arranca con `docker compose up` (puerto configurable con `API_DOCS_PORT`).
- **Swagger UI desde la API** (mismo spec embebido): `http://localhost:8080/api/docs/`.
- **Especificación**: [forms/docs/swagger.yaml](forms/docs/swagger.yaml) (OpenAPI 2.0, compatible con [go-swagger](https://github.com/go-swagger/go-swagger)).

Validar el YAML antes de commit o en CI:

```bash
cd forms
make swagger-validate
```

Instalación alternativa del CLI:

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
swagger validate docs/swagger.yaml
```

Con Docker/Podman:

```bash
docker run --rm -v "$(pwd):/work" -w /work quay.io/goswagger/swagger validate docs/swagger.yaml
```

En **Try it out**, para rutas admin use el botón **Authorize** e indique `Bearer <su_jwt>`.

### Variables de entorno (resumen)

| Variable | Descripción | Por defecto |
|----------|-------------|-------------|
| `DB_HOST` / `DB_PORT` | MySQL | `mysql` / `3306` |
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | Credenciales y BD | `root` / `secret` / `app_db` |
| `SERVER_PORT` | Puerto HTTP | `8080` |
| `EMAILS_DIR` / `IMAGES_DIR` | Almacenamiento local | `src/emails`, `src/img` |
| `JWT_SECRET` | **Debe coincidir** con users-api | — |
| `USERS_API_URL` | URL interna users-api | `http://users-api:3000` |

---

## Estructura del repositorio

```
forms-api/
├── forms/                 # Aplicación Go (punto de entrada: main.go)
│   ├── docs/              # swagger.yaml + Swagger UI (embebidos en el binario)
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── Makefile           # make swagger-validate
├── src/                   # Almacenamiento local (emails, imágenes de formularios)
│   ├── emails/
│   └── img/
├── SQL/                   # Scripts MySQL (ver SQL/README.md)
├── docker-compose.yml
└── README.md
```

---

## Autenticación y roles

1. El usuario inicia sesión en **users-api** y obtiene `accessToken`.
2. Peticiones admin: header `Authorization: Bearer <token>`.
3. Los **roles** se gestionan solo en users-api. Este servicio mapea, por ejemplo:
   - `super_admin`, `admin` → acceso admin
   - `editor` → edición de formularios / respuestas
   - `auditor`, `user` → lectura según middleware

---

## Endpoints (referencia)

**Base:** `/api`

### Públicos (sin JWT)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/health` | Salud del servicio |
| GET | `/public/forms` | Listado de formularios activos |
| GET | `/public/forms/{id}` | Detalle + schema |
| GET | `/public/forms/{id}/sections` | Secciones |
| GET | `/public/forms/{id}/customization` | Personalización |
| GET | `/images/{form_id}/{filename}` | Imagen de formulario |
| POST | `/forms/submit` | Enviar inscripción |
| POST | `/forms/auto-submit` | Auto-inscripción desde otro formulario mismo schema |
| POST | `/users/get` | Usuario por cédula + email |

### Admin (JWT)

Prefijo: `/api/admin/`

| Área | Ejemplos |
|------|----------|
| Formularios | `POST/GET /forms`, `GET /forms/{id}`, `PUT /forms/{id}/schema`, `POST /forms/{id}/close` |
| Secciones | `GET /forms/{form_id}/sections` |
| Plantillas | `GET /fields/templates` |
| Dashboard | `GET /dashboard/stats`, `forms`, `submissions`, `responses-timeline`, `top-active-forms`, `geographic-distribution`, `user-metrics` |
| Respuestas | `GET /submissions` |
| Asistencia | `POST /attendance/update` |
| Exportación | `POST /export/submissions` |
| Personalización | `GET|POST|PUT /customization/forms/{id}` |
| Imágenes | `POST|DELETE /forms/{id}/image` |
| Correos | `POST /email/bulk` (admin) |

---

## Arquitectura (flujo)

```
Frontend ──HTTP──► Forms API ◄──JWT compartido──► Users API
                      │
                      ▼
                   MySQL (app_db)
```

---

## Desarrollo local

```bash
cd forms
go mod download
go build -o app .
./app
```

Asegurar variables de entorno y MySQL accesible (o usar `docker compose`).

---

## Problemas frecuentes

| Síntoma | Comprobación |
|---------|--------------|
| 401 | Mismo `JWT_SECRET` que users-api; header `Authorization: Bearer …` |
| No conecta a MySQL | Contenedor `mysql` arriba; `DB_HOST` en red Docker |
| No llega a users-api | Red compartida; `USERS_API_URL` correcta |

---

## Seguridad

- Consultas parametrizadas contra inyección SQL.
- CORS y middleware de roles en rutas admin.
- No exponer `JWT_SECRET` ni credenciales en repositorio.

---

© Proyecto Forms. Todos los derechos reservados.
