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
		want []m.Triangle
	}{
		{
			obj:  `# empty file`,
			want: nil,
		},
		{
			obj: `# comment 
			v 1.0 -0.02 2.1754370e-002
			v 2 3 4
			v 4 5 6.0
			f 1 2 3`,
			want: []m.Triangle{
				m.NewTriangle(
					m.Vector{1.0, -0.02, 2.1754370e-002},
					m.Vector{2, 3, 4},
					m.Vector{4, 5, 6},
					&m.DiffuseMaterial{}),
				},
			},
	} {
		reader := strings.NewReader(tt.obj)
		scanner := bufio.NewScanner(reader)
		obj, err := loadObj(scanner, &m.DiffuseMaterial{})
		if err != nil {
			if tt.want == nil {
				continue
			}
			t.Errorf("%d): error in load: %s", i, err.Error())
			continue
		}
		got := obj.(*m.ComplexObject).Objects()
		if tt.want == nil && got != nil {
			t.Errorf("%d): expected nil, got: %s", i, got)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}
