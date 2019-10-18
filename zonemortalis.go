package main

import (
	"math"
	"math/rand"
	"time"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GenGeo/gen"
)

// 4x4 tiles
// each tile is 6x6 squares and 1' x 1'
// each square is ~50mm x 50mm
// squares can be walls in one of two directions or corners joining wall sections
// otherwise they are empty squares

var zm_alpha = [6][6]rune{
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', 'x', '-', '-', 'x', '-'},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', 'x', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
}

var zm_beta = [6][6]rune{
	{' ', ' ', ' ', ' ', '|', ' '},
	{' ', 'x', '-', '-', 'x', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', 'x', '-', '-', 'x', ' '},
	{' ', ' ', ' ', ' ', '|', ' '},
}

var zm_gamma = [6][6]rune{
	{' ', '|', ' ', ' ', '|', ' '},
	{' ', 'x', ' ', ' ', 'x', '-'},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', 'x', ' ', ' ', 'x', '-'},
	{' ', '|', ' ', ' ', '|', ' '},
}

var zm_delta = [6][6]rune{
	{' ', '|', ' ', ' ', '|', ' '},
	{'-', 'x', ' ', ' ', 'x', '-'},
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', ' '},
	{'-', 'x', ' ', ' ', 'x', '-'},
	{' ', '|', ' ', ' ', '|', ' '},
}

var zm_epsilon = [6][6]rune{
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', 'x', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', 'x', ' ', ' ', 'x', '-'},
	{' ', ' ', ' ', ' ', ' ', ' '},
}

var zm_zeta = [6][6]rune{
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', 'x', ' '},
	{' ', ' ', ' ', ' ', '|', ' '},
	{' ', ' ', ' ', ' ', '|', ' '},
	{' ', ' ', ' ', ' ', 'x', ' '},
	{' ', ' ', ' ', ' ', '|', ' '},
}

var zm_eta = [6][6]rune{
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', 'x', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
	{' ', 'x', '-', '-', 'x', ' '},
	{' ', '|', ' ', ' ', ' ', ' '},
}

var zm_theta = [6][6]rune{
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', ' '},
}

type ZoneMortalisParameters struct {
	floor    m.Object
	wall     m.Object
	corner   m.Object
	material m.Material
}

func NewZoneMortalis(p ZoneMortalisParameters) m.Object {
	alpha := newTile(zm_alpha, p)
	beta := newTile(zm_beta, p)
	gamma := newTile(zm_gamma, p)
	delta := newTile(zm_delta, p)
	epsilon := newTile(zm_epsilon, p)
	zeta := newTile(zm_zeta, p)
	eta := newTile(zm_eta, p)
	theta := newTile(zm_theta, p)

	board := []m.Object{
		alpha, alpha, beta, beta, gamma, gamma, delta, delta,
		epsilon, epsilon, zeta, zeta, eta, eta, theta, theta,
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	tiles := make([]m.Object, 16)
	perm := r.Perm(16)
	for i, randIndex := range perm {
		tiles[i] = board[randIndex]
	}

	i := 0
	for z := 0; z < 4; z++ {
		for x := 0; x < 4; x++ {
			transform := m.Translate(m.Vector{float32(x * 300), 0, float32(z * 300)})
			if n := r.Intn(4); n != 0 {
				transform = transform.Mul(m.Translate(m.Vector{150, 0, 150}))
				transform = transform.Mul(m.RotateY(float64(n) * math.Pi / 2.0).Mul(m.Translate(m.Vector{-150, 0, -150})))
			}
			tile := tiles[i]
			tiles[i] = m.NewSharedObject(tile, transform)
			i++
		}
	}
	return m.NewComplexObject(tiles)
}

func newTile(squaredef [6][6]rune, p ZoneMortalisParameters) m.Object {
	// extrude a tile
	t := []m.Vector{{0, 0, 300}, {300, 0, 300}, {300, 0, 0}, {0, 0, 0}}
	t1 := m.NewTriangle(t[0], t[1], t[2], p.material)
	t2 := m.NewTriangle(t[0], t[2], t[3], p.material)
	front := []m.Triangle{t1, t2}

	ef := gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{t},
		Material: p.material,
	}
	tile := ef.Extrude(m.Vector{0, 1, 0})

	// add squares
	squares := []m.Object{}
	for z := 0; z < 6; z++ {
		for x := 0; x < 6; x++ {
			switch squaredef[z][x] {
			case 'x':
				transform := m.Translate(m.Vector{float32(50 * x), 1, 250.0 - float32(50*z)})
				c := m.NewSharedObject(p.corner, transform)
				squares = append(squares, c)
			case '-':
				transform := m.Translate(m.Vector{float32(50 * x), 1, 250.0 - float32(50*z)})
				//transform = transform.Mul(m.Translate(m.Vector{25, 0, 25}))
				//transform = transform.Mul(m.RotateY(math.Pi / 2.0).Mul(m.Translate(m.Vector{-25, 0, -25})))
				w := m.NewSharedObject(p.wall, transform)
				squares = append(squares, w)
			case '|':
				transform := m.Translate(m.Vector{float32(50 * x), 1, 250.0 - float32(50*z)})
				transform = transform.Mul(m.Translate(m.Vector{25, 0, 25}))
				transform = transform.Mul(m.RotateY(math.Pi / 2.0).Mul(m.Translate(m.Vector{-25, 0, -25})))
				w := m.NewSharedObject(p.wall, transform)
				squares = append(squares, w)
			case ' ':
				transform := m.Translate(m.Vector{0.5 + float32(50*x), 1, 0.5 + 250.0 - float32(50*z)})
				s := m.NewSharedObject(p.floor, transform)
				squares = append(squares, s)
			}
		}
	}

	return m.NewComplexObject(append(squares, tile))
}
