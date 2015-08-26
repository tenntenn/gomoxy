package main

import (
	"image"
	"image/draw"
	"log"
	"math"
	"strings"

	"github.com/golang/freetype/truetype"

	sfont "golang.org/x/exp/shiny/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

type Label struct {
	*sprite.Node
	Text     string
	font     *truetype.Font
	fontSize float64
	rgba     *image.RGBA
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }

func NewLabel(sz size.Event, f *truetype.Font, fsize float64, bounds image.Rectangle) *Label {
	l := &Label{
		Node:     &sprite.Node{},
		font:     f,
		fontSize: fsize,
		rgba:     image.NewRGBA(bounds),
	}

	var last string
	l.Node.Arranger = arrangerFunc(func(e sprite.Engine, n *sprite.Node, t clock.Time) {
		if l.Text == last {
			return
		}
		last = l.Text
		if l.Text == "" {
			e.SetSubTex(l.Node, sprite.SubTex{})
			return
		}
		e.SetSubTex(l.Node, l.newTextTexture(e))
	})

	return l
}

func (l *Label) newTextTexture(eng sprite.Engine) sprite.SubTex {

	fg, bg := image.Black, image.White
	draw.Draw(l.rgba, l.rgba.Bounds(), bg, image.ZP, draw.Src)
	d := &sfont.Drawer{
		Dst: l.rgba,
		Src: fg,
		Face: truetype.NewFace(l.font, truetype.Options{
			Size:    l.fontSize,
			DPI:     72,
			Hinting: sfont.HintingFull,
		}),
	}

	spacing := 1.5
	dy := int(math.Ceil(l.fontSize * spacing))
	for i, s := range strings.Split(l.Text, "\n") {
		d.Dot = fixed.P(0, int(l.fontSize*0.8)+dy*i)
		d.DrawString(s)
	}

	t, err := eng.LoadTexture(l.rgba)
	if err != nil {
		log.Fatal(err)
	}

	return sprite.SubTex{t, l.rgba.Bounds()}
}
