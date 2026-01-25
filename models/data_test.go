package models

import "testing"

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

func TestSumData(t *testing.T) {
	data := []Data{
		{Name: "A", Value: "10"},
		{Name: "B", Value: "20"},
	}

	result := SumData(data, "")

	if len(result) != 3 {
		t.Fatalf("got %d items, want 3", len(result))
	}
	if result[2].Name != "sum" || result[2].Value != "30" {
		t.Errorf("sum = {%q, %q}, want {%q, %q}", result[2].Name, result[2].Value, "sum", "30")
	}
}

func TestSumData_WithSymbol(t *testing.T) {
	data := []Data{
		{Name: "A", Value: "10"},
		{Name: "B", Value: "20"},
	}

	result := SumData(data, "$")

	if len(result) != 4 {
		t.Fatalf("got %d items, want 4", len(result))
	}
	if result[2].Name != "sum" || result[2].Value != "30" {
		t.Errorf("sum = {%q, %q}, want {%q, %q}", result[2].Name, result[2].Value, "sum", "30")
	}
	if result[3].Name != "sum_symbol" || result[3].Value != "30$" {
		t.Errorf("sum_symbol = {%q, %q}, want {%q, %q}", result[3].Name, result[3].Value, "sum_symbol", "30$")
	}
}

func TestAddLines(t *testing.T) {
	result := AddLines(nil, []string{"foo", "bar"})

	if len(result) != 2 {
		t.Fatalf("got %d items, want 2", len(result))
	}
	if result[0].Value != "foo" || result[1].Value != "bar" {
		t.Errorf("got [%s, %s], want [foo, bar]", result[0].Value, result[1].Value)
	}
}
