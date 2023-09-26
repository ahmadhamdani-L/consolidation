package imports

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

func (h *handler) Import(dataImport abstraction.JsonDataImport) {
	cc := abstraction.Context{
		Auth: &abstraction.AuthContext{
			ID:        dataImport.UserID,
			CompanyID: dataImport.CompanyID,
		},
	}
	h.Service.ImportAll(&cc, &dataImport)
	return
}