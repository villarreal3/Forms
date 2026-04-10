import { ApiProperty } from '@nestjs/swagger';
import {
  IsEmail,
  IsString,
  IsOptional,
  IsUUID,
  IsArray,
  IsBoolean,
} from 'class-validator';

export class UpdateUserDto {
  @ApiProperty({ description: 'Nombre completo del usuario', required: false })
  @IsOptional()
  @IsString({ message: 'El nombre debe ser un texto' })
  fullName?: string;

  @ApiProperty({ description: 'Correo electrónico institucional', required: false })
  @IsOptional()
  @IsEmail({}, { message: 'Debe proporcionar un correo electrónico válido' })
  email?: string;

  @ApiProperty({ description: 'Cédula o identificación', required: false })
  @IsOptional()
  @IsString({ message: 'La identificación debe ser un texto' })
  identification?: string;

  @ApiProperty({ description: 'Teléfono de contacto', required: false })
  @IsOptional()
  @IsString({ message: 'El teléfono debe ser un texto' })
  phone?: string;

  @ApiProperty({ description: 'URL de avatar/foto', required: false })
  @IsOptional()
  @IsString({ message: 'La URL del avatar debe ser un texto' })
  avatarUrl?: string;

  @ApiProperty({ description: 'Usuario activo', required: false })
  @IsOptional()
  @IsBoolean({ message: 'isActive debe ser un valor booleano' })
  isActive?: boolean;

  @ApiProperty({ description: 'IDs de roles asignados', type: [String], required: false })
  @IsOptional()
  @IsArray({ message: 'Los roles deben ser un array' })
  @IsUUID('4', { each: true, message: 'Cada ID de rol debe ser un UUID válido' })
  roleIds?: string[];

  @ApiProperty({ description: 'ID del departamento', required: false })
  @IsOptional()
  @IsUUID('4', { message: 'El ID del departamento debe ser un UUID válido' })
  departmentId?: string;
}

