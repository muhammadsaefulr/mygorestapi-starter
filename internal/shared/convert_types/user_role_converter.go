package convert_types

import (
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/user_role/request"
	model "github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/model"
)

func CreateUserRoleToModel(req *request.CreateUserRole) *model.UserRole {
	permissions := make([]model.RolePermissions, len(req.Permission))
	for i, id := range req.Permission {
		permissions[i] = model.RolePermissions{ID: id}
	}

	return &model.UserRole{
		RoleName:    req.Name,
		Permissions: permissions,
	}
}

func UpdateUserRoleToModel(req *request.UpdateUserRole) *model.UserRole {

	permissions := make([]model.RolePermissions, len(req.Permission))
	for i, id := range req.Permission {
		permissions[i] = model.RolePermissions{ID: id}
	}

	return &model.UserRole{
		RoleName:    req.Name,
		Permissions: permissions,
	}
}
