package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const publicDir string = "/Users/ryan.farrell/go/src/github.com/rmfarrell/semistatic/"

type route struct {
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

type staticRouter struct{}

type routeHandler interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "Set the server port")
	flag.Parse()
}

// Wrap a route.handler in an adapter for mux handler
// func (rt *route) muxHandler() func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Write(rt.handler(r))
// 	}
// }

// Add a route to a mux.Route
func (rt *route) AddTo(m *mux.Router) *route {
	m.HandleFunc(rt.path, rt.handler)
	return rt
}

// Constructor for new route
func newRoute(p string, f func(http.ResponseWriter, *http.Request)) *route {
	// TODO: validate
	return &route{p, f}
}

func aHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

func aStaticHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm static dog!\n"))
}

// Compile a static route
func (rt *route) compile(req *http.Request) error {
	f, err := os.Create(fmt.Sprintf("%s%s", publicDir, rt.path))
	if err != nil {
		return err
	}
	defer f.Close()

	// fmt.Println()
	// w := bufio.NewWriter(f)
	rt.handler(w, req)
	// _, err = w.WriteString("string\n")
	// if err != nil {
	// 	return err
	// }
	w.Flush()
	return nil
}

func main() {
	static := mux.NewRouter()
	server := mux.NewRouter()
	newRoute("/", aHandler).AddTo(server)
	s := newRoute("index.html", aStaticHandler).AddTo(static)
	err := s.compile(static, &http.Request{})
	if err != nil {
		log.Println(err)
	}

	// err := compile(static, "public")
	// if err != nil {
	// 	return
	// }

	fmt.Println(fmt.Sprintf("Listening on port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
