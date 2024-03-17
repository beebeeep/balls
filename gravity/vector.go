package main

import "math"

type vector struct {
	x, y float64
}

func NewNormalized(x, y float64) vector {
	l := math.Sqrt(x*x + y*y)
	if l == 0 {
		return vector{0, 0}
	}
	return vector{x / l, y / l}
}

func (v vector) Add(u vector) vector {
	return vector{v.x + u.x, v.y + u.y}
}

func (v vector) Sub(u vector) vector {
	return vector{v.x - u.x, v.y - u.y}
}

func (v vector) Multiply(a float64) vector {
	return vector{v.x * a, v.y * a}
}

func (v vector) EntrywiseProduct(a vector) vector {
	return vector{v.x * a.x, v.y * a.y}
}

func (v vector) DotProduct(u vector) float64 {
	return v.x*u.x + v.y*u.y
}

func (v vector) Normalize() vector {
	return NewNormalized(v.x, v.y)
}

func (v vector) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v vector) Reflect(n vector) vector {
	// assuming v is normalized
	return v.Sub(n.Multiply(2.0 * v.DotProduct(n)))
}
