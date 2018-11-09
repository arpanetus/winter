package core

import (
	"net/http"
	"strconv"
)

var HTTPErrors = getHTTPErrorMap(
	// 4XX - Client Errors
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusPaymentRequired,
	http.StatusForbidden,
	http.StatusNotFound,
	http.StatusMethodNotAllowed,
	http.StatusNotAcceptable,
	http.StatusProxyAuthRequired,
	http.StatusRequestTimeout,
	http.StatusConflict,
	http.StatusGone,
	http.StatusLengthRequired,
	http.StatusPreconditionFailed,
	http.StatusRequestEntityTooLarge,
	http.StatusRequestURITooLong,
	http.StatusUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed,
	http.StatusTeapot,
	http.StatusMisdirectedRequest,
	http.StatusUnprocessableEntity,
	http.StatusLocked,
	http.StatusFailedDependency,
	http.StatusUpgradeRequired,
	http.StatusPreconditionRequired,
	http.StatusTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons,
	// 5XX - Server Errors
	http.StatusInternalServerError,
	http.StatusInternalServerError,
	http.StatusNotImplemented,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
	http.StatusHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates,
	http.StatusInsufficientStorage,
	http.StatusLoopDetected,
	http.StatusNotExtended,
	http.StatusNetworkAuthenticationRequired,
)

func NewErrorMap(errs ...BindError) *ErrorMap {
	errMap := HTTPErrors
	for _, n := range errs {
		(*errMap)[n.Code] = NewError(n.Status, n.Message)
	}
	return errMap
}

func NewError(status int, message interface{}) *Error {
	return &Error{
		Response: &Response{
			Status: status,
			Message: message,
		},
	}
}

func Err(code int, status int, message string) BindError {
	newErr := NewError(status, message)
	return BindError{
		Code: code,
		Error: newErr,
	}
}

func getHTTPErrorMap(httpErrorCodes ...int) *ErrorMap {
	httpErrors := &ErrorMap{}
	for _, n := range httpErrorCodes {
		(*httpErrors)[n] = NewError(n, http.StatusText(n))
	}
	return httpErrors
}

func (e *Error) Send(ctx *Context) {
	ctx.Status(e.Status).JSON(e)
}

func (e *Error) SetMessage(mess interface{}) {
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

func (e *ErrorMap) Set(code int, err *Error) {
	(*e)[code] = err
}
