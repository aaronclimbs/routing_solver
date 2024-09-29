package main

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

// Route represents a vehicle route
type Route struct {
	loads         []int // ids of loads in this route
	totalDistance float64
	lastDropoff   Point // need to keep track while merging
}
