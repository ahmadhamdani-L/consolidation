package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type FormatterDetail interface {
	Find(ctx *abstraction.Context, m *model.FormatterDetailFilterModel) (*[]model.FormatterDetailEntityModel, error)
	FindGroup(ctx *abstraction.Context, m *model.FormatterDetailFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailEntityModel, *abstraction.PaginationInfo, error)
}

type formatterdetail struct {
	abstraction.Repository
}

func NewFormatterDetail(db *gorm.DB) *formatterdetail {
	return &formatterdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *formatterdetail) Find(ctx *abstraction.Context, m *model.FormatterDetailFilterModel) (*[]model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailEntityModel

	query := conn.Model(&model.FormatterDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	query = query.Order("sort_id asc")

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *formatterdetail) FindGroup(ctx *abstraction.Context, m *model.FormatterDetailFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("code != 'ASET' AND code != 'ASET_LANCAR' AND code != 'KAS_DAN_SETARA_KAS' AND code != 'KAS_DI_TANGAN' AND code != 'KAS_SETARA_KAS' AND code != 'Pihak_Ketiga:~PIUTANG_LAIN~LAIN - JANGKA PENDEK' AND code != 'ASET TIDAK LANCAR' AND code != 'LIABILITAS LANCAR' AND code != 'LIABILITAS TIDAK LANCAR' AND code != 'LIABILITAS' AND code != 'BEBAN PEMASARAN DAN DISTRIBUSI' AND code != 'BEBAN UMUM DAN ADMINISTRASI' AND code != 'PENDAPATAN DAN BEBAN LAIN-LAIN' AND code != 'PENDAPATAN LAIN-LAIN' AND code != 'BEBAN LAIN-LAIN' AND code != 'KAS DI TANGAN' AND code != 'KAS DI BANK'")

	//sort
	if p.Sort == nil {
		sort := "asc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "sort_id"
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)

	//pagination
	if p.Page == nil {
		page := 1
		p.Page = &page
	}
	if p.PageSize == nil {
		pageSize := 10
		p.PageSize = &pageSize
	}
	info = abstraction.PaginationInfo{
		Pagination: p,
	}
	limit := *p.PageSize
	offset := limit * (*p.Page - 1)
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(datas) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}