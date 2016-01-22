package router

import "net/http"

type MiddlewareFunc func(http.ResponseWriter, *http.Request, Params, HandlerFunc)

func (m MiddlewareFunc) callHandler(_ http.ResponseWriter, _ *http.Request, _ Params) {}

func (m MiddlewareFunc) callMiddleware(res http.ResponseWriter, req *http.Request, params Params, next HandlerFunc) {
	m(res, req, params, next)
}

type middleware struct {
	handler iHandler
	next    *middleware
}

func (m *middleware) callHandler(res http.ResponseWriter, req *http.Request, params Params) {
	m.handler.callMiddleware(res, req, params, m.next.callHandler)
}

func (m *middleware) callMiddleware(_ http.ResponseWriter, _ *http.Request, _ Params, _ HandlerFunc) {}

func buildMiddleware(handlers ...iHandler) *middleware {

	if len(handlers) == 0 {
		return nil
	}
	return &middleware{handlers[0], buildMiddleware(handlers[1:]...)}
}
