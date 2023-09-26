package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type UserGetRequest struct {
	abstraction.Pagination
	model.UserFilterModel
}
type UserGetResponse struct {
	Datas          []model.UserEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type UserGetResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data []model.UserEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type UserGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type UserGetByIDResponse struct {
	model.UserEntityModel
}
type UserGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data UserGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type UserCreateRequest struct {
	Username     string `form:"username" validate:"required" example:"administrator"`
	Name         string `form:"name" validate:"required" example:"Lutfi Ramadhan"`
	Password     string `form:"password" validate:"required" gorm:"-" example:"nevemor3"`
	CompanyID    int    `form:"company_id" validate:"required" example:"1"`
	Email        string `form:"email" validate:"required" example:"admin@console.code"`
	RoleID       int    `form:"role_id" required:"required" example:"1"`
	IsActive     *bool  `form:"is_active" validate:"required" example:"true"`
	ImageProfile string `form:"-"`
}
type UserCreateResponse struct {
	model.UserEntityModel
}
type UserCreateResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data UserCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type UserUpdateRequest struct {
	ID           int    `param:"id" validate:"required,numeric"`
	Username     string `form:"username" example:"administrator"`
	Name         string `form:"name" example:"Lutfi Ramadhan"`
	Password     string `form:"password" gorm:"-" example:"nevemor3"`
	CompanyID    int    `form:"company_id" example:"1"`
	Email        string `form:"email" example:"admin@console.code"`
	RoleID       int    `form:"role_id" example:"1"`
	IsActive     *bool  `form:"is_active" example:"true"`
	ImageProfile string `form:"-"`
}
type UserUpdateResponse struct {
	model.UserEntityModel
}
type UserUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data UserUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type UserDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type UserDeleteResponse struct {
	model.UserEntityModel
}
type UserDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data UserDeleteResponse `json:"data"`
	} `json:"body"`
}

// User Forgot Password
type UserForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}
type UserForgotPasswordResponse struct {
	url string
}
type UserForgotPasswordResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data UserForgotPasswordRequest `json:"data"`
	} `json:"body"`
}

type UserResetPasswordRequest struct {
	ResetToken string `param:"resetToken"`
	Password   string `json:"password"`
}

type UserResetPasswordResponse struct {
	model.UserEntityModel
}

type UserResetPasswordResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data UserResetPasswordResponse `json:"data"`
	} `json:"body"`
}

type UserStatusRequest struct {
	UserID int `param:"user_id"`
}

type UserStatusResponse struct {
	model.UserEntityModel
}

type UserStatusResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data UserResetPasswordResponse `json:"data"`
	} `json:"body"`
}
