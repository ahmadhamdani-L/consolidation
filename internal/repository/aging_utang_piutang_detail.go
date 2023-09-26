package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutangDetail interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel) (*[]model.AgingUtangPiutangDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AgingUtangPiutangDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error)
}

type agingutangpiutangdetail struct {
	abstraction.Repository
}

func NewAgingUtangPiutangDetail(db *gorm.DB) *agingutangpiutangdetail {
	return &agingutangpiutangdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *agingutangpiutangdetail) Create(ctx *abstraction.Context, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return e, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *agingutangpiutangdetail) Find(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel) (*[]model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AgingUtangPiutangDetailEntityModel

	query := conn.Model(&model.AgingUtangPiutangDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *agingutangpiutangdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.AgingUtangPiutangDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
