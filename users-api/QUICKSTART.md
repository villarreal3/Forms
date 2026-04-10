# ⚡ Guía Rápida de Inicio - API de usuarios

## 🚀 Inicio Rápido (5 minutos)

### 1. Instalar dependencias

```bash
npm install
```

### 2. Configurar entorno

```bash
cp .env.example .env
```

Editar `.env` con tu configuración:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=app_user
DB_PASSWORD=tu_password
DB_DATABASE=users_db
```

### 3. Crear base de datos

```bash
# Conectar a PostgreSQL
psql -U postgres

# Crear base de datos y usuario
CREATE DATABASE users_db;
CREATE USER app_user WITH ENCRYPTED PASSWORD 'ChangeMeDb123';
GRANT ALL PRIVILEGES ON DATABASE users_db TO app_user;
\q
```

### 4. Ejecutar seed (datos iniciales)

```bash
npm run seed
```

Esto creará:
- ✅ 5 roles predefinidos (Super Admin, Admin, Editor, Auditor, Usuario)
- ✅ 5 departamentos (RH, IT, Dirección General, Legal, Comunicaciones)
- ✅ Usuario Super Admin (admin@example.com / ChangeMe123!)

### 5. Iniciar aplicación

```bash
npm run start:dev
```

La API estará disponible en: **http://localhost:3000/api/v1**

### 6. Acceder a documentación

Abrir en el navegador: **http://localhost:3000/docs**

---

## 🧪 Prueba Rápida

### Hacer Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
 -H "Content-Type: application/json" \
 -d '{
 "email": "admin@example.com",
 "password": "ChangeMe123!"
 }'
```

Copia el `accessToken` de la respuesta.

### Obtener tu perfil

```bash
curl -X GET http://localhost:3000/api/v1/users/me \
 -H "Authorization: Bearer TU_ACCESS_TOKEN_AQUI"
```

### Ver usuarios

```bash
curl -X GET http://localhost:3000/api/v1/users \
 -H "Authorization: Bearer TU_ACCESS_TOKEN_AQUI"
```

---

## 📊 Estructura de Roles

| Rol | Tipo | Permisos |
|-----|------|----------|
| Super Administrador | `super_admin` | Acceso completo (*) |
| Administrador | `admin` | CRUD usuarios |
| Editor | `editor` | Leer y actualizar |
| Auditor | `auditor` | Solo lectura y auditoría |
| Usuario | `user` | Gestionar su propio perfil |

---

## 🏢 Departamentos Predefinidos

- **RH** - Recursos Humanos
- **IT** - Tecnología e Informática
- **DG** - Dirección General
- **LEGAL** - Legal
- **COMM** - Comunicaciones

---

## 🔑 Credenciales Iniciales

**Email:** admin@example.com 
**Password:** ChangeMe123!

⚠️ **IMPORTANTE:** Cambiar la contraseña después del primer login:

```bash
curl -X POST http://localhost:3000/api/v1/auth/change-password \
 -H "Authorization: Bearer TU_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "currentPassword": "ChangeMe123!",
 "newPassword": "TuNuevaPasswordSegura123!"
 }'
```

---

## 📚 Endpoints Principales

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/auth/login` | POST | Iniciar sesión |
| `/auth/logout` | POST | Cerrar sesión |
| `/users` | GET | Listar usuarios |
| `/users` | POST | Crear usuario |
| `/users/me` | GET | Mi perfil |
| `/users/:id` | PATCH | Actualizar usuario |
| `/roles` | GET | Listar roles |
| `/departments` | GET | Listar departamentos |
| `/audit` | GET | Ver logs de auditoría |

---

## 🛠️ Scripts Disponibles

```bash
# Desarrollo
npm run start:dev # Iniciar en modo desarrollo (con hot-reload)
npm run start:debug # Iniciar con debugger

# Producción
npm run build # Compilar para producción
npm run start:prod # Iniciar en producción

# Base de datos
npm run seed # Ejecutar seed (datos iniciales)

# Testing
npm run test # Tests unitarios
npm run test:e2e # Tests end-to-end
npm run test:cov # Cobertura de tests

# Linting
npm run lint # Ejecutar linter
npm run format # Formatear código
```

---

## 🐛 Problemas Comunes

### Error: "Cannot connect to database"

**Solución:** Verificar que PostgreSQL esté corriendo:
```bash
sudo systemctl status postgresql
sudo systemctl start postgresql
```

### Error: "Port 3000 already in use"

**Solución:** Cambiar el puerto en `.env`:
```env
PORT=3001
```

### Error en seed: "relation does not exist"

**Solución:** Asegurar que `synchronize: true` en desarrollo (ya está configurado).

---

## 📖 Documentación Completa

- [README.md](./README.md) - Documentación completa
- [EXAMPLES.md](./EXAMPLES.md) - Ejemplos de código
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Guía de despliegue
- **Swagger UI:** http://localhost:3000/docs

---

## 📞 Soporte

**Email:** it@example.com

---

**¡Listo para usar! 🎉**

