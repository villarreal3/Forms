# 📘 Ejemplos de Uso - API Centralizada de Usuarios 

Este documento contiene ejemplos prácticos de cómo usar la API en diferentes escenarios.

## 📋 Tabla de Contenidos
1. [Autenticación](#autenticación)
2. [Gestión de Usuarios](#gestión-de-usuarios)
3. [Roles y Departamentos](#roles-y-departamentos)
4. [Auditoría](#auditoría)
5. [Integración Frontend](#integración-frontend)

---

## 🔐 Autenticación

### Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
 -H "Content-Type: application/json" \
 -d '{
 "email": "admin@example.com",
 "password": "ChangeMe123!"
 }'
```

**Respuesta exitosa:**
```json
{
 "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
 "refreshToken": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
 "user": {
 "id": "550e8400-e29b-41d4-a716-446655440000",
 "email": "admin@example.com",
 "fullName": "Administrador Sistema",
 "roles": [
 {
 "id": "...",
 "name": "Super Administrador",
 "type": "super_admin"
 }
 ],
 "department": {
 "id": "...",
 "name": "Tecnología e Informática",
 "code": "IT"
 }
 }
}
```

### Refrescar Token

```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
 -H "Content-Type: application/json" \
 -d '{
 "refreshToken": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
 }'
```

### Verificar Token

```bash
curl -X POST http://localhost:3000/api/v1/auth/verify \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Cerrar Sesión

```bash
curl -X POST http://localhost:3000/api/v1/auth/logout \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "refreshToken": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
 }'
```

### Cambiar Contraseña

```bash
curl -X POST http://localhost:3000/api/v1/auth/change-password \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "currentPassword": "ChangeMe123!",
 "newPassword": "NewSecurePass123!"
 }'
```

### Solicitar Reseteo de Contraseña

```bash
curl -X POST http://localhost:3000/api/v1/auth/request-password-reset \
 -H "Content-Type: application/json" \
 -d '{
 "email": "usuario@example.com"
 }'
```

### Restablecer Contraseña

```bash
curl -X POST http://localhost:3000/api/v1/auth/reset-password \
 -H "Content-Type: application/json" \
 -d '{
 "token": "TOKEN_RECIBIDO_POR_EMAIL",
 "newPassword": "NewSecurePass123!"
 }'
```

---

## 👥 Gestión de Usuarios

### Crear Usuario

```bash
curl -X POST http://localhost:3000/api/v1/users \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "fullName": "Juan Pérez",
 "email": "juan.perez@example.com",
 "password": "SecurePass123!",
 "identification": "8-123-4567",
 "phone": "+507 6000-0000",
 "roleIds": ["ROLE_UUID_HERE"],
 "departmentId": "DEPARTMENT_UUID_HERE"
 }'
```

### Listar Usuarios (con filtros y paginación)

```bash
# Listar todos
curl -X GET "http://localhost:3000/api/v1/users?page=1&limit=10" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Buscar por nombre/email
curl -X GET "http://localhost:3000/api/v1/users?search=juan" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Filtrar por departamento
curl -X GET "http://localhost:3000/api/v1/users?departmentId=DEPT_UUID" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Filtrar por estado activo
curl -X GET "http://localhost:3000/api/v1/users?isActive=true" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Combinar filtros
curl -X GET "http://localhost:3000/api/v1/users?search=juan&departmentId=DEPT_UUID&isActive=true&page=1&limit=20" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Obtener Usuario por ID

```bash
curl -X GET http://localhost:3000/api/v1/users/USER_UUID \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Obtener Perfil Propio

```bash
curl -X GET http://localhost:3000/api/v1/users/me \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Actualizar Usuario

```bash
curl -X PATCH http://localhost:3000/api/v1/users/USER_UUID \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "fullName": "Juan Carlos Pérez",
 "phone": "+507 6111-1111",
 "isActive": true
 }'
```

### Actualizar Perfil Propio

```bash
curl -X PATCH http://localhost:3000/api/v1/users/me \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "fullName": "Juan Carlos Pérez",
 "phone": "+507 6111-1111",
 "avatarUrl": "https://example.com/avatar.jpg"
 }'
```

### Desactivar Usuario (Soft Delete)

```bash
curl -X DELETE http://localhost:3000/api/v1/users/USER_UUID/soft \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Eliminar Usuario Permanentemente (Hard Delete)

```bash
curl -X DELETE http://localhost:3000/api/v1/users/USER_UUID/hard \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Restaurar Usuario Desactivado

```bash
curl -X POST http://localhost:3000/api/v1/users/USER_UUID/restore \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 🎭 Roles y Departamentos

### Crear Rol

```bash
curl -X POST http://localhost:3000/api/v1/roles \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "name": "Gerente",
 "type": "admin",
 "description": "Gerente de área",
 "permissions": ["users:read", "users:update"],
 "isActive": true
 }'
```

### Listar Roles

```bash
curl -X GET http://localhost:3000/api/v1/roles \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Obtener Rol por ID

```bash
curl -X GET http://localhost:3000/api/v1/roles/ROLE_UUID \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Actualizar Rol

```bash
curl -X PATCH http://localhost:3000/api/v1/roles/ROLE_UUID \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "description": "Nueva descripción",
 "permissions": ["users:read", "users:update", "users:create"]
 }'
```

### Crear Departamento

```bash
curl -X POST http://localhost:3000/api/v1/departments \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
 "name": "Finanzas",
 "code": "FIN",
 "description": "Departamento de Finanzas y Contabilidad",
 "contactEmail": "finanzas@example.com",
 "contactPhone": "+507 500-0000"
 }'
```

### Listar Departamentos

```bash
curl -X GET http://localhost:3000/api/v1/departments \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Obtener Departamento por ID

```bash
curl -X GET http://localhost:3000/api/v1/departments/DEPT_UUID \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 📊 Auditoría

### Listar Todos los Logs

```bash
curl -X GET "http://localhost:3000/api/v1/audit?page=1&limit=20" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Filtrar Logs por Usuario

```bash
curl -X GET "http://localhost:3000/api/v1/audit?userId=USER_UUID" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Filtrar Logs por Tipo de Entidad

```bash
curl -X GET "http://localhost:3000/api/v1/audit?entityType=User" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Filtrar Logs por Acción

```bash
curl -X GET "http://localhost:3000/api/v1/audit?action=login" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Obtener Logs de un Usuario Específico

```bash
curl -X GET "http://localhost:3000/api/v1/audit/user/USER_UUID?page=1&limit=10" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Obtener Logs de una Entidad Específica

```bash
curl -X GET "http://localhost:3000/api/v1/audit/entity/User/USER_UUID" \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 💻 Integración Frontend

### Ejemplo en JavaScript/TypeScript (React)

```typescript
// api.service.ts
import axios from 'axios';

const API_URL = 'http://localhost:3000/api/v1';

class ApiService {
 private accessToken: string | null = null;
 private refreshToken: string | null = null;

 constructor() {
 // Recuperar tokens del localStorage
 this.accessToken = localStorage.getItem('accessToken');
 this.refreshToken = localStorage.getItem('refreshToken');
 }

 // Login
 async login(email: string, password: string) {
 const response = await axios.post(`${API_URL}/auth/login`, {
 email,
 password,
 });

 const { accessToken, refreshToken, user } = response.data;

 // Guardar tokens
 this.accessToken = accessToken;
 this.refreshToken = refreshToken;
 localStorage.setItem('accessToken', accessToken);
 localStorage.setItem('refreshToken', refreshToken);

 return user;
 }

 // Logout
 async logout() {
 try {
 await axios.post(
 `${API_URL}/auth/logout`,
 { refreshToken: this.refreshToken },
 { headers: this.getHeaders() }
 );
 } finally {
 this.accessToken = null;
 this.refreshToken = null;
 localStorage.removeItem('accessToken');
 localStorage.removeItem('refreshToken');
 }
 }

 // Refrescar token
 async refreshAccessToken() {
 const response = await axios.post(`${API_URL}/auth/refresh`, {
 refreshToken: this.refreshToken,
 });

 const { accessToken, refreshToken } = response.data;
 this.accessToken = accessToken;
 this.refreshToken = refreshToken;
 localStorage.setItem('accessToken', accessToken);
 localStorage.setItem('refreshToken', refreshToken);

 return accessToken;
 }

 // Headers con autenticación
 private getHeaders() {
 return {
 Authorization: `Bearer ${this.accessToken}`,
 };
 }

 // Hacer petición con manejo automático de refresh token
 async request(method: string, url: string, data?: any) {
 try {
 const response = await axios({
 method,
 url: `${API_URL}${url}`,
 data,
 headers: this.getHeaders(),
 });
 return response.data;
 } catch (error: any) {
 // Si el token expiró, intentar refrescarlo
 if (error.response?.status === 401 && this.refreshToken) {
 await this.refreshAccessToken();
 // Reintentar la petición
 const response = await axios({
 method,
 url: `${API_URL}${url}`,
 data,
 headers: this.getHeaders(),
 });
 return response.data;
 }
 throw error;
 }
 }

 // Métodos específicos
 async getProfile() {
 return this.request('GET', '/users/me');
 }

 async getUsers(params?: any) {
 const queryString = new URLSearchParams(params).toString();
 return this.request('GET', `/users?${queryString}`);
 }

 async createUser(userData: any) {
 return this.request('POST', '/users', userData);
 }

 async updateUser(userId: string, userData: any) {
 return this.request('PATCH', `/users/${userId}`, userData);
 }

 async deleteUser(userId: string) {
 return this.request('DELETE', `/users/${userId}/soft`);
 }
}

export const apiService = new ApiService();
```

### Ejemplo de Uso en Componente React

```typescript
// LoginForm.tsx
import React, { useState } from 'react';
import { apiService } from './api.service';

const LoginForm: React.FC = () => {
 const [email, setEmail] = useState('');
 const [password, setPassword] = useState('');
 const [error, setError] = useState('');
 const [loading, setLoading] = useState(false);

 const handleLogin = async (e: React.FormEvent) => {
 e.preventDefault();
 setLoading(true);
 setError('');

 try {
 const user = await apiService.login(email, password);
 console.log('Usuario autenticado:', user);
 // Redirigir al dashboard
 window.location.href = '/dashboard';
 } catch (err: any) {
 setError(err.response?.data?.message || 'Error al iniciar sesión');
 } finally {
 setLoading(false);
 }
 };

 return (
 <form onSubmit={handleLogin}>
 <input
 type="email"
 value={email}
 onChange={(e) => setEmail(e.target.value)}
 placeholder="Email"
 required
 />
 <input
 type="password"
 value={password}
 onChange={(e) => setPassword(e.target.value)}
 placeholder="Contraseña"
 required
 />
 {error && <div className="error">{error}</div>}
 <button type="submit" disabled={loading}>
 {loading ? 'Cargando...' : 'Iniciar Sesión'}
 </button>
 </form>
 );
};

export default LoginForm;
```

### Ejemplo de Hook Personalizado para Usuarios

```typescript
// useUsers.ts
import { useState, useEffect } from 'react';
import { apiService } from './api.service';

export const useUsers = (filters?: any) => {
 const [users, setUsers] = useState([]);
 const [loading, setLoading] = useState(true);
 const [error, setError] = useState(null);
 const [pagination, setPagination] = useState({
 total: 0,
 page: 1,
 limit: 10,
 totalPages: 0,
 });

 const fetchUsers = async () => {
 setLoading(true);
 try {
 const response = await apiService.getUsers(filters);
 setUsers(response.data);
 setPagination(response.meta);
 } catch (err: any) {
 setError(err.message);
 } finally {
 setLoading(false);
 }
 };

 useEffect(() => {
 fetchUsers();
 }, [JSON.stringify(filters)]);

 const createUser = async (userData: any) => {
 await apiService.createUser(userData);
 await fetchUsers(); // Recargar lista
 };

 const updateUser = async (userId: string, userData: any) => {
 await apiService.updateUser(userId, userData);
 await fetchUsers();
 };

 const deleteUser = async (userId: string) => {
 await apiService.deleteUser(userId);
 await fetchUsers();
 };

 return {
 users,
 loading,
 error,
 pagination,
 createUser,
 updateUser,
 deleteUser,
 refresh: fetchUsers,
 };
};
```

---

## 🔍 Casos de Uso Comunes

### Caso 1: Sistema de Login Unificado

Todos los sistemas pueden usar esta API para autenticar usuarios:

```javascript
// En cualquier aplicación web interna
async function loginUser(email, password) {
 const response = await fetch('https://api.example.com/api/v1/auth/login', {
 method: 'POST',
 headers: { 'Content-Type': 'application/json' },
 body: JSON.stringify({ email, password })
 });
 
 if (response.ok) {
 const { accessToken, user } = await response.json();
 // Guardar token y redirigir
 sessionStorage.setItem('token', accessToken);
 return user;
 }
}
```

### Caso 2: Validación de Permisos

```javascript
async function checkUserPermission(permission) {
 const token = sessionStorage.getItem('token');
 const response = await fetch('https://api.example.com/api/v1/auth/verify', {
 method: 'POST',
 headers: { 'Authorization': `Bearer ${token}` }
 });
 
 const { user } = await response.json();
 const hasPermission = user.roles.some(role => 
 role.permissions.includes(permission) || role.permissions.includes('*')
 );
 
 return hasPermission;
}
```

### Caso 3: Directorio de Empleados

```javascript
async function getEmployeeDirectory(departmentCode) {
 const token = sessionStorage.getItem('token');
 const response = await fetch(
 `https://api.example.com/api/v1/users?departmentCode=${departmentCode}`,
 {
 headers: { 'Authorization': `Bearer ${token}` }
 }
 );
 
 const { data } = await response.json();
 return data;
}
```

---

## 📞 Soporte

Para más información, consultar:
- [README.md](./README.md)
- [DEPLOYMENT.md](./DEPLOYMENT.md)
- Documentación Swagger: `http://localhost:3000/docs`

**¡Desarrollado con ❤️ para el proyecto!**

