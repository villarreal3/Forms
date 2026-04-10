import { ApiProperty } from '@nestjs/swagger';
import { IsEmail, IsNotEmpty, IsString, MinLength, Matches } from 'class-validator';

export class RequestResetPasswordDto {
 @ApiProperty({ description: 'Correo electrónico del usuario', example: 'usuario@example.com' })
 @IsEmail({}, { message: 'Debe proporcionar un correo electrónico válido' })
 @IsNotEmpty({ message: 'El correo electrónico es obligatorio' })
 email: string;
}

export class ResetPasswordDto {
 @ApiProperty({ description: 'Token de reseteo enviado por correo' })
 @IsString({ message: 'El token debe ser un texto' })
 @IsNotEmpty({ message: 'El token es obligatorio' })
 token: string;

 @ApiProperty({ description: 'Nueva contraseña', example: 'NewPassword123!' })
 @IsString({ message: 'La nueva contraseña debe ser un texto' })
 @IsNotEmpty({ message: 'La nueva contraseña es obligatoria' })
 @MinLength(8, { message: 'La contraseña debe tener al menos 8 caracteres' })
 @Matches(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/, {
 message:
 'La contraseña debe contener al menos una mayúscula, una minúscula, un número y un carácter especial',
 })
 newPassword: string;
}

