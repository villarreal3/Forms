import { Entity, Column, ManyToOne, JoinColumn, Index } from 'typeorm';
import { ApiProperty } from '@nestjs/swagger';
import { BaseEntity } from '../../../common/entities/base.entity';
import { User } from '../../users/entities/user.entity';

export enum AuditAction {
  CREATE = 'create',
  UPDATE = 'update',
  DELETE = 'delete',
  LOGIN = 'login',
  LOGOUT = 'logout',
  PASSWORD_CHANGE = 'password_change',
  PASSWORD_RESET = 'password_reset',
  ROLE_CHANGE = 'role_change',
  DEPARTMENT_CHANGE = 'department_change',
}

@Entity('audit_logs')
export class AuditLog extends BaseEntity {
  @ApiProperty({ description: 'Usuario que realizó la acción', required: false, type: () => User })
  @ManyToOne(() => User, (user) => user.auditLogs, { nullable: true, onDelete: 'SET NULL' })
  @JoinColumn({ name: 'user_id' })
  user?: User;

  @ApiProperty({ description: 'Acción realizada', enum: AuditAction })
  @Column({ type: 'enum', enum: AuditAction })
  @Index()
  action: AuditAction;

  @ApiProperty({ description: 'Entidad afectada', example: 'User' })
  @Column({ name: 'entity_type', length: 50 })
  @Index()
  entityType: string;

  @ApiProperty({ description: 'ID de la entidad afectada' })
  @Column({ name: 'entity_id', type: 'uuid', nullable: true })
  @Index()
  entityId?: string;

  @ApiProperty({ description: 'Descripción de la acción' })
  @Column({ type: 'text' })
  description: string;

  @ApiProperty({ description: 'Datos anteriores (antes del cambio)', required: false })
  @Column({ name: 'old_values', type: 'jsonb', nullable: true })
  oldValues?: any;

  @ApiProperty({ description: 'Datos nuevos (después del cambio)', required: false })
  @Column({ name: 'new_values', type: 'jsonb', nullable: true })
  newValues?: any;

  @ApiProperty({ description: 'Dirección IP desde donde se realizó la acción' })
  @Column({ name: 'ip_address', length: 45, nullable: true })
  ipAddress?: string;

  @ApiProperty({ description: 'User Agent del navegador' })
  @Column({ name: 'user_agent', type: 'text', nullable: true })
  userAgent?: string;

  @ApiProperty({ description: 'Metadatos adicionales', required: false })
  @Column({ type: 'jsonb', nullable: true })
  metadata?: any;
}

