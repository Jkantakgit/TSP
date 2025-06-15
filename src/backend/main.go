package main

import (
	"math"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * Struct for storing City name and coordinates.
 */
type City struct {
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

/*
 * List of cities to be optimized.
 */
type CityRequest struct {
	Cities []City `json:"cities"`
}

/*
 * Response structure containing the optimized route.
 */
type CityResponse struct {
	Route []City `json:"route"`
}

func main() {
	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Define the POST endpoint for TSP optimization
	r.POST("/optimize", solveTSP)

	// Define the GET endpoint for health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "TSP backend is running"})
	})

	// Start the server on port 8080
	r.Run(":8080")
}

/*
* Function to solve the TSP problem by shuffling the cities.
 */
func solveTSP(c *gin.Context) {
	var req CityRequest
	// Print the incoming JSON for debugging

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	route := ACO(req.Cities)

	c.JSON(http.StatusOK, CityResponse{Route: route})
}

/*
* Function to solve the TSP problem using Ant Colony Optimization (ACO).
 */
func ACO(cities []City) []City {

	// Parameters for ACO
	n := len(cities)
	antCount := 20

	// pheromoneWeight := 1.0
	alpha := 1.0

	// heuristicWeight := 1.0
	beta := 2.0

	// evaporation rate of pheromones
	evaporationRate := 0.5

	// Q is a constant that influences the amount of pheromone deposited
	Q := 100.0

	// Number of iterations for the ACO algorithm
	iterations := 100

	pheromone := make([][]float64, n)

	// Initialize pheromone trail to small value
	for i := range pheromone {
		pheromone[i] = make([]float64, n)
		for j := range pheromone[i] {
			pheromone[i][j] = 1.0
		}
	}

	distances := make([][]float64, n)

	// Precalculate distances
	for i := range distances {
		distances[i] = make([]float64, n)
		for j := range distances[i] {
			if i == j {
				distances[i][j] = math.Inf(1)
			} else {
				distances[i][j] = distance(cities[i], cities[j]) + 0.0001
			}
		}
	}

	// Initialize best path and distance
	bestPath := make([]City, n)
	bestDistance := math.Inf(1)

	// Main loop for ACO iterations
	for range iterations {

		// Initialize paths and path lengths for each ant
		paths := make([][]City, antCount)
		pathLengths := make([]float64, antCount)

		// Each ant constructs a path
		for ant := range antCount {

			// Create a map to track allowed cities
			allowed := make(map[int]bool)
			for i := range n {
				allowed[i] = true
			}

			// Start from a random city
			start := rand.Intn(n)

			// Mark the starting city as visited
			allowed[start] = false

			// Initialize path with the starting city
			path := []int{start}

			// Construct the path until all cities are visited
			for len(path) < n {

				// Get the current city
				current := path[len(path)-1]

				// Calculate probabilities for next city
				probabilities := make([]float64, 0)
				candidates := make([]int, 0)

				// Initialize sum for probabilities
				var sum float64

				// Iterate over all cities to calculate probabilities
				for j := range n {
					// Only consider allowed cities
					if allowed[j] {
						// Calculate pheromone and heuristic values
						pher := math.Pow(pheromone[current][j], alpha)
						heur := math.Pow(1.0/distances[current][j], beta)
						w := pher * heur

						// Add the porbability to the sum
						sum += w

						// Store the probability and candidate city
						probabilities = append(probabilities, w)
						candidates = append(candidates, j)
					}
				}

				// Roulette Wheel Selection
				r := rand.Float64() * sum

				// Select the next city based on probabilities
				var selected bool
				for i, p := range probabilities {
					r -= p
					if r <= 0 {
						path = append(path, candidates[i])
						allowed[candidates[i]] = false
						selected = true
						break
					}
				}

				// If no city was selected (possible with floating point rounding),
				// select the first available one
				if !selected && len(candidates) > 0 {
					path = append(path, candidates[0])
					allowed[candidates[0]] = false
				}
			}

			// Calculate the tour length after the path is fully constructed
			tourLength := 0.0

			// Calculate the total distance of the path
			for i := 0; i < len(path)-1; i++ {
				tourLength += distances[path[i]][path[i+1]]
			}

			// Add the distance from the last city back to the first city
			tourLength += distances[path[len(path)-1]][path[0]]

			cityPath := make([]City, n)

			// Convert int path to City path
			for i, idx := range path {
				cityPath[i] = cities[idx]
			}

			paths[ant] = cityPath
			pathLengths[ant] = tourLength

			if tourLength < bestDistance {
				bestDistance = tourLength
				copy(bestPath, cityPath)
			}
		}

		// Evaporate pheromone
		for i := range pheromone {
			for j := range pheromone[i] {
				pheromone[i][j] *= (1 - evaporationRate)
				if pheromone[i][j] < 0.0001 {
					pheromone[i][j] = 0.0001
				}
			}
		}

		// Deposit pheromone on the paths found by ants
		for k := range antCount {
			// Find the indices of cities in the original array
			for i := 0; i < n-1; i++ {
				fromIdx := findCityIndex(cities, paths[k][i])
				toIdx := findCityIndex(cities, paths[k][i+1])
				pheromone[fromIdx][toIdx] += Q / pathLengths[k]
				pheromone[toIdx][fromIdx] += Q / pathLengths[k]
			}
			// Connect the last city with the first one
			fromIdx := findCityIndex(cities, paths[k][n-1])
			toIdx := findCityIndex(cities, paths[k][0])
			pheromone[fromIdx][toIdx] += Q / pathLengths[k]
			pheromone[toIdx][fromIdx] += Q / pathLengths[k]
		}
	}
	return bestPath
}

/*
* Helper function to find the index of a city in the cities array.
 */
func findCityIndex(cities []City, city City) int {
	for i, c := range cities {
		if c.Name == city.Name && c.X == city.X && c.Y == city.Y {
			return i
		}
	}
	return -1 // Should never happen if the city exists in the array
}

/*
* Function to calculate the Euclidean distance between two cities.
 */
func distance(city1, city2 City) float64 {
	// Euclidean distance should be the square root of sum of squares.
	return math.Sqrt((city1.X-city2.X)*(city1.X-city2.X) + (city1.Y-city2.Y)*(city1.Y-city2.Y))
}
