package repositories

import "forms/models"

// RoleRepository define las operaciones para la tabla roles (copia local para el frontend)
type RoleRepository interface {
	FindAll() ([]models.Role, error)
	Create(name, roleType string) (int64, error)
	ExistsByNameAndType(name, roleType string) (bool, error)
}
