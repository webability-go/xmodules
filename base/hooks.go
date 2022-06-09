package base

type Hook func(ds *Datasource, data ...interface{})

type Hooks struct {
	Hooks map[string]Hook
}

func NewHooks() *Hooks {
	return &Hooks{Hooks: map[string]Hook{}}
}

func (h *Hooks) Register(id string, f Hook) {
	h.Hooks[id] = f
}

func (h *Hooks) Call(ds *Datasource, data ...interface{}) {

	for _, f := range h.Hooks {
		f(ds, data...)
	}

}
