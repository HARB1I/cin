package cin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response определяет интерфейс для всех ответов
type Response interface {
	out(*Context)
}

// BaseResponse содержит общую логику для всех ответов
type BaseResponse[T any] struct {
	data       T
	statusCode int
}

func (resp BaseResponse[T]) out(ctx *Context) {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.WriteHeader(resp.statusCode)

	// Сериализуем ответ в JSON и отправляем его
	if err := json.NewEncoder(ctx.ResponseWriter).Encode(resp.data); err != nil {
		http.Error(ctx.ResponseWriter, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

type H map[string]string

func Resp[T any](code int, obj T) Response {
	return BaseResponse[T]{data: obj, statusCode: code}
}
