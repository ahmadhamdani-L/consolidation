package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"strings"

	"gorm.io/gorm"
)

type CoaDev interface {
	Find(ctx *abstraction.Context, m *model.CoaDevFilterModel, p *abstraction.Pagination) (*[]model.CoaDevEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.CoaDevEntityModel, error)
	Create(ctx *abstraction.Context, e *model.CoaDevEntityModel) (*model.CoaDevEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.CoaDevEntityModel) (*model.CoaDevEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.CoaDevEntityModel) (*model.CoaDevEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.CoaDevEntityModel, error)
	Export(ctx *abstraction.Context) (*[]model.CoaGroupEntityModel, error)
	CountDash(ctx *abstraction.Context) (int, error)
}

type coadev struct {
	abstraction.Repository
}

func NewCoaDev(db *gorm.DB) *coadev {
	return &coadev{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *coadev) Find(ctx *abstraction.Context, m *model.CoaDevFilterModel, p *abstraction.Pagination) (*[]model.CoaDevEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CoaDevEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.CoaDevEntityModel{})
	//filter
	tableName := model.CoaDevEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, model.CoaDevEntityModel{}.TableName())
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
	if len(*m.ArrCoaGroupID) > 0 {
		listCoaDevGroup := strings.Trim(strings.Join(strings.Split(fmt.Sprint(*m.ArrCoaGroupID), " "), ","), "[]")
		query = query.Where("m_coadev.coadev_type_id IN (?)", conn.Model(&model.CoaTypeEntityModel{}).Select("id").Where(fmt.Sprintf("coadev_group_id IN (%s)", listCoaDevGroup)))
	}

	if m.Search != nil {
		query = query.Where("(m_coadev.name ILIKE ? OR m_coadev.code ILIKE ?)", "%"+*m.Search+"%", "%"+*m.Search+"%")
	}

	//sort
	if p.Sort == nil {
		sort := "desc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "created_at"
		p.SortBy = &sortBy
	}

	tmpSortBy := p.SortBy
	if p.SortBy != nil && *p.SortBy == "coadev_type" {
		sortBy := "\"m_coadev_type\".name"
		p.SortBy = &sortBy
	} else if p.SortBy != nil && *p.SortBy == "coadev_group" {
		sortBy := "\"m_coadev_group\".name"
		p.SortBy = &sortBy
	} else {
		sortBy := "m_coadev." + *p.SortBy
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)
	p.SortBy = tmpSortBy

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

	if err := query.Joins("JOIN m_coadev_type ON m_coadev_type.id = m_coadev.coadev_type_id").Joins("JOIN m_coadev_group ON m_coadev_type.coadev_group_id = m_coadev_group.id").Preload("CoaDevType.CoaDevGroup").Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].UserCreatedString = v.UserCreated.Name
		datas[i].UserModifiedString = &v.UserModified.Name
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

func (r *coadev) FindByID(ctx *abstraction.Context, id *int) (*model.CoaDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaDevEntityModel
	if err := conn.Where("id = ?", &id).Preload("CoaDevType.CoaDevGroup").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *coadev) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.CoaDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.CoaDevEntityModel
	if err := conn.Where("code LIKE ?", *code+"%").Order("code asc").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *coadev) Create(ctx *abstraction.Context, e *model.CoaDevEntityModel) (*model.CoaDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *coadev) Update(ctx *abstraction.Context, id *int, e *model.CoaDevEntityModel) (*model.CoaDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *coadev) Delete(ctx *abstraction.Context, id *int, e *model.CoaDevEntityModel) (*model.CoaDevEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *coadev) Export(ctx *abstraction.Context) (*[]model.CoaGroupEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.CoaGroupEntityModel
	if err := conn.Model(&model.CoaGroupEntityModel{}).Preload("CoaDevType.CoaDev").Order("code ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *coadev) CountDash(ctx *abstraction.Context) (int, error) {
	conn := r.CheckTrx(ctx)
	var TotalDash int
	if err := conn.Model(&model.CoaDevEntityModel{}).Select("LENGTH(name) - LENGTH(REPLACE(name, '-', '')) as TotalDash").Order("TotalDash DESC").Limit(1).Scan(&TotalDash).Error; err != nil {
		return 0, err
	}
	return TotalDash, nil
}
