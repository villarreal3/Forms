import {
  Injectable,
  NotFoundException,
  BadRequestException,
  ConflictException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, Like, In } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { ConfigService } from '@nestjs/config';
import { User } from './entities/user.entity';
import { Role, RoleType } from '../roles/entities/role.entity';
import { Department } from '../departments/entities/department.entity';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { QueryUserDto } from './dto/query-user.dto';
import { AuditService } from '../audit/audit.service';
import { AuditAction } from '../audit/entities/audit-log.entity';

@Injectable()
export class UsersService {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
    @InjectRepository(Role)
    private roleRepository: Repository<Role>,
    @InjectRepository(Department)
    private departmentRepository: Repository<Department>,
    private configService: ConfigService,
    private auditService: AuditService,
  ) {}

  async create(createUserDto: CreateUserDto, createdBy?: string): Promise<User> {
    // Verificar si el email ya existe
    const existingUser = await this.userRepository.findOne({
      where: { email: createUserDto.email },
    });

    if (existingUser) {
      throw new ConflictException('El correo electrónico ya está registrado');
    }

    // Verificar si la identificación ya existe (si se proporciona)
    if (createUserDto.identification) {
      const existingIdentification = await this.userRepository.findOne({
        where: { identification: createUserDto.identification },
      });

      if (existingIdentification) {
        throw new ConflictException('La identificación ya está registrada');
      }
    }

    // Hashear contraseña
    const bcryptRounds = this.configService.get<number>('BCRYPT_ROUNDS') || 10;
    const hashedPassword = await bcrypt.hash(createUserDto.password, bcryptRounds);

    // Obtener roles
    let roles: Role[] = [];
    if (createUserDto.roleIds && createUserDto.roleIds.length > 0) {
      roles = await this.roleRepository.findBy({ id: In(createUserDto.roleIds) });
      if (roles.length !== createUserDto.roleIds.length) {
        throw new BadRequestException('Uno o más roles no existen');
      }
    } else {
      // Asignar rol de usuario por defecto
      const defaultRole = await this.roleRepository.findOne({ where: { type: RoleType.USER } });
      if (defaultRole) {
        roles = [defaultRole];
      }
    }

    // Obtener departamento
    let department: Department | undefined = undefined;
    if (createUserDto.departmentId) {
      const foundDepartment = await this.departmentRepository.findOne({
        where: { id: createUserDto.departmentId },
      });
      if (!foundDepartment) {
        throw new BadRequestException('El departamento no existe');
      }
      department = foundDepartment;
    }

    const user = this.userRepository.create({
      fullName: createUserDto.fullName,
      email: createUserDto.email,
      password: hashedPassword,
      identification: createUserDto.identification,
      phone: createUserDto.phone,
      avatarUrl: createUserDto.avatarUrl,
      roles,
      department,
      createdBy,
    });

    const savedUser = await this.userRepository.save(user);

    // Registrar auditoría
    await this.auditService.log({
      user: createdBy ? { id: createdBy } as User : null,
      action: AuditAction.CREATE,
      entityType: 'User',
      entityId: savedUser.id,
      description: `Usuario ${savedUser.email} creado`,
      newValues: { email: savedUser.email, fullName: savedUser.fullName },
    });

    return savedUser;
  }

  async findAll(query: QueryUserDto) {
    const { search, departmentId, roleId, isActive, page = 1, limit = 10 } = query;

    const queryBuilder = this.userRepository
      .createQueryBuilder('user')
      .leftJoinAndSelect('user.roles', 'role')
      .leftJoinAndSelect('user.department', 'department');

    // Búsqueda por texto
    if (search) {
      queryBuilder.andWhere(
        '(user.fullName ILIKE :search OR user.email ILIKE :search OR user.identification ILIKE :search)',
        { search: `%${search}%` },
      );
    }

    // Filtrar por departamento
    if (departmentId) {
      queryBuilder.andWhere('user.department.id = :departmentId', { departmentId });
    }

    // Filtrar por rol
    if (roleId) {
      queryBuilder.andWhere('role.id = :roleId', { roleId });
    }

    // Filtrar por estado activo
    if (isActive !== undefined) {
      queryBuilder.andWhere('user.isActive = :isActive', { isActive });
    }

    // Paginación
    const skip = (page - 1) * limit;
    queryBuilder.skip(skip).take(limit);

    // Ordenar por fecha de creación (más recientes primero)
    queryBuilder.orderBy('user.createdAt', 'DESC');

    const [users, total] = await queryBuilder.getManyAndCount();

    return {
      data: users,
      meta: {
        total,
        page,
        limit,
        totalPages: Math.ceil(total / limit),
      },
    };
  }

  async findOne(id: string): Promise<User> {
    const user = await this.userRepository.findOne({
      where: { id },
      relations: ['roles', 'department'],
    });

    if (!user) {
      throw new NotFoundException('Usuario no encontrado');
    }

    return user;
  }

  async update(id: string, updateUserDto: UpdateUserDto, updatedBy?: string): Promise<User> {
    const user = await this.findOne(id);

    const oldValues = { ...user };

    // Verificar si el email ya existe (si se está actualizando)
    if (updateUserDto.email && updateUserDto.email !== user.email) {
      const existingUser = await this.userRepository.findOne({
        where: { email: updateUserDto.email },
      });

      if (existingUser) {
        throw new ConflictException('El correo electrónico ya está registrado');
      }
    }

    // Verificar si la identificación ya existe (si se está actualizando)
    if (updateUserDto.identification && updateUserDto.identification !== user.identification) {
      const existingIdentification = await this.userRepository.findOne({
        where: { identification: updateUserDto.identification },
      });

      if (existingIdentification) {
        throw new ConflictException('La identificación ya está registrada');
      }
    }

    // Actualizar roles
    if (updateUserDto.roleIds) {
      const roles = await this.roleRepository.findBy({ id: In(updateUserDto.roleIds) });
      if (roles.length !== updateUserDto.roleIds.length) {
        throw new BadRequestException('Uno o más roles no existen');
      }
      user.roles = roles;
    }

    // Actualizar departamento
    if (updateUserDto.departmentId) {
      const department = await this.departmentRepository.findOne({
        where: { id: updateUserDto.departmentId },
      });
      if (!department) {
        throw new BadRequestException('El departamento no existe');
      }
      user.department = department;
    }

    // Actualizar campos
    Object.assign(user, {
      ...updateUserDto,
      updatedBy,
      roleIds: undefined,
      departmentId: undefined,
    });

    const updatedUser = await this.userRepository.save(user);

    // Registrar auditoría
    await this.auditService.log({
      user: updatedBy ? { id: updatedBy } as User : null,
      action: AuditAction.UPDATE,
      entityType: 'User',
      entityId: user.id,
      description: `Usuario ${user.email} actualizado`,
      oldValues: { email: oldValues.email, fullName: oldValues.fullName },
      newValues: { email: updatedUser.email, fullName: updatedUser.fullName },
    });

    return updatedUser;
  }

  async remove(id: string, deletedBy?: string): Promise<void> {
    const user = await this.findOne(id);

    // Soft delete
    await this.userRepository.softRemove(user);

    // Registrar auditoría
    await this.auditService.log({
      user: deletedBy ? { id: deletedBy } as User : null,
      action: AuditAction.DELETE,
      entityType: 'User',
      entityId: user.id,
      description: `Usuario ${user.email} desactivado (soft delete)`,
      oldValues: { email: user.email, fullName: user.fullName },
    });
  }

  async hardRemove(id: string, deletedBy?: string): Promise<void> {
    const user = await this.userRepository.findOne({
      where: { id },
      withDeleted: true,
    });

    if (!user) {
      throw new NotFoundException('Usuario no encontrado');
    }

    // Registrar auditoría antes de eliminar
    await this.auditService.log({
      user: deletedBy ? { id: deletedBy } as User : null,
      action: AuditAction.DELETE,
      entityType: 'User',
      entityId: user.id,
      description: `Usuario ${user.email} eliminado permanentemente (hard delete)`,
      oldValues: { email: user.email, fullName: user.fullName },
    });

    // Hard delete
    await this.userRepository.remove(user);
  }

  async restore(id: string): Promise<User> {
    const user = await this.userRepository.findOne({
      where: { id },
      withDeleted: true,
    });

    if (!user) {
      throw new NotFoundException('Usuario no encontrado');
    }

    if (!user.deletedAt) {
      throw new BadRequestException('El usuario no está desactivado');
    }

    user.deletedAt = undefined;
    return await this.userRepository.save(user);
  }
}

