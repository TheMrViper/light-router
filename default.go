package router

import (
	"fmt"
	"net/http"
)

////
type PanicHandler func(http.ResponseWriter, *http.Request, interface{})

func panicHandler(res http.ResponseWriter, req *http.Request, i interface{}) {

	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(http.StatusInternalServerError)

	res.Write([]byte("<center>"))
	err, ok := i.(error)
	if ok {
		res.Write([]byte("Panic: <b>" + err.Error() + "</b><br>"))
	} else {
		res.Write([]byte("Panic: <b>" + fmt.Sprint(err) + "</b><br>"))
	}
	res.Write([]byte("<h1>Server error</h1><br>"))
	res.Write([]byte("</center>"))
}

////
type NotFoundHandler func(http.ResponseWriter, *http.Request)

func notFoundHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(http.StatusNotFound)

	res.Write([]byte("<center>"))
	res.Write([]byte("Page: <b>" + req.URL.String() + "</b><br>"))
	res.Write([]byte("<h1>Not found</h1><br>"))
	res.Write([]byte("</center>"))
}

////
type MethodNotAllowedHandler func(http.ResponseWriter, *http.Request)

func methodNotAllowedHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(http.StatusMethodNotAllowed)

	res.Write([]byte("<center>"))
	res.Write([]byte("Method: <b>" + req.Method + "</b><br>"))
	res.Write([]byte("<h1>Not allowed</h1><br>"))
	res.Write([]byte("</center>"))
}
