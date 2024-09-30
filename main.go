package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	MAX_DISTANCE = 720.0

// COST_PENALTY = 500
)

var (
	DEPOT_LOCATION = Point{0, 0}
)

// computeSavings calculates the savings for all pairs of loads.
func computeSavings(loads []Load) []Saving {
	var savingsList []Saving

	// loop through each pair of loads (i, j) while i < j
	for i := 0; i < len(loads); i++ {
		for j := i + 1; j < len(loads); j++ {

			// this just utilizes the basic distance between the current pickup/dropoff
			// to the depot vs to the next pickup/dropoff
			distanceDepotToLoadI := Distance(DEPOT_LOCATION, loads[i].Pickup)
			distanceLoadJToDepot := Distance(loads[j].Dropoff, DEPOT_LOCATION)
			distanceLoadIToLoadJ := Distance(loads[i].Dropoff, loads[j].Pickup)

			savingValue := distanceDepotToLoadI + distanceLoadJToDepot - distanceLoadIToLoadJ

			// will need to confirm the distance from the next load back to depo is not greater
			// than max distance allowed

			// append savings value to list
			savingsList = append(savingsList, Saving{i: i, j: j, saving: savingValue})
		}
	}

	return savingsList
}

func attemptToMerge(route1, route2 *Route, loads []Load) (bool, *Route) {
	totalDistance := 0.0
	currentLocation := DEPOT_LOCATION

	totalDistance += CalculateRouteDistance(route1, loads, currentLocation)
	if len(route1.loads) > 0 {
		// update currentlocation to last dropoff location in route
		currentLocation = loads[route1.loads[len(route1.loads)-1]].Dropoff
	}

	// add distance from last dropoff of route1 to first pickup of route2
	if len(route2.loads) > 0 {
		firstPickup := loads[route2.loads[0]].Pickup
		totalDistance += Distance(currentLocation, firstPickup)
		currentLocation = firstPickup
	}

	totalDistance += CalculateRouteDistance(route2, loads, currentLocation)
	if len(route2.loads) > 0 {
		// update currentlocation to last dropoff location in route
		currentLocation = loads[route2.loads[len(route2.loads)-1]].Dropoff
	}

	// Add distance back to depot
	totalDistance += Distance(currentLocation, DEPOT_LOCATION)

	// separateCost := (2 * COST_PENALTY) + route1.totalDistance + route2.totalDistance

	// // if we merge them, we avoid one vehicle cost
	// mergedCost := COST_PENALTY + totalDistance

	// Check against MAX_DISTANCE
	if totalDistance <= MAX_DISTANCE {
		mergedLoads := append(route1.loads, route2.loads...)
		mergedRoute := &Route{
			loads:         mergedLoads,
			totalDistance: totalDistance,
			lastDropoff:   currentLocation,
		}
		return true, mergedRoute
	}
	return false, nil
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
		totalDistance := Distance(DEPOT_LOCATION, load.Pickup) +
			Distance(load.Pickup, load.Dropoff) +
			Distance(load.Dropoff, DEPOT_LOCATION)
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <input_file>")
		return
	}

	inputFile := os.Args[1]
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	loads := []Load{}

	scanner := bufio.NewScanner(file)
	gotHeader := false
	for scanner.Scan() {
		line := scanner.Text()
		if !gotHeader {
			gotHeader = true
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		splits := strings.Fields(line)
		if len(splits) != 3 {
			fmt.Println("Invalid line:", line)
			continue
		}

		id := splits[0]
		pickup := ParsePoint(splits[1])
		dropoff := ParsePoint(splits[2])

		load := Load{
			ID:      id,
			Pickup:  pickup,
			Dropoff: dropoff,
		}
		loads = append(loads, load)
	}

	// handle any scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	savings := computeSavings(loads)

	SortSavings(savings)

	finalRoutes := mergeRoutes(savings, loads)

	schedules := ExtractSchedules(finalRoutes, loads)
	PrintSchedules(schedules)
}
