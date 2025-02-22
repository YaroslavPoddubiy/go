package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var fuelWormth = []float64{20.47, 40.4, 33.08}
var fuelAsh = []float64{25.2, 0.15, 0.0}
var fuelAshCoef = []float64{0.8, 1.0, 0.0}
var combustibleSubstance = []float64{1.5, 0.0, 0.0}

func getEmissionIndicator(itemId int64) float64 {
	return (1000000 / fuelWormth[itemId]) * fuelAshCoef[itemId] *
		(fuelAsh[itemId] / (100 - combustibleSubstance[itemId])) * 0.015
}

func getEmission(itemId int64, fuelQuantity float64) float64 {
	emission := getEmissionIndicator(itemId)
	return 0.000001 * emission * fuelWormth[itemId] * fuelQuantity
}

type ViewData struct {
	Result       string
	FuelIndex    int64
	FuelQuantity string
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("static/index.html")
		tmpl.Execute(w, ViewData{})
	} else if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		fuelIndex, err := strconv.ParseInt(r.FormValue("fuel"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		fuelQuantityToShow := r.FormValue("quantity")

		fuelQuantity, err := strconv.ParseFloat(r.FormValue("quantity"), 64)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		result := getEmission(fuelIndex, fuelQuantity)

		data := ViewData{
			Result:       "Валовий викид дорівнює " + fmt.Sprintf("%.3f", result) + "т",
			FuelIndex:    fuelIndex,
			FuelQuantity: fuelQuantityToShow,
		}

		tmpl, _ := template.ParseFiles("static/index.html")
		tmpl.Execute(w, data)

	}

}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleRoot)

	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe("localhost:8000", nil)
}
