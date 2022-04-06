package rest

import (
	"polls/core"
	"polls/core/model"
)

// InternalApisHandler handles the rest internal APIs implementation
type InternalApisHandler struct {
	app    *core.Application
	config *model.Config
}
