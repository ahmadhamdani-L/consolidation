package repository

import (
	"notification/internal/abstraction"
	"notification/internal/model"

	"gorm.io/gorm"
)

type User interface {
	FindByID(ctx *abstraction.Context, id *int) (*model.UserEntityModel, error)
}

type user struct {
	abstraction.Repository
}

func NewUser(db *gorm.DB) *user {
	return &user{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *user) FindByID(ctx *abstraction.Context, id *int) (*model.UserEntityModel, error) {
	conn := r.checkTrx(ctx)

	var data model.UserEntityModel
	err := conn.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *user) checkTrx(ctx *abstraction.Context) *gorm.DB {
	if ctx.Trx != nil {
		return ctx.Trx.Db
	}
	return r.Db
}
