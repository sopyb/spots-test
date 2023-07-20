# Spots coding test

This is a coding test for Spots. The test is split in 2 parts:

## Part 1

The first part is a SQL query that returns spot information for those spots whose domain has occurred more than once in the data. 
The website field in the query is transformed such that it only contains the domain.
For the returned spots, the query also counts the occurrence of each domain. 
The final output of the query includes three columns - spot name, domain, and the count of each domain.

### My solution
```postgresql
SELECT
    name AS spot_name,
    SUBSTRING(website, '://([^/]+)') AS domain,
    COUNT(*) AS domain_count
FROM
    "MY_TABLE"
GROUP BY
    spot_name,
    SUBSTRING(website, '://([^/]+)')
HAVING
    COUNT(SUBSTRING(website, '://([^/]+)')) > 1;
```
also available in [part1.sql](part1.sql)

## Part 2
The second part of the task involves creating an endpoint in Golang that works as follows:

1. **Endpoint Creation:** An endpoint is to be created in Golang which can accept four parameters:
    - `Latitude`
    - `Longitude`
    - `Radius` (expressed in meters)
    - `Type` (The type could be a circle or a square)

2. **Spot Search:** The endpoint should utilise the received parameters to query all corresponding spots residing in the database (as defined in the "spots.sql" file).

3. **Result Sorting:** The returned results should be ordered based on their distance from a given spot. It's important to note a special requirement here, such that if the distance between two spots is less than 50 meters, the results should instead be ordered by the spots' respective ratings.

4. **Response:** The end response from the endpoint should be an array of objects where each object encapsulates all the fields present in the queried dataset.

### Solution File Structure
```
┌── .env (not included in repo, example in example.env) 
├── database
│   └── database.go (database connection)
├── go.mod
│   └── go.sum
├── handlers
│   └── handler.go (handler for /api/spots endpoint)
├── Main.go (entrypoint)
├── models
│   ├── nullstring.go (nullstring type for sql null values that can be marshalled to flat json)
│   └── spot.go (spot model)
└── spots.sql (sql file with provided spot data)
```
