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

	m.SetBackgroundColor(m.NewColor(10, 10, 10))

	l1 := m.NewDistantLight(m.Vector{1, -1, 1}, m.NewColor(255, 255, 255), 50)
	// l2 := m.NewPointLight(m.Vector{1, 2, 3}, m.NewColor(255, 255, 255), 200)
	scene.AddLights(l1)

	diffMat := &m.DiffuseMaterial{m.NewColor(250, 50, 50)}
	scene.Add(m.NewSphere(m.Vector{3, 1, 5}, 0.5, diffMat))

	// triangles
	r := m.NewQuadrilateral(
		m.Vector{0, 0, 6},
		m.Vector{4, 0, 3},
		m.Vector{0, 0, 0},
		m.Vector{-4, 0, 3},
		diffMat)

	diffMat = &m.DiffuseMaterial{m.NewColor(50, 150, 80)}
	terrain := gridToTriangles(perlinHeightMap(toPointGrid(r, 0.1)), diffMat)
	translation := m.Translate(m.Vector{0, 0, -1})
	scene.Add(m.NewSharedObject(terrain, translation))

	scene.Precompute()

	fmt.Println("Rendering...")

	// aw := render.NewAVI("out.avi", width, height)
	from, to := m.Vector{0, 1, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")

	// for i := 0; i < 30; i++ {
	// 	camera.LookAt(from, to, ey)
	// 	film := render.Render(scene, numWorkers)
	// 	render.AddToAVI(aw, film)
	// 	from = from.Add(m.Vector{0, 0, -0.05})
	// }
	// render.SaveAVI(aw)
}
