# 🏛️ API centralizada de usuarios

## 📋 Descripción

API RESTful segura, escalable y centralizada para la gestión, autenticación y autorización de usuarios. Actúa como **Single Source of Truth** (fuente única de verdad) para credenciales, roles y permisos de todos los sistemas institucionales.

## ✨ Características Principales

### 🔐 Autenticación & Seguridad
- ✅ Login institucional con JWT
- ✅ Refresh Tokens para renovación de sesiones
- ✅ Protección contra ataques (CSRF, XSS, bruteforce)
- ✅ Rate Limiting configurable
- ✅ CORS restrictivo con lista blanca de dominios/IPs
- ✅ Encriptación de contraseñas con bcrypt
- ✅ Bloqueo temporal de cuenta tras intentos fallidos
- ✅ Sistema de reseteo de contraseña

### 👥 Gestión de Usuarios
- ✅ CRUD completo (Create, Read, Update, Delete)
- ✅ Soft delete (desactivación)
- ✅ Hard delete (eliminación permanente)
- ✅ Restauración de usuarios desactivados
- ✅ Búsqueda y filtrado avanzado
- ✅ Paginación
- ✅ Gestión de perfiles personales

### 🎭 Roles y Permisos
- ✅ 5 roles predefinidos: Super Admin, Admin, Editor, Auditor, Usuario
- ✅ Permisos granulares por rol
- ✅ Asignación múltiple de roles
- ✅ Control de acceso basado en roles (RBAC)

### 🏢 Departamentos
- ✅ Gestión de departamentos organizacionales
- ✅ Asignación de usuarios a departamentos
- ✅ Información de contacto por departamento

### 📊 Auditoría
- ✅ Registro completo de todas las acciones
- ✅ Trazabilidad de cambios (antes/después)
- ✅ Registro de IP y User Agent
- ✅ Consultas de auditoría por usuario, entidad o acción

## 🛠️ Tecnologías Utilizadas

- **Framework**: NestJS
- **Lenguaje**: TypeScript
- **Base de Datos**: PostgreSQL
- **ORM**: TypeORM
- **Autenticación**: JWT + Passport
- **Validación**: class-validator, class-transformer
- **Documentación**: Swagger/OpenAPI
- **Seguridad**: Helmet, Throttler
- **Testing**: Jest

## 📦 Instalación

### Prerrequisitos

- Node.js >= 18.x
- PostgreSQL >= 14.x
- npm o yarn

### 1. Clonar el repositorio

```bash
git clone <repository-url>
cd users-api
```

### 2. Instalar dependencias

```bash
npm install
```

### 3. Configurar variables de entorno

Copiar el archivo de ejemplo y configurar:

```bash
cp .env.example .env
```

Editar `.env` con tus configuraciones:

```env
# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=app_user
DB_PASSWORD=tu_password_seguro
DB_DATABASE=users_db

# JWT
JWT_SECRET=tu_jwt_secret_muy_seguro_cambiar_en_produccion
JWT_REFRESH_SECRET=tu_refresh_secret_muy_seguro_cambiar_en_produccion

# Dominios permitidos (CORS)
ALLOWED_ORIGINS=http://localhost:3000,https://example.com

# Super Admin
SUPER_ADMIN_EMAIL=admin@example.com
SUPER_ADMIN_PASSWORD=CambiarEsto123!
```

### 4. Crear la base de datos

```bash
psql -U postgres
CREATE DATABASE users_db;
CREATE USER app_user WITH ENCRYPTED PASSWORD 'tu_password';
GRANT ALL PRIVILEGES ON DATABASE users_db TO app_user;
\q
```

### 5. Ejecutar seed (datos iniciales)

Este comando creará roles, departamentos y el usuario Super Admin:

```bash
npm run seed
```

### 6. Iniciar la aplicación

**Desarrollo:**
```bash
npm run start:dev
```

**Producción:**
```bash
npm run build
npm run start:prod
```

La API estará disponible en: `http://localhost:3000/api/v1`

## 📚 Documentación API

### Swagger UI

Una vez iniciada la aplicación, accede a la documentación interactiva en:

```
http://localhost:3000/docs
```

### Endpoints Principales

#### 🔐 Autenticación (`/api/v1/auth`)

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| POST | `/auth/login` | Iniciar sesión | No |
| POST | `/auth/refresh` | Refrescar token | No |
| POST | `/auth/logout` | Cerrar sesión | Sí |
| POST | `/auth/change-password` | Cambiar contraseña | Sí |
| POST | `/auth/request-password-reset` | Solicitar reseteo de contraseña | No |
| POST | `/auth/reset-password` | Restablecer contraseña | No |
| POST | `/auth/verify` | Verificar token | Sí |

#### 👥 Usuarios (`/api/v1/users`)

| Método | Endpoint | Descripción | Roles Requeridos |
|--------|----------|-------------|------------------|
| GET | `/users` | Listar usuarios | Admin, Auditor |
| GET | `/users/me` | Perfil propio | Cualquiera |
| GET | `/users/:id` | Obtener usuario | Admin, Auditor |
| POST | `/users` | Crear usuario | Super Admin, Admin |
| PATCH | `/users/me` | Actualizar perfil | Cualquiera |
| PATCH | `/users/:id` | Actualizar usuario | Super Admin, Admin |
| DELETE | `/users/:id/soft` | Desactivar usuario | Super Admin, Admin |
| DELETE | `/users/:id/hard` | Eliminar permanentemente | Super Admin |
| POST | `/users/:id/restore` | Restaurar usuario | Super Admin, Admin |

#### 🎭 Roles (`/api/v1/roles`)

| Método | Endpoint | Descripción | Roles Requeridos |
|--------|----------|-------------|------------------|
| GET | `/roles` | Listar roles | Admin |
| GET | `/roles/:id` | Obtener rol | Admin |
| POST | `/roles` | Crear rol | Super Admin |
| PATCH | `/roles/:id` | Actualizar rol | Super Admin |
| DELETE | `/roles/:id` | Eliminar rol | Super Admin |

#### 🏢 Departamentos (`/api/v1/departments`)

| Método | Endpoint | Descripción | Roles Requeridos |
|--------|----------|-------------|------------------|
| GET | `/departments` | Listar departamentos | Cualquiera |
| GET | `/departments/:id` | Obtener departamento | Cualquiera |
| POST | `/departments` | Crear departamento | Super Admin, Admin |
| PATCH | `/departments/:id` | Actualizar departamento | Super Admin, Admin |
| DELETE | `/departments/:id` | Eliminar departamento | Super Admin, Admin |

#### 📊 Auditoría (`/api/v1/audit`)

| Método | Endpoint | Descripción | Roles Requeridos |
|--------|----------|-------------|------------------|
| GET | `/audit` | Listar logs | Admin, Auditor |
| GET | `/audit/user/:userId` | Logs por usuario | Admin, Auditor |
| GET | `/audit/entity/:type/:id` | Logs por entidad | Admin, Auditor |

## 🔑 Autenticación

### Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
 -H "Content-Type: application/json" \
 -d '{
 "email": "admin@example.com",
 "password": "ChangeMe123!"
 }'
```

**Respuesta:**
```json
{
 "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
 "refreshToken": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
 "user": {
 "id": "uuid",
 "email": "admin@example.com",
 "fullName": "Administrador Sistema",
 "roles": [...],
 "department": {...}
 }
}
```

### Usar el Token

Incluir el `accessToken` en el header de las peticiones:

```bash
curl -X GET http://localhost:3000/api/v1/users/me \
 -H "Authorization: Bearer {accessToken}"
```

### Refrescar Token

```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
 -H "Content-Type: application/json" \
 -d '{
 "refreshToken": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
 }'
```

## 🚀 Despliegue en Producción

### Con PM2

```bash
# Instalar PM2
npm install -g pm2

# Compilar
npm run build

# Iniciar con PM2
pm2 start dist/main.js -name users-api

# Ver logs
pm2 logs users-api

# Monitorear
pm2 monit

# Guardar configuración
pm2 save
pm2 startup
```

### Con Nginx (Reverse Proxy)

```nginx
server {
 listen 80;
 server_name api.example.com;

 location / {
 proxy_pass http://localhost:3000;
 proxy_http_version 1.1;
 proxy_set_header Upgrade $http_upgrade;
 proxy_set_header Connection 'upgrade';
 proxy_set_header Host $host;
 proxy_set_header X-Real-IP $remote_addr;
 proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
 proxy_set_header X-Forwarded-Proto $scheme;
 proxy_cache_bypass $http_upgrade;
 }
}
```

### Variables de Entorno en Producción

**IMPORTANTE**: En producción, asegúrate de:

1. ✅ Cambiar `JWT_SECRET` y `JWT_REFRESH_SECRET`
2. ✅ Usar contraseñas fuertes para base de datos
3. ✅ Configurar `NODE_ENV=production`
4. ✅ Deshabilitar `synchronize: false` en TypeORM
5. ✅ Configurar dominios reales en `ALLOWED_ORIGINS`
6. ✅ Usar SSL/TLS (HTTPS)
7. ✅ Cambiar contraseña del Super Admin

## 🔒 Seguridad

### Buenas Prácticas Implementadas

- ✅ Contraseñas hasheadas con bcrypt (10 rounds)
- ✅ JWT con expiración corta (15 minutos)
- ✅ Refresh tokens con rotación
- ✅ Rate limiting por IP
- ✅ Validación estricta de datos de entrada
- ✅ Protección contra inyección SQL (TypeORM)
- ✅ Headers de seguridad con Helmet
- ✅ CORS restrictivo
- ✅ Bloqueo de cuenta tras intentos fallidos

### Validación de Contraseñas

Las contraseñas deben cumplir:
- Mínimo 8 caracteres
- Al menos una mayúscula
- Al menos una minúscula
- Al menos un número
- Al menos un carácter especial (@$!%*?&)

## 📊 Estructura del Proyecto

```
users-api/
├── src/
│ ├── common/
│ │ └── entities/
│ │ └── base.entity.ts
│ ├── modules/
│ │ ├── auth/
│ │ │ ├── decorators/
│ │ │ ├── dto/
│ │ │ ├── entities/
│ │ │ ├── guards/
│ │ │ ├── strategies/
│ │ │ ├── auth.controller.ts
│ │ │ ├── auth.module.ts
│ │ │ └── auth.service.ts
│ │ ├── users/
│ │ ├── roles/
│ │ ├── departments/
│ │ └── audit/
│ ├── database/
│ │ └── seed.ts
│ ├── app.module.ts
│ ├── app.controller.ts
│ └── main.ts
├── .env.example
├── .gitignore
├── nest-cli.json
├── package.json
├── tsconfig.json
└── README.md
```

## 🧪 Testing

```bash
# Unit tests
npm run test

# E2E tests
npm run test:e2e

# Test coverage
npm run test:cov
```

## 🤝 Integración con Otros Sistemas

Esta API está diseñada para integrarse fácilmente con cualquier sistema interno :

### Ejemplo de Integración (Frontend)

```javascript
// Login
const response = await fetch('http://api.example.com/api/v1/auth/login', {
 method: 'POST',
 headers: { 'Content-Type': 'application/json' },
 body: JSON.stringify({
 email: 'usuario@example.com',
 password: 'password123'
 })
});

const { accessToken, refreshToken, user } = await response.json();

// Guardar tokens
localStorage.setItem('accessToken', accessToken);
localStorage.setItem('refreshToken', refreshToken);

// Usar en peticiones
const userData = await fetch('http://api.example.com/api/v1/users/me', {
 headers: {
 'Authorization': `Bearer ${accessToken}`
 }
});
```

## 📝 Logs y Monitoreo

Los logs se gestionan automáticamente por NestJS y pueden visualizarse con:

```bash
# PM2 logs
pm2 logs users-api

# Logs en tiempo real
pm2 logs users-api -lines 100
```

## 🐛 Resolución de Problemas

### Error de conexión a base de datos

Verificar que PostgreSQL esté corriendo:
```bash
sudo systemctl status postgresql
```

### Error de permisos

Verificar permisos del usuario en PostgreSQL:
```bash
psql -U postgres
GRANT ALL PRIVILEGES ON DATABASE users_db TO app_user;
```

### Problemas con CORS

Agregar el dominio a `ALLOWED_ORIGINS` en `.env`:
```env
ALLOWED_ORIGINS=http://localhost:3000,https://tu-dominio.com
```

## 📞 Soporte

Para soporte técnico, contactar al equipo de IT :
- Email: it@example.com
- Extensión: XXXX

## 📄 Licencia

© 2025 Proyecto Forms

---

**Desarrollado con ❤️ para el proyecto**

