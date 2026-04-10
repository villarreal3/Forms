import { ApiProperty } from '@nestjs/swagger';
import { IsEmail, IsNotEmpty, IsString, MinLength } from 'class-validator';

export class LoginDto {
 @ApiProperty({
 description: 'Correo electrónico institucional',
 example: 'usuario@example.com',
 })
 @IsEmail({}, { message: 'Debe proporcionar un correo electrónico válido' })
 @IsNotEmpty({ message: 'El correo electrónico es obligatorio' })
 email: string;

 @ApiProperty({ description: 'Contraseña', example: 'Password123!', minLength: 6 })
 @IsString({ message: 'La contraseña debe ser un texto' })
 @IsNotEmpty({ message: 'La contraseña es obligatoria' })
 @MinLength(6, { message: 'La contraseña debe tener al menos 6 caracteres' })
 password: string;
}

