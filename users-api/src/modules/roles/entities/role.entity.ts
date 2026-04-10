import { Entity, Column, ManyToMany } from 'typeorm';
import { ApiProperty } from '@nestjs/swagger';
import { BaseEntity } from '../../../common/entities/base.entity';
import { User } from '../../users/entities/user.entity';

export enum RoleType {
  SUPER_ADMIN = 'super_admin',
  ADMIN = 'admin',
  EDITOR = 'editor',
  AUDITOR = 'auditor',
  USER = 'user',
}

@Entity('roles')
export class Role extends BaseEntity {
  @ApiProperty({ description: 'Nombre del rol', example: 'Administrador' })
  @Column({ unique: true, length: 50 })
  name: string;

  @ApiProperty({ description: 'Tipo de rol', enum: RoleType, example: RoleType.ADMIN })
  @Column({
    type: 'enum',
    enum: RoleType,
    default: RoleType.USER,
  })
  type: RoleType;

  @ApiProperty({ description: 'Descripción del rol', required: false })
  @Column({ type: 'text', nullable: true })
  description?: string;

  @ApiProperty({ description: 'Permisos del rol (array JSON)', type: [String] })
  @Column({ type: 'jsonb', default: [] })
  permissions: string[];

  @ApiProperty({ description: 'Rol activo', default: true })
  @Column({ default: true })
  isActive: boolean;

  @ManyToMany(() => User, (user) => user.roles)
  users: User[];
}

