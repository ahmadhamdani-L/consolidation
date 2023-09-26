package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type PembelianPenjualanBerelasi interface {
	Find(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*[]model.PembelianPenjualanBerelasiEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.PembelianPenjualanBerelasiEntityModel) (*model.PembelianPenjualanBerelasiEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.PembelianPenjualanBerelasiEntityModel) (*model.PembelianPenjualanBerelasiEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.PembelianPenjualanBerelasiEntityModel, error)
}

type pembelianpenjualanberelasi struct {
	abstraction.Repository
}

func NewPembelianPenjualanBerelasi(db *gorm.DB) *pembelianpenjualanberelasi {
	return &pembelianpenjualanberelasi{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *pembelianpenjualanberelasi) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.PembelianPenjualanBerelasiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.PembelianPenjualanBerelasiEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *pembelianpenjualanberelasi) Update(ctx *abstraction.Context, id *int, e *model.PembelianPenjualanBerelasiEntityModel) (*model.PembelianPenjualanBerelasiEntityModel, error) {
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

func (r *pembelianpenjualanberelasi) Create(ctx *abstraction.Context, e *model.PembelianPenjualanBerelasiEntityModel) (*model.PembelianPenjualanBerelasiEntityModel, error) {
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

func (r *pembelianpenjualanberelasi) Find(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*[]model.PembelianPenjualanBerelasiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.PembelianPenjualanBerelasiEntityModel

	query := conn.Model(&model.PembelianPenjualanBerelasiEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *pembelianpenjualanberelasi) GetCount(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.PembelianPenjualanBerelasiEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
