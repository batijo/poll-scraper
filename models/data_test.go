package models

import (
	"reflect"
	"testing"
)

func TestFilterData(t *testing.T) {
	data := []Data{
		{Name: "A", Value: "1"},
		{Name: "B", Value: "2"},
		{Name: "C", Value: "3"},
	}

	result := FilterData([]int{0, 2}, data)

	if len(result) != 2 {
		t.Fatalf("got %d items, want 2", len(result))
	}
	if result[0].Name != "A" || result[1].Name != "C" {
		t.Errorf("got [%s, %s], want [A, C]", result[0].Name, result[1].Name)
	}
}

func TestSumLines(t *testing.T) {
	data := []Data{
		{Name: "A", Value: "10"},
		{Name: "B", Value: "20"},
		{Name: "C", Value: "text"},
	}

	got := SumLines(data, "")

	want := []Data{{Name: "sum", Value: "30"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSumLines_WithSymbol(t *testing.T) {
	data := []Data{
		{Name: "A", Value: "10"},
		{Name: "B", Value: "20"},
	}

	got := SumLines(data, "$")

	want := []Data{{Name: "sum", Value: "30"}, {Name: "sum_symbol", Value: "30$"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAddSymbols(t *testing.T) {
	data := []Data{
		{Name: "A", Value: "1"},
		{Name: "B", Value: "2"},
		{Name: "C", Value: "3"},
	}

	got := AddSymbols(data, map[int]string{1: " eur", 3: "%", 9: "x"})

	want := []Data{
		{Name: "A", Value: "1"},
		{Name: "A_symbol", Value: "1 eur"},
		{Name: "B", Value: "2"},
		{Name: "C", Value: "3"},
		{Name: "C_symbol", Value: "3%"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAddLines(t *testing.T) {
	lines := []Data{
		{Name: "custom1", Value: "foo"},
		{Name: "custom2", Value: "bar"},
	}
	result := AddLines(nil, lines)

	if len(result) != 2 {
		t.Fatalf("got %d items, want 2", len(result))
	}
	if result[0].Name != "custom1" || result[0].Value != "foo" {
		t.Errorf("got {%s, %s}, want {custom1, foo}", result[0].Name, result[0].Value)
	}
	if result[1].Name != "custom2" || result[1].Value != "bar" {
		t.Errorf("got {%s, %s}, want {custom2, bar}", result[1].Name, result[1].Value)
	}
}
