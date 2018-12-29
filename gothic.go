package main

import (
	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GenGeo/gen"
)

// as per Generative Parametric Design of Gothic Window Tracery
// by Sven Havemann, Dieter W. Fellner
// (http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.97.8502&rep=rep1&type=pdf)

type archWindowWallParams struct {
	material m.Material
	// four points from llhc to ulhc, counterclockwise
	rectOutline m.Quadrilateral
	// ratio r / distance(pL, pR), excess >= 0.5
	excess float64
	// padding between outline and window
	xPadding      float64
	bottomPadding float64
	// depth of extrusion
	depth float64
	// height of upper segment endpoints pL and pR
	pLpRY float64
	// number of points on a full circle
	numPoints int
}

// returns a rectangular wall with an arch window in it
func archWindowWall(params archWindowWallParams) m.Object {
	// outline of front face, counterclockwise ordered
	llhc, lrhc, urhc, ulhc := params.rectOutline.P1, params.rectOutline.P2, params.rectOutline.P3, params.rectOutline.P4
	rect := []m.Vector{llhc, lrhc, urhc, ulhc}

	minX, maxX := llhc.X, lrhc.X
	leftPadding := minX + params.xPadding
	rightPadding := maxX - params.xPadding

	// inner line of the window, clockwise ordered
	pL, pR := m.Vector{leftPadding, params.pLpRY, 0}, m.Vector{rightPadding, params.pLpRY, 0}
	bpL, bpR := m.Vector{leftPadding, params.bottomPadding, 0}, m.Vector{rightPadding, params.bottomPadding, 0}
	// start with the lower box
	arch := []m.Vector{pR, bpR, bpL, pL}

	// then add the points on the circles of the actual arch
	dist := pR.Sub(pL).Length()
	r := params.excess * dist
	circle := gen.NewCircle(func(t float64) float64 { return r }, params.numPoints)

	mL := pL.Add(m.VectorFromTo(pL, pR).Times(params.excess))
	mR := pR.Add(m.VectorFromTo(pR, pL).Times(params.excess))

	// points returns a list of n points on the circle with radius r
	// around a given midpoint, starting from the right and going counterclockwise
	// with 8 points, [p0, p1, p2] describe the upper right arc
	// and [p2, p3, p4] describes the upper left arc
	// p0 = pR and p(N/2) = pL
	// assumption: numPoints is even
	// note: this only works for pointed archs because rest of the circle
	// overlaps with the wall anyways, but generates pointless triangles
	cL := circle.Points(mL, ex, ey, 0)
	upperLeftArc := make([]m.Vector, 0, params.numPoints/4)
	for i := params.numPoints / 2; i >= params.numPoints/4; i-- {
		upperLeftArc = append(upperLeftArc, cL[i])
		if i == params.numPoints/2 {
			continue
		}
		arch = append(arch, cL[i])
	}

	cR := circle.Points(mR, ex, ey, 0)
	upperRightArc := make([]m.Vector, 0, params.numPoints/4)
	for i := params.numPoints / 4; i >= 0; i-- {
		upperRightArc = append(upperRightArc, cR[i])
		if i == params.numPoints/4 || i == 0 {
			continue
		}
		arch = append(arch, cR[i])
	}

	// triangles of front face
	front := []m.Triangle{}
	t1, t2 := m.QuadrilateralToTriangles(llhc, lrhc, m.Vector{maxX, params.bottomPadding, 0}, m.Vector{minX, params.bottomPadding, 0}, params.material)
	front = append(front, t1, t2)
	t1, t2 = m.QuadrilateralToTriangles(m.Vector{minX, params.bottomPadding, 0}, bpL, pL, m.Vector{minX, params.pLpRY, 0}, params.material)
	front = append(front, t1, t2)
	t1, t2 = m.QuadrilateralToTriangles(bpR, m.Vector{maxX, params.bottomPadding, 0}, m.Vector{maxX, params.pLpRY, 0}, pR, params.material)
	front = append(front, t1, t2)

	// triangles to arch radiating from upper left/right hand corners
	topMidpoint := ulhc.Add(m.VectorFromTo(ulhc, urhc).Times(0.5))

	lPoints := append([]m.Vector{{minX, params.pLpRY, 0}}, upperLeftArc...)
	lPoints = append(lPoints, topMidpoint)
	for i, p1 := range lPoints[:len(lPoints)-1] {
		p2 := lPoints[i+1]
		t := m.NewTriangle(ulhc, p1, p2, params.material)
		front = append(front, t)
	}

	rPoints := append([]m.Vector{topMidpoint}, upperRightArc...)
	rPoints = append(rPoints, m.Vector{maxX, params.pLpRY, 0})
	for i, p1 := range rPoints[:len(rPoints)-1] {
		p2 := rPoints[i+1]
		t := m.NewTriangle(urhc, p1, p2, params.material)
		front = append(front, t)
	}

	ef := gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{rect},
		Inner:    [][]m.Vector{arch},
		Material: params.material,
	}
	return gen.Extrude(ef, m.Vector{0, 0, params.depth})
}

func roundedArchWindowWall(params archWindowWallParams) m.Object {
	params.excess = 0.5
	return archWindowWall(params)
}

func equilateralArchWindowWall(params archWindowWallParams) m.Object {
	params.excess = 1.0
	return archWindowWall(params)
}

type archWindowTraceryParams struct {
	material m.Material
	// ratio r / distance(pL, pR), excess >= 0.5
	excess float64
	// depth of extrusion
	depth float64
	// height of upper segment endpoints pL and pR
	pLpRY float64
	// number of points on a full circle
	numPoints int
}

// TODO: returns window tracery object
func archWindowTracery(params archWindowTraceryParams) m.Object {
	return m.Triangle{}
}
