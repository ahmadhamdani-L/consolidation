package export

import (
	"worker/internal/abstraction"
	"worker/internal/factory"
)

type handler struct {
	Service *service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		Service: NewService(f),
	}
}

func (h *handler) Export(data abstraction.JsonData) {
	cc := abstraction.Context{
		Auth: &abstraction.AuthContext{
			ID:        data.UserID,
			Name:      data.Name,
			CompanyID: data.CompanyID,
			Time:      data.Timestamp,
		},
	}
	h.Service.ExportAll(&cc, &data)

}

func (h *handler) ExportConsolidation(data abstraction.JsonData) {
	cc := abstraction.Context{
		Auth: &abstraction.AuthContext{
			ID:        data.UserID,
			Name:      data.Name,
			CompanyID: data.CompanyID,
			Time:      data.Timestamp,
		},
	}
	h.Service.ExportConsolidation(&cc, &data)

}
