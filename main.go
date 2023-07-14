package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type GeoJSON struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type IntersectionResult struct {
	LineID       string    `json:"lineId"`
	Intersection []float64 `json:"intersection"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/intersect", handleIntersect).Methods("POST")

	http.ListenAndServe(":8000", r)
}

func handleIntersect(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	authToken := r.Header.Get("Authorization")
	if !isValidAuthToken(authToken) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
		return
	}

	// Parse the GeoJSON linestring from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to read request body: %v", err)
		return
	}

	var linestring GeoJSON
	err = json.Unmarshal(body, &linestring)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to parse GeoJSON linestring: %v", err)
		return
	}

	// Calculate intersections
	intersections := []IntersectionResult{}

	for i := 1; i <= 50; i++ {
		line := getLineByID(fmt.Sprintf("L%02d", i))
		intersection := calculateLineIntersection(linestring.Coordinates, line.Coordinates)

		if intersection != nil {
			intersections = append(intersections, IntersectionResult{
				LineID:       fmt.Sprintf("L%02d", i),
				Intersection: intersection,
			})
		}
	}

	// Generate the response based on the intersections found
	response, err := json.Marshal(intersections)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to generate response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func isValidAuthToken(token string) bool {
	// Implement your authentication logic here
	// Verify the validity of the token or implement any other authentication method
	// Return true if the token is valid, otherwise return false

	// Example implementation: Accept any non-empty token as valid
	return token != ""
}

// Example function to retrieve line coordinates by ID
func getLineByID(lineID string) GeoJSON {
	// Implement your logic to retrieve line coordinates by ID
	// This is just an example, replace it with your actual data source or algorithm

	// Dummy line data
	lines := map[string][][]float64{
		"L01": {{0, 0}, {10, 10}},
		"L02": {{5, 5}, {15, 15}},
		// Add coordinates for other lines
	}

	return GeoJSON{
		Type:        "LineString",
		Coordinates: lines[lineID],
	}
}

// Calculate the intersection between the linestring and line
func calculateLineIntersection(linestring, line [][]float64) []float64 {
	intersections := []float64{}

	for i := 0; i < len(linestring)-1; i++ {
		pointA := linestring[i]
		pointB := linestring[i+1]

		for j := 0; j < len(line)-1; j++ {
			pointC := line[j]
			pointD := line[j+1]

			intersection := calculateIntersectionPoint(pointA, pointB, pointC, pointD)
			if intersection != nil {
				intersections = append(intersections, intersection[0], intersection[1])
			}
		}
	}

	return intersections
}

// Calculate the intersection point between two lines
func calculateIntersectionPoint(pointA, pointB, pointC, pointD []float64) []float64 {
	x1, y1 := pointA[0], pointA[1]
	x2, y2 := pointB[0], pointB[1]
	x3, y3 := pointC[0], pointC[1]
	x4, y4 := pointD[0], pointD[1]

	denominator := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if denominator == 0 {
		// Lines are parallel, no intersection
		return nil
	}

	x := ((x1*y2-y1*x2)*(x3-x4) - (x1-x2)*(x3*y4-y3*x4)) / denominator
	y := ((x1*y2-y1*x2)*(y3-y4) - (y1-y2)*(x3*y4-y3*x4)) / denominator

	// Check if intersection point is within the line segments
	if x < min(x1, x2) || x > max(x1, x2) || x < min(x3, x4) || x > max(x3, x4) ||
		y < min(y1, y2) || y > max(y1, y2) || y < min(y3, y4) || y > max(y3, y4) {
		return nil
	}

	return []float64{x, y}
}

// Utility function to calculate the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Utility function to calculate the maximum of two float64 values
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
