package route

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

// Route is struct which contains handlers for both static controllers and HTTP.ServeMux handlers
type Route struct {
	Path              string
	StaticHandlerFunc func(*http.Request) []byte
	HandlerFunc       func(http.ResponseWriter, *http.Request)
}

// NewRouteInput is the input struct for NewRoute constructor
// callback passes ResponseWriter which can be used to set headers in  server context
type NewRouteInput struct {
	Path              string
	StaticHandlerFunc func(*http.Request) []byte
	Callback          func(http.ResponseWriter)
}

// NewRoute is constructor for new route
func NewRoute(in *NewRouteInput) *Route {
	// TODO: validate
	route := Route{in.Path, in.StaticHandlerFunc, nil}
	route.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(route.StaticHandlerFunc(r))
		in.Callback(w)
	}
	return &route
}

// AddTo appends the route to an http.ServeMux
func (rt *Route) AddTo(m *http.ServeMux) *Route {
	m.HandleFunc(rt.Path, rt.HandlerFunc)
	return rt
}

// Compile a static route
// TODO: use server static path, if using gorilla mux, if possible.
func (rt *Route) Compile(req *http.Request, p string) error {
	f, err := os.Create(fmt.Sprintf("%s%s", p, rt.Path))
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.Write(rt.StaticHandlerFunc(req))
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
