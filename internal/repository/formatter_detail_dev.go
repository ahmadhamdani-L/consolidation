package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type FormatterDetailDev interface {
	Find(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailDevEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.FormatterDetailDevEntityModel, error)
	Create(ctx *abstraction.Context, e *model.FormatterDetailDevEntityModel) (*model.FormatterDetailDevEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.FormatterDetailDevEntityModel) (*model.FormatterDetailDevEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.FormatterDetailDevEntityModel) (*model.FormatterDetailDevEntityModel, error)
	FindByCode(ctx *abstraction.Context, formatterID *int, code *string) (*model.FormatterDetailDevEntityModel, error)
	FindSummary(ctx *abstraction.Context, formatterID *int) (*[]model.FormatterDetailDevEntityModel, error)
	Finds(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailDevFmtEntityModel, *abstraction.PaginationInfo, error)
	FindWithCriteria(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel) (*[]model.FormatterDetailDevEntityModel, error)
	FindByCriteriaControl(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error)
	FindCoaListing(ctx *abstraction.Context) (*[]model.CoaEntityModel, error)
	FindWithTotal(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailDevFmtEntityModel, *abstraction.PaginationInfo, error)
}

type formatterdetaildev struct {
	abstraction.Repository
}

func NewFormatterDetailDev(db *gorm.DB) *formatterdetaildev {
	return &formatterdetaildev{
		abstraction.Repository{
			Db: db,
		},
	}
}
func (r *formatterdetaildev) FindCoaListing(ctx *abstraction.Context) (*[]model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CoaEntityModel

	query := conn.Raw(`
	SELECT m_coa.code
	FROM m_coa
	LEFT JOIN m_formatter_detail ON 
		(LEFT(m_coa.code::text, 6) = m_formatter_detail.code::text OR 
		 LEFT(m_coa.code::text, 3) = m_formatter_detail.code::text OR 
		 LEFT(m_coa.code::text, 4) = m_formatter_detail.code::text)
		AND m_formatter_detail.formatter_dev_id = 3
	WHERE m_formatter_detail.code IS NULL
	AND LENGTH(m_coa.code) > 6;
	`)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
func (r *formatterdetaildev) FindByCriteriaControl(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error) {
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
func (r *formatterdetaildev) FindWithCriteria(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel) (*[]model.FormatterDetailDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailDevEntityModel

	query := conn.Model(&model.FormatterDetailDevEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	query = query.Order("sort_id asc")

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *formatterdetaildev) Finds(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailDevFmtEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailDevFmtEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterDetailDevFmtEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_dev_id = 3 AND is_total != true AND is_show_view = true")
	// query = query.Where("formatter_dev_id = 3 AND is_show_view = true")

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

func (r *formatterdetaildev) FindWithTotal(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailDevFmtEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailDevFmtEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterDetailDevFmtEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	// query = query.Where("formatter_dev_id = 3 AND is_total != true AND is_show_view = true")
	query = query.Where("formatter_dev_id = 3 AND is_show_view = true")

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

func (r *formatterdetaildev) Find(ctx *abstraction.Context, m *model.FormatterDetailDevFilterModel, p *abstraction.Pagination) (*[]model.FormatterDetailDevEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailDevEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterDetailDevEntityModel{})
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

func (r *formatterdetaildev) FindByID(ctx *abstraction.Context, id *int) (*model.FormatterDetailDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterDetailDevEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatterdetaildev) Create(ctx *abstraction.Context, e *model.FormatterDetailDevEntityModel) (*model.FormatterDetailDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *formatterdetaildev) Update(ctx *abstraction.Context, id *int, e *model.FormatterDetailDevEntityModel) (*model.FormatterDetailDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *formatterdetaildev) Delete(ctx *abstraction.Context, id *int, e *model.FormatterDetailDevEntityModel) (*model.FormatterDetailDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *formatterdetaildev) FindByCode(ctx *abstraction.Context, formatterID *int, code *string) (*model.FormatterDetailDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterDetailDevEntityModel
	if err := conn.Where("formatter_dev_id = ?", &formatterID).Where("code = ?", &code).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatterdetaildev) FindSummary(ctx *abstraction.Context, formatterID *int) (*[]model.FormatterDetailDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailDevEntityModel
	if err := conn.Where("(auto_summary = true OR is_total = true) AND formatter_dev_id = ?", &formatterID).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}
	return &datas, nil
}
