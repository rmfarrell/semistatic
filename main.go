package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	rt "github.com/rmfarrell/semistatic/route"
)

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

// Generic success handler
func success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// a handler
func aHandler(r *http.Request) []byte {
	return []byte("Hello World!\n")
}

func main() {
	static := mux.NewRouter()
	server := mux.NewRouter()
	s := rt.NewRoute(&rt.NewRouteInput{"/index.html", aHandler, success})
	s.AddTo(static)
	s.AddTo(server)
	err := s.Compile(&http.Request{}, dir)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(fmt.Sprintf("Listening on port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
