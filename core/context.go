package core

import (
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
)

func NullResponse() Response {
	return Response{}
}

func NewResponse(status int, message interface{}) Response {
	return Response{
		Status: status,
		Message: message,
	}
}

func NewSuccessResponse(message interface{}) Response {
	return Response{
		Status: 200,
		Message: message,
	}
}

func NewErrorResponse(err *Error) Response {
	return Response{
		Status: err.Status,
		Message: err.Message,
	}
}

func (c *Context) Send(msg []byte) {
	c.Response.Write(msg)
}

func (c *Context) JSON(mess interface{}) {
	json.NewEncoder(c.Response).Encode(mess)
}

func (c *Context) Status(code int) *Context {
	c.Response.WriteHeader(code)
	return c
}

func (c *Context) Header(key, val string) {
	c.Response.Header().Set(key, val)
}

func (c *Context) SendError(err Error) {
	if err == (Error{}) {
		c.SendResponse(err.Status, Error{
			&Response{
				500,
				"Unknown error",
			},
		})
		return
	}

	c.SendResponse(err.Status, err)
}

func (c *Context) SendSuccess(message interface{}) {
	c.SendResponse(200, message)
}

func (c *Context) SendResponse(status int, message interface{}) {
	c.Status(status).JSON(Response{
		status,
		message,
	})
}

func (c *Context) GetParam(key string) (string, bool) {
	param, ok := c.GetParams()[key]
	return param, ok
}

func (c *Context) GetParams() map[string]string {
	params := mux.Vars(c.Request)
	return params
}

func (c *Context) GetBody(body interface{}) error {
	b, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, body)
	if err != nil {
		return err
	}

	return nil
}

func (m *MiddlewareContext) Next() {
	m.handler.ServeHTTP(m.Response, m.Request)
}
