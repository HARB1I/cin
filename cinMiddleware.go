package cin

import "net/http"

// Group представляет группу маршрутов
type Group struct {
	prefix     string
	middleware []HandlerFunc
	router     *Router
}

// Group создает новую группу маршрутов
func (r *Router) Group(prefix string, middleware ...HandlerFunc) *Group {
	return &Group{
		prefix:     prefix,
		middleware: append(r.middleware, middleware...),
		router:     r,
	}
}

// Group создает новую группу маршрутов
// Handle регистрирует обработчик для указанного пути в группе
func (g *Group) Handle(method, path string, handler HandlerFunc) {
	fullPath := g.prefix + path

	// Создаем цепочку middleware
	wrappedHandler := chainMiddleware(handler, g.middleware...)

	g.router.mux.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		// Создаем контекст
		ctx := NewContext(r.Context(), w, r)
		defer ctx.Request.Body.Close()

		// Вызываем обработчик
		response := wrappedHandler(ctx)
		if response != nil {
			response.out(ctx)
		}
	})
}

// chainMiddleware объединяет middleware в одну цепочку
func chainMiddleware(final HandlerFunc, middleware ...HandlerFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		mw := middleware[i]
		next := final
		final = func(ctx *Context) Response {
			// Вызываем middleware
			resp := mw(ctx)
			if resp != nil {
				// Если middleware вернул ответ, завершаем цепочку
				return resp
			}
			// Иначе передаем управление следующему middleware или финальному обработчику
			return next(ctx)
		}
	}
	return final
}

func (g *Group) GET(path string, handler HandlerFunc) {
	g.Handle(http.MethodGet, path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) {
	g.Handle(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler HandlerFunc) {
	g.Handle(http.MethodPut, path, handler)
}

func (g *Group) DELETE(path string, handler HandlerFunc) {
	g.Handle(http.MethodDelete, path, handler)
}

func (g *Group) PATCH(path string, handler HandlerFunc) {
	g.Handle(http.MethodPatch, path, handler)
}
