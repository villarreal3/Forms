# 🚀 Guía de Despliegue - API Centralizada de Usuarios 

## 📋 Tabla de Contenidos
1. [Preparación del Servidor](#preparación-del-servidor)
2. [Instalación de Dependencias](#instalación-de-dependencias)
3. [Configuración de PostgreSQL](#configuración-de-postgresql)
4. [Despliegue con PM2](#despliegue-con-pm2)
5. [Despliegue con Docker](#despliegue-con-docker)
6. [Configuración de Nginx](#configuración-de-nginx)
7. [SSL/HTTPS](#ssl-https)
8. [Mantenimiento](#mantenimiento)

---

## 🖥️ Preparación del Servidor

### Requisitos del Sistema

- **OS**: Ubuntu 20.04/22.04 LTS o similar
- **RAM**: Mínimo 2GB (Recomendado 4GB)
- **CPU**: Mínimo 2 cores
- **Disco**: Mínimo 20GB libres
- **Red**: Puertos 80, 443 y 3000 disponibles

### Actualizar el Sistema

```bash
sudo apt update && sudo apt upgrade -y
```

---

## 📦 Instalación de Dependencias

### 1. Instalar Node.js 18.x

```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
node -v # Verificar versión
npm -v
```

### 2. Instalar PostgreSQL

```bash
sudo apt install postgresql postgresql-contrib -y
sudo systemctl start postgresql
sudo systemctl enable postgresql
sudo systemctl status postgresql
```

### 3. Instalar PM2 (Opcional, para despliegue tradicional)

```bash
sudo npm install -g pm2
```

### 4. Instalar Nginx

```bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

---

## 🗄️ Configuración de PostgreSQL

### 1. Crear Usuario y Base de Datos

```bash
sudo -u postgres psql

-- En el prompt de PostgreSQL:
CREATE DATABASE users_db;
CREATE USER app_user WITH ENCRYPTED PASSWORD 'TU_PASSWORD_SUPER_SEGURO';
GRANT ALL PRIVILEGES ON DATABASE users_db TO app_user;
ALTER DATABASE users_db OWNER TO app_user;
\q
```

### 2. Configurar Acceso Remoto (Opcional)

Editar `/etc/postgresql/14/main/postgresql.conf`:
```conf
listen_addresses = 'localhost' # Para acceso local solamente
# O '*' para acceso desde cualquier IP
```

Editar `/etc/postgresql/14/main/pg_hba.conf`:
```conf
# Permitir acceso local
local all all peer
host all all 127.0.0.1/32 md5
```

Reiniciar PostgreSQL:
```bash
sudo systemctl restart postgresql
```

---

## 🚀 Despliegue con PM2

### 1. Clonar el Repositorio

```bash
cd /var/www
sudo git clone <repository-url> users-api
cd users-api
sudo chown -R $USER:$USER .
```

### 2. Instalar Dependencias

```bash
npm install
```

### 3. Configurar Variables de Entorno

```bash
cp .env.example .env
nano .env
```

**Configuración recomendada para producción:**

```env
NODE_ENV=production
PORT=3000
API_PREFIX=api/v1

# Base de Datos
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=app_user
DB_PASSWORD=TU_PASSWORD_SUPER_SEGURO
DB_DATABASE=users_db

# JWT - CAMBIAR ESTOS VALORES
JWT_SECRET=GENERA_UN_SECRET_ALEATORIO_MUY_LARGO_Y_SEGURO_123456789
JWT_EXPIRATION=15m
JWT_REFRESH_SECRET=GENERA_OTRO_SECRET_DIFERENTE_PARA_REFRESH_987654321
JWT_REFRESH_EXPIRATION=7d

# CORS
ALLOWED_ORIGINS=https://example.com,https://www.example.com
ALLOWED_IPS=127.0.0.1,::1

# Rate Limiting
THROTTLE_TTL=60
THROTTLE_LIMIT=10

# Logging
LOG_LEVEL=info

# Seguridad
BCRYPT_ROUNDS=10

# Super Admin
SUPER_ADMIN_EMAIL=admin@example.com
SUPER_ADMIN_PASSWORD=CAMBIAR_ESTO_POR_PASSWORD_SEGURO_123!
SUPER_ADMIN_NAME=Administrador Sistema
```

### 4. Compilar la Aplicación

```bash
npm run build
```

### 5. Ejecutar Seed (Primera vez)

```bash
npm run seed
```

### 6. Iniciar con PM2

```bash
pm2 start ecosystem.config.js
pm2 save
pm2 startup
```

### 7. Verificar Estado

```bash
pm2 status
pm2 logs users-api
```

---

## 🐳 Despliegue con Docker

### 1. Instalar Docker y Docker Compose

```bash
# Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Instalar Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verificar instalación
docker -version
docker-compose -version
```

### 2. Configurar Variables de Entorno

```bash
cp .env.example .env
nano .env
# Editar las variables necesarias
```

### 3. Construir e Iniciar Contenedores

```bash
docker-compose up -d -build
```

### 4. Verificar Contenedores

```bash
docker-compose ps
docker-compose logs -f api
```

### 5. Ejecutar Seed en Docker

```bash
docker-compose exec api npm run seed
```

### 6. Comandos Útiles

```bash
# Ver logs
docker-compose logs -f

# Reiniciar servicios
docker-compose restart

# Detener servicios
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v
```

---

## 🌐 Configuración de Nginx

### 1. Crear Configuración de Nginx

```bash
sudo nano /etc/nginx/sites-available/users-api
```

**Contenido:**

```nginx
server {
 listen 80;
 server_name api.example.com;

 # Logs
 access_log /var/log/nginx/-api-access.log;
 error_log /var/log/nginx/-api-error.log;

 # Rate Limiting
 limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
 limit_req zone=api_limit burst=20 nodelay;

 # Proxy Headers
 proxy_set_header Host $host;
 proxy_set_header X-Real-IP $remote_addr;
 proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
 proxy_set_header X-Forwarded-Proto $scheme;

 location / {
 proxy_pass http://localhost:3000;
 proxy_http_version 1.1;
 proxy_set_header Upgrade $http_upgrade;
 proxy_set_header Connection 'upgrade';
 proxy_cache_bypass $http_upgrade;
 }

 # Health check
 location /health {
 proxy_pass http://localhost:3000;
 access_log off;
 }
}
```

### 2. Habilitar Sitio

```bash
sudo ln -s /etc/nginx/sites-available/users-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 🔒 SSL/HTTPS con Let's Encrypt

### 1. Instalar Certbot

```bash
sudo apt install certbot python3-certbot-nginx -y
```

### 2. Obtener Certificado SSL

```bash
sudo certbot -nginx -d api.example.com
```

Seguir las instrucciones en pantalla.

### 3. Renovación Automática

Certbot configura automáticamente la renovación. Verificar:

```bash
sudo certbot renew -dry-run
```

### 4. Configuración Final de Nginx con SSL

Certbot actualizará automáticamente el archivo de configuración. Verificar:

```bash
sudo nano /etc/nginx/sites-available/users-api
```

Debe incluir:

```nginx
server {
 listen 443 ssl http2;
 server_name api.example.com;

 ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
 ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;
 
 # SSL Configuration
 ssl_protocols TLSv1.2 TLSv1.3;
 ssl_ciphers HIGH:!aNULL:!MD5;
 ssl_prefer_server_ciphers on;

 # ... resto de la configuración
}

# Redirect HTTP to HTTPS
server {
 listen 80;
 server_name api.example.com;
 return 301 https://$server_name$request_uri;
}
```

---

## 🔧 Mantenimiento

### Backup de Base de Datos

```bash
# Crear backup
sudo -u postgres pg_dump users_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Restaurar backup
sudo -u postgres psql users_db < backup_20250101_120000.sql
```

### Actualizar la Aplicación

```bash
cd /var/www/users-api

# Detener la aplicación
pm2 stop users-api

# Actualizar código
git pull origin main

# Instalar dependencias
npm install

# Compilar
npm run build

# Reiniciar
pm2 restart users-api
```

### Monitoreo

```bash
# Ver logs de PM2
pm2 logs users-api

# Monitoreo en tiempo real
pm2 monit

# Ver estado
pm2 status

# Ver logs de Nginx
sudo tail -f /var/log/nginx/-api-access.log
sudo tail -f /var/log/nginx/-api-error.log
```

### Cambiar Contraseña del Super Admin

1. Iniciar sesión como Super Admin
2. Hacer POST a `/api/v1/auth/change-password`:

```bash
curl -X POST https://api.example.com/api/v1/auth/change-password \
 -H "Authorization: Bearer YOUR_JWT_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "currentPassword": "ChangeMe123!",
 "newPassword": "NuevaPasswordSegura123!"
 }'
```

---

## 🆘 Solución de Problemas

### La aplicación no inicia

```bash
# Verificar logs
pm2 logs users-api -lines 100

# Verificar que PostgreSQL esté corriendo
sudo systemctl status postgresql

# Verificar conectividad a base de datos
psql -h localhost -U app_user -d users_db
```

### Error 502 Bad Gateway en Nginx

```bash
# Verificar que la aplicación esté corriendo
pm2 status
curl http://localhost:3000

# Verificar configuración de Nginx
sudo nginx -t
sudo systemctl restart nginx
```

### Problemas de memoria

```bash
# Aumentar límite de memoria en ecosystem.config.js
max_memory_restart: '2G'

# Reiniciar
pm2 restart users-api
```

---

## 📞 Contacto

Para soporte técnico:
- **Email**: it@example.com
- **Documentación**: [README.md](./README.md)

---

**¡Despliegue completado! 🎉**

