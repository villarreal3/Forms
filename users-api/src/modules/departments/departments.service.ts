import { Injectable, NotFoundException, ConflictException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Department } from './entities/department.entity';
import { CreateDepartmentDto } from './dto/create-department.dto';
import { UpdateDepartmentDto } from './dto/update-department.dto';

@Injectable()
export class DepartmentsService {
  constructor(
    @InjectRepository(Department)
    private departmentRepository: Repository<Department>,
  ) {}

  async create(createDepartmentDto: CreateDepartmentDto): Promise<Department> {
    // Verificar si el nombre o código ya existe
    const existingDepartment = await this.departmentRepository.findOne({
      where: [
        { name: createDepartmentDto.name },
        { code: createDepartmentDto.code },
      ],
    });

    if (existingDepartment) {
      throw new ConflictException('El departamento o código ya existe');
    }

    const department = this.departmentRepository.create(createDepartmentDto);
    return await this.departmentRepository.save(department);
  }

  async findAll(): Promise<Department[]> {
    return await this.departmentRepository.find({
      order: { name: 'ASC' },
    });
  }

  async findOne(id: string): Promise<Department> {
    const department = await this.departmentRepository.findOne({
      where: { id },
      relations: ['users'],
    });

    if (!department) {
      throw new NotFoundException('Departamento no encontrado');
    }

    return department;
  }

  async update(id: string, updateDepartmentDto: UpdateDepartmentDto): Promise<Department> {
    const department = await this.findOne(id);

    // Verificar si el nuevo nombre o código ya existe
    if (updateDepartmentDto.name || updateDepartmentDto.code) {
      const existingDepartment = await this.departmentRepository.findOne({
        where: [
          { name: updateDepartmentDto.name },
          { code: updateDepartmentDto.code },
        ],
      });

      if (existingDepartment && existingDepartment.id !== id) {
        throw new ConflictException('El nombre o código del departamento ya existe');
      }
    }

    Object.assign(department, updateDepartmentDto);
    return await this.departmentRepository.save(department);
  }

  async remove(id: string): Promise<void> {
    const department = await this.findOne(id);
    await this.departmentRepository.softRemove(department);
  }
}

