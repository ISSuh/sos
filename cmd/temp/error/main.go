package main

import (
	"errors"
	"fmt"
)

var (
	notFound = NewNotFoundError(nil)
)

type Error struct {
	code int
	err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Err: %v", e.code, e.err)
}

// Unwrap 메서드는 errors.Unwrap 함수와 호환되도록 합니다.
func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.code == t.code
}

func (e *Error) ToError() error {
	return e
}

type NotFound struct {
	Error
}

func NewNotFoundError(err error) error {
	errsa := NotFound{}
	return errsa.ToError()
}

func Testfunc() error {
	return NewNotFoundError(fmt.Errorf("Resource not found"))
}

func main() {
	err := Testfunc()
	fmt.Println(err)

	if errors.Is(err, notFound) {
		fmt.Println("Error wraps the system error")
	} else {
		fmt.Println("Error does not wrap the system error")
	}
}
