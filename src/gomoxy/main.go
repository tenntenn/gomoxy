package main

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
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

	l := NewLabel("hoge", texs)
	eng.Register(l.Node)
	eng.SetTransform(l.Node, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})
	scene.AppendChild(l.Node)
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

func newTextTexture(eng sprite.Engine, dst image.Image, size int, font *truetype.Font, text string) sprite.SubTex {

	fg, bg := image.Black, image.White
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	d := &font.Drawer{
		Dst: dst,
		Src: fg,
		Face: truetype.NewFace(f, truetype.Options{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingFull,
		}),
	}
	d.Dot = fixed.P(0, size*0.8)
	d.DrawString(text)

	t, err := eng.LoadTexture(img)
	if err != nil {
		log.Fatal(err)
	}

	r := dst.Bounds()
	return sprite.SubTex{t, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)}
}
