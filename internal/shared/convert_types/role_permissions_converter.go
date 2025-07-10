package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/role_permissions/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateRolePermissionsToModel(req *request.CreateRolePermissions) *model.RolePermissions {
	return &model.RolePermissions{
		PermissionName: req.Name,
	}
}

func UpdateRolePermissionsToModel(req *request.UpdateRolePermissions) *model.RolePermissions {
	return &model.RolePermissions{
		PermissionName: req.Name,
	}
}
