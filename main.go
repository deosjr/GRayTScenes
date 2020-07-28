package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/pzsz/voronoi"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
	"github.com/deosjr/GenGeo/gen"
)

var (
	width      uint = 1600
	height     uint = 1200
	numWorkers      = 10
	numSamples      = 10

	ex = m.Vector{1, 0, 0}
	ey = m.Vector{0, 1, 0}
	ez = m.Vector{0, 0, 1}
)

func main() {

	fmt.Println("Creating scene...")
	m.SIMD_ENABLED = true
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)
	scene := m.NewScene(camera)

	pointLight := m.NewPointLight(m.Vector{0, 10, -100}, m.NewColor(255, 255, 255), 50000000)
	scene.AddLights(pointLight)

	//l1 := m.NewDistantLight(m.Vector{-1, -1, 1}, m.NewColor(255, 255, 255), 20)
	//l2 := m.NewDistantLight(m.Vector{1, -1, 1}, m.NewColor(255, 255, 255), 20)
	// l2 := m.NewPointLight(m.Vector{-2, 4.5, 7}, m.NewColor(255, 255, 255), 500)
	//scene.AddLights(l1, l2)

	m.SetBackgroundColor(m.NewColor(15, 200, 215))

	q := m.Quadrilateral{P1: m.Vector{-5.0, -5.0, 0.0}, P2: m.Vector{5.0, -5.0, 0.0}, P3: m.Vector{5.0, 5.0, 0.0}, P4: m.Vector{-5.0, -5.0, 0.0}}
	points := poisson(q, 1.0)
	sites := make([]voronoi.Vertex, len(points))
	for i, p := range points {
		sites[i] = voronoi.Vertex{float64(p.X), float64(p.Y)}
	}
	bbox := voronoi.NewBBox(-5.0, 5.0, -5.0, 5.0)
	diagram := voronoi.ComputeDiagram(sites, bbox, true)
	cells := make([][]m.Vector, len(diagram.Cells))
	for i, d := range diagram.Cells {
		cell := make([]m.Vector, len(d.Halfedges))
		he0 := d.Halfedges[0].Edge
		he1 := d.Halfedges[1].Edge

		// all of this shit because although halfedges are ordered,
		// their Va and Vb vertices are not...

		e0va := m.Vector{float32(he0.Va.Vertex.X), float32(he0.Va.Vertex.Y), 0}
		e0vb := m.Vector{float32(he0.Vb.Vertex.X), float32(he0.Vb.Vertex.Y), 0}
		e1va := m.Vector{float32(he1.Va.Vertex.X), float32(he1.Va.Vertex.Y), 0}
		e1vb := m.Vector{float32(he1.Vb.Vertex.X), float32(he1.Vb.Vertex.Y), 0}

		a0a1 := e0va.Sub(e1va).Length()
		a0b1 := e0va.Sub(e1vb).Length()
		b0a1 := e0vb.Sub(e1va).Length()
		b0b1 := e0vb.Sub(e1vb).Length()

		var prev m.Vector

		// find the first 2 vertices by finding the duplicate between va/vb of first 2 halfedges

		if a0a1 == 0 {
			cell[0] = e0vb
			cell[1] = e0va
			prev = e0va
		} else if a0b1 == 0 {
			cell[0] = e0vb
			cell[1] = e0va
			prev = e0va
		} else if b0a1 == 0 {
			cell[0] = e0va
			cell[1] = e0vb
			prev = e0vb
		} else if b0b1 == 0 {
			cell[0] = e0va
			cell[1] = e0vb
			prev = e0vb
		}

		for j, he := range d.Halfedges[2:] {
			va := m.Vector{float32(he.Edge.Va.Vertex.X), float32(he.Edge.Va.Vertex.Y), 0}
			vb := m.Vector{float32(he.Edge.Vb.Vertex.X), float32(he.Edge.Vb.Vertex.Y), 0}

			valen := prev.Sub(va).Length()
			vblen := prev.Sub(vb).Length()
			next := va
			if vblen != 0 && vblen < valen {
				next = vb
			}

			cell[j+2] = next
			prev = next
		}
		cells[i] = cell
	}

	for _, cell := range cells {
		mat := &m.DiffuseMaterial{Color: m.NewColor(uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)))}
		depth := -5 * rand.Float32()
		esf := gen.ExtrudeSolidFace(cell, m.Vector{0, 0, depth}, mat)
		scene.Add(esf)
	}

	radmat := &m.RadiantMaterial{Color: m.NewColor(176, 237, 255)}
	skybox := m.NewCuboid(m.NewAABB(m.Vector{-1000, -1000, -1000}, m.Vector{1000, 1000, 1000}), radmat)
	triangles := skybox.TesselateInsideOut()
	skyboxObject := m.NewTriangleComplexObject(triangles)
	scene.Add(skyboxObject)
	scene.Emitters = triangles
	scene.Precompute()

	fmt.Println("Rendering...")

	//	from, to := m.Vector{25, 150, -50}, m.Vector{25, 0, 150}
	from, to := m.Vector{0, 1, -10}, m.Vector{0, 0, 10}
	camera.LookAt(from, to, ey)

	params := render.Params{
		Scene:        scene,
		NumWorkers:   numWorkers,
		NumSamples:   numSamples,
		AntiAliasing: true,
		//TracerType: 	m.WhittedStyle,
		TracerType:   m.PathNextEventEstimate,
	}
	film := render.Render(params)
	film.SaveAsPNG("out.png")
}
