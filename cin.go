package cin

import (
	"net/http"
)

// CustomRouter представляет обертку над стандартным HTTP-роутером
type Router struct {
	mux        *http.ServeMux
	middleware []HandlerFunc
}

// NewCustomRouter создает новый экземпляр CustomRouter
func New() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

type HandlerFunc func(ctx *Context) Response

// Handle регистрирует обработчик для указанного пути
func (r *Router) registerHandler(method string, pattern string, handler HandlerFunc) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Создаем кастомный контекст
		ctx := NewContext(req.Context(), w, req)
		defer ctx.Request.Body.Close()

		// Вызываем пользовательский обработчик
		response := handler(ctx)
		if response != nil {
			response.out(ctx)
		}
	})
}

// GET регистрирует обработчик для GET-запросов
func (r *Router) GET(path string, handler HandlerFunc) {
	r.registerHandler(http.MethodGet, path, handler)
}

// POST регистрирует обработчик для POST-запросов
func (r *Router) POST(path string, handler HandlerFunc) {
	r.registerHandler(http.MethodPost, path, handler)
}

// PUT регистрирует обработчик для PUT-запросов
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.registerHandler(http.MethodPut, path, handler)
}

// DELETE регистрирует обработчик для DELETE-запросов
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.registerHandler(http.MethodDelete, path, handler)
}

func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.registerHandler(http.MethodPatch, path, handler)
}

// ServeHTTP делает CustomRouter совместимым с интерфейсом http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}
