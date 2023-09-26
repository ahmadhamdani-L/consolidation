package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type PembelianPenjualanBerelasiDetail interface {
	Find(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiDetailFilterModel, p *abstraction.Pagination) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.PembelianPenjualanBerelasiDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.PembelianPenjualanBerelasiDetailEntityModel) (*model.PembelianPenjualanBerelasiDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.PembelianPenjualanBerelasiDetailEntityModel) (*model.PembelianPenjualanBerelasiDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.PembelianPenjualanBerelasiDetailEntityModel) (*model.PembelianPenjualanBerelasiDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, error)
	GetTotal(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiDetailFilterModel) (float64, float64, error)
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

func (r *pembelianpenjualanberelasidetail) Find(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiDetailFilterModel, p *abstraction.Pagination) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.PembelianPenjualanBerelasiDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.PembelianPenjualanBerelasiDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	//sort
	if p.Sort == nil {
		sort := "asc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "id"
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

	if err := query.Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *pembelianpenjualanberelasidetail) FindByID(ctx *abstraction.Context, id *int) (*model.PembelianPenjualanBerelasiDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.PembelianPenjualanBerelasiDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *pembelianpenjualanberelasidetail) Create(ctx *abstraction.Context, e *model.PembelianPenjualanBerelasiDetailEntityModel) (*model.PembelianPenjualanBerelasiDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("PembelianPenjualanBerelasi").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *pembelianpenjualanberelasidetail) Update(ctx *abstraction.Context, id *int, e *model.PembelianPenjualanBerelasiDetailEntityModel) (*model.PembelianPenjualanBerelasiDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Preload("PembelianPenjualanBerelasi").WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("PembelianPenjualanBerelasi").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *pembelianpenjualanberelasidetail) Delete(ctx *abstraction.Context, id *int, e *model.PembelianPenjualanBerelasiDetailEntityModel) (*model.PembelianPenjualanBerelasiDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *pembelianpenjualanberelasidetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.PembelianPenjualanBerelasiDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.PembelianPenjualanBerelasiDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *pembelianpenjualanberelasidetail) GetTotal(ctx *abstraction.Context, m *model.PembelianPenjualanBerelasiDetailFilterModel) (float64, float64, error) {
	conn := r.CheckTrx(ctx)

	type Total struct {
		TotalPembelian float64
		TotalPenjualan float64
	}
	var data Total

	query := conn.Model(&model.PembelianPenjualanBerelasiDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Select("SUM(bought_amount) total_pembelian, SUM(sales_amount) total_penjualan").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return 0, 0, err
	}

	return data.TotalPembelian, data.TotalPenjualan, nil
}
