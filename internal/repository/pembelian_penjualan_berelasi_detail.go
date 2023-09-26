package repository

import (
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type PembelianPenjualanBerelasiDetail interface {
	Find(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiDetailFilterModel) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, error)
}

type pembelianpenjualanberelasidetail struct {
	abstraction.Repository
}

func NewPembelianPenjualanBerelasiDetail(db *gorm.DB) *pembelianpenjualanberelasidetail {
	return &pembelianpenjualanberelasidetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *pembelianpenjualanberelasidetail) Find(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiDetailFilterModel) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.PembelianPenjualanBerelasiDetailEntityModel

	query := conn.Model(&model.PembelianPenjualanBerelasiDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *pembelianpenjualanberelasidetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.PembelianPenjualanBerelasiDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
