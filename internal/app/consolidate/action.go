package consolidate

import (
	"worker-consol/internal/abstraction"
)

func (h *handler) Action(act string, data abstraction.JsonData) {
	switch act {
	case "CONSOLIDATE":
		h.Consolidate(data)
	}

}
