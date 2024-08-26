package libsqlvectorgorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	libsqlvector "github.com/ryanskidmore/libsql-vector-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Vector struct {
	v *libsqlvector.Vector
}

func (v *Vector) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "vector(?)",
		Vars: []interface{}{v.v.FormatFloats()},
	}
}

func NewVector(vec []float32) Vector {
	v := libsqlvector.NewVector(vec)
	return Vector{
		v: &v,
	}
}

func (v Vector) Slice() []float32 {
	return v.v.Slice()
}

func (v Vector) FormatFloats() string {
	return ""
}

// String returns a string representation of the vector
func (v Vector) String() string {
	return v.v.String()
}

func (v *Vector) Parse(s string) error {
	return v.v.Parse(s)
}

// EncodeBinary encodes a binary representation of the vector.
func (v Vector) EncodeBinary(buf []byte) (newBuf []byte, err error) {
	return v.v.EncodeBinary(buf)
}

// DecodeBinary decodes a binary representation of a vector.
func (v *Vector) DecodeBinary(buf []byte) error {
	return v.v.DecodeBinary(buf)
}

// statically assert that Vector implements sql.Scanner.
var _ sql.Scanner = (*Vector)(nil)

// Scan implements the sql.Scanner interface.
func (v *Vector) Scan(src interface{}) (err error) {
	return v.v.Scan(src)
}

// statically assert that Vector implements driver.Valuer.
var _ driver.Valuer = (*Vector)(nil)

// Value implements the driver.Valuer interface.
func (v Vector) Value() (driver.Value, error) {
	return v.v.String(), nil
}

// statically assert that Vector implements json.Marshaler.
var _ json.Marshaler = (*Vector)(nil)

// MarshalJSON implements the json.Marshaler interface.
func (v Vector) MarshalJSON() ([]byte, error) {
	return v.v.MarshalJSON()
}

// statically assert that Vector implements json.Unmarshaler.
var _ json.Unmarshaler = (*Vector)(nil)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Vector) UnmarshalJSON(data []byte) error {
	return v.v.UnmarshalJSON(data)
}
