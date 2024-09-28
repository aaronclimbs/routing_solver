package main

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

func TestComputeSavings(t *testing.T) {
	loads := []Load{
		{Pickup: Point{x: 1, y: 1}, Dropoff: Point{x: 2, y: 2}},
		{Pickup: Point{x: 2, y: 3}, Dropoff: Point{x: 3, y: 4}},
		{Pickup: Point{x: 4, y: 5}, Dropoff: Point{x: 6, y: 7}},
	}

	expected := []Saving{
   {0, 1, 5.414213562373095},
   {0, 2, 7.028206744201993},
   {1, 2, 11.41088217038378},
	}

	savings := computeSavings(loads)

	for i := range savings {
		if savings[i].i != expected[i].i || savings[i].j != expected[i].j || !almostEqual(savings[i].saving, expected[i].saving) {
			t.Errorf("expected saving: %v, got: %v", expected[i], savings[i])
		}
	}
}

// https://stackoverflow.com/questions/47969385/go-float-comparison
func almostEqual(a, b float64) bool {
    return math.Abs(a - b) <= float64EqualityThreshold
}
