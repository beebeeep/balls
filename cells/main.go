package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth   int16   = 800
	winHeight  int16   = 800
	cellSize   int16   = 40
	eps        float64 = 0.01
	speedScale float64 = 30.0
)

var (
	cells      [winWidth / cellSize][winHeight / cellSize]bool
	leftColor  = sdl.Color{R: 0x17, G: 0xcf, B: 0xff, A: 0x80}
	rightColor = sdl.Color{R: 0xff, G: 0x89, B: 0x17, A: 0x80}
)

type ball struct {
	isLeft bool
	x, y   float64
	dx, dy float64
	r      float64
}

func fpsleep(start time.Time) {
	delay := 16*time.Millisecond - time.Now().Sub(start)
	if delay < 0 {
		delay = 0
	}
	sdl.Delay(uint32(delay.Milliseconds()))
}

func (b *ball) move() {
	b.x += b.dx
	b.y += b.dy
}

func (b ball) getCell() (int, int) {
	row := math.Floor(b.x / float64(cellSize))
	column := math.Floor(b.y / float64(cellSize))
	return int(row), int(column)
}

func (b *ball) collide() {
	if b.x-b.r < 0 || b.x+b.r > float64(winWidth) {
		fmt.Println("collide x")
		b.dx *= -1.0
		return
	}
	if b.y-b.r < 0 || b.y+b.r > float64(winHeight) {
		fmt.Println("collide y")
		b.dy *= -1.0
		return
	}

	cx, cy := b.getCell()
	// fmt.Printf("ball %v, row %d column %d\n", b.isLeft, r, c)
	if b.dx > 0 && cx <= len(cells)-2 {
		// check cell at the right
		if cells[cx+1][cy] != b.isLeft && b.x+b.r >= float64(cx+1)*float64(cellSize) {
			cells[cx+1][cy] = b.isLeft
			b.dx *= -1.0
			fmt.Println("boop cell right", cy, cx)
			return
		}

	}
	if b.dx < 0 && cx >= 2 {
		// check cell at the left
		if cells[cx-1][cy] != b.isLeft && b.x-b.r <= float64(cx)*float64(cellSize) {
			cells[cx-1][cy] = b.isLeft
			b.dx *= -1.0
			fmt.Println("boop cell left", cy, cx)
			return
		}

	}
	if b.dy > 0 && cy <= len(cells[0])-2 {
		// check cell below
		if cells[cx][cy+1] != b.isLeft && b.y+b.r >= float64(cy+1)*float64(cellSize) {
			cells[cx][cy+1] = b.isLeft
			b.dy *= -1.0
			fmt.Println("boop cell below", cy, cx)
			return
		}

	}
	if b.dy < 0 && cy >= 2 {
		// check cell above
		if cells[cx][cy-1] != b.isLeft && b.y-b.r <= float64(cy)*float64(cellSize) {
			cells[cx][cy-1] = b.isLeft
			b.dy *= -1.0
			fmt.Println("boop cell above", cy, cx)
			return
		}
	}

}

func drawCells(r *sdl.Renderer) {
	for i := range cells {
		for j := range cells[i] {
			x := int16(i) * cellSize
			y := int16(j) * cellSize
			color := rightColor
			if cells[i][j] {
				color = leftColor
			}
			gfx.FilledPolygonColor(r,
				[]int16{x, x + cellSize, x + cellSize, x},
				[]int16{y + cellSize, y + cellSize, y, y},
				color,
			)
		}
	}
}

func (b ball) render(r *sdl.Renderer) {
	color := leftColor
	if b.isLeft {
		color = rightColor
	}
	color.A = 0xff
	gfx.FilledCircleColor(r, int32(b.x), int32(b.y), int32(b.r), color)
}

func run(r *sdl.Renderer) {
	running := true
	paused := false
	balls := []ball{
		{isLeft: true, r: float64(cellSize) / 2,
			x:  float64(winWidth) / 4.0,
			y:  float64(winHeight) / 2.0,
			dx: speedScale * (rand.Float64() - 0.5),
			dy: speedScale * (rand.Float64() * 0.5),
		},
		{isLeft: false, r: float64(cellSize) / 2,
			x:  3 * float64(winWidth) / 4.0,
			y:  float64(winHeight) / 2.0,
			dx: speedScale * (rand.Float64() - 0.5),
			dy: speedScale * (rand.Float64() - 0.5),
		},
	}

	for i := range cells {
		for j := range cells[i] {
			if i < len(cells)/2 {
				cells[i][j] = true
			} else {
				cells[i][j] = false
			}
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
					fmt.Printf("%v, %v\n", math.Floor(float64(ev.Y)/float64(cellSize)), math.Floor(float64(ev.X)/float64(cellSize)))
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
		drawCells(r)
		for i := range balls {
			if !paused {
				balls[i].move()
				balls[i].collide()
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
