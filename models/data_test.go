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

	result := SumData(data)

	if len(result) != 3 {
		t.Fatalf("got %d items, want 3", len(result))
	}
	if result[2].Name != "sum" || result[2].Value != "30" {
		t.Errorf("sum = {%q, %q}, want {%q, %q}", result[2].Name, result[2].Value, "sum", "30")
	}
}

func TestAddLines(t *testing.T) {
	t.Setenv("PS_ADD_LINES", "[foo,bar]")

	result := AddLines(nil)

	if len(result) != 2 {
		t.Fatalf("got %d items, want 2", len(result))
	}
	if result[0].Value != "foo" || result[1].Value != "bar" {
		t.Errorf("got [%s, %s], want [foo, bar]", result[0].Value, result[1].Value)
	}
}
