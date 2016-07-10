package listener

import "fmt"

type Registry struct {
	listeners map[interface{}]func()
}

func NewRegistry() *Registry {
	return &Registry{
		listeners: make(map[interface{}]func()),
	}
}

func (r *Registry) Add(key interface{}, listener func()) {
	if key == nil {
		key = new(int)
	}
	if _, ok := r.listeners[key]; ok {
		panic(fmt.Sprintf("listener with key already exists: %v", key))
	}
	r.listeners[key] = listener
}

func (r *Registry) Remove(key interface{}) {
	delete(r.listeners, key)
}

func (r *Registry) Fire() {
	for _, l := range r.listeners {
		l()
	}
}
