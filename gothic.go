package main

import (
	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GenGeo/gen"
)

// as per Generative Parametric Design of Gothic Window Tracery
// by Sven Havemann, Dieter W. Fellner
// (http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.97.8502&rep=rep1&type=pdf)

// returns a rectangular wall with an arch window in it
// params:
// - excess: ratio r / distance(pL, pR) >= 0.5
func simpleArchWindow(excess float64, mat m.Material) m.Object {
	// outline of front face, counterclockwise ordered
	llhc, lrhc, urhc, ulhc := m.Vector{0, 0, 0}, m.Vector{1, 0, 0}, m.Vector{1, 2, 0}, m.Vector{0, 2, 0}
	rect := []m.Vector{llhc, lrhc, urhc, ulhc}

	// inner line of the window, clockwise ordered
	pL, pR := m.Vector{0.25, 4.0 / 3.0, 0}, m.Vector{0.75, 4.0 / 3.0, 0}
	bpL, bpR := m.Vector{0.25, 0.25, 0}, m.Vector{0.75, 0.25, 0}
	// start with the lower box
	arch := []m.Vector{pR, bpR, bpL, pL}

	// then add the points on the circles of the actual arch
	dist := pR.Sub(pL).Length()
	r := excess * dist
	numPoints := 100
	circle := gen.NewCircle(func(t float64) float64 { return r }, numPoints)

	mL := pL.Add(m.VectorFromTo(pL, pR).Times(excess))
	mR := pR.Add(m.VectorFromTo(pR, pL).Times(excess))

	// points returns a list of n points on the circle with radius r
	// around a given midpoint, starting from the right and going counterclockwise
	// with 8 points, [p0, p1, p2] describe the upper right arc
	// and [p2, p3, p4] describes the upper left arc
	// p0 = pR and p(N/2) = pL
	// assumption: numPoints is even
	cL := circle.Points(mL, ex, ey, 0)
	upperLeftArc := make([]m.Vector, 0, numPoints/4)
	for i := numPoints / 2; i >= numPoints/4; i-- {
		upperLeftArc = append(upperLeftArc, cL[i])
		if i == numPoints/2 {
			continue
		}
		arch = append(arch, cL[i])
	}

	cR := circle.Points(mR, ex, ey, 0)
	upperRightArc := make([]m.Vector, 0, numPoints/4)
	for i := numPoints / 4; i >= 0; i-- {
		upperRightArc = append(upperRightArc, cR[i])
		if i == numPoints/4 || i == 0 {
			continue
		}
		arch = append(arch, cR[i])
	}

	// triangles of front face
	front := []m.Triangle{}
	t1, t2 := m.QuadrilateralToTriangles(m.Vector{0, 0, 0}, m.Vector{1, 0, 0}, m.Vector{1, 0.25, 0}, m.Vector{0, 0.25, 0}, mat)
	front = append(front, t1, t2)
	t1, t2 = m.QuadrilateralToTriangles(m.Vector{0, 0.25, 0}, bpL, pL, m.Vector{0, 4.0 / 3.0, 0}, mat)
	front = append(front, t1, t2)
	t1, t2 = m.QuadrilateralToTriangles(bpR, m.Vector{1, 0.25, 0}, m.Vector{1, 4.0 / 3.0, 0}, pR, mat)
	front = append(front, t1, t2)

	lPoints := append([]m.Vector{{0, 4.0 / 3.0, 0}}, upperLeftArc...)
	lPoints = append(lPoints, m.Vector{0.5, 2, 0})
	for i, p1 := range lPoints[:len(lPoints)-1] {
		p2 := lPoints[i+1]
		t := m.NewTriangle(m.Vector{0, 2, 0}, p1, p2, mat)
		front = append(front, t)
	}

	rPoints := append([]m.Vector{{0.5, 2, 0}}, upperRightArc...)
	rPoints = append(rPoints, m.Vector{1, 4.0 / 3.0, 0})
	for i, p1 := range rPoints[:len(rPoints)-1] {
		p2 := rPoints[i+1]
		t := m.NewTriangle(m.Vector{1, 2, 0}, p1, p2, mat)
		front = append(front, t)
	}

	ef := gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{rect},
		Inner:    [][]m.Vector{arch},
		Material: mat,
	}
	return gen.Extrude(ef, m.Vector{0, 0, 0.5})
}

func roundedArchWindow(mat m.Material) m.Object {
	return simpleArchWindow(0.5, mat)
}

func equilateralArchWindow(mat m.Material) m.Object {
	return simpleArchWindow(1.0, mat)
}
