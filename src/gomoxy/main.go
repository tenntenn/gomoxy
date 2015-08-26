package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/golang/freetype/truetype"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
)

var (
	requests []*http.Request

	startTime = time.Now()
	eng       = glsprite.Engine()
	scene     *sprite.Node
	font      *truetype.Font
	label     *Label
)

func main() {
	var proxy *goproxy.ProxyHttpServer
	app.Main(func(a app.App) {
		var sz size.Event
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case lifecycle.Event:
				if e.Crosses(lifecycle.StageAlive) == lifecycle.CrossOn && proxy == nil {
					proxy = goproxy.NewProxyHttpServer()
					//proxy.Verbose = true
					re := regexp.MustCompile(`.*`)
					proxy.OnResponse(goproxy.UrlMatches(re)).DoFunc(
						func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
							if label != nil {
								label.Text = fmt.Sprintf("%s\n%s\n", ctx.Req.URL, label.Text)
							}
							return res
						})
					go func() {
						log.Fatal(http.ListenAndServe(":8888", proxy))
					}()
				}
			case paint.Event:
				onPaint(sz)
				a.EndPaint(e)
			case size.Event:
				sz = e
			}
		}
	})
}

func onPaint(sz size.Event) {
	if scene == nil {
		loadScene(sz)
	}
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, sz)
}

func loadScene(sz size.Event) {
	font = loadFont()

	scene = &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	label = NewLabel(sz, font, 12, image.Rect(0, 0, 400, 400))
	eng.Register(label.Node)
	eng.SetTransform(label.Node, f32.Affine{
		{400, 0, 0},
		{0, 400, 0},
	})
	scene.AppendChild(label.Node)
}

func loadFont() *truetype.Font {
	ttf, err := asset.Open("luximr.ttf")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(ttf)
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	return f
}
