package consolidate

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/factory"
)

type handler struct {
	Service *service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		Service: NewService(f),
	}
}

func (h *handler) Consolidate(data abstraction.JsonData) {
	cc := abstraction.Context{
		Auth: &abstraction.AuthContext{
			ID:        data.UserID,
			Name:      data.Name,
			CompanyID: data.CompanyID,
			Time:      data.Timestamp,
		},
	}
	h.Service.Process(&cc, &data)

}
