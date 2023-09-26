package repository

import (
	"errors"
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type EmployeeBenefitDetail interface {
	Find(ctx *abstraction.Context, m *model.EmployeeBenefitDetailFilterModel) (*[]model.EmployeeBenefitDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.EmployeeBenefitDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.EmployeeBenefitDetailFilterModel) (data *model.EmployeeBenefitDetailEntityModel, err error)
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

func (r *employeebenefitdetail) Find(ctx *abstraction.Context, m *model.EmployeeBenefitDetailFilterModel) (*[]model.EmployeeBenefitDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.EmployeeBenefitDetailEntityModel

	query := conn.Model(&model.EmployeeBenefitDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "AGING-UTANG-PIUTANG").Where("trx_ref_id = ?", m.EmployeeBenefitID))
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

func (r *employeebenefitdetail) FindByCriteria(ctx *abstraction.Context, filter *model.EmployeeBenefitDetailFilterModel) (data *model.EmployeeBenefitDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.EmployeeBenefitDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
