package repositories

import (
	"forms/database"
	"forms/models"
)

type mysqlRoleRepository struct{}

func NewRoleRepository() RoleRepository {
	return &mysqlRoleRepository{}
}

func (r *mysqlRoleRepository) FindAll() ([]models.Role, error) {
	rows, err := database.DB.Query("SELECT id, name, type, created_at FROM v_roles ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Type, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, role)
	}
	return list, nil
}

func (r *mysqlRoleRepository) Create(name, roleType string) (int64, error) {
	res, err := database.DB.Exec(
		"INSERT INTO roles (name, type) VALUES (?, ?)",
		name, roleType,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *mysqlRoleRepository) ExistsByNameAndType(name, roleType string) (bool, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(1) FROM roles WHERE name = ? AND type = ?", name, roleType).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
