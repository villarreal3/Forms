# 🔒 Guía de Seguridad - API de usuarios

## ⚠️ Checklist de Seguridad para Producción

Antes de desplegar en producción, asegúrate de completar estos pasos:

### 🔑 1. Secretos y Contraseñas

- [ ] **JWT_SECRET**: Generar un secret aleatorio de al menos 64 caracteres
- [ ] **JWT_REFRESH_SECRET**: Generar otro secret diferente de al menos 64 caracteres
- [ ] **DB_PASSWORD**: Usar una contraseña fuerte para PostgreSQL
- [ ] **SUPER_ADMIN_PASSWORD**: Cambiar la contraseña por defecto inmediatamente

#### Generar Secretos Seguros

```bash
# En Linux/Mac
openssl rand -base64 64

# O con Node.js
node -e "console.log(require('crypto').randomBytes(64).toString('base64'))"
```

### 🌐 2. Configuración de Red

- [ ] **CORS**: Configurar solo los dominios específicos 
- [ ] **ALLOWED_IPS**: Limitar IPs si es posible
- [ ] **Firewall**: Configurar firewall del servidor (UFW, iptables)
- [ ] **Puertos**: Exponer solo los puertos necesarios

```env
# Ejemplo de configuración restrictiva
ALLOWED_ORIGINS=https://example.com,https://www.example.com,https://admin.example.com
ALLOWED_IPS=192.168.1.0/24,10.0.0.0/8
```

### 🔐 3. Base de Datos

- [ ] Crear usuario específico para la aplicación (no usar postgres)
- [ ] Usar contraseñas fuertes
- [ ] Configurar PostgreSQL para aceptar solo conexiones locales o VPN
- [ ] Habilitar SSL/TLS en conexiones a base de datos
- [ ] Configurar backups automáticos diarios

```sql
-- Configuración segura de usuario PostgreSQL
CREATE USER app_user WITH ENCRYPTED PASSWORD 'PASSWORD_MUY_SEGURO_AQUI';
GRANT CONNECT ON DATABASE users_db TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;
```

### 🛡️ 4. Servidor

- [ ] Mantener el sistema operativo actualizado
- [ ] Instalar y configurar firewall (UFW)
- [ ] Configurar fail2ban para protección contra bruteforce
- [ ] Deshabilitar login root por SSH
- [ ] Usar autenticación con claves SSH
- [ ] Configurar actualizaciones automáticas de seguridad

```bash
# Configurar UFW (Firewall)
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Instalar fail2ban
sudo apt install fail2ban -y
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

### 🔒 5. SSL/TLS (HTTPS)

- [ ] Instalar certificado SSL (Let's Encrypt recomendado)
- [ ] Forzar redirección HTTP a HTTPS
- [ ] Configurar HSTS (HTTP Strict Transport Security)
- [ ] Usar TLS 1.2 o superior
- [ ] Configurar renovación automática de certificados

```bash
# Instalar certificado con Let's Encrypt
sudo apt install certbot python3-certbot-nginx -y
sudo certbot -nginx -d api.example.com
```

### 📝 6. Logs y Monitoreo

- [ ] Configurar rotación de logs
- [ ] Monitorear logs de errores regularmente
- [ ] Configurar alertas para eventos críticos
- [ ] Revisar logs de auditoría periódicamente
- [ ] Configurar monitoreo de recursos (CPU, RAM, Disco)

```bash
# Configurar logrotate para PM2
sudo nano /etc/logrotate.d/pm2

# Contenido:
/var/www/users-api/logs/*.log {
 daily
 rotate 14
 compress
 delaycompress
 missingok
 notifempty
 create 0640 www-data www-data
 sharedscripts
 postrotate
 pm2 reloadLogs
 endscript
}
```

### 🔄 7. Actualizaciones

- [ ] Mantener dependencias actualizadas
- [ ] Revisar vulnerabilidades con `npm audit`
- [ ] Aplicar parches de seguridad regularmente
- [ ] Probar actualizaciones en ambiente de staging primero

```bash
# Revisar vulnerabilidades
npm audit

# Corregir vulnerabilidades automáticamente
npm audit fix

# Actualizar dependencias
npm update
```

### 👤 8. Gestión de Usuarios

- [ ] Cambiar contraseña del Super Admin
- [ ] Eliminar o desactivar usuarios de prueba
- [ ] Revisar permisos de usuarios regularmente
- [ ] Implementar política de contraseñas fuertes
- [ ] Revisar logs de acceso periódicamente

### 🚫 9. Rate Limiting

- [ ] Ajustar límites según necesidad
- [ ] Monitorear intentos de abuso
- [ ] Configurar diferentes límites por endpoint

```env
# Configuración de Rate Limiting
THROTTLE_TTL=60 # Ventana de tiempo en segundos
THROTTLE_LIMIT=10 # Número de requests permitidos
```

### 🔐 10. Variables de Entorno

- [ ] Nunca commitear archivos .env
- [ ] Usar variables de entorno del servidor en producción
- [ ] Rotar secretos periódicamente
- [ ] Almacenar backups de .env de forma segura

---

## 🚨 Procedimientos de Emergencia

### Si se Compromete un Token JWT

```bash
# 1. Cambiar JWT_SECRET en .env
nano .env

# 2. Reiniciar la aplicación
pm2 restart users-api

# 3. Todos los usuarios necesitarán hacer login nuevamente
```

### Si se Detecta Acceso No Autorizado

```bash
# 1. Revisar logs de auditoría
curl -X GET "https://api.example.com/api/v1/audit?action=login&page=1&limit=100" \
 -H "Authorization: Bearer SUPER_ADMIN_TOKEN"

# 2. Desactivar usuario comprometido
curl -X DELETE "https://api.example.com/api/v1/users/USER_ID/soft" \
 -H "Authorization: Bearer SUPER_ADMIN_TOKEN"

# 3. Revisar todos los cambios recientes del usuario
curl -X GET "https://api.example.com/api/v1/audit/user/USER_ID" \
 -H "Authorization: Bearer SUPER_ADMIN_TOKEN"

# 4. Revocar todos los tokens del usuario
# (Automático al desactivar usuario)

# 5. Cambiar contraseña si es necesario
```

### Backup de Emergencia

```bash
# Backup completo de base de datos
sudo -u postgres pg_dump users_db | gzip > backup_emergency_$(date +%Y%m%d_%H%M%S).sql.gz

# Copiar a ubicación segura
scp backup_emergency_*.sql.gz usuario@servidor-backup:/backups/
```

---

## 📊 Monitoreo de Seguridad

### Métricas a Monitorear

1. **Intentos de Login Fallidos**
 ```sql
 SELECT COUNT(*) FROM audit_logs 
 WHERE action = 'login' 
 AND created_at > NOW() - INTERVAL '1 hour'
 AND metadata->>'success' = 'false';
 ```

2. **Usuarios Bloqueados**
 ```sql
 SELECT email, locked_until 
 FROM users 
 WHERE locked_until > NOW();
 ```

3. **Cambios de Roles Recientes**
 ```sql
 SELECT * FROM audit_logs 
 WHERE action = 'role_change' 
 AND created_at > NOW() - INTERVAL '24 hours';
 ```

4. **Nuevos Usuarios Creados**
 ```sql
 SELECT * FROM audit_logs 
 WHERE action = 'create' 
 AND entity_type = 'User'
 AND created_at > NOW() - INTERVAL '24 hours';
 ```

### Alertas Recomendadas

Configurar alertas para:

- ✅ Más de 10 intentos de login fallidos desde una IP en 1 hora
- ✅ Creación de usuarios con rol Super Admin
- ✅ Eliminación permanente de usuarios (hard delete)
- ✅ Cambios en roles de Super Admin
- ✅ Uso de CPU/RAM superior al 80%
- ✅ Disco con menos de 20% disponible
- ✅ Errores 500 frecuentes

---

## 🔍 Auditoría Regular

### Checklist Mensual

- [ ] Revisar logs de auditoría
- [ ] Verificar usuarios activos
- [ ] Revisar permisos y roles
- [ ] Actualizar dependencias
- [ ] Verificar backups
- [ ] Revisar logs de errores
- [ ] Verificar certificados SSL (renovación)
- [ ] Revisar métricas de uso
- [ ] Probar procedimiento de restauración

### Checklist Trimestral

- [ ] Auditoría completa de seguridad
- [ ] Revisión de políticas de acceso
- [ ] Pruebas de penetración
- [ ] Rotación de secretos (JWT, DB passwords)
- [ ] Revisión de logs de firewall
- [ ] Actualización de documentación
- [ ] Capacitación de usuarios administradores

---

## 🛡️ Mejores Prácticas

### Para Desarrolladores

1. **Nunca** commitear archivos .env o secretos
2. **Siempre** validar datos de entrada
3. **Usar** DTOs con validación estricta
4. **Implementar** rate limiting en endpoints sensibles
5. **Registrar** todas las acciones críticas en auditoría

### Para Administradores

1. **Cambiar** contraseñas regularmente
2. **Revisar** logs de auditoría frecuentemente
3. **Mantener** backups actualizados
4. **Actualizar** el sistema regularmente
5. **Limitar** acceso al servidor solo a personal autorizado

### Para Usuarios

1. **Usar** contraseñas fuertes y únicas
2. **No compartir** credenciales
3. **Cerrar sesión** al terminar
4. **Reportar** actividad sospechosa
5. **Cambiar contraseña** si se sospecha compromiso

---

## 📞 Reporte de Vulnerabilidades

Si descubre una vulnerabilidad de seguridad, por favor:

1. **NO** la divulgue públicamente
2. Contacte inmediatamente al equipo de IT: it@example.com
3. Proporcione detalles específicos y pasos para reproducir
4. Espere confirmación antes de divulgar

---

## 📚 Referencias

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [NestJS Security](https://docs.nestjs.com/security/authentication)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/auth-pg-hba-conf.html)
- [Let's Encrypt](https://letsencrypt.org/)

---

**Mantén esta API segura. La seguridad es responsabilidad de todos. 🔒**

