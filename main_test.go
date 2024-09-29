package main

import (
	"math"
	"testing"
)
const float64EqualityThreshold = 1e-9

// source: https://stackoverflow.com/a/47969546
func almostEqual(a, b float64) bool {
    return math.Abs(a - b) <= float64EqualityThreshold
}

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

func TestAttemptMerge(t *testing.T) {
  // all should be 1.414 distance
	loads := []Load{
		{ID: "A", Pickup: Point{0, 0}, Dropoff: Point{1, 1}},   
		{ID: "B", Pickup: Point{1, 1}, Dropoff: Point{2, 2}},
		{ID: "C", Pickup: Point{2, 2}, Dropoff: Point{3, 3}},
		{ID: "D", Pickup: Point{3, 3}, Dropoff: Point{4, 4}},
	}

	route1 := &Route{loads: []int{0, 1}, totalDistance: 4.242} // A -> B
	route2 := &Route{loads: []int{2, 3}, totalDistance: 4.242} // C -> D

	success, mergedRoute := attemptToMerge(route1, route2, loads)

	if !success {
		t.Errorf("expected merge success, got failure")
	}

	expectedLoads := []int{0, 1, 2, 3}
	if len(mergedRoute.loads) != len(expectedLoads) {
		t.Errorf("expected %d loads in merged route, got %d", len(expectedLoads), len(mergedRoute.loads))
	}
	for i, loadIdx := range expectedLoads {
		if mergedRoute.loads[i] != loadIdx {
			t.Errorf("expected load %d at position %d, got %d", loadIdx, i, mergedRoute.loads[i])
		}
	}

  // should be the sum of both floats distances
  // also floats are a pain in tests
	expectedDistance := 8.485281374 
	if !almostEqual(mergedRoute.totalDistance, expectedDistance) {
		t.Errorf("expected total distance %.9f, got %.9f", expectedDistance, mergedRoute.totalDistance)
	}
}

func TestMergeRoutes(t *testing.T) {
	loads := []Load{
		{ID: "A", Pickup: Point{0, 0}, Dropoff: Point{1, 1}},
		{ID: "B", Pickup: Point{1, 1}, Dropoff: Point{2, 2}},
		{ID: "C", Pickup: Point{2, 2}, Dropoff: Point{3, 3}},
		{ID: "D", Pickup: Point{3, 3}, Dropoff: Point{4, 4}},
	}

	savings := []Saving{
		{i: 0, j: 1, saving: 2.0},
		{i: 1, j: 2, saving: 1.5},
		{i: 2, j: 3, saving: 1.2},
	}

	finalRoutes := mergeRoutes(savings, loads)

	expectedNumRoutes := 1 // all loads in this case should merge into a single route
	if len(finalRoutes) != expectedNumRoutes {
		t.Errorf("expected %d routes, got %d", expectedNumRoutes, len(finalRoutes))
	}

	expectedLoads := []int{0, 1, 2, 3}
	for i, loadIdx := range expectedLoads {
		if finalRoutes[0].loads[i] != loadIdx {
			t.Errorf("expected load %d at position %d, got %d", loadIdx, i, finalRoutes[0].loads[i])
		}
	}

	if len(finalRoutes[0].loads) != len(expectedLoads) {
		t.Errorf("expected %d loads in the merged route, got %d", len(expectedLoads), len(finalRoutes[0].loads))
	}
}


func TestExtractSchedules(t *testing.T) {
	loads := []Load{
		{ID: "A", Pickup: Point{0, 0}, Dropoff: Point{1, 1}},
		{ID: "B", Pickup: Point{1, 1}, Dropoff: Point{2, 2}},
		{ID: "C", Pickup: Point{2, 2}, Dropoff: Point{3, 3}},
	}

	routes := []*Route{
		{loads: []int{0, 1}}, // Route 1 with load indices 0, 1
		{loads: []int{2}},    // Route 2 with load index 2
	}

	result := extractSchedules(routes, loads)

	expected := [][]string{
		{"A", "B"}, // Route 1: Loads A and B
		{"C"},      // Route 2: Load C
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d schedules, got %d", len(expected), len(result))
	}

	for i, schedule := range result {
		if len(schedule) != len(expected[i]) {
			t.Errorf("expected schedule %v, got %v", expected[i], schedule)
		}
		for j, id := range schedule {
			if id != expected[i][j] {
				t.Errorf("expected load ID %s at position %d, got %s", expected[i][j], j, id)
			}
		}
	}
}
