package handler

import (
	"database/sql"
	// core packages
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// local packages
	"SpotsTest/database"
	"SpotsTest/models"
)

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
				rating
			FROM public."MY_TABLE"
			WHERE ST_DWithin(
			    ST_MakePoint($1, $2)::geography, 
			    coordinates,
			    $3
			);
		`
	} else if recruitmentType == "square" {
		query = `
			SELECT
				id,
				name,
				website,
				coordinates,
				description,
				rating
			FROM 
			    public."MY_TABLE",
			    (SELECT
					$1::double precision AS long, 
					$2::double precision AS lat,
					$3::double precision AS dist, 
					111319.0 AS metersPerDegree) AS c
			WHERE 
				coordinates::geometry &&
				ST_MakeEnvelope(
					long - (dist / (metersPerDegree * COS(RADIANS(lat)))),
					lat  - (dist / metersPerDegree), 
					long + (dist / (metersPerDegree * COS(RADIANS(lat)))), 
					lat  + (dist / metersPerDegree), 
					4326
				)`
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

	// Iterate through the result set and populate the spots array
	for rows.Next() {
		var spot models.Spot
		err := rows.Scan(&spot.ID, &spot.Name, &spot.Website, &spot.Coordinates, &spot.Description, &spot.Rating)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Println("error scanning row: ", err)
			return
		}
		spots = append(spots, spot)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating result set", http.StatusInternalServerError)
		return
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
