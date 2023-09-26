package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type Notification interface {
	Create(ctx *abstraction.Context, e *model.NotificationEntityModel) (*model.NotificationEntityModel, error)
}

type notification struct {
	abstraction.Repository
}

func NewNotification(db *gorm.DB) *notification {
	return &notification{
		abstraction.Repository{
			Db: db,
		},
	}
}
func (r *notification) Create(ctx *abstraction.Context, e *model.NotificationEntityModel) (*model.NotificationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}
