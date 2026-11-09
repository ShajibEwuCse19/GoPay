package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalMiddlewares []Middleware
}

func NewManager() *Manager {
	return &Manager{
		globalMiddlewares: make([]Middleware, 0),
	}
}

func (m *Manager) AppendGlobalMiddlewares(middlewares ...Middleware) {
	m.globalMiddlewares = append(m.globalMiddlewares, middlewares...)
}

func (m *Manager) WrapMuxWithGlobalMiddleware(next http.Handler, middlewares ...Middleware) http.Handler {
	handler := next

	for i := len(m.globalMiddlewares) - 1; i >= 0; i-- {
		handler = m.globalMiddlewares[i](handler)
	}

	return handler
}

func (m *Manager) AddMiddleware(next http.Handler, middlewares ...Middleware) http.Handler {
	handler := next

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}