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
// - all points in the left arc
// - all points in the right arc
type arch struct {
	left  []m.Vector
	right []m.Vector
}

// createArch returns an arch with endpoints pL, pR and midpoints to the arcs mL and mR.
// mL and mR are parameters instead of calculated from excess here, so that this function
// can be used for calculating offset for an arch too
func createArch(pL, pR, mL, mR m.Vector, numPoints int) arch {
	// generate left and right arc up until their intersection point (top)
	// the top of the arch is given by translating the middle of line pLpR up with y
	// where y is sqrt(pLpR * (r-(pLpR/4))), where r is the circle radius
	// arc angle alpha: tan alpha = ||top,pM|| / || pM,mL||

	// assumption: length of mR-pR given as equal
	r := m.VectorFromTo(mL, pL).Length()

	pLpR := m.VectorFromTo(pL, pR).Length()
	pM := pL.Add(m.VectorFromTo(pL, pR).Times(0.5))
	y := math.Sqrt(pLpR * (r - (pLpR / 4.0)))
	top := pM.Add(m.Vector{0, y, 0})
	alpha := math.Atan2(m.VectorFromTo(top, pM).Length(), m.VectorFromTo(pM, mL).Length())

	c := gen.NewCircle(mL, r)
	upperLeftArc := c.PointsPhaseRange(math.Pi, math.Pi-alpha, numPoints, m.Vector{1, 0, 0}, m.Vector{0, 1, 0})
	c = gen.NewCircle(mR, r)
	upperRightArc := c.PointsPhaseRange(alpha, 0, numPoints, m.Vector{1, 0, 0}, m.Vector{0, 1, 0})

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
	// number of points on a quarter circle arc
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

type emptyArchWindowTraceryParams struct {
	material m.Material
	// ratio r / distance(pL, pR), excess >= 0.5
	excess float64
	// width of tracery, extends inwards
	offset float64
	// depth of extrusion
	depth float64
	// pL and pR, left and right endpoints of the arch
	pL m.Vector
	pR m.Vector
	// bpL and bpR, left and right bottom points of the window
	bpL m.Vector
	bpR m.Vector
	// number of points on a quarter circle arc
	numPoints int
}

// emptyArchWindowTracery returns window tracery object which is a simple outline
func emptyArchWindowTracery(params emptyArchWindowTraceryParams) m.Object {
	mL := params.pL.Add(m.VectorFromTo(params.pL, params.pR).Times(params.excess))
	mR := params.pR.Add(m.VectorFromTo(params.pR, params.pL).Times(params.excess))

	outerArch := createArch(params.pL, params.pR, mL, mR, params.numPoints)
	outerFrame := append([]m.Vector{params.pR, params.bpR, params.bpL, params.pL}, outerArch.left...)
	outerFrame = append(outerFrame, outerArch.right...)

	ipL := m.Vector{params.pL.X + params.offset, params.pL.Y, params.pL.Z}
	ipR := m.Vector{params.pR.X - params.offset, params.pR.Y, params.pR.Z}
	innerArch := createArch(ipL, ipR, mL, mR, params.numPoints)

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

type archWindowTraceryParams struct {
	material m.Material
	// ratio r / distance(pL, pR), excess >= 0.5
	excess float64
	// outerWidth is the width of outermost arch
	// innerWidth is the width of the inner tracery
	outerWidth float64
	innerWidth float64
	// offset between main and subarch points
	verticalOffset float64
	// depth of extrusion
	depth float64
	// pL and pR, left and right endpoints of the arch
	pL m.Vector
	pR m.Vector
	// bpL and bpR, left and right bottom points of the window
	bpL m.Vector
	bpR m.Vector
	// number of points on a quarter circle arc
	numPoints int
}

// archWindowTracery returns window tracery object
func archWindowTracery(params archWindowTraceryParams) m.Object {
	eparams := emptyArchWindowTraceryParams{
		material:  params.material,
		excess:    params.excess,
		offset:    params.outerWidth,
		depth:     params.depth,
		pL:        params.pL,
		pR:        params.pR,
		bpL:       params.bpL,
		bpR:       params.bpR,
		numPoints: params.numPoints,
	}
	mainArch := emptyArchWindowTracery(eparams)

	eparams.offset = params.innerWidth
	innerOffset := m.Vector{params.outerWidth - params.innerWidth, 0, 0}
	innerWidth := m.Vector{params.innerWidth, 0, 0}
	verticalOffset := m.Vector{0, -params.verticalOffset, 0}
	pM := params.pL.Add(m.VectorFromTo(params.pL, params.pR).Times(0.5)).Add(verticalOffset)
	bpM := params.bpL.Add(m.VectorFromTo(params.bpL, params.bpR).Times(0.5))

	// left sub arch
	eparams.pL = params.pL.Add(verticalOffset).Add(innerOffset)
	eparams.bpL = params.bpL.Add(innerOffset)
	eparams.pR = pM.Add(innerWidth.Times(0.5))
	eparams.bpR = bpM.Add(innerWidth.Times(0.5))
	leftArch := emptyArchWindowTracery(eparams)

	// rosette
	mR := params.pR.Add(m.VectorFromTo(params.pR, params.pL).Times(params.excess))
	rR := m.VectorFromTo(mR, params.pR).Length() - params.outerWidth
	mLR := eparams.pR.Add(m.VectorFromTo(eparams.pR, eparams.pL).Times(params.excess))
	rLR := m.VectorFromTo(mLR, eparams.pR).Length()

	// intersection of axis of symmetry and ellipse (mR,mLR,rR+rLR)
	center := mR.Add(m.VectorFromTo(mR, mLR).Times(0.5))
	a := (rR + rLR) / 2.0
	c := m.VectorFromTo(mR, mLR).Length() * 0.5
	b := math.Sqrt(a*a - c*c)
	y := (a/b)*math.Sqrt(b*b-math.Pow(pM.X-center.X, 2)) + center.Y
	mC := m.Vector{pM.X, y, 0}

	r := rR - m.VectorFromTo(mC, mR).Length()
	innerCircle := gen.NewCircle(mC, r)
	ip := innerCircle.PointsPhaseRange(0, 2*math.Pi, params.numPoints, ex, ey)
	outerCircle := gen.NewCircle(mC, r+params.innerWidth)
	op := outerCircle.PointsPhaseRange(0, 2*math.Pi, params.numPoints, ex, ey)
	ipRev := make([]m.Vector, len(ip))
	for i, v := range ip {
		ipRev[len(op)-1-i] = v
	}

	triangles := gen.JoinPoints([][]m.Vector{ip, op}, params.material)
	rface := make([]m.Triangle, len(triangles))
	for i, o := range triangles {
		rface[len(triangles)-1-i] = o.(m.Triangle)
	}
	ef := gen.ExtrusionFace{
		Front:    rface,
		Outer:    [][]m.Vector{op},
		Inner:    [][]m.Vector{ipRev},
		Material: params.material,
	}
	rosette := gen.Extrude(ef, m.Vector{0, 0, params.depth})

	// right sub arch
	eparams.pL = pM.Sub(innerWidth.Times(0.5))
	eparams.bpL = bpM.Sub(innerWidth.Times(0.5))
	eparams.pR = params.pR.Add(verticalOffset).Sub(innerOffset)
	eparams.bpR = params.bpR.Sub(innerOffset)
	rightArch := emptyArchWindowTracery(eparams)

	return m.NewComplexObject([]m.Object{mainArch, leftArch, rightArch, rosette})
}
