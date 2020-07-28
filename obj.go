package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GenGeo/gen"
)

// SaveObj builds a string in .obj format
// representing the list of triangles
// NOTE: .obj vertex count is 1-based
// NOTE: .obj line values are whitespace-separated
func SaveObj(o m.Object) string {
	triangles := trianglesFromObject(o)
	vertices := []m.Vector{}
	vertexMap := map[m.Vector]int64{}
	faces := make([]m.Face, len(triangles))
	for i, t := range triangles {
		v0, ok := vertexMap[t.P0]
		if !ok {
			v0 = int64(len(vertexMap)) + 1
			vertexMap[t.P0] = v0
			vertices = append(vertices, t.P0)
		}
		v1, ok := vertexMap[t.P1]
		if !ok {
			v1 = int64(len(vertexMap)) + 1
			vertexMap[t.P1] = v1
			vertices = append(vertices, t.P1)
		}
		v2, ok := vertexMap[t.P2]
		if !ok {
			v2 = int64(len(vertexMap)) + 1
			vertexMap[t.P2] = v2
			vertices = append(vertices, t.P2)
		}
		faces[i] = m.Face{v0, v1, v2}
	}

	s := ""
	for _, v := range vertices {
		s += fmt.Sprintf("v %f %f %f\n", v.X, v.Y, v.Z)
	}
	for _, f := range faces {
		s += fmt.Sprintf("f %d %d %d\n", f.V0, f.V1, f.V2)
	}
	return s
}

func trianglesFromObject(objects ...m.Object) []m.Triangle {
	triangles := []m.Triangle{}
	for _, o := range objects {
		switch t := o.(type) {
		case m.Triangle:
			triangles = append(triangles, t)
		case *m.ComplexObject:
			triangles = append(triangles, trianglesFromObject(t.Objects()...)...)
		case *m.SharedObject:
			trs := trianglesFromObject(t.Object)
			for _, tr := range trs {
				p0 := t.ObjectToWorld.Point(tr.P0)
				p1 := t.ObjectToWorld.Point(tr.P1)
				p2 := t.ObjectToWorld.Point(tr.P2)
				newTr := m.NewTriangle(p0, p1, p2, tr.Material)
				triangles = append(triangles, newTr)
			}
		}
	}
	return triangles
}

// LoadObj assumes filename contains one triangle mesh object
func LoadObj(filename string, mat m.Material) (m.Object, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return loadObj(scanner, mat)
}

func loadObj(scanner *bufio.Scanner, mat m.Material) (m.Object, error) {
	var vertices []m.Vector
	var faces []m.Face
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key, values := fields[0], fields[1:]
		switch key {
		case "#":
			continue
		case "v":
			vertex, err := readVertex(values)
			if err != nil {
				return nil, err
			}
			vertices = append(vertices, vertex)
		case "f":
			face, err := readFace(values, int64(len(vertices)))
			if err != nil {
				return nil, err
			}
			faces = append(faces, face)
		default:
			fmt.Printf("Unexpected line: %s", line)
		}
	}
	return toObject(vertices, faces, mat)
}

func readVertex(coordinates []string) (m.Vector, error) {
	if len(coordinates) != 3 {
		return m.Vector{}, fmt.Errorf("Invalid coordinates: %v", coordinates)
	}
	p1, err := strconv.ParseFloat(coordinates[0], 32)
	if err != nil {
		return m.Vector{}, err
	}
	p2, err := strconv.ParseFloat(coordinates[1], 32)
	if err != nil {
		return m.Vector{}, err
	}
	p3, err := strconv.ParseFloat(coordinates[2], 32)
	if err != nil {
		return m.Vector{}, err
	}
	return m.Vector{float32(p1), float32(p2), float32(p3)}, nil
}

func readFace(indices []string, numVertices int64) (m.Face, error) {
	if len(indices) != 3 {
		return m.Face{}, fmt.Errorf("Invalid indices: %v", indices)
	}
	i1, err := strconv.ParseInt(indices[0], 10, 64)
	if err != nil {
		return m.Face{}, err
	}

	if i1 < 1 || numVertices < i1 {
		return m.Face{}, fmt.Errorf("Invalid index: %d #indices: %d", i1, numVertices)
	}
	i2, err := strconv.ParseInt(indices[1], 10, 64)
	if err != nil {
		return m.Face{}, err
	}
	if i2 < 1 || numVertices < i2 {
		return m.Face{}, fmt.Errorf("Invalid index: %d #indices: %d", i2, numVertices)
	}
	i3, err := strconv.ParseInt(indices[2], 10, 64)
	if err != nil {
		return m.Face{}, err
	}
	if i3 < 1 || numVertices < i3 {
		return m.Face{}, fmt.Errorf("Invalid index: %d #indices: %d", i3, numVertices)
	}
	return m.NewFace(i1-1, i2-1, i3-1), nil
}

func toObject(vertices []m.Vector, faces []m.Face, mat m.Material) (m.Object, error) {
	if len(faces) == 0 {
		return nil, errors.New("Object list empty")
	}
	vertices = gen.CenterPointsOnOrigin(vertices)
	return m.NewTriangleMesh(vertices, faces, mat), nil
}
