package services

import (
	"log"
	"net/http"

	"forms/handlers/helpers"
	"forms/handlers/repositories"
)

// CreateRoleRequest body para crear o sincronizar un rol
type CreateRoleRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// GetRoles lista roles desde la tabla local (para cargar en el select del frontend)
func GetRoles(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodGet) {
		return
	}

	list, err := repositories.RoleRepo.FindAll()
	if err != nil {
		log.Printf("Error listando roles: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo roles")
		return
	}

	out := make([]map[string]interface{}, 0, len(list))
	for _, role := range list {
		out = append(out, map[string]interface{}{
			"id":   role.ID,
			"name": role.Name,
			"type": role.Type,
		})
	}

	helpers.RespondJSON(w, http.StatusOK, out)
}

// CreateRole crea un rol en la tabla local (sincronización desde users-api)
func CreateRole(w http.ResponseWriter, r *http.Request) {
	if !helpers.EnsureMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateRoleRequest
	if !helpers.DecodeJSON(w, r, &req) {
		return
	}
	if req.Name == "" || req.Type == "" {
		helpers.RespondError(w, http.StatusBadRequest, "name y type son obligatorios")
		return
	}

	exists, err := repositories.RoleRepo.ExistsByNameAndType(req.Name, req.Type)
	if err != nil {
		log.Printf("Error comprobando rol: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error al crear rol")
		return
	}
	if exists {
		helpers.RespondError(w, http.StatusConflict, "El rol ya existe")
		return
	}

	_, err = repositories.RoleRepo.Create(req.Name, req.Type)
	if err != nil {
		log.Printf("Error creando rol: %v", err)
		helpers.RespondError(w, http.StatusInternalServerError, "Error al crear rol")
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Rol creado",
	})
}
