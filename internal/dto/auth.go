package dto

import (
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Login
type AuthLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type AuthLoginResponse struct {
	Token                      string `json:"token"`
	NotificationAuthToken      string `json:"na_token"`
	NotificationSubscribeToken string `json:"ns_token"`
	model.UserEntityModel
}
type AuthLoginResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data AuthLoginResponse `json:"data"`
	} `json:"body"`
}

// Register
type AuthRegisterRequest struct {
	model.UserEntity
}
type AuthRegisterResponse struct {
	model.UserEntityModel
}
type AuthRegisterResponseDoc struct {
	Body struct {
		Meta res.Meta             `json:"meta"`
		Data AuthRegisterResponse `json:"data"`
	} `json:"body"`
}

// CheckAuth
type CheckAuthResponse struct {
	User        model.UserEntityModel              `json:"user"`
	Permission  *[]model.RolePermissionEntityModel `json:"permission"`
	AccessScope *model.AccessScopeEntityModel      `json:"access_scope"`
}
type CheckAuthResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data CheckAuthResponse `json:"data"`
	} `json:"body"`
}

//Change Password

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type GetNotificationTokenResponse struct {
	NotificationAuthToken      string `json:"na_token"`
	NotificationSubscribeToken string `json:"ns_token"`
}
