package validation

import (
	"worker-validation/internal/abstraction"
)

func (h *handler) Action(act string, data abstraction.JsonData) {
	switch act {
	case "VALIDATE":
		h.Validate(data)
	}

}
