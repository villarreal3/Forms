import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsNotEmpty, IsEnum, IsOptional, IsArray, IsBoolean } from 'class-validator';
import { RoleType } from '../entities/role.entity';

export class CreateRoleDto {
  @ApiProperty({ description: 'Nombre del rol', example: 'Administrador' })
  @IsString({ message: 'El nombre debe ser un texto' })
  @IsNotEmpty({ message: 'El nombre es obligatorio' })
  name: string;

  @ApiProperty({ description: 'Tipo de rol', enum: RoleType, example: RoleType.ADMIN })
  @IsEnum(RoleType, { message: 'El tipo de rol debe ser válido' })
  @IsNotEmpty({ message: 'El tipo de rol es obligatorio' })
  type: RoleType;

  @ApiProperty({ description: 'Descripción del rol', required: false })
  @IsOptional()
  @IsString({ message: 'La descripción debe ser un texto' })
  description?: string;

  @ApiProperty({ description: 'Permisos del rol', type: [String], required: false })
  @IsOptional()
  @IsArray({ message: 'Los permisos deben ser un array' })
  @IsString({ each: true, message: 'Cada permiso debe ser un texto' })
  permissions?: string[];

  @ApiProperty({ description: 'Rol activo', default: true })
  @IsOptional()
  @IsBoolean({ message: 'isActive debe ser un valor booleano' })
  isActive?: boolean;
}

