package repository

import (
	"errors"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type EmployeeBenefit interface {
	Find(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*[]model.EmployeeBenefitEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.ExportFilter) (data *model.EmployeeBenefitEntityModel, err error)
}

type employeebenefit struct {
	abstraction.Repository
}

func NewEmployeeBenefit(db *gorm.DB) *employeebenefit {
	return &employeebenefit{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *employeebenefit) Find(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*[]model.EmployeeBenefitEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.EmployeeBenefitEntityModel

	query := conn.Model(&model.EmployeeBenefitEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Company").Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *employeebenefit) FindByCriteria(ctx *abstraction.Context, filter *model.ExportFilter) (data *model.EmployeeBenefitEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.EmployeeBenefitEntityModel{})
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Version).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
