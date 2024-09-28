package main

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// Point represents a 2D point with x, y coordinates
type Point struct {
	x, y float64
}

// Load represents a task with a pickup and dropoff location
type Load struct {
	ID      string
	Pickup  Point
	Dropoff Point
}

// Saving represents the savings between two loads
type Saving struct {
	i, j   int     
	saving float64
}

// distance calculates the euclidean distance between two points
func distance(p1, p2 Point) float64 {
	xDiff := p1.x - p2.x
	yDiff := p1.y - p2.y
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

// parsePoint converts a string in the format "(x,y)" into a Point struct
func parsePoint(s string) Point {
	cleanedString := strings.TrimSpace(s)
	cleanedString = strings.TrimPrefix(cleanedString, "(")
	cleanedString = strings.TrimSuffix(cleanedString, ")")
	coords := strings.Split(cleanedString, ",")

  // this shouldn't happen but i've see one too many NPEs from filereaders
	if len(coords) != 2 {
		return Point{0, 0} // Return (0, 0) for invalid input
	}

	x, errX := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
	y, errY := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)

  // this shouldn't happen
	if errX != nil || errY != nil {
		return Point{0, 0}
	}

	return Point{x, y}
}

// computeSavings calculates the savings for all pairs of loads.
func computeSavings(loads []Load) []Saving {
	depot := Point{0, 0}
	var savingsList []Saving

	// loop through each pair of loads (i, j) while i < j
	for i := 0; i < len(loads); i++ {
		for j := i + 1; j < len(loads); j++ {
      
      // this just utilizes the basic distance between the current pickup/dropoff
      // to the depot vs to the next pickup/dropoff
			distanceDepotToLoadI := distance(depot, loads[i].Pickup)
			distanceLoadJToDepot := distance(loads[j].Dropoff, depot)
			distanceLoadIToLoadJ := distance(loads[i].Dropoff, loads[j].Pickup)

			savingValue := distanceDepotToLoadI + distanceLoadJToDepot - distanceLoadIToLoadJ

      // will need to confirm the distance from the next load back to depo is not greater
      // than max distance allowed

      // append savings value to list
			savingsList = append(savingsList, Saving{i: i, j: j, saving: savingValue})
		}
	}

	return savingsList
}

// sortSavings sorts the savings in descending order
func sortSavings(savings []Saving) {
	sort.Slice(savings, func(a, b int) bool {
		return savings[a].saving > savings[b].saving
	})
}
