import { Entity, Column, OneToMany } from 'typeorm';
import { ApiProperty } from '@nestjs/swagger';
import { BaseEntity } from '../../../common/entities/base.entity';
import { User } from '../../users/entities/user.entity';

@Entity('departments')
export class Department extends BaseEntity {
  @ApiProperty({ description: 'Nombre del departamento', example: 'Recursos Humanos' })
  @Column({ unique: true, length: 100 })
  name: string;

  @ApiProperty({ description: 'Código del departamento', example: 'RH' })
  @Column({ unique: true, length: 20 })
  code: string;

  @ApiProperty({ description: 'Descripción del departamento', required: false })
  @Column({ type: 'text', nullable: true })
  description?: string;

  @ApiProperty({ description: 'Email de contacto del departamento', required: false })
  @Column({ length: 100, nullable: true })
  contactEmail?: string;

  @ApiProperty({ description: 'Teléfono de contacto del departamento', required: false })
  @Column({ length: 20, nullable: true })
  contactPhone?: string;

  @ApiProperty({ description: 'Departamento activo', default: true })
  @Column({ default: true })
  isActive: boolean;

  @OneToMany(() => User, (user) => user.department)
  users: User[];
}

