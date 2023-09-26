package imports

import (
	"worker/internal/abstraction"
)

func (h *handler) ActionImport(act string, dataImport abstraction.JsonDataImport) {
	switch act {
	case "IMPORT":
		h.Import(dataImport)
		break
	}

}
