package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
)

func main() {
	var proxy *goproxy.ProxyHttpServer
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case lifecycle.Event:
				if e.Crosses(lifecycle.StageAlive) == lifecycle.CrossOn && proxy == nil {
					proxy = goproxy.NewProxyHttpServer()
					proxy.Verbose = true
					re := regexp.MustCompile(`.*`)
					proxy.OnRequest(goproxy.UrlMatches(re)).DoFunc(
						func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
							log.Println(r.URL)
							return r, nil
						})
					log.Fatal(http.ListenAndServe(":8888", proxy))
				}
			}
		}
	})
}
