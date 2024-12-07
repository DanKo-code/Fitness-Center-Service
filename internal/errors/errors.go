package errors

import "errors"

var (
	VoidServiceData          = errors.New("void service data")
	ServiceAlreadyExists     = errors.New("service already exists")
	ServiceNotFound          = errors.New("service not found")
	CoachNotFound            = errors.New("coach not found")
	InternalCoachServerError = errors.New("internal coach server error")
)
