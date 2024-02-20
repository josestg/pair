package pair

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func ExampleOf() {
	p1 := Of("first", "second")
	fmt.Println(p1)
	fmt.Println(p1.First())
	fmt.Println(p1.Second())

	p2 := Of(42, errors.New("something went wrong"))
	fmt.Println(p2)

	// Output:
	// ("first", "second")
	// first
	// second
	// (42, &errors.errorString{s:"something went wrong"})
}

func ExamplePair_MarshalJSON() {
	p := Of(42, "hello")
	b, err := p.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	// Output:
	// {"key":42,"value":"hello"}
}

func ExamplePair_UnmarshalJSON() {
	raw := []byte(`{"key":42,"value":"hello"}`)
	var p Pair[int, string]
	if err := json.Unmarshal(raw, &p); err != nil {
		panic(err)
	}

	fmt.Println(p)
	// Output:
	// (42, "hello")
}

func TestOf(t *testing.T) {
	p := Of(1, "hello")
	if p.First() != 1 {
		t.Errorf("expected 1, got %#v", p.First())
	}
	if p.Second() != "hello" {
		t.Errorf("expected hello, got %#v", p.Second())
	}
}

func TestString(t *testing.T) {
	p := Of(1, "hello")
	e := "(1, \"hello\")"
	if p.String() != e {
		t.Errorf("expected %q, got %q", e, p.String())
	}
}

func TestMarshalJSON(t *testing.T) {
	p := Of(1, "hello")
	e := `{"key":1,"value":"hello"}`

	t.Run("direct", func(t *testing.T) {
		b, err := p.MarshalJSON()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(b) != e {
			t.Errorf("expected %q, got %q", e, string(b))
		}
	})

	t.Run("un-direct", func(t *testing.T) {
		b, err := json.Marshal(p)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(b) != e {
			t.Errorf("expected %q, got %q", e, string(b))
		}
	})

	t.Run("pointer", func(t *testing.T) {
		b, err := json.Marshal(&p)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(b) != e {
			t.Errorf("expected %q, got %q", e, string(b))
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	e := Of(1, "hello")
	b := []byte(`{"key":1,"value":"hello"}`)

	t.Run("direct", func(t *testing.T) {
		var p Pair[int, string]
		if err := p.UnmarshalJSON(b); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if p != e {
			t.Errorf("expected %#v, got %#v", e, p)
		}
	})

	t.Run("un-direct", func(t *testing.T) {
		var p Pair[int, string]
		if err := json.Unmarshal(b, &p); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var p Pair[int, string]
		if err := p.UnmarshalJSON([]byte(`{"key":1,"value":2}`)); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestValue(t *testing.T) {
	p := Of(1, "hello")
	e := `{"key":1,"value":"hello"}`
	b, err := p.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(b.([]byte)) != e {
		t.Errorf("expected %q, got %q", e, string(b.([]byte)))
	}
}

func TestScan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var p Pair[int, string]
		if err := p.Scan(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var q Pair[int, string]
		if p != q {
			t.Errorf("expected %#v, got %#v", q, p)
		}
	})

	t.Run("unexpected type", func(t *testing.T) {
		var p Pair[int, string]
		if err := p.Scan(123); err == nil {
			t.Errorf("expected error, got nil")
		}

		var q Pair[int, string]
		if p != q {
			t.Errorf("expected %#v, got %#v", q, p)
		}
	})

	t.Run("from string", func(t *testing.T) {
		var p Pair[int, string]
		if err := p.Scan(`{"key":1,"value":"hello"}`); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if p.First() != 1 {
			t.Errorf("expected 1, got %#v", p.First())
		}

		if p.Second() != "hello" {
			t.Errorf("expected hello, got %#v", p.Second())
		}
	})

	t.Run("from []byte", func(t *testing.T) {
		var p Pair[int, string]
		if err := p.Scan([]byte(`{"key":1,"value":"hello"}`)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if p.First() != 1 {
			t.Errorf("expected 1, got %#v", p.First())
		}

		if p.Second() != "hello" {
			t.Errorf("expected hello, got %#v", p.Second())
		}
	})
}
