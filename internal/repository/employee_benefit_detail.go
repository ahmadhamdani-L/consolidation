package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type EmployeeBenefitDetail interface {
	Find(ctx *abstraction.Context, m *model.EmployeeBenefitDetailFilterModel) (*[]model.EmployeeBenefitDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.EmployeeBenefitDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.EmployeeBenefitDetailEntityModel) (*model.EmployeeBenefitDetailEntityModel, error)
}

type employeebenefitdetail struct {
	abstraction.Repository
}

func NewEmployeeBenefitDetail(db *gorm.DB) *employeebenefitdetail {
	return &employeebenefitdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *employeebenefitdetail) Create(ctx *abstraction.Context, e *model.EmployeeBenefitDetailEntityModel) (*model.EmployeeBenefitDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *employeebenefitdetail) Find(ctx *abstraction.Context, m *model.EmployeeBenefitDetailFilterModel) (*[]model.EmployeeBenefitDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.EmployeeBenefitDetailEntityModel

	query := conn.Model(&model.EmployeeBenefitDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *employeebenefitdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.EmployeeBenefitDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.EmployeeBenefitDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
