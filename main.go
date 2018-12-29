package main

import (
	"fmt"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
	"github.com/deosjr/GenGeo/gen"
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

	m.SetBackgroundColor(m.NewColor(100, 100, 100))

	mat := &m.DiffuseMaterial{m.NewColor(200, 0, 0)}

	// base case: a cuboid
	box := m.NewAABB(m.Vector{0, 0, 0}, m.Vector{1, 1, 1})
	c := m.NewCuboid(box, mat).Tesselate()
	rotation := m.RotateY(-math.Pi / 8)
	translation := m.Translate(m.Vector{-0.5, 0, 2})
	sharedC := m.NewSharedObject(c, translation.Mul(rotation))
	scene.Add(sharedC)

	// extruded quadrilateral
	mat = &m.DiffuseMaterial{m.NewColor(0, 0, 200)}
	q := []m.Vector{{0, 0, 0}, {1, 0, 0}, {1, 1, 0}, {0.5, 1.5, 0}, {0, 1, 0}}
	extruded := gen.ExtrudeSolidFace(q, ez, mat)

	translation = m.Translate(m.Vector{-1.5, 0, 3})
	rotation = m.RotateY(math.Pi)
	sharedQ := m.NewSharedObject(extruded, translation.Mul(rotation))
	scene.Add(sharedQ)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
