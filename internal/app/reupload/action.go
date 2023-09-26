package reupload

import (
	"worker/internal/abstraction"
)

func (h *handler) ActionImport(act string, dataReUpload abstraction.JsonDataReUpload) {
	switch act {
	case "REUPLOAD":
		h.ReUpload(dataReUpload)
		break
	}

}
