# RestAPI for Spatial Data
A web API developed using golang to find neighbours for a country or search for a country name.
Data for countries taken from [datahub](https://datahub.io/core/geo-countries#resource-geo-countries_zip)


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

