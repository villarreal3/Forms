import {
 Entity,
 Column,
 ManyToMany,
 JoinTable,
 ManyToOne,
 JoinColumn,
 OneToMany,
 Index,
} from 'typeorm';
import { ApiProperty } from '@nestjs/swagger';
import { Exclude } from 'class-transformer';
import { BaseEntity } from '../../../common/entities/base.entity';
import { Role } from '../../roles/entities/role.entity';
import { Department } from '../../departments/entities/department.entity';
import { RefreshToken } from '../../auth/entities/refresh-token.entity';
import { AuditLog } from '../../audit/entities/audit-log.entity';

@Entity('users')
export class User extends BaseEntity {
 @ApiProperty({ description: 'Nombre completo del usuario', example: 'Juan Pérez' })
 @Column({ name: 'full_name', length: 100 })
 fullName: string;

 @ApiProperty({
 description: 'Correo electrónico institucional',
 example: 'juan.perez@example.com',
 })
 @Column({ unique: true, length: 100 })
 @Index()
 email: string;

 @ApiProperty({ description: 'Contraseña cifrada' })
 @Column()
 @Exclude()
 password: string;

 @ApiProperty({ description: 'Cédula o identificación', example: '8-123-4567', required: false })
 @Column({ unique: true, length: 50, nullable: true })
 @Index()
 identification?: string;

 @ApiProperty({ description: 'Teléfono de contacto', required: false })
 @Column({ length: 20, nullable: true })
 phone?: string;

 @ApiProperty({ description: 'URL de avatar/foto', required: false })
 @Column({ type: 'text', nullable: true })
 avatarUrl?: string;

 @ApiProperty({ description: 'Usuario activo', default: true })
 @Column({ default: true })
 isActive: boolean;

 @ApiProperty({ description: 'Email verificado', default: false })
 @Column({ name: 'email_verified', default: false })
 emailVerified: boolean;

 @ApiProperty({ description: 'Fecha de último login', required: false })
 @Column({ name: 'last_login_at', type: 'timestamp', nullable: true })
 lastLoginAt?: Date;

 @ApiProperty({ description: 'Intentos de login fallidos', default: 0 })
 @Column({ name: 'failed_login_attempts', default: 0 })
 failedLoginAttempts: number;

 @ApiProperty({ description: 'Fecha de bloqueo de cuenta', required: false })
 @Column({ name: 'locked_until', type: 'timestamp', nullable: true })
 lockedUntil?: Date;

 @ApiProperty({ description: 'Token de reseteo de contraseña', required: false })
 @Column({ name: 'password_reset_token', length: 255, nullable: true })
 @Exclude()
 passwordResetToken?: string;

 @ApiProperty({ description: 'Expiración del token de reseteo', required: false })
 @Column({ name: 'password_reset_expires', type: 'timestamp', nullable: true })
 passwordResetExpires?: Date;

 @ApiProperty({ description: 'Roles del usuario', type: () => [Role] })
 @ManyToMany(() => Role, (role) => role.users, { eager: true })
 @JoinTable({
 name: 'user_roles',
 joinColumn: { name: 'user_id', referencedColumnName: 'id' },
 inverseJoinColumn: { name: 'role_id', referencedColumnName: 'id' },
 })
 roles: Role[];

 @ApiProperty({ description: 'Departamento del usuario', type: () => Department, required: false })
 @ManyToOne(() => Department, (department) => department.users, { eager: true })
 @JoinColumn({ name: 'department_id' })
 department?: Department;

 @OneToMany(() => RefreshToken, (token) => token.user)
 refreshTokens: RefreshToken[];

 @OneToMany(() => AuditLog, (log) => log.user)
 auditLogs: AuditLog[];

 @ApiProperty({ description: 'Usuario que creó este registro', required: false })
 @Column({ name: 'created_by', type: 'uuid', nullable: true })
 createdBy?: string;

 @ApiProperty({ description: 'Usuario que actualizó este registro', required: false })
 @Column({ name: 'updated_by', type: 'uuid', nullable: true })
 updatedBy?: string;
}

