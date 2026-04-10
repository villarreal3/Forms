import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsString } from 'class-validator';

export class RefreshTokenDto {
  @ApiProperty({ description: 'Token de refresco' })
  @IsString({ message: 'El token debe ser un texto' })
  @IsNotEmpty({ message: 'El token de refresco es obligatorio' })
  refreshToken: string;
}

