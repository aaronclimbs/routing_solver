package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	MAX_DISTANCE = 720.0
)

var (
	DEPOT_LOCATION = Point{0, 0}
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

// Route represents a vehicle route
type Route struct {
	loads         []int // ids of loads in this route
	totalDistance float64
	lastDropoff   Point // ultimately will be depot but need to keep track
}

// attemptToMerge tries to merge two routes if the total distance remains within the allowed limit
// it returns true and the new merged route if successful
func attemptToMerge(route1, route2 *Route, loads []Load) (bool, *Route) {
	totalDistance := 0.0
	currentLocation := DEPOT_LOCATION

	// calculate the total distance for route1, then for route2, and back to depot
	totalDistance = calculateRouteDistance(route1, loads, currentLocation)

	// get the last dropoff from route1 and the first pickup from route2
	if len(route1.loads) > 0 && len(route2.loads) > 0 {
		lastDropoff := loads[route1.loads[len(route1.loads)-1]].Dropoff
		firstPickup := loads[route2.loads[0]].Pickup
		// add the distance between the last dropoff of route1 and the first pickup of route2
		totalDistance += distance(lastDropoff, firstPickup)
		currentLocation = firstPickup // update the current location to the first pickup of route2
	}

	totalDistance += calculateRouteDistance(route2, loads, currentLocation)
	totalDistance += distance(currentLocation, DEPOT_LOCATION)

  fmt.Printf("totalDistance: %.3f\n", totalDistance)

	if totalDistance <= MAX_DISTANCE {
		// then we can merge
		mergedLoads := append(route1.loads, route2.loads...)
		mergedRoute := &Route{
			loads:         mergedLoads,
			totalDistance: totalDistance,
			lastDropoff:   currentLocation,
		}
		return true, mergedRoute
	}

	// false if the merge would exceed max distance allowed
	return false, nil
}

// calculateRouteDistance calculates the distance traveled for a given route.
func calculateRouteDistance(route *Route, loads []Load, startLocation Point) float64 {
	totalDistance := 0.0
	currentLocation := startLocation

	for _, loadIdx := range route.loads {
		// add the distance to the pickup point and then to the dropoff point
		totalDistance += distance(currentLocation, loads[loadIdx].Pickup)
		totalDistance += distance(loads[loadIdx].Pickup, loads[loadIdx].Dropoff)
		// update the current location to the dropoff point
		currentLocation = loads[loadIdx].Dropoff
		fmt.Printf("current_load: %d, totalDistance: %f\n", loadIdx, totalDistance)
	}

	return totalDistance
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

// computeSavings calculates the savings for all pairs of loads.
func computeSavings(loads []Load) []Saving {
	var savingsList []Saving

	// loop through each pair of loads (i, j) while i < j
	for i := 0; i < len(loads); i++ {
		for j := i + 1; j < len(loads); j++ {

			// this just utilizes the basic distance between the current pickup/dropoff
			// to the depot vs to the next pickup/dropoff
			distanceDepotToLoadI := distance(DEPOT_LOCATION, loads[i].Pickup)
			distanceLoadJToDepot := distance(loads[j].Dropoff, DEPOT_LOCATION)
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

// mergeRoutes merges routes based on savings
func mergeRoutes(savings []Saving, loads []Load) []*Route {
	routes := make(map[int]*Route)
	routeIndices := make(map[int]int)

	nextRouteID := 0
	// create an individual route for each load
	for idx, load := range loads {
		// assigning a single vehicle for every load, calculating the distance from the depot,
		// to the load's pickup, then to the dropoff, and back to the depot.
		// this works to start, but would obviously make the cost much higher as it would use
		// the maximum number of vehicles
		totalDistance := distance(DEPOT_LOCATION, load.Pickup) +
			distance(load.Pickup, load.Dropoff) +
			distance(load.Dropoff, DEPOT_LOCATION)
		route := &Route{
			loads:         []int{idx},
			totalDistance: totalDistance,
			lastDropoff:   load.Dropoff,
		}
		routes[nextRouteID] = route     // assign a route to the current load
		routeIndices[idx] = nextRouteID // map the load index to the route ID
		nextRouteID++
	}

	for _, saving := range savings {
		// current route IDs for the two loads in this saving
		routeIDi := routeIndices[saving.i]
		routeIDj := routeIndices[saving.j]

		// only attempt to merge if the loads are in different routes
		// if this is already merged, it could duplicate loads across routes
		if routeIDi != routeIDj {
			routeI := routes[routeIDi]
			routeJ := routes[routeIDj]

			canMerge, newRoute := attemptToMerge(routeI, routeJ, loads)
			if canMerge {
				// merge successful, remove the old routes
				delete(routes, routeIDi)
				delete(routes, routeIDj)
				routes[nextRouteID] = newRoute

				// all loads in newly merged route need to be mapped to correct routeID
				for _, loadIdx := range newRoute.loads {
					routeIndices[loadIdx] = nextRouteID
				}

				// increment the routeID for next potential merge
				nextRouteID++
			}
		}
	}

	// collect all remaining routes into the final result
	finalRoutes := make([]*Route, 0, len(routes))
	for _, route := range routes {
		finalRoutes = append(finalRoutes, route)
	}

	return finalRoutes
}

// extractSchedules generates a list of schedules from the routes
func extractSchedules(routes []*Route, loads []Load) [][]string {
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

// printSchedules formats and prints the schedules
func printSchedules(schedules [][]string) {
	for _, schedule := range schedules {
		fmt.Printf("[%s]\n", strings.Join(schedule, ","))
	}
}

