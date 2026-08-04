// Package pool реализует пул объектов с автоматическим сбросом
package pool

import "sync"

type Resetter interface {
	Reset()
}

type Pool[T Resetter] struct {
	pool    sync.Pool
	newFunc func() T
}

func New[T Resetter](newFunc func() T) *Pool[T] {
	p := &Pool[T]{
		newFunc: newFunc,
	}

	p.pool.New = func() interface{} {
		return p.newFunc()
	}

	return p
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.pool.Put(obj)
}
