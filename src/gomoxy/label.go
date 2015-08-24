package main

import (
	"image"

	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

type Label struct {
	*sprite.Node
	Text string
	t    map[rune]sprite.SubTex
	rect image.Rectangle
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }

func NewLabel(s string) *Label {
	l := &Label{&sprite.Node{}, s, t}

	var last string
	l.Node.Arranger = arrangerFunc(func(e sprite.Engine, n *sprite.Node, t clock.Time) {
		if l.Text == "" || l.Text == last {
			return
		}
		last = l.Text
		e.SetSubTex(l.Node, tex)
	})

	return l
}
