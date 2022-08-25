# RestAPI for Spatial Data
A web API developed using golang to find neighbours for a country or search for a country name.
Data for countries taken from [datahub](https://datahub.io/core/geo-countries#resource-geo-countries_zip)

## Installation and Execution
-Clone the repository in the desired directory
    ```
    >git clone https://github.com/GarvitSadhwani/spatialDataRestAPI
    ```    
-Install [gdal](https://gdal.org/download.html) if not present
-Credentials to the postGIS database are present in ```docker-compose.yml``` file
-Start Docker
    ```
    >docker compose up
    ```
-Run the following command to import countries.geojason data to the PostGIS database
    ```
    >ogr2ogr -f "PostgreSQL" PG:"dbname=spatialdata user=pixxeldb password=pixxeldb" "countries.geojson" -nln spatialdatadb -append
    ```
 This creates a table 'spatialdatadb' in the PostGIS database which is used by the code
-Execute the ```main.go``` file
    ```
    >go run main.go
    ```
-The API will be hosted on localhost:8080


## Features
- Search for a country
    ![homepage](https://raw.githubusercontent.com/GarvitSadhwani/spatialDataRestAPI/main/templates/searchcountry.JPG)
    ![homepage](https://raw.githubusercontent.com/GarvitSadhwani/spatialDataRestAPI/main/templates/showcountry.JPG)
- Find neighbours for a country on the data list
    ![homepage](https://raw.githubusercontent.com/GarvitSadhwani/spatialDataRestAPI/main/templates/searchnghbr.JPG)
    ![homepage](https://raw.githubusercontent.com/GarvitSadhwani/spatialDataRestAPI/main/templates/shownghbr.JPG)
- Add or Delete any country
- More features coming soon!

## Dependencies
This app runs with the help of [go-chi](https://github.com/go-chi) and jack's [driver](https://github.com/jackc/pgx) for PostgreSQL

