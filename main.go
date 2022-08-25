package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"math"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	wkb "github.com/twpayne/go-geom/encoding/wkb"
)

type Template struct {
	HTMLTpl *template.Template
}

type Country struct {
	Name string
	Geo  wkb.Geom
}

func StaticHandler(w http.ResponseWriter, file string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles(file)
	if err != nil {
		fmt.Printf("error parsing")
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Printf("error executing")
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func findn(refname string, cnt []float64) []string {
	var ans []string
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
		a := len(cntBoundary)
		b := len(cnt)
		neigh := false
		for i := 0; i < a-1; i += 2 {
			for j := 0; j < b-1; j += 2 {
				if cnt[j] == cntBoundary[i] && cnt[j+1] == cntBoundary[i+1] {
					neigh = true
					break
				}
			}
			if neigh {
				break
			}
		}
		if neigh && cntry.Name != refname {
			ans = append(ans, cntry.Name+string(','))
		}
	}
	return ans

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	StaticHandler(w, "home.gohtml")
}

func neighbourHandler(w http.ResponseWriter, r *http.Request) {
	StaticHandler(w, "neighbour.gohtml")
}

func neighbourUtil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	db, err := sql.Open("pgx", "host=localhost port=5432 user=pixxeldb password=pixxeldb dbname=spatialdata sslmode=disable")
	if err != nil {
		fmt.Printf("error connecting to db")
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("error communicating with db")
	}
	defer db.Close()
	res, err := db.Query("select admin,ST_AsBinary(wkb_geometry) from spatialdatadb where admin=$1", r.FormValue("country"))

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
	neighbours := findn(cntry.Name, cntBoundary)
	result := "The countries intersecting with " + r.FormValue("country") + " are:<br><br>"
	fmt.Fprint(w, result)
	fmt.Fprint(w, neighbours)
	fmt.Fprint(w, "<br><br><a href=\"/\">home</a>")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	StaticHandler(w, "search.gohtml")
}

func searchUtil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var result []string
	db, err := sql.Open("pgx", "host=localhost port=5432 user=pixxeldb password=pixxeldb dbname=spatialdata sslmode=disable")
	if err != nil {
		fmt.Printf("error connecting to db")
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("error communicating with db")
	}
	defer db.Close()
	per := '%'
	res, err := db.Query("select admin from spatialdatadb where admin like $1 ", r.FormValue("country")+string(per))

	if err != nil {
		fmt.Println("error running query")
	}
	defer res.Close()
	var temp string
	for res.Next() {
		err = res.Scan(&temp)
		if err != nil {
			fmt.Println("error retrieving data from row")
		}
		result = append(result, temp+string(','))
	}
	fmt.Fprint(w, "You might be looking for: <br><br>")
	fmt.Fprint(w, result)
	fmt.Fprint(w, "<br><br><a href=\"/\">home</a>")
}

func addCountryHandler(w http.ResponseWriter, r *http.Request) {
	StaticHandler(w, "newcountry.gohtml")
}

func newCountryUtil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	db, err := sql.Open("pgx", "host=localhost port=5432 user=pixxeldb password=pixxeldb dbname=spatialdata sslmode=disable")
	if err != nil {
		fmt.Printf("error connecting to db")
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("error communicating with db")
	}
	defer db.Close()
	qry := "insert into spatialdatadb (admin,iso_a3,wkb_geometry) values ('" + r.FormValue("name") + "','" + r.FormValue("abbr") + "',ST_AsBinary(ST_GeomFromText('POLYGON((" + r.FormValue("coord") + "))',4326)));"
	_, err = db.Exec(qry)
	if err != nil {
		fmt.Println("error running query")
	}

	fmt.Fprint(w, "Country added<br><br><a href=\"/\">home</a>")
}

func deleteCountryHandler(w http.ResponseWriter, r *http.Request) {
	StaticHandler(w, "removecountry.gohtml")
}

func deleteCountryUtil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	db, err := sql.Open("pgx", "host=localhost port=5432 user=pixxeldb password=pixxeldb dbname=spatialdata sslmode=disable")
	if err != nil {
		fmt.Printf("error connecting to db")
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("error communicating with db")
	}
	defer db.Close()
	_, err = db.Exec("delete from spatialdatadb where admin=$1", r.FormValue("name"))
	if err != nil {
		fmt.Println("error running query")
	}

	fmt.Fprint(w, "Country deleted<br><br><a href=\"/\">home</a>")
}

func main() {
	rt := chi.NewRouter()
	rt.Get("/", homeHandler)
	rt.Get("/searchcountry", searchHandler)
	rt.Get("/findneighbour", neighbourHandler)
	rt.Post("/search", searchUtil)
	rt.Post("/showneighbour", neighbourUtil)
	rt.Get("/addcountry", addCountryHandler)
	rt.Post("/newcountry", newCountryUtil)
	rt.Get("/deletecountry", deleteCountryHandler)
	rt.Post("/erasecountry", deleteCountryUtil)
	rt.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "page not found")
	})
	fmt.Println("starting server on 8080")
	http.ListenAndServe(":8080", rt)
}
