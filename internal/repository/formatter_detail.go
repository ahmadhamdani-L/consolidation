package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type FormatterDetail interface {
	Find(ctx *abstraction.Context, m *model.FormatterDetailFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.FormatterDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.FormatterDetailEntityModel) (*model.FormatterDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.FormatterDetailEntityModel) (*model.FormatterDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.FormatterDetailEntityModel) (*model.FormatterDetailEntityModel, error)
	FindByCode(ctx *abstraction.Context, formatterID *int, code *string) (*model.FormatterDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, formatterID *int) (*[]model.FormatterDetailEntityModel, error)
	Finds(ctx *abstraction.Context, m *model.FormatterDetailFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailFmtEntityModel, *abstraction.PaginationInfo, error)
	FindWithCriteria(ctx *abstraction.Context, m *model.FormatterDetailFilterModel) (*[]model.FormatterDetailEntityModel, error)
	FindByCriteriaControl(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error)
	FindCoaListing(ctx *abstraction.Context) (*[]model.CoaEntityModel, error)
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
func (r *formatterdetail) FindCoaListing(ctx *abstraction.Context) (*[]model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CoaEntityModel

	query := conn.Raw(`
	SELECT m_coa.code
	FROM m_coa
	LEFT JOIN m_formatter_detail ON 
		(LEFT(m_coa.code::text, 6) = m_formatter_detail.code::text OR 
		 LEFT(m_coa.code::text, 3) = m_formatter_detail.code::text OR 
		 LEFT(m_coa.code::text, 4) = m_formatter_detail.code::text)
		AND m_formatter_detail.formatter_id = 3
	WHERE m_formatter_detail.code IS NULL
	AND LENGTH(m_coa.code) > 6;
	`)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
func (r *formatterdetail) FindByCriteriaControl(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ControllerEntityModel

	query := conn.Model(&model.ControllerEntityModel{})

	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
func (r *formatterdetail) FindWithCriteria(ctx *abstraction.Context, m *model.FormatterDetailFilterModel) (*[]model.FormatterDetailEntityModel, error) {
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

func (r *formatterdetail) Finds(ctx *abstraction.Context, m *model.FormatterDetailFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailFmtEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailFmtEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterDetailFmtEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_id = 3 AND is_total != true AND is_show_view = true")

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

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *formatterdetail) Find(ctx *abstraction.Context, m *model.FormatterDetailFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

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

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *formatterdetail) FindByID(ctx *abstraction.Context, id *int) (*model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatterdetail) Create(ctx *abstraction.Context, e *model.FormatterDetailEntityModel) (*model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *formatterdetail) Update(ctx *abstraction.Context, id *int, e *model.FormatterDetailEntityModel) (*model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *formatterdetail) Delete(ctx *abstraction.Context, id *int, e *model.FormatterDetailEntityModel) (*model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *formatterdetail) FindByCode(ctx *abstraction.Context, formatterID *int, code *string) (*model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterDetailEntityModel
	if err := conn.Where("formatter_id = ?", &formatterID).Where("code = ?", &code).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatterdetail) FindSummary(ctx *abstraction.Context, formatterID *int) (*[]model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailEntityModel
	if err := conn.Where("(auto_summary = true OR is_total = true) AND formatter_id = ?", &formatterID).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}
	return &datas, nil
}
