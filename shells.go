package main

import (
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GenGeo/gen"
)

// D. M. Raup - Geometric Analysis of Shell Coiling: General Problems (1966)
// variable naming after R. Dawkins - Climbing Mount Improbable (1996)

// For extra inspiration maybe check
// Bernard Tursch - Spiral Growth: The 'Museum of All Shells' Revisited (1996)

// TODO: spire=0 gives stack overflow errors in bvh.go... ?
// TODO: param definitions between different papers seem all over the place.

// flare: Raup's W
// If w1 is the distance between a point P1 on the generating spiral and that spiral's axis,
// and w2 the distance from the corresponding point P2 one turn later, then w2/w1=W, which is
// strictly greater than 1
// verm:  Raup's D
// If d1 is the distance between the spiral axis and the inner edge of the shell cavity,
// and d2 the distance between the spiral axis and the outer edge on the same turn, then
// d1/d2 = D, which is greater than or equal to 0 and strictly less than 1
// spire: Raup's T
// If the acute angle between the line P1P2 and the horizontal is alpha, then tan alpha = T.
func generateShell(flare, verm, spire float64, numWindings int) m.Object {
	aFunc := func(t float64) float64 {
		return math.Pow(flare, (t/(2*math.Pi))) - 1
	}
	// pitch is 2pi*b, meaning a full rotation gains 2pi*b in height
	// T = tan alpha = b / delta a
	// b = T * delta a
	// for a full rotation to change exactly b in height, we correct by 2pi
	bFunc := func(t float64) float64 {
		aNow := aFunc(t)
		aPrev := 0.0
		if t > 2*math.Pi {
			aPrev = aFunc(t - 2*math.Pi)
		}
		aDiff := math.Abs(aNow - aPrev)
		return (spire * aDiff) / (2 * math.Pi)
	}
	helix := gen.NewHelix(aFunc, bFunc)

	// |           a             |
	// C------------------|------P------|
	// |       a - r      |   r  |   r  |
	// |            a + r               |

	// C is center of spiral, P point on helix
	// We have radius a and want to know r
	// verm (D) is ratio of (a-r) to (a+r), so
	// (a - r) / (a + r) = D
	// a - r = (a + r) * D = Da + Dr
	// a - Da = Dr + r = r(D + 1)
	// r = (a - Da) / (D + 1) = -(D - 1) * a / (D + 1)

	generatingCurve := gen.NewCircle(func(t float64) float64 {
		return (-verm + 1) * aFunc(t) / (verm + 1)
	}, 100)
	numSteps := 64 * numWindings
	stepSize := math.Pi / 32.0
	mat := &m.DiffuseMaterial{m.NewColor(200, 100, 0)}

	po := gen.NewParametricObject(helix, generatingCurve, numSteps, stepSize, mat)
	return po.Build()
}
