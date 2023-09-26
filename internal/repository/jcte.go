package repository

import (
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"
)

type Jcte interface {
	Find(ctx *abstraction.Context, m *model.JcteFilterModel) (*[]model.JcteEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JcteEntityModel, error)
	Export(ctx *abstraction.Context, e *model.JcteFilterModel) (*model.JcteEntityModel, error)
}

type jcte struct {
	abstraction.Repository
}

func NewJcte(db *gorm.DB) *jcte {
	return &jcte{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jcte) Find(ctx *abstraction.Context, m *model.JcteFilterModel) (*[]model.JcteEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JcteEntityModel

	query := conn.Model(&model.JcteEntityModel{})

	// filter
	query = r.Filter(ctx, query, *m)
	err := query.Preload("JcteDetail").Find(&datas).Error
	if err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jcte) FindByID(ctx *abstraction.Context, id *int) (*model.JcteEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JcteEntityModel

	err := conn.Where("id = ?", id).Preload("JcteDetail").First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jcte) Export(ctx *abstraction.Context, e *model.JcteFilterModel) (*model.JcteEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JcteEntityModel
	query := conn.Model(&model.JcteEntityModel{}).Preload("Company").Preload("JcteDetail").Where("tb_id = ?", &e.TbID).Where("period = ?", &e.Period).Where("company_id = ?", &e.CompanyID).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}
