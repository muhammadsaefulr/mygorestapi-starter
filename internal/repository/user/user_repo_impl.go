package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type newUserRepositryImpl struct {
	DB *gorm.DB
}

func NewUserRepositryImpl(db *gorm.DB) UserRepo {
	return &newUserRepositryImpl{
		DB: db,
	}
}

// GetAllUser implements UserRepo.
func (n *newUserRepositryImpl) GetAllUser(ctx context.Context, param *request.QueryUser) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := n.DB.WithContext(ctx).Model(&model.User{}).Order("created_at asc")
	offset := (param.Page - 1) * param.Limit

	if param.Search != "" {
		searchLike := "%" + param.Search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR role LIKE ?", searchLike, searchLike, searchLike)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Page).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (n *newUserRepositryImpl) CreateUser(ctx context.Context, user *model.User) error {
	result := n.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateUser implements UserRepo.
func (n *newUserRepositryImpl) UpdateUser(ctx context.Context, user *model.User) error {

	users := &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
	}

	result := n.DB.WithContext(ctx).Where("id = ?", user.ID).Updates(users)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserByEmail implements UserRepo.
func (n *newUserRepositryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)

	result := n.DB.WithContext(ctx).Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// GetUserByID implements UserRepo.
func (n *newUserRepositryImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := new(model.User)

	result := n.DB.WithContext(ctx).Where("id = ?", id).Preload("UserRole").
		Preload("UserRole.Permissions").First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// DeleteUser implements UserRepo.
func (n *newUserRepositryImpl) DeleteUser(ctx context.Context, id string) error {
	err := n.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	return nil
}
