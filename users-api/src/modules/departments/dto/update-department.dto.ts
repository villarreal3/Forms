import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsEmail, IsOptional, IsBoolean } from 'class-validator';

export class UpdateDepartmentDto {
  @ApiProperty({ description: 'Nombre del departamento', required: false })
  @IsOptional()
  @IsString({ message: 'El nombre debe ser un texto' })
  name?: string;

  @ApiProperty({ description: 'Código del departamento', required: false })
  @IsOptional()
  @IsString({ message: 'El código debe ser un texto' })
  code?: string;

  @ApiProperty({ description: 'Descripción del departamento', required: false })
  @IsOptional()
  @IsString({ message: 'La descripción debe ser un texto' })
  description?: string;

  @ApiProperty({ description: 'Email de contacto del departamento', required: false })
  @IsOptional()
  @IsEmail({}, { message: 'Debe proporcionar un correo electrónico válido' })
  contactEmail?: string;

  @ApiProperty({ description: 'Teléfono de contacto del departamento', required: false })
  @IsOptional()
  @IsString({ message: 'El teléfono debe ser un texto' })
  contactPhone?: string;

  @ApiProperty({ description: 'Departamento activo', required: false })
  @IsOptional()
  @IsBoolean({ message: 'isActive debe ser un valor booleano' })
  isActive?: boolean;
}

