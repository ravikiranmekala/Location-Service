package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravikiranmekala/MapUp-Backend-assessment/database"
	"github.com/ravikiranmekala/MapUp-Backend-assessment/models"
)

func GetLocationController(c *gin.Context) {
	startTime := time.Now()

	category := c.Param("category")

	var locations []models.Location
	database.DB.Db.Where("category = ?", category).Find(&locations)

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Nanoseconds()

	c.JSON(200, gin.H{"locations": locations, "time_ns": elapsedTime})
}

func CreateLocationController(c *gin.Context) {
	startTime := time.Now()

	location := models.Location{}
	fmt.Println(location)
	c.BindJSON(&location)
	fmt.Println(location)

	database.DB.Db.Create(&location)

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Nanoseconds()

	response := map[string]interface{}{
		"status":  "success",
		"id":      location.Id,
		"time_ns": elapsedTime,
	}
	c.JSON(200, response)
}

type TripCostRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type TripCostResponse struct {
	TotalCost float64 `json:"total_cost"`
	FuelCost  float64 `json:"fuel_cost"`
	TollCost  float64 `json:"toll_cost"`
}

func TripCostController(c *gin.Context) {
	startTime := time.Now()

	// Extract location ID from URL parameters
	locationID := c.Param("location_id")

	// Parse user's current location from the request body
	var request TripCostRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Fetch destination location from the database
	var destination models.Location
	if err := database.DB.Db.First(&destination, locationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	// Fetch trip cost using TollGuru API
	tripCost, err := getTripCost(request, destination)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trip cost"})
		return
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Nanoseconds()

	// Add time_ns to the tripCost response
	response := gin.H{
		"total_cost": tripCost.TotalCost,
		"fuel_cost":  tripCost.FuelCost,
		"toll_cost":  tripCost.TollCost,
		"time_ns":    elapsedTime,
	}

	// Return JSON response
	c.JSON(http.StatusOK, response)

}

func getTripCost(currentLocation TripCostRequest, destination models.Location) (TripCostResponse, error) {
	tollGuruURL := "https://apis.tollguru.com/toll/v2/origin-destination-waypoints"

	requestPayload := map[string]interface{}{
		"from": map[string]float64{
			"lat": currentLocation.Latitude,
			"lng": currentLocation.Longitude,
		},
		"to": map[string]float64{
			"lat": destination.Latitude,
			"lng": destination.Longitude,
		},
	}

	// Marshal the request payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return TripCostResponse{}, err
	}

	// Create an HTTP request to the TollGuru API
	req, err := http.NewRequest("POST", tollGuruURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return TripCostResponse{}, err
	}

	// Set the API key in the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "dgBRrqtqhdB2842gJ4tMrRJhTD7nQ3gh")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TripCostResponse{}, err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var tollGuruResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&tollGuruResponse)
	if err != nil {
		return TripCostResponse{}, err
	}
	fmt.Println("TollGuru response:")
	fmt.Println(tollGuruResponse)

	// Calculate total fuel and toll costs by traversing the response structure
	var totalFuelCost, totalTollCost float64
	routes, ok := tollGuruResponse["routes"].([]interface{})
	if !ok || len(routes) == 0 {
		return TripCostResponse{}, fmt.Errorf("No valid routes found")
	}

	fmt.Println("Routes:")
	fmt.Println(routes)

	// Pick the first route
	firstRoute, ok := routes[0].(map[string]interface{})
	if !ok {
		return TripCostResponse{}, fmt.Errorf("Invalid route format")
	}

	costs, ok := firstRoute["costs"].(map[string]interface{})
	if !ok {
		return TripCostResponse{}, fmt.Errorf("Costs not found in route")
	}

	// Set fuelCost to zero if 'fuel' is not available or is nil
	if fuel, ok := costs["fuel"].(float64); ok {
		totalFuelCost += fuel
	}

	// Set tollCost to zero if 'tag' is not available or is nil
	if tag, ok := costs["tag"].(float64); ok {
		totalTollCost += tag
	}

	fmt.Println("Total fuel cost:", totalFuelCost)
	fmt.Println("Total toll cost:", totalTollCost)

	// Create and return the response
	response := TripCostResponse{
		TotalCost: totalFuelCost + totalTollCost,
		FuelCost:  totalFuelCost,
		TollCost:  totalTollCost,
	}

	return response, nil
}

func SearchLocationController(c *gin.Context) {
	startTime := time.Now()
	var searchParams struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		Category  string  `json:"category" binding:"required"`
		RadiusKM  float64 `json:"radius_km" binding:"required"`
	}

	if err := c.BindJSON(&searchParams); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	minLat, maxLat, minLon, maxLon := calculateBoundingBox(
		searchParams.Latitude,
		searchParams.Longitude,
		searchParams.RadiusKM,
	)

	var locations []models.Location
	database.DB.Db.Where("category = ? AND latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
		searchParams.Category, minLat, maxLat, minLon, maxLon).Find(&locations)

	// Filter locations based on the exact distance within the specified radius
	filteredLocations := filterLocationsByDistance(
		searchParams.Latitude,
		searchParams.Longitude,
		searchParams.RadiusKM,
		locations,
	)

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Nanoseconds()

	// Return JSON response
	c.JSON(200, gin.H{"locations": filteredLocations, "time_ns": elapsedTime})

}

// calculateBoundingBox calculates the bounding box coordinates based on the radius
func calculateBoundingBox(lat, lon, radius float64) (minLat, maxLat, minLon, maxLon float64) {
	// Earth radius in kilometers
	earthRadius := 6371.0

	// Calculate angular distance in radians on the earth's surface
	angularDist := radius / earthRadius

	// Calculate bounding box coordinates
	minLat = lat - radToDeg(angularDist)
	maxLat = lat + radToDeg(angularDist)
	minLon = lon - radToDeg(angularDist)
	maxLon = lon + radToDeg(angularDist)

	return minLat, maxLat, minLon, maxLon
}

// degToRad converts degrees to radians
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

// radToDeg converts radians to degrees
func radToDeg(rad float64) float64 {
	return rad * (180.0 / math.Pi)
}

// filterLocationsByDistance filters locations based on the exact distance within the specified radius
func filterLocationsByDistance(lat, lon, radius float64, locations []models.Location) []gin.H {
	var filteredLocations []gin.H

	// Earth radius in kilometers
	earthRadius := 6371.0

	// Convert latitude and longitude from degrees to radians
	latRad := degToRad(lat)
	lonRad := degToRad(lon)

	// Calculate maximum distance in radians for filtering
	maxDistRad := radius / earthRadius

	for _, loc := range locations {
		// Convert location coordinates to radians
		locLatRad := degToRad(loc.Latitude)
		locLonRad := degToRad(loc.Longitude)

		// Calculate the distance between the points using the Haversine formula
		distRad := haversine(latRad, lonRad, locLatRad, locLonRad)

		// Check if the distance is within the specified radius
		if distRad <= maxDistRad {
			// Calculate the distance in kilometers
			distanceKM := radToDeg(distRad) * earthRadius

			// Add the location to the filtered list with distance
			filteredLocation := gin.H{
				"id":       loc.ID,
				"name":     loc.Name,
				"address":  loc.Address,
				"distance": distanceKM,
				"category": loc.Category,
			}
			filteredLocations = append(filteredLocations, filteredLocation)
		}
	}

	return filteredLocations
}

// haversine calculates the Haversine distance between two points given their coordinates in radians
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Haversine formula
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c
}
