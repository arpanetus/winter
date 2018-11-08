package core

import (
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
)

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
		c.Status(err.Status).JSON(Error{
			&Response{
				500,
				"Unknown error",
			},
		})
		return
	}

	c.Status(err.Status).JSON(err)
}

func (c *Context) SendNewResponse() {
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
	m.Handler.ServeHTTP(m.Response, m.Request)
}
