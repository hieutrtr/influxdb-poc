package measurementdb

import "fmt"

type dbError struct {
	err  string
	code int
}

func NewDBError(err string, code int) error {
	return &dbError{
		err:  err,
		code: code,
	}
}

func (e *dbError) Error() string {
	return fmt.Sprintf("dberror: %s", e.err)
}

func (e *dbError) Code() int {
	return e.code
}

type ErrMeasurementNotFound struct {
	err  string
	code int
}

func NewErrMeasurementNotFound(err string, code int) error {
	return &ErrMeasurementNotFound{
		err:  err,
		code: code,
	}
}

func (e *ErrMeasurementNotFound) Error() string {
	return fmt.Sprintf("ErrMeasurementNotFound: %s", e.err)
}

func (e *ErrMeasurementNotFound) Code() int {
	return e.code
}
