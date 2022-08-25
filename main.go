package main

import (
	// "html/template"
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type Template interface {
	Execute(w http.ResponseWriter, data interface{})
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "my restAPI")
}

func main() {
	rt := chi.NewRouter()
	rt.Get("/", StaticHandler)
	rt.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "page not found")
	})
	fmt.Println("starting server on 8080")
	http.ListenAndServe(":8080", rt)
}
