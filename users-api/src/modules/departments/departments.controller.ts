import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  UseGuards,
  HttpCode,
  HttpStatus,
} from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
  ApiParam,
} from '@nestjs/swagger';
import { DepartmentsService } from './departments.service';
import { CreateDepartmentDto } from './dto/create-department.dto';
import { UpdateDepartmentDto } from './dto/update-department.dto';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { RolesGuard } from '../auth/guards/roles.guard';
import { Roles } from '../auth/decorators/roles.decorator';
import { RoleType } from '../roles/entities/role.entity';
import { Department } from './entities/department.entity';

@ApiTags('Departments')
@Controller('departments')
@UseGuards(JwtAuthGuard, RolesGuard)
@ApiBearerAuth('JWT-auth')
export class DepartmentsController {
  constructor(private readonly departmentsService: DepartmentsService) {}

  @Post()
  @Roles(RoleType.SUPER_ADMIN, RoleType.ADMIN)
  @ApiOperation({ summary: 'Crear nuevo departamento' })
  @ApiResponse({ status: 201, description: 'Departamento creado exitosamente', type: Department })
  @ApiResponse({ status: 409, description: 'El departamento ya existe' })
  async create(@Body() createDepartmentDto: CreateDepartmentDto) {
    return this.departmentsService.create(createDepartmentDto);
  }

  @Get()
  @ApiOperation({ summary: 'Listar todos los departamentos' })
  @ApiResponse({ status: 200, description: 'Lista de departamentos', type: [Department] })
  async findAll() {
    return this.departmentsService.findAll();
  }

  @Get(':id')
  @ApiOperation({ summary: 'Obtener departamento por ID' })
  @ApiParam({ name: 'id', description: 'ID del departamento' })
  @ApiResponse({ status: 200, description: 'Departamento encontrado', type: Department })
  @ApiResponse({ status: 404, description: 'Departamento no encontrado' })
  async findOne(@Param('id') id: string) {
    return this.departmentsService.findOne(id);
  }

  @Patch(':id')
  @Roles(RoleType.SUPER_ADMIN, RoleType.ADMIN)
  @ApiOperation({ summary: 'Actualizar departamento' })
  @ApiParam({ name: 'id', description: 'ID del departamento' })
  @ApiResponse({ status: 200, description: 'Departamento actualizado', type: Department })
  @ApiResponse({ status: 404, description: 'Departamento no encontrado' })
  async update(@Param('id') id: string, @Body() updateDepartmentDto: UpdateDepartmentDto) {
    return this.departmentsService.update(id, updateDepartmentDto);
  }

  @Delete(':id')
  @Roles(RoleType.SUPER_ADMIN, RoleType.ADMIN)
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiOperation({ summary: 'Eliminar departamento' })
  @ApiParam({ name: 'id', description: 'ID del departamento' })
  @ApiResponse({ status: 204, description: 'Departamento eliminado exitosamente' })
  @ApiResponse({ status: 404, description: 'Departamento no encontrado' })
  async remove(@Param('id') id: string) {
    await this.departmentsService.remove(id);
  }
}

