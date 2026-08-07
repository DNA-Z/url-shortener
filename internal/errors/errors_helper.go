// Package errors содержит переопределенный метод для возврата ошибки
package errors

import "fmt"

type ConflictError struct {
	Status int
	URL    string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("URL already exists (%d)", e.Status)
}
