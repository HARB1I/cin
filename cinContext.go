package cin

import (
	"context"
	"encoding/json"
	"net/http"
)

// CustomContext представляет кастомный контекст
type Context struct {
	context.Context                     // Встраиваем стандартный context.Context
	ResponseWriter  http.ResponseWriter // Публичное поле для ResponseWriter
	Request         *http.Request       // Публичное поле для Request
}

// NewCustomContext создает новый экземпляр CustomContext
func NewContext(parent context.Context, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Context:        parent,
		ResponseWriter: w,
		Request:        r,
	}
}

func (c Context) BindJSON(obj any) error {
	return json.NewDecoder(c.Request.Body).Decode(obj)
}

func (c Context) PathValue(name string) string {
	return c.Request.PathValue(name)
}
