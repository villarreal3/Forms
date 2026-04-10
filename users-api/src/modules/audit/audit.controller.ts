import { Controller, Get, Query, UseGuards, Param } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiParam } from '@nestjs/swagger';
import { AuditService } from './audit.service';
import { QueryAuditDto } from './dto/query-audit.dto';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { RolesGuard } from '../auth/guards/roles.guard';
import { Roles } from '../auth/decorators/roles.decorator';
import { RoleType } from '../roles/entities/role.entity';
import { AuditLog } from './entities/audit-log.entity';

@ApiTags('Audit')
@Controller('audit')
@UseGuards(JwtAuthGuard, RolesGuard)
@ApiBearerAuth('JWT-auth')
export class AuditController {
  constructor(private readonly auditService: AuditService) {}

  @Get()
  @Roles(RoleType.SUPER_ADMIN, RoleType.ADMIN, RoleType.AUDITOR)
  @ApiOperation({ summary: 'Listar todos los logs de auditoría con filtros' })
  @ApiResponse({ status: 200, description: 'Lista de logs de auditoría' })
  async findAll(@Query() query: QueryAuditDto) {
    return this.auditService.findAll(query);
  }

  @Get('user/:userId')
  @Roles(RoleType.SUPER_ADMIN, RoleType.ADMIN, RoleType.AUDITOR)
  @ApiOperation({ summary: 'Obtener logs de auditoría de un usuario específico' })
  @ApiParam({ name: 'userId', description: 'ID del usuario' })
  @ApiResponse({ status: 200, description: 'Logs de auditoría del usuario', type: [AuditLog] })
  async findByUser(@Param('userId') userId: string, @Query('page') page?: number, @Query('limit') limit?: number) {
    return this.auditService.findByUser(userId, page, limit);
  }

  @Get('entity/:entityType/:entityId')
  @Roles(RoleType.SUPER_ADMIN, RoleType.ADMIN, RoleType.AUDITOR)
  @ApiOperation({ summary: 'Obtener logs de auditoría de una entidad específica' })
  @ApiParam({ name: 'entityType', description: 'Tipo de entidad (User, Role, Department)' })
  @ApiParam({ name: 'entityId', description: 'ID de la entidad' })
  @ApiResponse({ status: 200, description: 'Logs de auditoría de la entidad', type: [AuditLog] })
  async findByEntity(
    @Param('entityType') entityType: string,
    @Param('entityId') entityId: string,
    @Query('page') page?: number,
    @Query('limit') limit?: number,
  ) {
    return this.auditService.findByEntity(entityType, entityId, page, limit);
  }
}

