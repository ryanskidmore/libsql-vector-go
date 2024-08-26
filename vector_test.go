package libsqlvector

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestNewVector(t *testing.T) {
	vec := []float32{1.1, 2.2, 3.3}
	vector := NewVector(vec)
	if !reflect.DeepEqual(vector.vec, vec) {
		t.Errorf("NewVector(%v) = %v, want %v", vec, vector.vec, vec)
	}
}

func TestVectorSlice(t *testing.T) {
	vec := []float32{1.1, 2.2, 3.3}
	v := Vector{vec: vec}
	if !reflect.DeepEqual(v.Slice(), vec) {
		t.Errorf("Vector.Slice() = %v, want %v", v.Slice(), vec)
	}
}

func TestVectorFormatFloats(t *testing.T) {
	tests := []struct {
		name     string
		vec      []float32
		expected string
	}{
		{
			name:     "Empty vector",
			vec:      []float32{},
			expected: "[]",
		},
		{
			name:     "Single element vector",
			vec:      []float32{1.5},
			expected: "[1.5]",
		},
		{
			name:     "Multiple element vector",
			vec:      []float32{1.5, 2.75, -3.25, 0},
			expected: "[1.5,2.75,-3.25,0]",
		},
		{
			name:     "Vector with very small number",
			vec:      []float32{1e-7},
			expected: "[0.0000001]",
		},
		{
			name:     "Vector with very large number",
			vec:      []float32{1e7},
			expected: "[10000000]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vec := Vector{vec: tt.vec}
			result := vec.FormatFloats()
			if result != tt.expected {
				t.Errorf("formatFloats() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVectorString(t *testing.T) {
	tests := []struct {
		name string
		vec  []float32
		want string
	}{
		{"Empty", []float32{}, "vector('[]')"},
		{"Single", []float32{1.0}, "vector('[1]')"},
		{"Multiple", []float32{1.0, 2.0, 3.0}, "vector('[1,2,3]')"},
		{"Decimal", []float32{1.1, 2.2, 3.3}, "vector('[1.1,2.2,3.3]')"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{vec: tt.vec}
			if got := v.String(); got != tt.want {
				t.Errorf("Vector.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []float32
		wantErr bool
	}{
		{"Empty", "vector('[]')", []float32{}, false},
		{"Single", "vector('[1]')", []float32{1.0}, false},
		{"Multiple", "vector('[1,2,3]')", []float32{1.0, 2.0, 3.0}, false},
		{"Decimal", "vector('[1.1,2.2,3.3]')", []float32{1.1, 2.2, 3.3}, false},
		{"Invalid", "vector('[1,2,a]')", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{}
			err := v.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(v.vec, tt.want) {
				t.Errorf("Vector.Parse() = %v, want %v", v.vec, tt.want)
			}
		})
	}
}

func TestVectorEncodeBinary(t *testing.T) {
	tests := []struct {
		name    string
		vec     []float32
		wantLen int
	}{
		{"Empty", []float32{}, 4},
		{"Single", []float32{1.0}, 8},
		{"Multiple", []float32{1.0, 2.0, 3.0}, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{vec: tt.vec}
			buf, err := v.EncodeBinary(nil)
			if err != nil {
				t.Errorf("Vector.EncodeBinary() error = %v", err)
				return
			}
			if len(buf) != tt.wantLen {
				t.Errorf("Vector.EncodeBinary() len = %v, want %v", len(buf), tt.wantLen)
			}
		})
	}
}

func TestVectorDecodeBinary(t *testing.T) {
	tests := []struct {
		name    string
		vec     []float32
		wantErr bool
	}{
		{"Empty", []float32{}, false},
		{"Single", []float32{1.0}, false},
		{"Multiple", []float32{1.0, 2.0, 3.0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := Vector{vec: tt.vec}
			buf, _ := v1.EncodeBinary(nil)

			v2 := &Vector{}
			err := v2.DecodeBinary(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector.DecodeBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(v1.vec, v2.vec) {
				t.Errorf("Vector.DecodeBinary() = %v, want %v", v2.vec, v1.vec)
			}
		})
	}
}

func TestVectorScan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []float32
		wantErr bool
	}{
		{"String", "vector('[1,2,3]')", []float32{1.0, 2.0, 3.0}, false},
		{"Bytes", []byte("vector('[1,2,3]')"), []float32{1.0, 2.0, 3.0}, false},
		{"Invalid", 123, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{}
			err := v.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(v.vec, tt.want) {
				t.Errorf("Vector.Scan() = %v, want %v", v.vec, tt.want)
			}
		})
	}
}

func TestVectorValue(t *testing.T) {
	tests := []struct {
		name    string
		vec     []float32
		want    driver.Value
		wantErr bool
	}{
		{"Empty", []float32{}, "vector('[]')", false},
		{"Single", []float32{1.0}, "vector('[1]')", false},
		{"Multiple", []float32{1.0, 2.0, 3.0}, "vector('[1,2,3]')", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{vec: tt.vec}
			got, err := v.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Vector.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		vec     []float32
		want    string
		wantErr bool
	}{
		{"Empty", []float32{}, "[]", false},
		{"Single", []float32{1.0}, "[1]", false},
		{"Multiple", []float32{1.0, 2.0, 3.0}, "[1,2,3]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{vec: tt.vec}
			got, err := v.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("Vector.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestVectorUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []float32
		wantErr bool
	}{
		{"Empty", "[]", []float32{}, false},
		{"Single", "[1]", []float32{1.0}, false},
		{"Multiple", "[1,2,3]", []float32{1.0, 2.0, 3.0}, false},
		{"Invalid", "[1,2,a]", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{}
			err := v.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Vector.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(v.vec, tt.want) {
				t.Errorf("Vector.UnmarshalJSON() = %v, want %v", v.vec, tt.want)
			}
		})
	}
}
