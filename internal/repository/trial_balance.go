package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type TrialBalance interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.TrialBalanceEntityModel, error)
	Create(ctx *abstraction.Context, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error)
	Get(ctx *abstraction.Context, id int) (*model.TrialBalanceEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*int64, error)
	GetVersion(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.GetVersionModel, error)
	FindByCriteria(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error)
	Finds(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
}

type trialbalance struct {
	abstraction.Repository
}

func NewTrialBalance(db *gorm.DB) *trialbalance {
	return &trialbalance{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *trialbalance) Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
    conn := r.CheckTrx(ctx)

    var datas []model.TrialBalanceEntityModel
    var info abstraction.PaginationInfo

    query := conn.Model(&model.TrialBalanceEntityModel{})
    //filter
    tableName := model.TrialBalanceEntityModel{}.TableName()
    query = r.FilterTable(ctx, query, *m, tableName)
    query = r.AllowedCompany(ctx, query, tableName)

    //filter custom
    tmp1 := m.CompanyCustomFilter
    queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
    query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, tableName)
    query = r.FilterMultiVersion(ctx, query, m.TrialBalanceFilter)
    queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
    query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
    query = query.Where("trial_balance.status != 0")

    if m.Search != nil {
        query = query.Where("trial_balance.created_by IN (?)", conn.Model(&model.UserEntityModel{}).Select("id").Where("name ILIKE ?", "%"+*m.Search+"%"))
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
    if p.SortBy != nil && *p.SortBy == "company" {
        sortBy := "\"Company\".name"
        p.SortBy = &sortBy
    } else if p.SortBy != nil && *p.SortBy == "user" {
        sortBy := "\"UserCreated\".name"
        p.SortBy = &sortBy
    }
    if p.SortBy != nil && (tmpSortBy != nil && *tmpSortBy != "company" && *tmpSortBy != "user") {
        sortBy := fmt.Sprintf("trial_balance.%s", *p.SortBy)
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
	if m.ArrVersions == nil || len(*m.ArrVersions) == 0 {
		query = query.Joins(`INNER JOIN "imported_worksheet" "Import" ON "Import"."company_id" = "trial_balance"."company_id" 
		AND "Import"."period" = "trial_balance"."period" 
		AND "Import"."versions" = "trial_balance"."versions" 
		AND "Import"."status" NOT IN (0, 1)`)
	}
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
    if int(totalData) > *p.PageSize {
        info.MoreRecords = true
        // info.Count -= 1
        // datas = datas[:len(datas)-1]
    }

    return &datas, &info, nil
}

func (r *trialbalance) Finds(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
    conn := r.CheckTrx(ctx)

    var datas []model.TrialBalanceEntityModel
    var info abstraction.PaginationInfo

    query := conn.Model(&model.TrialBalanceEntityModel{})
    //filter
    tableName := model.TrialBalanceEntityModel{}.TableName()
    query = r.FilterTable(ctx, query, *m, tableName)
    query = r.AllowedCompany(ctx, query, tableName)

    //filter custom
    tmp1 := m.CompanyCustomFilter
    queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
    query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, tableName)
    query = r.FilterMultiVersion(ctx, query, m.TrialBalanceFilter)
    queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
    query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
    query = query.Where("status != 0")

    if m.Search != nil {
        query = query.Where("trial_balance.created_by IN (?)", conn.Model(&model.UserEntityModel{}).Select("id").Where("name ILIKE ?", "%"+*m.Search+"%"))
    }

	// query = query.Joins(`INNER JOIN "imported_worksheet" "Import" ON "Import"."company_id" = "trial_balance"."company_id" 
	// AND "Import"."period" = "trial_balance"."period" 
	// AND "Import"."versions" = "trial_balance"."versions" 
	// AND "Import"."status" NOT IN (0, 1)`)
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
        sortBy := fmt.Sprintf("trial_balance.%s", *p.SortBy)
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
    if int(totalData) > *p.PageSize {
        info.MoreRecords = true
        // info.Count -= 1
        // datas = datas[:len(datas)-1]
    }

    return &datas, &info, nil
}

func (r *trialbalance) FindByID(ctx *abstraction.Context, id *int) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *trialbalance) Create(ctx *abstraction.Context, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
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

func (r *trialbalance) Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)
	err := conn.Model(e).Where("id = ?", id).Updates(e).WithContext(ctx.Request().Context()).Error
	if err != nil {
		return nil, err
	}

	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil

}

func (r *trialbalance) Destroy(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *trialbalance) Delete(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Update("status", 4).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *trialbalance) Get(ctx *abstraction.Context, id int) (*model.TrialBalanceEntityModel, error) {
	var datas model.TrialBalanceEntityModel

	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter

	if err := query.Where("id = ?", id).Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}

	var formatterBridges []model.FormatterBridgesEntityModel

	query = conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ?", datas.ID).Where("source = ?", "TRIAL-BALANCE")
	if err := query.Find(&formatterBridges).Error; err != nil {
		return &datas, err
	}
	datas.FormatterBridges = formatterBridges

	return &datas, nil
}

func (r *trialbalance) GetCount(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).WithContext(ctx.Request().Context()).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *trialbalance) GetVersion(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.GetVersionModel, error) {
	var data []model.TrialBalanceEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	query = r.Filter(ctx, query, *m)
	query = query.Where("status != 0")

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")
	query = r.FilterMultiStatus(ctx, query, *m)

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

func (r *trialbalance) FindByCriteria(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error) {
	var data model.TrialBalanceEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&data).Preload("Company").WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
