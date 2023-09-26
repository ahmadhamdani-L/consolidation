package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type PembelianPenjualanBerelasi interface {
	Find(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*[]model.PembelianPenjualanBerelasiEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*int64, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.PembelianPenjualanBerelasiEntityModel, error)
	Export(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*model.PembelianPenjualanBerelasiEntityModel, error)
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

func (r *pembelianpenjualanberelasi) FindByID(ctx *abstraction.Context, id *int) (*model.PembelianPenjualanBerelasiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas model.PembelianPenjualanBerelasiEntityModel

	query := conn.Model(&model.PembelianPenjualanBerelasiEntityModel{}).Where("id = ?", &id)

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

func (r *pembelianpenjualanberelasi) Export(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiFilterModel) (*model.PembelianPenjualanBerelasiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.PembelianPenjualanBerelasiEntityModel
	query := conn.Model(&data)
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Company").Preload("PembelianPenjualanBerelasiDetail.Company").Find(&data).Error; err != nil {
		return &data, err
	}

	return &data, nil
}
