import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { AuditLog, AuditAction } from './entities/audit-log.entity';
import { User } from '../users/entities/user.entity';
import { QueryAuditDto } from './dto/query-audit.dto';

export interface LogAuditParams {
  user?: User | null;
  action: AuditAction;
  entityType: string;
  entityId?: string;
  description: string;
  oldValues?: any;
  newValues?: any;
  ipAddress?: string;
  userAgent?: string;
  metadata?: any;
}

@Injectable()
export class AuditService {
  constructor(
    @InjectRepository(AuditLog)
    private auditLogRepository: Repository<AuditLog>,
  ) {}

  async log(params: LogAuditParams): Promise<AuditLog> {
    const auditLog = this.auditLogRepository.create({
      user: params.user || undefined,
      action: params.action,
      entityType: params.entityType,
      entityId: params.entityId,
      description: params.description,
      oldValues: params.oldValues,
      newValues: params.newValues,
      ipAddress: params.ipAddress,
      userAgent: params.userAgent,
      metadata: params.metadata,
    });

    return await this.auditLogRepository.save(auditLog);
  }

  async findAll(query: QueryAuditDto) {
    const { userId, entityType, entityId, action, page = 1, limit = 20 } = query;

    const queryBuilder = this.auditLogRepository
      .createQueryBuilder('auditLog')
      .leftJoinAndSelect('auditLog.user', 'user');

    // Filtrar por usuario
    if (userId) {
      queryBuilder.andWhere('user.id = :userId', { userId });
    }

    // Filtrar por tipo de entidad
    if (entityType) {
      queryBuilder.andWhere('auditLog.entityType = :entityType', { entityType });
    }

    // Filtrar por ID de entidad
    if (entityId) {
      queryBuilder.andWhere('auditLog.entityId = :entityId', { entityId });
    }

    // Filtrar por acción
    if (action) {
      queryBuilder.andWhere('auditLog.action = :action', { action });
    }

    // Paginación
    const skip = (page - 1) * limit;
    queryBuilder.skip(skip).take(limit);

    // Ordenar por fecha de creación (más recientes primero)
    queryBuilder.orderBy('auditLog.createdAt', 'DESC');

    const [auditLogs, total] = await queryBuilder.getManyAndCount();

    return {
      data: auditLogs,
      meta: {
        total,
        page,
        limit,
        totalPages: Math.ceil(total / limit),
      },
    };
  }

  async findByUser(userId: string, page: number = 1, limit: number = 20) {
    return this.findAll({ userId, page, limit });
  }

  async findByEntity(entityType: string, entityId: string, page: number = 1, limit: number = 20) {
    return this.findAll({ entityType, entityId, page, limit });
  }
}

