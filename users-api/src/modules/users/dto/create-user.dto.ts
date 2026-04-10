import { ApiProperty } from '@nestjs/swagger';
import {
 IsEmail,
 IsNotEmpty,
 IsString,
 MinLength,
 Matches,
 IsOptional,
 IsUUID,
 IsArray,
} from 'class-validator';

export class CreateUserDto {
 @ApiProperty({ description: 'Nombre completo del usuario', example: 'Juan Pérez' })
 @IsString({ message: 'El nombre debe ser un texto' })
 @IsNotEmpty({ message: 'El nombre es obligatorio' })
 fullName: string;

 @ApiProperty({
 description: 'Correo electrónico institucional',
 example: 'juan.perez@example.com',
 })
 @IsEmail({}, { message: 'Debe proporcionar un correo electrónico válido' })
 @IsNotEmpty({ message: 'El correo electrónico es obligatorio' })
 email: string;

 @ApiProperty({
 description: 'Contraseña (mínimo 8 caracteres, al menos una mayúscula, una minúscula, un número y un carácter especial)',
 example: 'Password123!',
 })
 @IsString({ message: 'La contraseña debe ser un texto' })
 @IsNotEmpty({ message: 'La contraseña es obligatoria' })
 @MinLength(8, { message: 'La contraseña debe tener al menos 8 caracteres' })
 @Matches(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/, {
 message:
 'La contraseña debe contener al menos una mayúscula, una minúscula, un número y un carácter especial',
 })
 password: string;

 @ApiProperty({
 description: 'Cédula o identificación',
 example: '8-123-4567',
 required: false,
 })
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

