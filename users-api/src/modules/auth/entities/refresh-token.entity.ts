import { Entity, Column, ManyToOne, JoinColumn, Index } from 'typeorm';
import { ApiProperty } from '@nestjs/swagger';
import { BaseEntity } from '../../../common/entities/base.entity';
import { User } from '../../users/entities/user.entity';

@Entity('refresh_tokens')
export class RefreshToken extends BaseEntity {
  @ApiProperty({ description: 'Token de refresco' })
  @Column({ type: 'text', unique: true })
  @Index()
  token: string;

  @ApiProperty({ description: 'Usuario propietario del token', type: () => User })
  @ManyToOne(() => User, (user) => user.refreshTokens, { onDelete: 'CASCADE' })
  @JoinColumn({ name: 'user_id' })
  user: User;

  @ApiProperty({ description: 'Fecha de expiración del token' })
  @Column({ name: 'expires_at', type: 'timestamp' })
  expiresAt: Date;

  @ApiProperty({ description: 'Token revocado', default: false })
  @Column({ name: 'is_revoked', default: false })
  isRevoked: boolean;

  @ApiProperty({ description: 'Dirección IP desde donde se generó el token', required: false })
  @Column({ name: 'ip_address', length: 45, nullable: true })
  ipAddress?: string;

  @ApiProperty({ description: 'User Agent del navegador', required: false })
  @Column({ name: 'user_agent', type: 'text', nullable: true })
  userAgent?: string;
}

