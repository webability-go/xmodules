package base

import (
	"github.com/webability-go/xamboo/applications"
)

type Hook func(ds applications.Datasource, data ...interface{})

type Hooks struct {
	Hooks map[string]Hook
}

func NewHooks() *Hooks {
	return &Hooks{Hooks: map[string]Hook{}}
}

func (h *Hooks) Register(id string, f Hook) {
	h.Hooks[id] = f
}

func (h *Hooks) Call(ds applications.Datasource, data ...interface{}) {

	for _, f := range h.Hooks {
		f(ds, data...)
	}

}
