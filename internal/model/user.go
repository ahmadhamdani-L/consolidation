package model

import (
	"encoding/base64"
	"mcash-finance-console-core/configs"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserEntity struct {
	Username     string `json:"username" validate:"required" example:"administrator"`
	Name         string `json:"name" validate:"required" example:"Lutfi Ramadhan"`
	Password     string `json:"password" validate:"required" gorm:"-" example:"nevemor3"`
	ImageProfile string `json:"image_profile" validate:"required"`
	CompanyID    int    `json:"company_id" validate:"required" example:"1"`
	PasswordHash string `json:"-" gorm:"column:password"`
	Email        string `json:"email" validate:"required" example:"admin@console.code"`
	RoleID       int    `json:"role_id" required:"required" example:"1"`
	IsActive     *bool  `json:"is_active" validate:"required" gorm:"default:true" example:"true"`
}

type UserFilter struct {
	Username          *string    `query:"username" filter:"ILIKE"`
	Name              *string    `query:"name" filter:"ILIKE"`
	Email             *string    `query:"email" filter:"ILIKE" example:"admin@console.code"`
	RoleID            *int       `query:"role_id" example:"1"`
	IsActive          *bool      `query:"is_active"`
	ArrRoleID         *[]int     `filter:"CUSTOM"`
	CreatedAt         *time.Time `query:"created_at" filter:"DATE" example:"2022-08-17T15:04:05Z"`
	CreatedBy         *int       `query:"created_by" example:"1"`
	UserCreatedString *string    `query:"user_created" filter:"CUSTOM" example:"Lutfi Ramadhan"`
	Search            *string    `query:"search" filter:"CUSTOM" example:"Lutfi Ramadhan"`
}

type UserEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	CreatedAt         time.Time        `json:"created_at"`
	UserCreated       *UserEntityModel `json:"-" gorm:"foreignKey:CreatedBy"`
	UserCreatedString string           `json:"user_created" gorm:"-"`
	CreatedBy         int              `json:"created_by"`

	// entity
	UserEntity
	Role    *RoleEntityModel    `json:"role" gorm:"foreignKey:RoleID"`
	Company *CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type UserFilterModel struct {
	// abstraction
	// abstraction.Filter

	// filter
	UserFilter
	CompanyCustomFilter
}

func (UserEntityModel) TableName() string {
	return "users"
}

func (m *UserEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = 1

	m.hashPassword()
	m.Password = ""
	return
}

func (m *UserEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.Password != "" {
		m.hashPassword()
		m.Password = ""
	}
	return
}

func (m *UserEntityModel) hashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	m.PasswordHash = string(bytes)
}

func (m *UserEntityModel) GenerateToken() (string, error) {
	jwtKey := configs.Jwt().SecretKey()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": m.ID,
		// "username":   m.Username,
		"cid": m.CompanyID,
		// "name":       m.Name,
		"rid": m.RoleID,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, err
}

func (m *UserEntityModel) GenerateTokenResetPassword() (string, error) {
	jwtKey := configs.Jwt().SecretKey()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    m.ID,
		"email": m.Email,
		"name":  m.Name,
		// "phone":      m.Phone,
		"company_id": m.CompanyID,
		"exp":        time.Now().Add(time.Minute * 59).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}
