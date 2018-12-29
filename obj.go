package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	m "github.com/deosjr/GRayT/src/model"
)

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
	var triangles []m.Object
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
			face, err := readTriangle(values, vertices, mat)
			if err != nil {
				return nil, err
			}
			triangles = append(triangles, face)
		default:
			fmt.Printf("Unexpected line: %s", line)
		}
	}
	return toObject(triangles)
}

func readVertex(coordinates []string) (m.Vector, error) {
	if len(coordinates) != 3 {
		return m.Vector{}, fmt.Errorf("Invalid coordinates: %v", coordinates)
	}
	p1, err := strconv.ParseFloat(coordinates[0], 64)
	if err != nil {
		return m.Vector{}, err
	}
	p2, err := strconv.ParseFloat(coordinates[1], 64)
	if err != nil {
		return m.Vector{}, err
	}
	p3, err := strconv.ParseFloat(coordinates[2], 64)
	if err != nil {
		return m.Vector{}, err
	}
	return m.Vector{p1, p2, p3}, nil
}

func readTriangle(indices []string, vertices []m.Vector, mat m.Material) (m.Triangle, error) {
	if len(indices) != 3 {
		return m.Triangle{}, fmt.Errorf("Invalid indices: %v", indices)
	}
	i1, err := strconv.ParseInt(indices[0], 10, 64)
	if err != nil {
		return m.Triangle{}, err
	}

	numVertices := int64(len(vertices))
	if i1 < 1 || numVertices < i1 {
		return m.Triangle{}, fmt.Errorf("Invalid index: %d #indices: %d", i1, numVertices)
	}
	i2, err := strconv.ParseInt(indices[1], 10, 64)
	if err != nil {
		return m.Triangle{}, err
	}
	if i2 < 1 || numVertices < i2 {
		return m.Triangle{}, fmt.Errorf("Invalid index: %d #indices: %d", i2, numVertices)
	}
	i3, err := strconv.ParseInt(indices[2], 10, 64)
	if err != nil {
		return m.Triangle{}, err
	}
	if i3 < 1 || numVertices < i3 {
		return m.Triangle{}, fmt.Errorf("Invalid index: %d #indices: %d", i3, numVertices)
	}
	// TODO: coordinate handedness!
	return m.NewTriangle(vertices[i3-1], vertices[i2-1], vertices[i1-1], mat), nil
}

func toObject(triangles []m.Object) (m.Object, error) {
	if len(triangles) == 0 {
		return nil, errors.New("Object list empty")
	}
	centerTrianglesOnOrigin(triangles)
	return m.NewComplexObject(triangles), nil
}

func centerTrianglesOnOrigin(triangles []m.Object) {
	b := m.ObjectsBound(triangles, m.ScaleUniform(1.0))
	objectToOrigin := m.Translate(b.Centroid()).Inverse()

	for i, tobj := range triangles {
		t := tobj.(m.Triangle)
		triangles[i] = m.NewTriangle(
			objectToOrigin.Point(t.P0),
			objectToOrigin.Point(t.P1),
			objectToOrigin.Point(t.P2),
			t.Material)
	}
}
