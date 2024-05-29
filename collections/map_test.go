package collections_test

import (
	"testing"

	zerocollections "github.com/zerogo-hub/zero-helper/collections"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		m    map[int]string
		want int
	}{
		{
			name: "Non-empty map",
			m:    map[int]string{1: "a", 2: "b", 3: "c"},
			want: 3,
		},
		{
			name: "Empty map",
			m:    map[int]string{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zerocollections.Keys(tt.m); len(got) != tt.want {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeysNonEmptyMap(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	got := zerocollections.Keys(m)
	if len(m) != 3 {
		t.Errorf("Unexcept got: %v", got)
	}
}

func TestKeysEmptyMap(t *testing.T) {
	m := map[int]string{}
	got := zerocollections.Keys(m)
	if got != nil {
		t.Errorf("Unexcept got: %v", got)
	}
}

func TestValues(t *testing.T) {
	// Test case 1: Empty map
	t.Run("EmptyMap", func(t *testing.T) {
		m := make(map[int]string)
		result := zerocollections.Values(m)
		if len(result) != 0 {
			t.Errorf("Expected an empty slice, but got %v", result)
		}
	})

	// Test case 2: Map with string values
	t.Run("StringValues", func(t *testing.T) {
		m := map[int]string{1: "apple", 2: "banana", 3: "cherry"}
		expected := []string{"apple", "banana", "cherry"}
		result := zerocollections.Values(m)
		if len(result) != len(expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// Test case 3: Map with int values
	t.Run("IntValues", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		expected := []int{1, 2, 3}
		result := zerocollections.Values(m)
		if len(result) != len(expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})
}
