package request

type CreateUser struct {
	Name     string `json:"name" validate:"required,max=50" example:"fake name"`
	Email    string `json:"email" validate:"required,email,max=50" example:"fake@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20,password" example:"password1"`
	RoleId   uint   `json:"role_id" validate:"required,number,max=50" example:"1"`
	Role     string `json:"role" validate:"required,oneof=user admin,max=50" example:"user"`
}

type UpdateUser struct {
	ID       string `json:"-"`
	Name     string `json:"name,omitempty" validate:"omitempty,max=50" example:"fake name"`
	Email    string `json:"email" validate:"omitempty,email,max=50" example:"fake@example.com"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=user admin,max=50" example:"user"`
	RoleId   uint   `json:"role_id,omitempty" validate:"omitempty,number,max=50" example:"1"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8,max=20,password" example:"password1"`
}

type GoogleLogin struct {
	Name          string `json:"name" validate:"required,max=50"`
	Email         string `json:"email" validate:"required,email,max=50"`
	VerifiedEmail bool   `json:"verified_email" validate:"required"`
}

type UpdatePassOrVerify struct {
	Password      string `json:"password,omitempty" validate:"omitempty,min=8,max=20,password" example:"password1"`
	VerifiedEmail bool   `json:"verified_email" swaggerignore:"true" validate:"omitempty,boolean"`
}

type QueryUser struct {
	Page   int    `validate:"omitempty,number,max=50"`
	Limit  int    `validate:"omitempty,number,max=50"`
	Search string `validate:"omitempty,max=50"`
}
