package libsqlvectorgorm

import (
	"context"
	"gorm.io/gorm/clause"
	"testing"
)

func TestVectorGormValue(t *testing.T) {
	testCases := []struct {
		name     string
		vector   Vector
		expected clause.Expr
	}{
		{
			name:   "Empty vector",
			vector: NewVector([]float32{}),
			expected: clause.Expr{
				SQL:  "vector(?)",
				Vars: []interface{}{"[]"},
			},
		},
		{
			name:   "Vector with single element",
			vector: NewVector([]float32{1.5}),
			expected: clause.Expr{
				SQL:  "vector(?)",
				Vars: []interface{}{"[1.5]"},
			},
		},
		{
			name:   "Vector with multiple elements",
			vector: NewVector([]float32{1.5, -2.75, 3.0}),
			expected: clause.Expr{
				SQL:  "vector(?)",
				Vars: []interface{}{"[1.5,-2.75,3]"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.vector.GormValue(context.Background(), nil)

			if result.SQL != tc.expected.SQL {
				t.Errorf("Expected SQL %q, but got %q", tc.expected.SQL, result.SQL)
			}

			if len(result.Vars) != len(tc.expected.Vars) {
				t.Errorf("Expected %d vars, but got %d", len(tc.expected.Vars), len(result.Vars))
			} else if result.Vars[0] != tc.expected.Vars[0] {
				t.Errorf("Expected var %v, but got %v", tc.expected.Vars[0], result.Vars[0])
			}
		})
	}
}
