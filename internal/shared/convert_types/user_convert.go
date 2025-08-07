package convert_types

import (
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/user/request"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/user/response"
	model "github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/model"
)

func CreateUserToUserModel(user *request.CreateUser) *model.User {
	return &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		RoleId:   user.RoleId,
	}
}

func UpdateUserToUserModel(user *request.UpdateUser) *model.User {
	return &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		RoleId:   user.RoleId,
		Password: user.Password,
	}
}

func UpdatePassOrVerifyToUserModel(user *request.UpdatePassOrVerify) *model.User {
	return &model.User{
		Password:      user.Password,
		VerifiedEmail: user.VerifiedEmail,
	}
}

func UserResponseToUserModel(user *model.User) *model.User {
	return &model.User{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Role:          user.Role,
		VerifiedEmail: user.VerifiedEmail,
	}
}

func UserModelToUserResponse(user *model.User) *response.GetUsersResponse {
	return &response.GetUsersResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Role:            user.Role,
		Roles:           user.UserRole,
		IsEmailVerified: user.VerifiedEmail,
	}
}
