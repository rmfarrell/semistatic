package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type route struct {
	path       string
	handler    func(*http.Request) []byte
	muxHandler func(http.ResponseWriter, *http.Request)
}

type newRouteInput struct {
	path     string
	handler  func(*http.Request) []byte
	callback func(http.ResponseWriter)
}

var (
	port int
	dir  string
)

func init() {
	flag.IntVar(&port, "port", 8000, "Set the server port")
	flag.Parse()
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir = currentDir
}

// Add a route to a router
func (rt *route) AddTo(m *mux.Router) *route {
	m.HandleFunc(rt.path, rt.muxHandler)
	return rt
}

// Constructor for new lib route
func newRoute(in *newRouteInput) *route {
	// TODO: validate
	route := route{
		path:    in.path,
		handler: in.handler,
	}
	route.muxHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write(route.handler(r))
		in.callback(w)
	}
	return &route
}

// Generic success handler
func success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// a handler
func aHandler(r *http.Request) []byte {
	return []byte("Hello World!\n")
}

// Compile a static route
func (rt *route) compile(req *http.Request) error {
	f, err := os.Create(fmt.Sprintf("%s%s", dir, rt.path))
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.Write(rt.handler(req))
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func main() {
	static := mux.NewRouter()
	server := mux.NewRouter()
	s := newRoute(&newRouteInput{"/index.html", aHandler, success})
	s.AddTo(static)
	s.AddTo(server)
	err := s.compile(&http.Request{})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(fmt.Sprintf("Listening on port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
