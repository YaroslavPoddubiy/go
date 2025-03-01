package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

func erf(x float64, terms int) float64 {
	result := 0.0
	for n := 0; n < terms; n++ {
		term := math.Pow(-1, float64(n)) * math.Pow(x, float64(2*n+1)) / factorial(n) / float64(2*n+1)
		result += term
	}
	return (2 / math.Sqrt(math.Pi)) * result
}

func factorial(n int) float64 {
	if n == 0 {
		return 1.0
	}
	return float64(n) * factorial(n-1)
}

func noImbalance(dailyPower, stdDev float64) float64 {
	error := 0.05 * dailyPower
	noImbalanceCoef := (erf(error/(stdDev*math.Sqrt(2.0)), 100)/2 - erf(-error/(stdDev*math.Sqrt(2.0)), 100)/2)
	return math.Round(noImbalanceCoef*100) / 100
}

func calculateProfit(dailyPower, stdDeviation, energyCost float64) float64 {
	noImbalanceCoef := noImbalance(dailyPower, stdDeviation)
	profit := dailyPower * energyCost * 24.0 * noImbalanceCoef
	tax := dailyPower * energyCost * 24.0 * (1.0 - noImbalanceCoef)
	return profit - tax
}

type ViewData struct {
	Result    string
	Power     string
	Deviation string
	Price     string
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

		power, err := strconv.ParseFloat(r.FormValue("power"), 64)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		powerToShow := r.FormValue("power")

		deviation, err := strconv.ParseFloat(r.FormValue("deviation"), 64)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		deviationToShow := r.FormValue("deviation")

		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		priceToShow := r.FormValue("price")

		result := calculateProfit(power, deviation, price*1000)

		data := ViewData{
			Result:    "Очікуваний прибуток: ₴" + fmt.Sprintf("%.2f", result),
			Power:     powerToShow,
			Deviation: deviationToShow,
			Price:     priceToShow,
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
