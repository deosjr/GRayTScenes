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
	// outline of front face, counterclockwise ordered
	innerCircle := gen.NewCircle(func(t float64) float64 { return 1 }, 8)
	outerCircle := gen.NewCircle(func(t float64) float64 { return 2 }, 8)
	innerPoints := innerCircle.Points(m.Vector{1, 1, 0}, ex, ey, 0)
	outerPoints := outerCircle.Points(m.Vector{1, 1, 0}, ex, ey, 0)
	// triangles of front face
	front := []m.Triangle{}
	for _, o := range gen.JoinPoints([][]m.Vector{innerPoints, outerPoints}, mat) {
		t := o.(m.Triangle)
		front = append(front, t)
	}

	ef := gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{outerPoints},
		Inner:    [][]m.Vector{innerPoints},
		Material: mat,
	}
	extruded := gen.Extrude(ef, ez)

	translation := m.Translate(m.Vector{0, 0, 3})
	rotation := m.RotateY(math.Pi / 8)
	sharedQ := m.NewSharedObject(extruded, translation.Mul(rotation))
	scene.Add(sharedQ)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 2, 0}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
