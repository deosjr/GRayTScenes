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
	numSamples      = 100

	ex = m.Vector{1, 0, 0}
	ey = m.Vector{0, 1, 0}
	ez = m.Vector{0, 0, 1}
)

func main() {

	fmt.Println("Creating scene...")
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)
	scene := m.NewScene(camera)

	//	l1 := m.NewDistantLight(m.Vector{-1, -1, 1}, m.NewColor(255, 255, 255), 20)
	//	l2 := m.NewDistantLight(m.Vector{1, -1, 1}, m.NewColor(255, 255, 255), 20)
	// l2 := m.NewPointLight(m.Vector{-2, 4.5, 7}, m.NewColor(255, 255, 255), 500)
	//	scene.AddLights(l1, l2)

	m.SetBackgroundColor(m.NewColor(0, 0, 0))
	mat := &m.DiffuseMaterial{Color: m.NewColor(100, 100, 100)}

	// extrude a wall
	awwparams := archWindowWallParams{
		rectOutline:   m.NewQuadrilateral(m.Vector{0, 0, 0}, m.Vector{50, 0, 0}, m.Vector{50, 100, 0}, m.Vector{0, 100, 0}, mat),
		excess:        1.25,
		xPadding:      15,
		bottomPadding: 20,
		depth:         10,
		pLpRY:         2 * (100.0 / 3.0),
		numPoints:     100,
		material:      mat,
	}
	wallOrigin := archWindowWall(awwparams)
	transform := m.Translate(m.Vector{0, 0, 20})
	wall := m.NewSharedObject(wallOrigin, transform)

	// extrude a corner
	t := []m.Vector{{0, 0, 50}, {50, 0, 50}, {50, 0, 0}, {0, 0, 0}}
	t1 := m.NewTriangle(t[0], t[1], t[2], mat)
	t2 := m.NewTriangle(t[0], t[2], t[3], mat)
	front := []m.Triangle{t1, t2}

	ef := gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{t},
		Material: mat,
	}
	corner := ef.Extrude(m.Vector{0, 100, 0})

	// extrude a square
	t = []m.Vector{{0, 0, 49}, {49, 0, 49}, {49, 0, 0}, {0, 0, 0}}
	t1 = m.NewTriangle(t[0], t[1], t[2], mat)
	t2 = m.NewTriangle(t[0], t[2], t[3], mat)
	front = []m.Triangle{t1, t2}

	ef = gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{t},
		Material: mat,
	}
	square := ef.Extrude(m.Vector{0, 0.3, 0})

	zmp := ZoneMortalisParameters{
		floor:    square,
		wall:     wall,
		corner:   corner,
		material: mat,
	}

	board := NewZoneMortalis(zmp)
	scene.Add(board)

	radmat := &m.RadiantMaterial{Color: m.NewColor(176, 237, 255)}
	skybox := m.NewCuboid(m.NewAABB(m.Vector{-1000, -1000, -1000}, m.Vector{1000, 1000, 1000}), radmat)
	scene.Add(skybox.TesselateInsideOut())

	scene.Precompute()

	fmt.Println("Rendering...")

	//	from, to := m.Vector{25, 150, -50}, m.Vector{25, 0, 150}
	from, to := m.Vector{600, 250, -50}, m.Vector{600, 0, 250}
	camera.LookAt(from, to, ey)
	film := render.RenderWithPathTracer(scene, numWorkers, numSamples)
	//film := render.RenderNaive(scene, numWorkers)
	film.SaveAsPNG("out.png")
}
