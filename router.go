package router

import "net/http"

type iHandler interface {
	callHandler(http.ResponseWriter, *http.Request, Params)
	callMiddleware(http.ResponseWriter, *http.Request, Params, HandlerFunc)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, Params)

func (h HandlerFunc) callHandler(res http.ResponseWriter, req *http.Request, params Params) {
	h(res, req, params)
}

func (h HandlerFunc) callMiddleware(res http.ResponseWriter, req *http.Request, params Params, _ HandlerFunc) {
	h(res, req, params)
}

type iRouter interface {
	//method, path, handler
	handle(string, string, iHandler)
}

//////

type Router struct {
	*router
	hostsTree map[string]*router
}

func New() *Router {
	rout := &Router{

		router:    newRouter(),
		hostsTree: make(map[string]*router),
	}
	return rout
}

func (rout *Router) SubDomain(domain string) *router {
	domainRouter, ok := rout.hostsTree[domain]
	if !ok {
		domainRouter = newRouter()
		rout.hostsTree[domain] = domainRouter
	}
	return domainRouter
}

func (rout *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	domainRouter, ok := rout.hostsTree[req.URL.Host]
	if ok {

		if domainRouter.panicHandler != nil {
			defer func() {
				if err := recover(); err != nil {
					domainRouter.panicHandler(res, req, err)
				}
			}()
		} else if rout.panicHandler != nil {
			defer func() {
				if err := recover(); err != nil {
					rout.panicHandler(res, req, err)
				}
			}()
		} else {
			defer func() {
				if err := recover(); err != nil {
					panicHandler(res, req, err)
				}
			}()
		}

		methodRouter, ok := domainRouter.methodsTree[req.Method]
		if ok {
			if handler, params := methodRouter.FindRoute(req.URL.Path); handler != nil {
				handler.callHandler(res, req, params)
				return
			}
		}
	} else {

		if rout.panicHandler != nil {
			defer func() {
				if err := recover(); err != nil {
					rout.panicHandler(res, req, err)
				}
			}()
		} else {
			defer func() {
				if err := recover(); err != nil {
					panicHandler(res, req, err)
				}
			}()
		}
	}
	methodRouter, ok := rout.router.methodsTree[req.Method]
	if ok {
		if handler, params := methodRouter.FindRoute(req.URL.Path); handler != nil {
			handler.callHandler(res, req, params)
			return
		} else {
			domainRouter, ok := rout.hostsTree[req.URL.Host]
			if ok {
				if domainRouter.notFoundHandler != nil {
					domainRouter.notFoundHandler(res, req)
					return
				}
			}
			if rout.notFoundHandler != nil {
				rout.notFoundHandler(res, req)
				return
			} else {
				notFoundHandler(res, req)
				return
			}
		}
	} else {
		domainRouter, ok := rout.hostsTree[req.URL.Host]
		if ok {
			if domainRouter.methodNotAllowedHandler != nil {
				domainRouter.methodNotAllowedHandler(res, req)
				return
			}
		}
		if rout.methodNotAllowedHandler != nil {
			rout.methodNotAllowedHandler(res, req)
			return
		} else {
			methodNotAllowedHandler(res, req)
			return
		}
	}

}

///////

type router struct {
	*routerTree

	panicHandler            PanicHandler
	notFoundHandler         NotFoundHandler
	methodNotAllowedHandler MethodNotAllowedHandler

	methodsTree map[string]*tree
}

func newRouter() *router {
	rout := &router{

		methodsTree: make(map[string]*tree),
	}
	rout.routerTree = newRouterTree("", rout)

	return rout
}

func (rout *router) handle(method string, path string, handler iHandler) {
	currentTree, ok := rout.methodsTree[method]
	if !ok {
		currentTree = newTree()
		rout.methodsTree[method] = currentTree
	}

	currentTree.AddRoute(path, handler)
}

func (rout *router) PanicHandler(handler PanicHandler) {
	rout.panicHandler = handler
}

func (rout *router) NotFoundHandler(handler NotFoundHandler) {
	rout.notFoundHandler = handler
}

func (rout *router) MethodNotAllowedHandler(handler MethodNotAllowedHandler) {
	rout.methodNotAllowedHandler = handler
}

////////
type routerTree struct {
	path   string
	parent iRouter

	middlewares []iHandler
}

func newRouterTree(path string, parent iRouter) *routerTree {
	rout := &routerTree{
		path:   path,
		parent: parent,

		middlewares: make([]iHandler, 0),
	}

	return rout
}

func (rout *routerTree) SubRouter(path string) *routerTree {
	return newRouterTree(path, rout)
}

func (rout *routerTree) handle(method string, path string, handler iHandler) {
	rout.parent.handle(method, rout.path+path, buildMiddleware(append(rout.middlewares, handler)...))
}

func (rout *routerTree) Use(middleware MiddlewareFunc) {
	rout.middlewares = append(rout.middlewares, middleware)
}

func (rout *routerTree) Get(path string, handler HandlerFunc) {
	rout.handle("GET", path, handler)
}
func (rout *routerTree) Put(path string, handler HandlerFunc) {
	rout.handle("PUT", path, handler)
}
func (rout *routerTree) Post(path string, handler HandlerFunc) {
	rout.handle("POST", path, handler)
}
func (rout *routerTree) Head(path string, handler HandlerFunc) {
	rout.handle("HEAD", path, handler)
}
func (rout *routerTree) Patch(path string, handler HandlerFunc) {
	rout.handle("PATCH", path, handler)
}
func (rout *routerTree) Delete(path string, handler HandlerFunc) {
	rout.handle("DELETE", path, handler)
}
func (rout *routerTree) Options(path string, handler HandlerFunc) {
	rout.handle("OPTIONS", path, handler)
}

/////////
