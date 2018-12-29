package main

import (
	"fmt"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
)

var (
	width      uint = 1600
	height     uint = 1200
	numWorkers      = 10

	ex = m.Vector{1, 0, 0}
	ey = m.Vector{0, 1, 0}
	ez = m.Vector{0, 0, 1}
)

func main() {
	fmt.Println("Creating scene...")
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)
	scene := m.NewScene(camera)

	l1 := m.NewDistantLight(m.Vector{1, -1, 1}, m.NewColor(255, 255, 255), 50)
	scene.AddLights(l1)

	m.SetBackgroundColor(m.NewColor(0, 0, 100))

	mat := &m.DiffuseMaterial{m.NewColor(100, 100, 100)}
	q := m.NewQuadrilateral(m.Vector{0, 0, 0}, m.Vector{1, 0, 0}, m.Vector{1, 2, 0}, m.Vector{0, 2, 0}, mat)
	w := equilateralArchWindow(q, mat)

	for y := 0; y < 3; y += 2 {
		for x := -5; x < 10; x++ {
			translation := m.Translate(m.Vector{float64(x), float64(y), 3})
			shared := m.NewSharedObject(w, translation)
			scene.Add(shared)
		}
	}

	mat = &m.DiffuseMaterial{m.NewColor(50, 150, 0)}
	q = m.NewQuadrilateral(m.Vector{-5, 0, 0}, m.Vector{5, 0, 0}, m.Vector{5, 0, 3}, m.Vector{-5, 0, 3}, mat)
	scene.Add(q.Tesselate())

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
