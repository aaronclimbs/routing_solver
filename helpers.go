package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Distance calculates the euclidean distance between two points
func Distance(p1, p2 Point) float64 {
	xDiff := p1.x - p2.x
	yDiff := p1.y - p2.y
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

// CalculateRouteDistance calculates the distance traveled for a given route.
func CalculateRouteDistance(route *Route, loads []Load, startLocation Point) float64 {
	totalDistance := 0.0
	currentLocation := startLocation

	for _, loadIdx := range route.loads {
		// add the distance to the pickup point and then to the dropoff point
		totalDistance += Distance(currentLocation, loads[loadIdx].Pickup)
		totalDistance += Distance(loads[loadIdx].Pickup, loads[loadIdx].Dropoff)
		// update the current location to the dropoff point
		currentLocation = loads[loadIdx].Dropoff
		// fmt.Printf("current_load: %d, totalDistance: %f\n", loadIdx, totalDistance)
	}

	return totalDistance
}

// ParsePoint converts a string in the format "(x,y)" into a Point struct
func ParsePoint(s string) Point {
	cleanedString := strings.TrimSpace(s)
	cleanedString = strings.TrimPrefix(cleanedString, "(")
	cleanedString = strings.TrimSuffix(cleanedString, ")")
	coords := strings.Split(cleanedString, ",")

	// this shouldn't happen but I've see one too many NPEs from filereaders
	if len(coords) != 2 {
		return DEPOT_LOCATION
	}

	x, errX := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
	y, errY := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)

	// this shouldn't happen
	if errX != nil || errY != nil {
		return DEPOT_LOCATION
	}

	return Point{x, y}
}

// SortSavings sorts the savings in descending order
// might as well save some space by mutating existing slice
func SortSavings(savings []Saving) {
	sort.Slice(savings, func(a, b int) bool {
		return savings[a].saving > savings[b].saving
	})
}

// ExtractSchedules generates a list of schedules from the routes
func ExtractSchedules(routes []*Route, loads []Load) [][]string {
	schedules := make([][]string, 0, len(routes))

	for _, route := range routes {
		schedule := make([]string, 0, len(route.loads))

		for _, loadIdx := range route.loads {
			schedule = append(schedule, loads[loadIdx].ID)
		}

		schedules = append(schedules, schedule)
	}

	return schedules
}

// PrintSchedules formats and prints the schedules
func PrintSchedules(schedules [][]string) {
	for _, schedule := range schedules {
		fmt.Printf("[%s]\n", strings.Join(schedule, ","))
	}
}
