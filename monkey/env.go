package monkey

type Env struct {
	store map[string]Value
	outer *Env
}

func (e *Env) Get(name string) (Value, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		val, ok = e.outer.Get(name)
	}
	return val, ok
}

func (e *Env) Set(name string, val Value) {
	e.store[name] = val
}

func NewEnv(outer *Env) *Env {
	return &Env{
		store: make(map[string]Value),
		outer: outer,
	}
}
