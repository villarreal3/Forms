import {
  Injectable,
  UnauthorizedException,
  BadRequestException,
  NotFoundException,
} from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { ConfigService } from '@nestjs/config';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, MoreThan } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { v4 as uuidv4 } from 'uuid';
import { User } from '../users/entities/user.entity';
import { RefreshToken } from './entities/refresh-token.entity';
import { ChangePasswordDto } from './dto/change-password.dto';
import { ResetPasswordDto } from './dto/reset-password.dto';
import { AuditService } from '../audit/audit.service';
import { AuditAction } from '../audit/entities/audit-log.entity';

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
    @InjectRepository(RefreshToken)
    private refreshTokenRepository: Repository<RefreshToken>,
    private jwtService: JwtService,
    private configService: ConfigService,
    private auditService: AuditService,
  ) {}

  async validateUser(email: string, password: string): Promise<any> {
    const user = await this.userRepository.findOne({
      where: { email },
      relations: ['roles', 'department'],
    });

    if (!user) {
      return null;
    }

    // Verificar si la cuenta está bloqueada
    if (user.lockedUntil && user.lockedUntil > new Date()) {
      throw new UnauthorizedException('Cuenta bloqueada temporalmente debido a múltiples intentos fallidos');
    }

    // Validar contraseña
    const isPasswordValid = await bcrypt.compare(password, user.password);

    if (!isPasswordValid) {
      // Incrementar intentos fallidos
      user.failedLoginAttempts += 1;

      // Bloquear cuenta si hay más de 5 intentos fallidos
      if (user.failedLoginAttempts >= 5) {
        user.lockedUntil = new Date(Date.now() + 15 * 60 * 1000); // 15 minutos
      }

      await this.userRepository.save(user);
      return null;
    }

    if (!user.isActive) {
      throw new UnauthorizedException('Usuario inactivo');
    }

    // Resetear intentos fallidos y actualizar último login
    user.failedLoginAttempts = 0;
    user.lockedUntil = undefined;
    user.lastLoginAt = new Date();
    await this.userRepository.save(user);

    const { password: _, ...result } = user;
    return result;
  }

  async login(user: User, ipAddress?: string, userAgent?: string) {
    const payload = {
      sub: user.id,
      email: user.email,
      roles: user.roles.map((role) => role.type),
    };

    const accessToken = this.jwtService.sign(payload, {
      secret: this.configService.get('JWT_SECRET'),
      expiresIn: this.configService.get('JWT_EXPIRATION'),
    });

    const refreshToken = await this.createRefreshToken(user, ipAddress, userAgent);

    // Registrar auditoría
    await this.auditService.log({
      user,
      action: AuditAction.LOGIN,
      entityType: 'User',
      entityId: user.id,
      description: `Usuario ${user.email} ha iniciado sesión`,
      ipAddress,
      userAgent,
    });

    return {
      accessToken,
      refreshToken: refreshToken.token,
      user: {
        id: user.id,
        email: user.email,
        fullName: user.fullName,
        roles: user.roles,
        department: user.department,
      },
    };
  }

  async createRefreshToken(user: User, ipAddress?: string, userAgent?: string): Promise<RefreshToken> {
    const token = uuidv4();
    const expiresIn = this.configService.get('JWT_REFRESH_EXPIRATION') || '7d';
    const expiresAt = new Date();

    // Calcular fecha de expiración
    const match = expiresIn.match(/(\d+)([dhms])/);
    if (match) {
      const value = parseInt(match[1]);
      const unit = match[2];
      switch (unit) {
        case 'd':
          expiresAt.setDate(expiresAt.getDate() + value);
          break;
        case 'h':
          expiresAt.setHours(expiresAt.getHours() + value);
          break;
        case 'm':
          expiresAt.setMinutes(expiresAt.getMinutes() + value);
          break;
        case 's':
          expiresAt.setSeconds(expiresAt.getSeconds() + value);
          break;
      }
    }

    const refreshToken = this.refreshTokenRepository.create({
      token,
      user,
      expiresAt,
      ipAddress,
      userAgent,
    });

    return await this.refreshTokenRepository.save(refreshToken);
  }

  async refreshAccessToken(refreshTokenValue: string, ipAddress?: string, userAgent?: string) {
    const refreshToken = await this.refreshTokenRepository.findOne({
      where: { token: refreshTokenValue },
      relations: ['user', 'user.roles', 'user.department'],
    });

    if (!refreshToken) {
      throw new UnauthorizedException('Token de refresco inválido');
    }

    if (refreshToken.isRevoked) {
      throw new UnauthorizedException('Token de refresco revocado');
    }

    if (refreshToken.expiresAt < new Date()) {
      throw new UnauthorizedException('Token de refresco expirado');
    }

    const user = refreshToken.user;

    if (!user.isActive) {
      throw new UnauthorizedException('Usuario inactivo');
    }

    // Generar nuevo access token
    const payload = {
      sub: user.id,
      email: user.email,
      roles: user.roles.map((role) => role.type),
    };

    const accessToken = this.jwtService.sign(payload, {
      secret: this.configService.get('JWT_SECRET'),
      expiresIn: this.configService.get('JWT_EXPIRATION'),
    });

    // Opcional: rotar el refresh token (crear uno nuevo y revocar el anterior)
    const newRefreshToken = await this.createRefreshToken(user, ipAddress, userAgent);
    refreshToken.isRevoked = true;
    await this.refreshTokenRepository.save(refreshToken);

    return {
      accessToken,
      refreshToken: newRefreshToken.token,
    };
  }

  async logout(userId: string, refreshTokenValue?: string) {
    if (refreshTokenValue) {
      const refreshToken = await this.refreshTokenRepository.findOne({
        where: { token: refreshTokenValue },
      });

      if (refreshToken) {
        refreshToken.isRevoked = true;
        await this.refreshTokenRepository.save(refreshToken);
      }
    }

    // Revocar todos los tokens activos del usuario (opcional)
    await this.refreshTokenRepository.update(
      { user: { id: userId }, isRevoked: false },
      { isRevoked: true },
    );

    const user = await this.userRepository.findOne({ where: { id: userId } });

    // Registrar auditoría
    await this.auditService.log({
      user,
      action: AuditAction.LOGOUT,
      entityType: 'User',
      entityId: userId,
      description: `Usuario ${user?.email} ha cerrado sesión`,
    });

    return { message: 'Sesión cerrada exitosamente' };
  }

  async changePassword(userId: string, changePasswordDto: ChangePasswordDto) {
    const user = await this.userRepository.findOne({ where: { id: userId } });

    if (!user) {
      throw new NotFoundException('Usuario no encontrado');
    }

    // Verificar contraseña actual
    const isPasswordValid = await bcrypt.compare(changePasswordDto.currentPassword, user.password);

    if (!isPasswordValid) {
      throw new BadRequestException('Contraseña actual incorrecta');
    }

    // Verificar que la nueva contraseña sea diferente
    const isSamePassword = await bcrypt.compare(changePasswordDto.newPassword, user.password);

    if (isSamePassword) {
      throw new BadRequestException('La nueva contraseña debe ser diferente a la actual');
    }

    // Hashear nueva contraseña
    const bcryptRounds = this.configService.get<number>('BCRYPT_ROUNDS') || 10;
    user.password = await bcrypt.hash(changePasswordDto.newPassword, bcryptRounds);

    await this.userRepository.save(user);

    // Revocar todos los refresh tokens
    await this.refreshTokenRepository.update(
      { user: { id: userId }, isRevoked: false },
      { isRevoked: true },
    );

    // Registrar auditoría
    await this.auditService.log({
      user,
      action: AuditAction.PASSWORD_CHANGE,
      entityType: 'User',
      entityId: userId,
      description: `Usuario ${user.email} cambió su contraseña`,
    });

    return { message: 'Contraseña cambiada exitosamente' };
  }

  async requestPasswordReset(email: string) {
    const user = await this.userRepository.findOne({ where: { email } });

    if (!user) {
      // No revelamos si el usuario existe o no por seguridad
      return { message: 'Si el correo existe, recibirá instrucciones para restablecer su contraseña' };
    }

    // Generar token de reseteo
    const resetToken = uuidv4();
    user.passwordResetToken = resetToken;
    user.passwordResetExpires = new Date(Date.now() + 3600000); // 1 hora

    await this.userRepository.save(user);

    // TODO: Aquí deberías enviar un correo con el token de reseteo
    // Por ahora solo retornamos el token (en producción, NO hacer esto)
    console.log(`Token de reseteo para ${email}: ${resetToken}`);

    return { message: 'Si el correo existe, recibirá instrucciones para restablecer su contraseña' };
  }

  async resetPassword(resetPasswordDto: ResetPasswordDto) {
    const user = await this.userRepository.findOne({
      where: {
        passwordResetToken: resetPasswordDto.token,
        passwordResetExpires: MoreThan(new Date()),
      },
    });

    if (!user) {
      throw new BadRequestException('Token de reseteo inválido o expirado');
    }

    // Hashear nueva contraseña
    const bcryptRounds = this.configService.get<number>('BCRYPT_ROUNDS') || 10;
    user.password = await bcrypt.hash(resetPasswordDto.newPassword, bcryptRounds);
    user.passwordResetToken = undefined;
    user.passwordResetExpires = undefined;

    await this.userRepository.save(user);

    // Revocar todos los refresh tokens
    await this.refreshTokenRepository.update(
      { user: { id: user.id }, isRevoked: false },
      { isRevoked: true },
    );

    // Registrar auditoría
    await this.auditService.log({
      user,
      action: AuditAction.PASSWORD_RESET,
      entityType: 'User',
      entityId: user.id,
      description: `Usuario ${user.email} restableció su contraseña`,
    });

    return { message: 'Contraseña restablecida exitosamente' };
  }
}

