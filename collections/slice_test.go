package collections_test

import (
	"fmt"
	"reflect"
	"testing"

	zerocollections "github.com/zerogo-hub/zero-helper/collections"
)

func TestContains_WhenTargetExists_ReturnsTrue(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	target := 3
	result := zerocollections.Contains(data, target)
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestContains_WhenTargetDoesNotExist_ReturnsFalse(t *testing.T) {
	data := []string{"apple", "banana", "orange"}
	target := "grape"
	result := zerocollections.Contains(data, target)
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestContains_WhenEmptySlice_ReturnsFalse(t *testing.T) {
	var data []float64
	target := 5.5
	result := zerocollections.Contains(data, target)
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestContains_WhenTargetIsFirstElement_ReturnsTrue(t *testing.T) {
	data := []int{10, 20, 30, 40, 50}
	target := 10
	result := zerocollections.Contains(data, target)
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestContains_WhenTargetIsLastElement_ReturnsTrue(t *testing.T) {
	data := []int{10, 20, 30, 40, 50}
	target := 50
	result := zerocollections.Contains(data, target)
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		l    []string
		sep  string
		want string
	}{
		{"JoinStrings", []string{"a", "b", "c"}, ",", "a,b,c"},
		{"JoinInts", []string{"1", "2", "3"}, "-", "1-2-3"},
		{"JoinEmpty", []string{}, ",", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zerocollections.Join(tt.l, tt.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSlice_index(b *testing.B) {
	n := 1000000
	l := make([]int, n)
	for i := 0; i < n; i++ {
		l = append(l, i)
	}
	b.ResetTimer()

	l2 := make([]string, len(l))
	for i, t := range l {
		l2[i] = fmt.Sprintf("%v", t)
	}
	fmt.Println(len(l2))
}

func BenchmarkSlice_append(b *testing.B) {
	n := 1000000
	l := make([]int, n)
	for i := 0; i < n; i++ {
		l = append(l, i)
	}
	b.ResetTimer()

	l2 := make([]string, 0, len(l))
	for _, t := range l {
		l2 = append(l2, fmt.Sprintf("%v", t))
	}
	fmt.Println(len(l2))
}

func TestSum_Integer(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := 15
	result := zerocollections.Sum(input)
	if result != expected {
		t.Errorf("Sum of %v = %d; want %d", input, result, expected)
	}
}

func TestSum_Float(t *testing.T) {
	input := []float64{1.5, 2.5, 3.5, 4.5, 5.5}
	expected := 17.5
	result := zerocollections.Sum(input)
	if result != expected {
		t.Errorf("Sum of %v = %f; want %f", input, result, expected)
	}
}

func TestSum_EmptySlice(t *testing.T) {
	input := []int{}
	expected := 0
	result := zerocollections.Sum(input)
	if result != expected {
		t.Errorf("Sum of %v = %d; want %d", input, result, expected)
	}
}

func TestSum_NegativeNumbers(t *testing.T) {
	input := []int{-1, -2, -3, -4, -5}
	expected := -15
	result := zerocollections.Sum(input)
	if result != expected {
		t.Errorf("Sum of %v = %d; want %d", input, result, expected)
	}
}

func TestUnique(t *testing.T) {
	// Test case 1: Empty input slice
	input1 := []int{}
	expected1 := []int{}
	if output1 := zerocollections.Unique(input1); !reflect.DeepEqual(output1, expected1) {
		t.Errorf("Unique(%v) = %v, want %v", input1, output1, expected1)
	}

	// Test case 2: Input slice with single element
	input2 := []string{"a"}
	expected2 := []string{"a"}
	if output2 := zerocollections.Unique(input2); !reflect.DeepEqual(output2, expected2) {
		t.Errorf("Unique(%v) = %v, want %v", input2, output2, expected2)
	}

	// Test case 3: Input slice with multiple duplicate elements
	input3 := []int{1, 1, 2, 3, 3}
	expected3 := []int{1, 2, 3}
	if output3 := zerocollections.Unique(input3); len(output3) != len(expected3) {
		t.Errorf("Unique(%v) = %v, want %v", input3, output3, expected3)
	}

	// Test case 4: Input slice with all unique elements
	input4 := []float64{1.5, 2.3, 3.7}
	expected4 := []float64{1.5, 2.3, 3.7}
	if output4 := zerocollections.Unique(input4); len(output4) != len(expected4) {
		t.Errorf("Unique(%v) = %v, want %v", input4, output4, expected4)
	}

	// Test case 5: Input slice with mixed types
	input5 := []interface{}{"a", "b", "a", "c", "b"}
	expected5 := []interface{}{"a", "b", "c"}
	if output5 := zerocollections.Unique(input5); len(output5) != len(expected5) {
		t.Errorf("Unique(%v) = %v, want %v", input5, output5, expected5)
	}
}

func TestDifferenceWithDifferentDataTypes(t *testing.T) {
	type testCase struct {
		a        []interface{}
		b        []interface{}
		expected []interface{}
	}

	testCases := []testCase{
		{a: []interface{}{1, 2, 3}, b: []interface{}{1, 3}, expected: []interface{}{2}},
		{a: []interface{}{1, 2, 3}, b: []interface{}{}, expected: []interface{}{1, 2, 3}},
	}

	for _, tc := range testCases {
		result := zerocollections.Difference(tc.a, tc.b)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Difference(%v, %v) = %v; want %v", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestMaxCount(t *testing.T) {
	t.Run("First list longer", func(t *testing.T) {
		a := []int{1, 2, 3, 4}
		b := []int{5, 6}
		expected := 4
		result := zerocollections.MaxCount(a, b)
		if result != expected {
			t.Errorf("Expected %d, but got %d", expected, result)
		}
	})

	t.Run("Second list longer", func(t *testing.T) {
		a := []string{"a", "b"}
		b := []string{"c", "d", "e"}
		expected := 3
		result := zerocollections.MaxCount(a, b)
		if result != expected {
			t.Errorf("Expected %d, but got %d", expected, result)
		}
	})

	t.Run("Both lists equal length", func(t *testing.T) {
		a := []float64{1.1, 2.2, 3.3}
		b := []float64{4.4, 5.5, 6.6}
		expected := 3
		result := zerocollections.MaxCount(a, b)
		if result != expected {
			t.Errorf("Expected %d, but got %d", expected, result)
		}
	})
}
