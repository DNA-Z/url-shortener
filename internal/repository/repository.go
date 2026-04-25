package repository

type URL[T any] struct {
	value T
}

type URLRepository[T any] interface {
	WriteURL(url URL[T]) func() string
	ReadURL(url URL[T]) func() string
}
