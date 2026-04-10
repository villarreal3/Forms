# Forms — Plataforma de formularios

Monorepositorio con **tres piezas** que trabajan juntas: una API de usuarios (autenticación y autorización), una API de formularios dinámicos (Go + MySQL) y un **frontend estático** para el catálogo público y el panel de administración.

## Componentes

| Parte | Descripción | Documentación |
|--------|-------------|----------------|
| **users-api** | API REST en **NestJS** + **PostgreSQL**: usuarios, roles, departamentos, JWT, auditoría. | [users-api/README.md](users-api/README.md) |
| **forms-api** | API REST en **Go**: formularios table-per-form, envíos, dashboard, exportación, personalización; valida el mismo JWT que users-api. | [forms-api/README.md](forms-api/README.md) |
| **fontend** | HTML/CSS/JS: formularios públicos, login y administración (formularios, respuestas, correos, etc.). | [fontend/README.md](fontend/README.md) |

## Flujo general

1. El usuario inicia sesión contra **users-api** y obtiene un token JWT.
2. Las rutas administrativas de **forms-api** exigen `Authorization: Bearer <token>`; el secreto JWT debe ser **idéntico** en ambas APIs.
3. El frontend consume las URLs configuradas (por ejemplo en `fontend/services/config.js`).

## Requisitos

- **Docker** (o Podman) + Compose para bases de datos y servicios, o entornos locales equivalentes.
- **Node.js** y npm para users-api y herramientas del frontend.
- **Go 1.21+** para compilar y ejecutar forms-api fuera de Docker.

## Puesta en marcha (orientación)

1. **Red Docker compartida** (si usas compose en varios proyectos): la documentación de forms-api describe la red `forms-users-shared` para enlazar con users-api.
2. **Variables de entorno**: copia los `env.template` de cada servicio a `.env` y ajusta credenciales y `JWT_SECRET` compartido. **No subas `.env` al repositorio** (está ignorado en `.gitignore`).
3. Arranca **users-api** y **forms-api** según sus README (puertos habituales: users `3000`, forms `8080`).
4. Sirve **fontend** con un servidor HTTP estático (por ejemplo `python3 -m http.server`) y revisa `config.js` para las URLs de las APIs.

## Estructura del repositorio

```
Forms/
├── users-api/          # API de usuarios (NestJS)
├── forms-api/          # API de formularios (Go) + SQL MySQL + Docker
│   ├── forms/          # Código fuente Go
│   └── SQL/            # Scripts de creación y parches
├── fontend/            # Frontend estático
├── LICENSE
├── README.md           # Este archivo
└── .gitignore
```

## Documentación adicional

- Esquema MySQL: [forms-api/SQL/README.md](forms-api/SQL/README.md)
- Swagger embebido en forms-api: `http://localhost:8080/api/docs/` (tras levantar el servicio)

## Licencia

Ver el archivo [LICENSE](LICENSE) en la raíz del repositorio.
