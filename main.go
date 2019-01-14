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

	l1 := m.NewDistantLight(m.Vector{1, -1, 1}, m.NewColor(255, 255, 255), 20)
	l2 := m.NewPointLight(m.Vector{-2, 4.5, 7}, m.NewColor(255, 255, 255), 500)
	scene.AddLights(l1, l2)

	m.SetBackgroundColor(m.NewColor(150, 100, 20))

	mat := &m.DiffuseMaterial{m.NewColor(100, 100, 100)}
	awwparams := archWindowWallParams{
		rectOutline:   m.NewQuadrilateral(m.Vector{0, 0, 0}, m.Vector{1, 0, 0}, m.Vector{1, 2, 0}, m.Vector{0, 2, 0}, mat),
		excess:        1.25,
		xPadding:      0.25,
		bottomPadding: 0.25,
		depth:         0.25,
		pLpRY:         4.0 / 3.0,
		numPoints:     100,
		material:      mat,
	}
	w := archWindowWall(awwparams)

	for y := 0; y < 3; y += 2 {
		for x := -5; x < 10; x++ {
			translation := m.Translate(m.Vector{float64(x), float64(y), 3})
			shared := m.NewSharedObject(w, translation)
			scene.Add(shared)
		}
	}

	mat = &m.DiffuseMaterial{m.NewColor(50, 50, 50)}
	awtparams := archWindowTraceryParams{
		material:       mat,
		excess:         1.25,
		outerWidth:     0.05,
		innerWidth:     0.03,
		verticalOffset: 0.25,
		depth:          0.05,
		pL:             m.Vector{0.25, 4.0 / 3.0, 0},
		pR:             m.Vector{0.75, 4.0 / 3.0, 0},
		bpL:            m.Vector{0.25, 0.25, 0},
		bpR:            m.Vector{0.75, 0.25, 0},
		numPoints:      100,
		numFoils:       3,
	}
	t := archWindowTracery(awtparams)

	for y := 0; y < 3; y += 2 {
		for x := -5; x < 10; x++ {
			translation := m.Translate(m.Vector{float64(x), float64(y), 3.10})
			shared := m.NewSharedObject(t, translation)
			scene.Add(shared)
		}
	}

	//outside floor
	mat = &m.DiffuseMaterial{m.NewColor(50, 150, 0)}
	q := m.NewQuadrilateral(m.Vector{-5, 0, 0}, m.Vector{5, 0, 0}, m.Vector{5, 0, 3}, m.Vector{-5, 0, 3}, mat)
	scene.Add(q.Tesselate())

	//inside floor
	mat = &m.DiffuseMaterial{m.NewColor(250, 150, 0)}
	q = m.NewQuadrilateral(m.Vector{-25, 0, 3}, m.Vector{25, 0, 3}, m.Vector{25, 0, 10}, m.Vector{-25, 0, 10}, mat)
	scene.Add(q.Tesselate())

	//inside wall
	mat = &m.DiffuseMaterial{m.NewColor(100, 100, 100)}
	q = m.NewQuadrilateral(m.Vector{25, 0, 10}, m.Vector{25, 10, 10}, m.Vector{-25, 10, 10}, m.Vector{-25, 0, 10}, mat)
	scene.Add(q.Tesselate())

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
