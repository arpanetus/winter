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
	c.Status(err.Status).JSON(err)
}

func (c *Context) GetParam(key string) string {
	param, ok := c.GetParams()[key]
	if !ok {
		return ""
	}

	return param
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
