package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	winWidth  int16 = 800
	winHeight int16 = 800
	ballColor       = sdl.Color{R: 0x17, G: 0xcf, B: 0xff, A: 0x80}
	massScale       = 1.0
)

type ball struct {
	r     vector
	v     vector
	mass  float64
	color sdl.Color
}

func (b ball) render(r *sdl.Renderer) {
	gfx.FilledCircleColor(r, int32(b.r.x), int32(b.r.y), int32(massScale*b.mass), b.color)
	gfx.CircleColor(r, int32(b.r.x), int32(b.r.y), int32(massScale*b.mass), b.color)
}

func (b *ball) move() {
	b.r.x += b.v.x
	b.r.y += b.v.y

	if b.r.x > float64(winWidth) {
		b.r.x = 0
	}
	if b.r.x < 0 {
		b.r.x = float64(winWidth)
	}
	if b.r.y > float64(winHeight) {
		b.r.y = 0
	}
	if b.r.y < 0 {
		b.r.y = float64(winHeight)
	}

}

func (b *ball) touches(t ball) bool {
	d := b.r.Sub(t.r).Length()
	return d <= massScale*(b.mass+t.mass)
}

func fpsleep(start time.Time) {
	delay := 16*time.Millisecond - time.Since(start)
	if delay < 0 {
		delay = 0
	}
	sdl.Delay(uint32(delay.Milliseconds()))
}

func run(r *sdl.Renderer) {
	running := true
	paused := false
	balls := make([]ball, 100)
	for i := range balls {
		balls[i].r.x = rand.Float64() * float64(winWidth)
		balls[i].r.y = rand.Float64() * float64(winHeight)
		balls[i].v.x = 10 * (0.5 - rand.Float64())
		balls[i].v.y = 10 * (0.5 - rand.Float64())
		balls[i].mass = rand.Float64() * 10
		balls[i].color = sdl.Color{
			R: uint8(rand.Uint32() % 0xff),
			G: uint8(rand.Uint32() % 0xff),
			B: uint8(rand.Uint32() % 0xff),
			A: 0xa0,
		}
	}

	for running {
		startT := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if ev.Type == sdl.MOUSEBUTTONUP {
					//
				}
			case *sdl.KeyboardEvent:
				if ev.Type == sdl.KEYUP {
					if ev.Keysym.Sym == sdl.K_SPACE {
						paused = !paused
					}
				}
			}
		}

		r.SetDrawColor(0, 0, 0, 255)
		r.Clear()
		for i := range balls {
			if !paused {
				balls[i].move()
			}
			for j := i; j < len(balls); j++ {
				if i == j {
					continue
				}
				if balls[i].touches(balls[j]) {
					b1 := &balls[i]
					b2 := &balls[j]

					n := b1.r.Sub(b2.r).Normalize()
					b1.v = b1.v.Sub(n.Multiply(2.0 * b1.v.DotProduct(n) / (n.Length() * n.Length())))
					b2.v = b2.v.Sub(n.Multiply(2.0 * b2.v.DotProduct(n) / (n.Length() * n.Length())))
				}
			}

			balls[i].render(r)
		}
		r.Present()

		fpsleep(startT)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	window, err := sdl.CreateWindow("balls", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("creating window: %s", err)
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("creating renderer: %s", err)
	}
	defer renderer.Destroy()
	run(renderer)
}
