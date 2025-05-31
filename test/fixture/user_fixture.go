package fixture

import (
	user_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model/user"

	"github.com/google/uuid"
)

var UserOne = &user_model.User{
	ID:            uuid.New(),
	Name:          "Test1",
	Email:         "test1@gmail.com",
	Password:      "password1",
	Role:          "user",
	VerifiedEmail: false,
}

var UserTwo = &user_model.User{
	ID:            uuid.New(),
	Name:          "Test2",
	Email:         "test2@gmail.com",
	Password:      "password1",
	Role:          "user",
	VerifiedEmail: false,
}

var Admin = &user_model.User{
	ID:            uuid.New(),
	Name:          "Admin",
	Email:         "admin@gmail.com",
	Password:      "password1",
	Role:          "admin",
	VerifiedEmail: false,
}
