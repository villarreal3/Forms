import { DataSource } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { config } from 'dotenv';
import { Role, RoleType } from '../modules/roles/entities/role.entity';
import { Department } from '../modules/departments/entities/department.entity';
import { User } from '../modules/users/entities/user.entity';

config();

const AppDataSource = new DataSource({
 type: 'postgres',
 host: process.env.DB_HOST || 'localhost',
 port: parseInt(process.env.DB_PORT || '5432'),
 username: process.env.DB_USERNAME || 'app_user',
 password: process.env.DB_PASSWORD || '',
 database: process.env.DB_DATABASE || 'users_db',
 entities: [__dirname + '/../**/*.entity{.ts,.js}'],
 synchronize: false,
});

async function seed() {
 console.log('🌱 Iniciando seed de base de datos...');

 try {
 await AppDataSource.initialize();
 console.log('✅ Conectado a la base de datos');

 // Crear roles
 console.log('📝 Creando roles...');
 const roleRepository = AppDataSource.getRepository(Role);

 const roles = [
 {
 name: 'Super Administrador',
 type: RoleType.SUPER_ADMIN,
 description: 'Acceso completo al sistema',
 permissions: ['*'],
 },
 {
 name: 'Administrador',
 type: RoleType.ADMIN,
 description: 'Gestión de usuarios y configuración',
 permissions: ['users:create', 'users:read', 'users:update', 'users:delete'],
 },
 {
 name: 'Editor',
 type: RoleType.EDITOR,
 description: 'Puede editar información',
 permissions: ['users:read', 'users:update'],
 },
 {
 name: 'Auditor',
 type: RoleType.AUDITOR,
 description: 'Solo lectura y auditoría',
 permissions: ['users:read', 'audit:read'],
 },
 {
 name: 'Usuario',
 type: RoleType.USER,
 description: 'Usuario estándar',
 permissions: ['profile:read', 'profile:update'],
 },
 ];

 const createdRoles: Role[] = [];
 for (const roleData of roles) {
 const existingRole = await roleRepository.findOne({ where: { type: roleData.type } });
 if (!existingRole) {
 const role = roleRepository.create(roleData);
 const savedRole = await roleRepository.save(role);
 createdRoles.push(savedRole);
 console.log(` ✓ Rol creado: ${roleData.name}`);
 } else {
 createdRoles.push(existingRole);
 console.log(` → Rol ya existe: ${roleData.name}`);
 }
 }

 // Crear departamentos
 console.log('🏢 Creando departamentos...');
 const departmentRepository = AppDataSource.getRepository(Department);

 const departments = [
 {
 name: 'Recursos Humanos',
 code: 'RH',
 description: 'Gestión de personal y recursos humanos',
 contactEmail: 'rh@example.com',
 },
 {
 name: 'Tecnología e Informática',
 code: 'IT',
 description: 'Departamento de tecnología y sistemas',
 contactEmail: 'it@example.com',
 },
 {
 name: 'Dirección General',
 code: 'DG',
 description: 'Dirección General',
 contactEmail: 'direccion@example.com',
 },
 {
 name: 'Legal',
 code: 'LEGAL',
 description: 'Departamento Legal',
 contactEmail: 'legal@example.com',
 },
 {
 name: 'Comunicaciones',
 code: 'COMM',
 description: 'Departamento de Comunicación y Prensa',
 contactEmail: 'comunicaciones@example.com',
 },
 ];

 const createdDepartments: Department[] = [];
 for (const deptData of departments) {
 const existingDept = await departmentRepository.findOne({ where: { code: deptData.code } });
 if (!existingDept) {
 const dept = departmentRepository.create(deptData);
 const savedDept = await departmentRepository.save(dept);
 createdDepartments.push(savedDept);
 console.log(` ✓ Departamento creado: ${deptData.name}`);
 } else {
 createdDepartments.push(existingDept);
 console.log(` → Departamento ya existe: ${deptData.name}`);
 }
 }

 // Crear usuario super admin
 console.log('👤 Creando usuario Super Admin...');
 const userRepository = AppDataSource.getRepository(User);

 const superAdminEmail = process.env.SUPER_ADMIN_EMAIL || 'admin@example.com';
 const existingSuperAdmin = await userRepository.findOne({ where: { email: superAdminEmail } });

 if (!existingSuperAdmin) {
 const password = process.env.SUPER_ADMIN_PASSWORD || 'ChangeMe123!';
 const bcryptRounds = parseInt(process.env.BCRYPT_ROUNDS || '10');
 const hashedPassword = await bcrypt.hash(password, bcryptRounds);

 const superAdminRole = createdRoles.find((r) => r.type === RoleType.SUPER_ADMIN);
 const itDepartment = createdDepartments.find((d) => d.code === 'IT');

 const superAdmin = userRepository.create({
 fullName: process.env.SUPER_ADMIN_NAME || 'Administrador Sistema',
 email: superAdminEmail,
 password: hashedPassword,
 isActive: true,
 emailVerified: true,
 roles: superAdminRole ? [superAdminRole] : [],
 department: itDepartment,
 });

 await userRepository.save(superAdmin);
 console.log(` ✓ Super Admin creado: ${superAdminEmail}`);
 console.log(` ⚠️ Contraseña por defecto: ${password}`);
 console.log(' ⚠️ ¡IMPORTANTE! Cambie la contraseña después del primer login');
 } else {
 console.log(` → Super Admin ya existe: ${superAdminEmail}`);
 }

 console.log('\n✅ Seed completado exitosamente');
 await AppDataSource.destroy();
 process.exit(0);
 } catch (error) {
 console.error('❌ Error durante el seed:', error);
 await AppDataSource.destroy();
 process.exit(1);
 }
}

seed();

