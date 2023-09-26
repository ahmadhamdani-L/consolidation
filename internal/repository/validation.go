package repository

import (
	"errors"
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"
	"worker-validation/pkg/constant"

	"gorm.io/gorm"
)

type Validation interface {
	Create(ctx *abstraction.Context, payload *model.ValidationDetailEntityModel) (data *model.ValidationDetailEntityModel, err error)
	FindByCriteria(ctx *abstraction.Context, filter *model.ValidationDetailFilterModel) (data *model.ValidationDetailEntityModel, err error)
	FindByID(ctx *abstraction.Context, id *int) (data *model.ValidationDetailEntityModel, err error)
	Update(ctx *abstraction.Context, id *int, payload *model.ValidationDetailEntityModel) (data *model.ValidationDetailEntityModel, err error)
	UpdateByCriteria(ctx *abstraction.Context, payload *model.ValidationDetailFilterModel, data *model.ValidationDetailEntityModel) (err error)
	FirstOrCreate(ctx *abstraction.Context, payload *model.ValidationDetailEntityModel) (data *model.ValidationDetailEntityModel, err error)
	UpdateStatus(ctx *abstraction.Context, tableName *string, id *int, status *int) error
	UpdateJurnal(ctx *abstraction.Context, trialBalanceID *int, status *int) error
	CountByCriteria(ctx *abstraction.Context, filter *model.ValidationDetailFilterModel) (count int64, err error)
}

type validation struct {
	abstraction.Repository
}

func NewValidation(db *gorm.DB) *validation {
	return &validation{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *validation) Create(ctx *abstraction.Context, payload *model.ValidationDetailEntityModel) (data *model.ValidationDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)

	err = conn.Create(&payload).Error
	if err != nil {
		return
	}
	err = conn.Model(&payload).First(&payload).Error
	if err != nil {
		return
	}
	return
}

func (r *validation) Update(ctx *abstraction.Context, id *int, payload *model.ValidationDetailEntityModel) (data *model.ValidationDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)

	err = conn.Model(&model.ValidationDetailEntityModel{}).Where("id = ?", &id).Updates(&payload).Error
	if err != nil {
		return
	}
	err = conn.Model(&payload).Where("id = ?", &id).First(&data).Error
	if err != nil {
		return
	}
	return
}

func (r *validation) UpdateByCriteria(ctx *abstraction.Context, payload *model.ValidationDetailFilterModel, data *model.ValidationDetailEntityModel) (err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ValidationDetailEntityModel{})
	query = r.Filter(ctx, query, *payload)
	err = query.Updates(&payload).Error
	if err != nil {
		return
	}
	return
}

func (r *validation) FindByCriteria(ctx *abstraction.Context, filter *model.ValidationDetailFilterModel) (data *model.ValidationDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ValidationDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *validation) FindByID(ctx *abstraction.Context, id *int) (data *model.ValidationDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ValidationDetailEntityModel{})
	// err = query.Where("id = ?", id).First(&data).Error
	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *validation) FirstOrCreate(ctx *abstraction.Context, payload *model.ValidationDetailEntityModel) (data *model.ValidationDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ValidationDetailEntityModel{})
	query = query.Where("name = ?", payload.Name).Where("company_id = ?", payload.CompanyID).Where("period = ?", payload.Period).Where("versions = ?", payload.Versions)
	if err = query.FirstOrCreate(&payload).First(&data).Error; err != nil {
		return
	}

	return
}

func (r *validation) UpdateStatus(ctx *abstraction.Context, tableName *string, id *int, status *int) error {
	conn := r.CheckTrx(ctx)
	sql := fmt.Sprintf("UPDATE %s SET status = ? WHERE id = ?", *tableName)
	if *tableName == "trial_balance" {
		if *status == constant.MODUL_STATUS_DRAFT {
			sql = fmt.Sprintf("UPDATE %s SET status = ?, validation_note = '%s' WHERE id = ?", *tableName, "Imbalance")
		} else if *status == constant.MODUL_STATUS_VALIDATED {
			sql = fmt.Sprintf("UPDATE %s SET status = ?, validation_note = '%s' WHERE id = ?", *tableName, "Balance")
		}
	}
	if err := conn.Exec(sql, status, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *validation) CountByCriteria(ctx *abstraction.Context, filter *model.ValidationDetailFilterModel) (count int64, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ValidationDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.Count(&count).Error; err != nil {
		return
	}
	return
}

func (r *validation) UpdateJurnal(ctx *abstraction.Context, trialBalanceID *int, status *int) error {
	conn := r.CheckTrx(ctx)

	err := conn.Model(&model.AdjustmentEntityModel{}).Where("tb_id = ?", &trialBalanceID).Update("status = ?", &status).Error
	if err != nil {
		return err
	}
	return nil
}
