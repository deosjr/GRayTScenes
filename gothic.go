package main

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GenGeo/gen"
)

// as per Generative Parametric Design of Gothic Window Tracery
// by Sven Havemann, Dieter W. Fellner
// (http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.97.8502&rep=rep1&type=pdf)

// an arch consists of two lists of points:
// - all points in the left arch (including the top point)
// - all points in the right arch (excluding top point)
type arch struct {
	left  []m.Vector
	right []m.Vector
}

// createArch returns an arch with endpoints pL, pR and midpoints to the arcs mL and mR.
// mL and mR are parameters instead of calculated from excess here, so that this function
// can be used for calculating offset for an arch too
func createArch(pL, pR, mL, mR m.Vector, numPoints int) arch {
	// points returns a list of n points on the circle with radius r
	// around a given midpoint, starting from the right and going counterclockwise
	// the first quarter of points form the upper right quarter of the circle
	// and the second quarter of points form the upper left quarter
	// the left and right arc of the arch are subsets of these points,
	// since the circles meet earlier if the excess is greater than 0.5

	// the top of the arch is given by translating the middle of line pLpR up with y
	// where y is sqrt(pLpR * (r-(pLpR/4))), where r is the circle radius

	// assumption: length of mR-pR given as equal
	r := m.VectorFromTo(mL, pL).Length()
	circle := gen.NewCircle(func(t float64) float64 { return r }, numPoints)

	pLpR := m.VectorFromTo(pL, pR).Length()
	middle := pL.Add(m.VectorFromTo(pL, pR).Times(0.5))
	y := math.Sqrt(pLpR * (r - (pLpR / 4.0)))
	top := middle.Add(m.Vector{0, y, 0})

	// TODO: generate the quarter circle instead of picking points from the full one
	cL := circle.Points(mL, ex, ey, 0)
	upperLeftArc := []m.Vector{}
	for i := numPoints / 2; i >= numPoints/4; i-- {
		p := cL[i]
		if p.X > top.X {
			break
		}
		upperLeftArc = append(upperLeftArc, p)
	}
	upperLeftArc = append(upperLeftArc, top)

	cR := circle.Points(mR, ex, ey, 0)
	upperRightArc := []m.Vector{top}
	for i := numPoints / 4; i >= 0; i-- {
		p := cR[i]
		if p.X <= top.X {
			continue
		}
		upperRightArc = append(upperRightArc, p)
	}

	return arch{
		left:  upperLeftArc,
		right: upperRightArc,
	}
}

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

// archWindowWall returns a rectangular wall with an arch window in it
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
	wholeArch := []m.Vector{pR, bpR, bpL, pL}

	// then add the points on the circles of the actual arch
	mL := pL.Add(m.VectorFromTo(pL, pR).Times(params.excess))
	mR := pR.Add(m.VectorFromTo(pR, pL).Times(params.excess))

	arch := createArch(pL, pR, mL, mR, params.numPoints)
	wholeArch = append(wholeArch, arch.left...)
	wholeArch = append(wholeArch, arch.right...)

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

	lPoints := append([]m.Vector{{minX, params.pLpRY, 0}}, arch.left...)
	lPoints = append(lPoints, topMidpoint)
	for i, p1 := range lPoints[:len(lPoints)-1] {
		p2 := lPoints[i+1]
		t := m.NewTriangle(ulhc, p1, p2, params.material)
		front = append(front, t)
	}

	rPoints := append([]m.Vector{topMidpoint}, arch.right...)
	rPoints = append(rPoints, m.Vector{maxX, params.pLpRY, 0})
	for i, p1 := range rPoints[:len(rPoints)-1] {
		p2 := rPoints[i+1]
		t := m.NewTriangle(urhc, p1, p2, params.material)
		front = append(front, t)
	}

	ef := gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{rect},
		Inner:    [][]m.Vector{wholeArch},
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
	// width of tracery
	offset float64
	// depth of extrusion
	depth float64
	// pL and pR, left and right endpoints of the arch
	pL m.Vector
	pR m.Vector
	// bpL and bpR, left and right bottom points of the window
	bpL m.Vector
	bpR m.Vector
	// number of points on a full circle
	numPoints int
}

// archWindowTracery returns window tracery object
func archWindowTracery(params archWindowTraceryParams) m.Object {
	mL := params.pL.Add(m.VectorFromTo(params.pL, params.pR).Times(params.excess))
	mR := params.pR.Add(m.VectorFromTo(params.pR, params.pL).Times(params.excess))

	outerArch := createArch(params.pL, params.pR, mL, mR, params.numPoints)
	outerFrame := append([]m.Vector{params.pR, params.bpR, params.bpL, params.pL}, outerArch.left...)
	outerFrame = append(outerFrame, outerArch.right...)

	ipL := m.Vector{params.pL.X + params.offset, params.pL.Y, params.pL.Z}
	ipR := m.Vector{params.pR.X - params.offset, params.pR.Y, params.pR.Z}
	// TODO: note this *2 is so that innerarch has at least as many points as outerarch
	// otherwise, gen.JoinPoints will throw an indexoutofbounds exception
	// to solve this issue, fix the above TODO in createArch()
	innerArch := createArch(ipL, ipR, mL, mR, params.numPoints*2)

	ibpL := m.Vector{params.bpL.X + params.offset, params.bpL.Y + params.offset, params.bpL.Z}
	ibpR := m.Vector{params.bpR.X - params.offset, params.bpR.Y + params.offset, params.bpR.Z}
	innerFrame := append([]m.Vector{ipR, ibpR, ibpL, ipL}, innerArch.left...)
	innerFrame = append(innerFrame, innerArch.right...)

	triangles := gen.JoinPoints([][]m.Vector{outerFrame, innerFrame}, params.material)
	front := make([]m.Triangle, len(triangles))
	for i, o := range triangles {
		t := o.(m.Triangle)
		front[i] = t
	}

	ef := gen.ExtrusionFace{
		Front:    front,
		Outer:    [][]m.Vector{outerFrame},
		Inner:    [][]m.Vector{innerFrame},
		Material: params.material,
	}
	return gen.Extrude(ef, m.Vector{0, 0, params.depth})
}
