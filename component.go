package main

import (
	"bytes"
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"io"
	"net/http"
)

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo
	text       string
	backendURL string
	clicks     int
}

func (h *hello) OnMount(ctx app.Context) {
	h.text = "Hello World!"
	h.backendURL = "/post"
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello) Render() app.UI {
	return app.Div().Body(app.Button().ID("reqbut").Text(h.text).OnClick(h.onClick),
		app.Button().ID("urlbut").Text(h.backendURL).OnClick(h.onSwitchBackendURL))
}

func (h *hello) onSwitchBackendURL(ctx app.Context, ev app.Event) {
	if h.backendURL == "/post" {
		h.backendURL = "/empty"
	} else {
		h.backendURL = "/post"
	}
	app.Log("Switched to endpoint " + h.backendURL)
}

func (h *hello) onClick(ctx app.Context, ev app.Event) {
	app.Log("Clicked start")
	h.clicks++
	h.text = fmt.Sprintf("I got clicked %d times", h.clicks)
	app.Logf("Clicked %d times", h.clicks)

	ctx.Async(func() {
		app.Log("Start Backend Request")
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPost, h.backendURL, bytes.NewBuffer([]byte("Test")))
		if err != nil {
			panic(err)
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/text; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// we still read the body of the answer
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		app.Logf("Status: %d Body: '%s'", resp.StatusCode, string(buf))
		app.Log("Backend Request done")
	})
	app.Log("Clicked end")
}
