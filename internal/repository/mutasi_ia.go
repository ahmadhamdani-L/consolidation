package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type MutasiIa interface {
	Find(ctx *abstraction.Context, m *model.MutasiIaFilterModel, p *abstraction.Pagination) (*[]model.MutasiIaEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.MutasiIaEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*int64, error)
	GetVersion(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*model.GetVersionModel, error)
	FindByCriteria(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*model.MutasiIaEntityModel, error)
}

type mutasiia struct {
	abstraction.Repository
}

func NewMutasiIa(db *gorm.DB) *mutasiia {
	return &mutasiia{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasiia) Find(ctx *abstraction.Context, m *model.MutasiIaFilterModel, p *abstraction.Pagination) (*[]model.MutasiIaEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiIaEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.MutasiIaEntityModel{})
	//filter
	tableName := model.MutasiIaEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, tableName)
	query = r.AllowedCompany(ctx, query, tableName)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, tableName)
	query = r.FilterMultiVersion(ctx, query, m.MutasiIaFilter)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
	query = query.Where("mutasi_ia.status != 0")
	query = query.Joins(`INNER JOIN "imported_worksheet" "Import" ON "Import"."company_id" = "mutasi_ia"."company_id" 
	AND "Import"."period" = "mutasi_ia"."period" 
	AND "Import"."versions" = "mutasi_ia"."versions" 
	AND "Import"."status" NOT IN (0, 1)`)
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
	if p.SortBy != nil && *p.SortBy == "company" {
		sortBy := "\"Company\".name"
		p.SortBy = &sortBy
	} else if p.SortBy != nil && *p.SortBy == "user" {
		sortBy := "\"UserCreated\".name"
		p.SortBy = &sortBy
	}
	if p.SortBy != nil && (tmpSortBy != nil && *tmpSortBy != "company" && *tmpSortBy != "user") {
		sortBy := fmt.Sprintf("mutasi_ia.%s", *p.SortBy)
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

	// if err := query.Preload("Company").Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	if err := query.Joins("Company").Joins("UserCreated", func(db *gorm.DB) *gorm.DB {
		db = db.Select("id, name")
		return db
	}).Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *mutasiia) FindByID(ctx *abstraction.Context, id *int) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiIaEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *mutasiia) Create(ctx *abstraction.Context, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *mutasiia) Update(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil

}

func (r *mutasiia) Destroy(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *mutasiia) Delete(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Update("status", 4).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *mutasiia) GetCount(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiIaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).WithContext(ctx.Request().Context()).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *mutasiia) GetVersion(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*model.GetVersionModel, error) {
	var data []model.MutasiIaEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiIaEntityModel{})
	query = r.Filter(ctx, query, *m)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")

	query = query.Select("versions").Where("status != 0").Group("versions").Order("versions ASC")

	if err := query.Find(&data).Error; err != nil {
		return &model.GetVersionModel{}, err
	}

	var result model.GetVersionModel
	tmp := []map[string]string{}
	for _, v := range data {
		tmp1 := map[string]string{
			"value": fmt.Sprintf("%d", v.Versions),
			"label": fmt.Sprintf("Version %d", v.Versions),
		}
		tmp = append(tmp, tmp1)
	}
	result.Version = tmp
	return &result, nil
}

func (r *mutasiia) FindByCriteria(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiIaEntityModel
	query := conn.Model(&model.MutasiIaEntityModel{})
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Company").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
