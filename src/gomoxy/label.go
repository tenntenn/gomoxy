package main

import (
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

type Label struct {
	*sprite.Node
	Text string
	t    map[rune]sprite.SubTex
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }

func NewLabel(s string, t map[rune]sprite.SubTex) *Label {
	l := &Label{&sprite.Node{}, s, t}

	var children []*sprite.Node
	var last string
	l.Node.Arranger = arrangerFunc(func(e sprite.Engine, n *sprite.Node, t clock.Time) {
		if l.Text == "" || l.Text == last {
			return
		}
		last = l.Text

		for _, n := range children {
			l.Node.RemoveChild(n)
		}

		children = make([]*sprite.Node, len(l.Text))

		var x float32 = 0
		for i, r := range l.Text {
			n := &sprite.Node{}
			e.Register(n)
			t, ok := l.t[r]
			if !ok {
				t = l.t['?']
			}
			e.SetSubTex(n, t)
			w, h := float32(t.R.Max.X-t.R.Min.X), float32(t.R.Max.Y-t.R.Min.Y)
			e.SetTransform(n, f32.Affine{
				{w, 0, x + w},
				{0, h, 0},
			})
			x += w
			children[i] = n
			l.Node.AppendChild(n)
		}
	})

	return l
}
