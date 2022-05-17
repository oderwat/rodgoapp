//go:build wasm

package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func testing() {}
func backend() {}
func frontend() {
	app.RunWhenOnBrowser()
}
