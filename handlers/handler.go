package handler

import (
	"database/sql"
	"math"
	"sort"

	// core packages
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// local packages
	"SpotsTest/database"
	"SpotsTest/models"
)

type Point struct {
	first  models.Spot
	second float64
}

func GetSpotsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	latStr := r.URL.Query().Get("latitude")
	lonStr := r.URL.Query().Get("longitude")
	radiusStr := r.URL.Query().Get("radius")
	recruitmentType := r.URL.Query().Get("type")

	// Convert parameters to appropriate types
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	// Prepare the query based on the recruitment type
	var query string
	if recruitmentType == "circle" {
		query = `
			SELECT
				id,
				name, 
				website, 
				coordinates, 
				description, 
				rating,
				st_distance(coordinates, ST_MakePoint($1, $2)::geography) AS distance
			FROM public."MY_TABLE"
			WHERE ST_DWithin(
			    ST_MakePoint($1, $2)::geography, 
			    coordinates,
			    $3
			)
			ORDER BY distance;`
	} else if recruitmentType == "square" {
		query = `
			SELECT
				id,
				name,
				website,
				coordinates,
				description,
				rating,
				st_distance(coordinates, given_point) AS distance
			FROM 
			    public."MY_TABLE",
			    (SELECT
					$1::double precision AS long, 
					$2::double precision AS lat,
					ST_MakePoint($1, $2)::geography AS given_point,
					$3::double precision AS radius,
					111319.0 AS metersPerDegree) AS c
			WHERE 
				coordinates::geometry &&
				ST_MakeEnvelope(
					long - (radius / (metersPerDegree * COS(RADIANS(lat)))),
					lat  - (radius / metersPerDegree), 
					long + (radius / (metersPerDegree * COS(RADIANS(lat)))), 
					lat  + (radius / metersPerDegree), 
					4326
				)
			ORDER BY distance;`
	} else {
		http.Error(w, "Invalid recruitment type", http.StatusBadRequest)
		return
	}

	// Execute the query
	rows, err := database.DB.Query(query, lon, lat, radius)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Printf("error executing query %q: %v", query, err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("error closing rows: ", err)
		}
	}(rows)

	// Create an array to store the spots
	var spots []models.Spot
	var pairs []Point

	// Iterate through the result set and populate the spots array
	for rows.Next() {
		var spot models.Spot
		var point Point
		err := rows.Scan(&spot.ID, &spot.Name, &spot.Website, &spot.Coordinates, &spot.Description, &spot.Rating, &point.second)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Println("error scanning row: ", err)
			return
		}
		point.first = spot
		pairs = append(pairs, point)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating result set", http.StatusInternalServerError)
		return
	}

	// sort spots by rating of distance <50m
	sort.Slice(pairs, func(i, j int) bool {
		if math.Abs(pairs[i].second-pairs[j].second) < 50 {
			return pairs[i].first.Rating > pairs[j].first.Rating
		}

		return pairs[i].second < pairs[j].second
	})

	// loop through pairs and append to spots
	for _, pair := range pairs {
		log.Printf("spot: %v, Dist: %f", pair.first, pair.second)
		spots = append(spots, pair.first)
	}

	// Convert the spots array to JSON
	respJSON, err := json.Marshal(spots)
	if err != nil {
		http.Error(w, "JSON marshaling error", http.StatusInternalServerError)
		return
	}

	// Set the response content type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(spots) == 0 {
		_, _ = w.Write([]byte("[]")) // return empty array if no spots found
		return
	}
	_, err = w.Write(respJSON)
	if err != nil {
		return
	}
}
