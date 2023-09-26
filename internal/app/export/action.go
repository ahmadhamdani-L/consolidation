package export

import (
	"worker/internal/abstraction"
)

func (h *handler) Action(act string, data abstraction.JsonData) {
	switch act {
	case "EXPORT":
		h.Export(data)
	case "EXPORT_CONSOLIDATION":
		h.ExportConsolidation(data)
	}

}
