package main

import (
	"math"
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
