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

	m.SetBackgroundColor(m.NewColor(15, 200, 215))

	mat := &m.DiffuseMaterial{m.NewColor(255, 255, 255)}
	c := gen.NewCircle(func(t float64) float64 { return 0.03 }, 10)
	c1 := c.Points(m.Vector{0, 0, 0}, m.Vector{1, 0, 0}, m.Vector{0, 0, 1}, 0)
	c2 := c.Points(m.Vector{0, 1, 0}, m.Vector{1, 0, 0}, m.Vector{0, 0, 1}, 0)
	trunkTriangles := gen.JoinPoints([][]m.Vector{c2, c1}, mat)
	trunk := m.NewComplexObject(trunkTriangles)

	mat = &m.DiffuseMaterial{m.NewColor(255, 180, 40)}
	leaves := m.NewSphere(m.Vector{0, 1.5, 0}, 0.5, mat)
	treeObject := m.NewComplexObject([]m.Object{trunk, leaves})

	// poisson disc sampling allows a random distribution
	// which enforces distance at least r between points
	q := m.Quadrilateral{P1: m.Vector{-5.0, -5.0, 0.0}, P3: m.Vector{-2.0, 5.0, 0.0}}
	points := poisson(q, 1.0)
	for _, p := range points {
		translation := m.Translate(m.Vector{p.X, 0, p.Y})
		rotation := m.RotateY(math.Pi / 4.0)
		tree := m.NewSharedObject(treeObject, translation.Mul(rotation))
		scene.Add(tree)
	}

	q = m.Quadrilateral{P1: m.Vector{2.0, -5.0, 0.0}, P3: m.Vector{5.0, 5.0, 0.0}}
	points = poisson(q, 1.0)
	for _, p := range points {
		translation := m.Translate(m.Vector{p.X, 0, p.Y})
		rotation := m.RotateY(math.Pi / 4.0)
		tree := m.NewSharedObject(treeObject, translation.Mul(rotation))
		scene.Add(tree)
	}

	flatColor := &m.DiffuseMaterial{m.NewColor(255, 180, 40)}
	steepColor := &m.DiffuseMaterial{m.NewColor(50, 50, 50)}
	fmat := &m.PosFuncMat{
		Func: func(si *m.SurfaceInteraction, l m.Light) m.Color {
			if si.Normal.Dot(ey) > 0.99 {
				return flatColor.GetColor(si, l)
			}
			return steepColor.GetColor(si, l)
		},
	}
	q = m.NewQuadrilateral(
		m.Vector{-5.0, 0.0, 5.0},
		m.Vector{5.0, 0.0, 5.0},
		m.Vector{5.0, 0.0, -5.0},
		m.Vector{-5.0, 0.0, -5.0},
		fmat,
	)
	grid := toPointGrid(q, 0.1)
	// perlin := perlinHeightMap(grid, 3, []float64{0.5, 0.7, 0.25, 0.15}, 3.75)
	perlin := perlinHeightMap(grid, 3, []float64{0.5, 0.7, 0.25, 0.15}, 1.75)
	ground := gridToTriangles(perlin, fmat)
	scene.Add(ground)

	scene.Precompute()

	fmt.Println("Rendering...")

	from, to := m.Vector{0, 1, -2}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)
	film := render.Render(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
