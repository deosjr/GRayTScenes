package main

import (
	"math"
	"math/rand"
	"time"

	perlin "github.com/aquilax/go-perlin"
	"github.com/fogleman/poissondisc"

	"github.com/deosjr/GRayT/src/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func perlinHeightMap(grid [][]model.Vector, n int, weights []float64, pow float64) [][]model.Vector {
	xSize, ySize := len(grid), len(grid[0])
	// alpha, beta, n iterations, random seed
	p := perlin.NewPerlin(2, 2, 3, rand.Int63())
	for y, row := range grid {
		for x, _ := range row {
			nx := float64(x)/float64(xSize) - 0.5
			ny := float64(y)/float64(ySize) - 0.5
			noise := weights[0] * p.Noise2D(nx, ny)
			sum := weights[0]
			for i := 0; i < n; i++ {
				exp := math.Pow(2, float64(i+1))
				noise += weights[i+1] * p.Noise2D(exp*nx, exp*ny)
				sum += weights[i+1]
			}
			// normalize
			noise = noise / sum
			// map from [-1,1] to [0,1]
			noise = (noise + 1) / 2
			noise = math.Pow(noise, pow)
			grid[y][x].Y = float32(noise)
		}
	}
	return grid
}

// assumption: r is a rectangle
func toPointGrid(r model.Quadrilateral, roughSize float64) [][]model.Vector {
	xlen := float64(model.VectorFromTo(r.P1, r.P2).Length())
	ylen := float64(model.VectorFromTo(r.P1, r.P4).Length())
	numDivisionsX := math.Ceil(xlen / roughSize)
	numDivisionsY := math.Ceil(ylen / roughSize)
	pointSizeX := xlen / numDivisionsX
	pointSizeY := ylen / numDivisionsY
	xVector := model.VectorFromTo(r.P1, r.P2).Normalize().Times(float32(pointSizeX))
	yVector := model.VectorFromTo(r.P1, r.P4).Normalize().Times(float32(pointSizeY))

	numPointsX := int(numDivisionsX) + 1
	numPointsY := int(numDivisionsY) + 1

	grid := make([][]model.Vector, numPointsY)
	for y := 0; y < numPointsY; y++ {
		row := make([]model.Vector, numPointsX)
		for x := 0; x < numPointsX; x++ {
			row[x] = r.P1.Add(xVector.Times(float32(x))).Add(yVector.Times(float32(y)))
		}
		grid[y] = row
	}

	return grid
}

func gridToTriangles(grid [][]model.Vector, mat model.Material) model.Object {
	ylen := len(grid)
	xlen := len(grid[0])
	triangles := []model.Object{}
	for y := 0; y < ylen-1; y++ {
		for x := 0; x < xlen-1; x++ {
			p1 := grid[y][x]
			p2 := grid[y][x+1]
			p3 := grid[y+1][x+1]
			p4 := grid[y+1][x]
			t1 := model.NewTriangle(p1, p4, p2, mat)
			t2 := model.NewTriangle(p4, p3, p2, mat)
			triangles = append(triangles, t1, t2)
		}
	}
	return model.NewComplexObject(triangles)
}

// assumption: q is perpendicular to z
func poisson(q model.Quadrilateral, r float64) []model.Vector {
	x0, y0 := float64(q.P1.X), float64(q.P1.Y)
	x1, y1 := float64(q.P3.X), float64(q.P3.Y)
	k := 30
	points := poissondisc.Sample(x0, y0, x1, y1, r, k, nil)
	vectors := make([]model.Vector, len(points))
	for i, p := range points {
		vectors[i] = model.Vector{float32(p.X), float32(p.Y), 0.0}
	}
	return vectors
}
