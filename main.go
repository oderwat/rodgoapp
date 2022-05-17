package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

//go:generate sh -c "GOARCH=wasm GOOS=js go build -o web/app.wasm"
//go:generate go build -o tester

const siteURL = "localhost:8000"

func main() {
	// exists in front and backend because of possible SEO
	app.Route("/", &hello{})
	frontend()
	go backend()

	testing()
}
