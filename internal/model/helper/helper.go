package helper

import (
	"errors"
	"mcash-finance-console-core/internal/model"
	"net/url"
	"strconv"
)

func MultiCompanyFilter(queryParam url.Values) (model.CompanyCustomFilter, error) {
	payload := model.CompanyCustomFilter{}
	companyFilter := queryParam["company_id[]"]
	if len(companyFilter) > 0 {
		var company []int
		for _, v := range companyFilter {
			companyInt, err := strconv.Atoi(v)
			if err != nil {
				return model.CompanyCustomFilter{}, errors.New("Not Integer Value!")
			}
			company = append(company, companyInt)
		}
		payload.ArrCompanyID = &company
	}

	companyStringFilter := queryParam["company[]"]
	companyOperatorFilter := queryParam["company_operator[]"]
	if len(companyStringFilter) > 0 && len(companyOperatorFilter) > 0 && len(companyStringFilter) == len(companyOperatorFilter) {
		var companyString []string
		var companyOperator []string
		for i, v := range companyStringFilter {
			companyString = append(companyString, v)
			companyOperator = append(companyOperator, string(companyOperatorFilter[i]))
		}
		payload.ArrCompanyString = &companyString
		payload.ArrCompanyOperator = &companyOperator
	} else if len(companyStringFilter) != len(companyOperatorFilter) {
		return model.CompanyCustomFilter{}, errors.New("Missing Parameter!")
	}
	return payload, nil
}
