//go:build !wasm

package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"strconv"
	"time"
)

func testing() {
	l := launcher.New().
		Headless(false).
		Devtools(true)

	defer l.Cleanup() // remove launcher.FlagUserDataDir

	url := l.MustLaunch()
	browser := rod.New().
		ControlURL(url).
		//Trace(false).
		//SlowMotion(2000 * time.Millisecond).
		MustConnect()

	defer browser.MustClose()
	page := browser.MustPage()

	logConsole := true
	// Listen for all events of console output.
	go page.EachEvent(
		func(e *proto.RuntimeConsoleAPICalled) {
			if logConsole {
				prefix := log.Prefix()
				log.SetPrefix("[con] ")
				log.Printf("%s", page.MustObjectsToJSON(e.Args))
				log.SetPrefix(prefix)
			}
		},
		func(e *proto.RuntimeExceptionThrown) {
			prefix := log.Prefix()
			log.SetPrefix("[exc] ")
			log.Printf("%#v", e.ExceptionDetails)
			log.SetPrefix(prefix)
		},
	)()

	router := browser.HijackRequests()
	defer router.MustStop()
	// without some route defined it hangs here
	router.MustAdd("not-to-happen", func(*rod.Hijack) {})
	go router.Run()

	wait := page.WaitEvent(&proto.PageLoadEventFired{})
	page.MustNavigate("http://" + siteURL)
	wait()

	// Wait till the loader has vanished
	page.MustWait(`() => document.getElementById("app-wasm-loader") == null`)

	// Click without modifying the result
	fmt.Println("\n")
	log.Println("Clicking the button without hijacking active")
	wait = page.MustWaitRequestIdle()
	page.MustElement("#reqbut").MustClick()
	wait()

	fmt.Println("\n")
	log.Println("Clicking the button with hijacking 'Data' active")
	router.MustAdd("*/post", injectData)
	wait = page.MustWaitRequestIdle()
	page.MustElement("#reqbut").MustClick()
	wait()
	// remove that hijacker
	router.MustRemove("*/post")

	// switching the endpoint url
	fmt.Println("\n")
	log.Println("Clicking the button to change the endpoint url")
	page.MustElement("#urlbut").MustClick()

	log.Println("Clicking the button, button replies with an empty body")
	router.MustAdd("*/post", injectData)

	wait = page.MustWaitRequestIdle()
	page.MustElement("#reqbut").MustClick()
	wait()
	log.Println("Notice: The frontend got the empty body response just as expected.")

	// remove that hijacker
	router.MustRemove("*/post")

	// This is the Error case which actually should work
	fmt.Println("\n")
	log.Println("Clicking the button with hijacking 'EmptyBody' active")
	router.MustAdd("*/empty", injectEmpty)
	wait = page.MustWaitRequestIdle()
	page.MustElement("#reqbut").MustClick()
	wait()
	log.Println("Notice: The javascript console ([con]) says the backend request started, but now it hangs. The reasons seems to be that there is no (empty) body.")

	time.Sleep(5 * time.Second)
	//utils.Pause()
}

func injectData(ctx *rod.Hijack) {
	log.Println("Start: Intercept and modifying POST with data")
	//ctx.MustLoadResponse()
	body := "Injected response"
	ctx.Response.SetBody(body).SetHeader(
		"Access-Control-Allow-Origin", siteURL,
		"Access-Control-Expose-Headers", "Link",
		"Cache-Control", "public, max-age=60",
		"Content-Length", strconv.Itoa(len(body)),
		"Content-Type", "application/text",
		//"Date", "Mon, 16 May 2022 20:53:07 GMT",
		"Vary", "Origin",
	)
	log.Println("End: Intercept and modifying POST with data")
	//ctx.Response.Fail(proto.NetworkErrorReasonAddressUnreachable)
}

func injectEmpty(ctx *rod.Hijack) {
	log.Println("Start: Intercept and modifying POST with empty BODY")
	//ctx.MustLoadResponse()
	ctx.Response.SetBody([]byte{}).SetHeader(
		"Access-Control-Allow-Origin", siteURL,
		"Access-Control-Expose-Headers", "Link",
		"Cache-Control", "public, max-age=60",
		"Content-Length", "0",
		"Content-Type", "application/text",
		//"Date", "Mon, 16 May 2022 20:53:07 GMT",
		"Vary", "Origin",
	)
	log.Println("End: Intercept and modifying POST with empty BODY")
	//ctx.Response.Fail(proto.NetworkErrorReasonAddressUnreachable)
}
