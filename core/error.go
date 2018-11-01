package core

import "strconv"

func NewErrorMap(errs ...BindError) ErrorMap {
	errMap := ErrorMap{}
	for _, n := range errs {
		errMap[n.Code] = NewError(n.Status, n.Message)
	}
	return errMap
}

func NewError(status int, message string) *Error {
	return &Error{
		Status: status,
		Message: message,
	}
}

func Err(code int, status int, message string) BindError {
	newErr := NewError(status, message)
	return BindError{
		Code: code,
		Error: newErr,
	}
}

func (e *Error) Send(ctx *Context) {
	ctx.Status(e.Status).JSON(e)
}

func (e *Error) SetMessage(mess string) {
	e.Message = mess
}

func (e *Error) SetStatus(status int) {
	e.Status = status
}

func (e *ErrorMap) Get(code int) *Error {
	err, ok := (*e)[code]
	if !ok {
		routerLogger.Warn("Missin error with code " + strconv.Itoa(code))
		return &Error{}
	}
	return err
}

func (e *ErrorMap) Set(code int, err Error) {
	(*e)[code] = &err
}
