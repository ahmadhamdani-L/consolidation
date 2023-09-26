package repository

import (
	"errors"
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutang interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*[]model.AgingUtangPiutangEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.AgingUtangPiutangEntityModel, err error)
	Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
}

type agingutangpiutang struct {
	abstraction.Repository
}

func NewAgingUtangPiutang(db *gorm.DB) *agingutangpiutang {
	return &agingutangpiutang{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *agingutangpiutang) Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*[]model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AgingUtangPiutangEntityModel

	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *agingutangpiutang) FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.AgingUtangPiutangEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Versions).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *agingutangpiutang) Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(e)
	fmt.Println(query)
	if err := query.Where("id = ?", id).Updates(&e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(&e).First(&e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
