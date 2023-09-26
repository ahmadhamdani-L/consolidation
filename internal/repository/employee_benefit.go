package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type EmployeeBenefit interface {
	Find(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*[]model.EmployeeBenefitEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.EmployeeBenefitEntityModel) (*model.EmployeeBenefitEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.EmployeeBenefitEntityModel) (*model.EmployeeBenefitEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.EmployeeBenefitEntityModel, error)
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

func (r *employeebenefit) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.EmployeeBenefitEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.EmployeeBenefitEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *employeebenefit) Update(ctx *abstraction.Context, id *int, e *model.EmployeeBenefitEntityModel) (*model.EmployeeBenefitEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil

}

func (r *employeebenefit) Create(ctx *abstraction.Context, e *model.EmployeeBenefitEntityModel) (*model.EmployeeBenefitEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name

	return e, nil
}

func (r *employeebenefit) Find(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*[]model.EmployeeBenefitEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.EmployeeBenefitEntityModel

	query := conn.Model(&model.EmployeeBenefitEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *employeebenefit) GetCount(ctx *abstraction.Context, m *model.EmployeeBenefitFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.EmployeeBenefitEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
