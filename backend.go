//go:build !wasm

package main

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"log"
	"net/http"
)

func frontend() {}
func backend() {
	mux := http.NewServeMux()

	mux.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	mux.HandleFunc("/post", func(res http.ResponseWriter, req *http.Request) {
		log.Println("Backend Route /post called. We send data!")
		_, _ = fmt.Fprint(res, "Backend Response")
	})

	mux.HandleFunc("/empty", func(res http.ResponseWriter, req *http.Request) {
		log.Println("Backend Route /empty called. We send nothing!")
	})

	fmt.Println("Open browser at " + siteURL)
	if err := http.ListenAndServe(siteURL, mux); err != nil {
		log.Fatal(err)
	}
}
