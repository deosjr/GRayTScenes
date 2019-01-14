package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	m "github.com/deosjr/GRayT/src/model"
)

func TestLoadObj(t *testing.T) {
	for i, tt := range []struct {
		obj  string
		want m.ComplexObject
	}{
		{
			obj:  `# empty file`,
			want: nil,
		},
		{
			obj: `# note reversed vertex order
			v 1.0 -0.02 2.1754370e-002
			v 2 3 4
			v 4 5 6.0
			f 1 2 3`,
			want: []m.Object{
				m.NewTriangle(
					m.Vector{4, 5, 6},
					m.Vector{2, 3, 4},
					m.Vector{1.0, -0.02, 2.1754370e-002},
					&m.DiffuseMaterial{}),
			},
		},
	} {
		reader := strings.NewReader(tt.obj)
		scanner := bufio.NewScanner(reader)
		got, err := loadObj(scanner, &m.DiffuseMaterial{})
		if err != nil {
			if tt.want == nil {
				continue
			}
			t.Errorf("%d): error in load: %s", i, err.Error())
			continue
		}
		if tt.want == nil && got != nil {
			t.Errorf("%d): expected nil, got: %s", i, got)
			continue
		}
		want := m.NewComplexObject(tt.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%d): got %v want %v", i, got, want)
		}
	}
}
