package pair

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Pair of two elements of different/same types.
type Pair[F, S any] struct {
	f F
	s S
}

type jsonValue[F, S any] struct {
	Key   F `json:"key"`
	Value S `json:"value"`
}

// Of returns a new pair with the given elements.
func Of[F, S any](f F, s S) Pair[F, S] {
	return Pair[F, S]{f: f, s: s}
}

// First returns the first element of the pair.
func (p Pair[F, S]) First() F { return p.f }

// Second returns the second element of the pair.
func (p Pair[F, S]) Second() S { return p.s }

// String returns a string representation of the pair.
func (p Pair[F, S]) String() string {
	return fmt.Sprintf("(%#v, %#v)", p.f, p.s)
}

// MarshalJSON implements the json.Marshaler.
func (p Pair[F, S]) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonValue[F, S]{Key: p.f, Value: p.s})
}

// UnmarshalJSON implements the json.Unmarshaler.
func (p *Pair[F, S]) UnmarshalJSON(data []byte) error {
	var jv jsonValue[F, S]
	if err := json.Unmarshal(data, &jv); err != nil {
		return err
	}
	p.f = jv.Key
	p.s = jv.Value
	return nil
}

// Value implements the driver.Valuer.
func (p Pair[F, S]) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner.
func (p *Pair[F, S]) Scan(src any) error {
	if src == nil {
		return nil
	}

	var b []byte
	switch t := src.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		return fmt.Errorf("pair: unexpected source data type: %T", src)
	}

	return json.Unmarshal(b, p)
}
