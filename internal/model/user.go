package model

import (
	"os"
	"time"
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/constant"
	"worker-validation/pkg/util/date"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserEntity struct {
	Username     string `json:"username" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Password     string `json:"password" validate:"required" gorm:"-"`
	ImageProfile string `json:"image_profile" validate:"required"`
	CompanyID    int    `json:"company_id" validate:"required"`
	PasswordHash string `json:"-" gorm:"column:password"`
}

type UserEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`

	// entity
	UserEntity

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

func (UserEntityModel) TableName() string {
	return "users"
}

func (m *UserEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = constant.DB_DEFAULT_CREATED_BY

	m.hashPassword()
	m.Password = ""
	return
}

// func (m *UserEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
// 	m.ModifiedAt = date.DateTodayLocal()
// 	m.ModifiedBy = &m.Context.Auth.ID
// 	return
// }

func (m *UserEntityModel) hashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	m.PasswordHash = string(bytes)
}

func (m *UserEntityModel) GenerateToken() (string, error) {
	var (
		jwtKey = os.Getenv("JWT_KEY")
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         m.ID,
		"username":   m.Username,
		"company_id": m.CompanyID,
		"name":       m.Name,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, err
}
