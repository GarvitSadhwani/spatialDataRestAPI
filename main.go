package main

import (
	// "html/template"
	"database/sql"
	"fmt"
	"math"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	wkb "github.com/twpayne/go-geom/encoding/wkb"
)

type Template interface {
	Execute(w http.ResponseWriter, data interface{})
}

type Country struct {
	Name string
	Geo  wkb.Geom
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "my restAPI")
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func findn(cnt []float64) []string {
	var ans []string
	//sort.Float64s(cnt)
	db, err := sql.Open("pgx", "host=localhost port=5432 user=pixxeldb password=pixxeldb dbname=spatialdata sslmode=disable")
	if err != nil {
		fmt.Printf("error connecting to db")
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("error communicating with db")
	}
	defer db.Close()

	res, err := db.Query("select admin,ST_AsBinary(wkb_geometry) from spatialdatadb")

	if err != nil {
		fmt.Println("error running query")
	}
	defer res.Close()
	var cntry Country
	for res.Next() {
		err = res.Scan(&cntry.Name, &cntry.Geo)
		if err != nil {
			fmt.Println("error retrieving data from row")
		}
		var cntBoundary []float64
		for _, e := range cntry.Geo.FlatCoords() {
			cntBoundary = append(cntBoundary, toFixed(e, 1))
		}
		//sort.Float64s(cntBoundary)
		a := len(cntBoundary)
		b := len(cnt)
		if a > b {
			a = b
		}
		neigh := false
		for i := 0; i < a-1; i += 2 {
			for j := 0; j < b-1; j += 2 {
				if cnt[j] == cntBoundary[i] && cnt[j+1] == cntBoundary[i+1] {
					neigh = true
					break
				}
			}

		}
		if neigh {
			ans = append(ans, cntry.Name)
		}
	}
	return ans

}

func main() {
	rt := chi.NewRouter()
	db, err := sql.Open("pgx", "host=localhost port=5432 user=pixxeldb password=pixxeldb dbname=spatialdata sslmode=disable")
	if err != nil {
		fmt.Printf("error connecting to db")
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("error communicating with db")
	}
	defer db.Close()
	//ST_AsBinary
	res, err := db.Query("select admin,ST_AsBinary(wkb_geometry) from spatialdatadb where admin='India'")

	if err != nil {
		fmt.Println("error running query")
	}
	defer res.Close()
	var cntry Country
	for res.Next() {
		err = res.Scan(&cntry.Name, &cntry.Geo)
		if err != nil {
			fmt.Println("error retrieving data from row")
		}
	}
	var cntBoundary []float64
	for _, e := range cntry.Geo.FlatCoords() {
		cntBoundary = append(cntBoundary, toFixed(e, 1))
	}

	var neighbours []string
	neighbours = findn(cntBoundary)
	fmt.Printf("Neighbours: %v", neighbours)

	rt.Get("/", StaticHandler)
	rt.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "page not found")
	})
	fmt.Println("starting server on 8080")
	http.ListenAndServe(":8080", rt)
}
