package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Task1ViewData struct {
	Voltage       string
	Current       string
	FaultDuration string
	Load          string
	Result        string
}

type Task2ViewData struct {
	Power  string
	Result string
}

type Task3ViewData struct {
	ModeIndex int64
	Result    string
}

func calculate(voltage, shortCircuitCurrent, faultDuration, load float64) float64 {
	im := (load / 2.0) / (math.Sqrt(3.0) * voltage)
	sek := im / 1.4
	smin := (shortCircuitCurrent * 1000 * math.Sqrt(faultDuration)) / 92.0
	if smin <= sek {
		smin = sek
	}
	return math.Ceil(smin/10) * 10
}

func handleTask1(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("static/task1.html")
		tmpl.Execute(w, Task1ViewData{})
	} else if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		voltageToShow := r.FormValue("voltage")
		currentToShow := r.FormValue("current")
		faultDurationToShow := r.FormValue("faultDuration")
		loadToShow := r.FormValue("load")

		voltage, _ := strconv.ParseFloat(voltageToShow, 64)
		current, _ := strconv.ParseFloat(currentToShow, 64)
		faultDuration, _ := strconv.ParseFloat(faultDurationToShow, 64)
		load, _ := strconv.ParseFloat(loadToShow, 64)

		result := calculate(voltage, current, faultDuration, load)

		data := Task1ViewData{
			voltageToShow,
			currentToShow,
			faultDurationToShow,
			loadToShow,
			"Рекомендований переріз кабеля: " + fmt.Sprintf("%.2f", result) + "мм^2",
		}

		tmpl, _ := template.ParseFiles("static/task1.html")
		tmpl.Execute(w, data)
	}
}

func calculatePower(power float64) float64 {
	return 10.5 / (math.Sqrt(3.0) * (math.Pow(10.5, 2)/power + math.Pow(10.5, 3)/630))
}

func handleTask2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("static/task2.html")
		tmpl.Execute(w, Task2ViewData{})
	} else if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		powerToShow := r.FormValue("power")

		power, _ := strconv.ParseFloat(powerToShow, 64)

		result := calculatePower(power)

		data := Task2ViewData{
			powerToShow,
			"Струм КЗ на шинах 10 кВ ГПП: " + fmt.Sprintf("%.2f", result) + "кА",
		}

		tmpl, _ := template.ParseFiles("static/task2.html")
		tmpl.Execute(w, data)
	}
}

func calculateCurrency(rResistance, xResistance float64, phases int) float64 {
	rc := rResistance * 0.009
	x := (xResistance + 233) * 0.009
	z := math.Sqrt(math.Pow(rc, 2) + math.Pow(x, 2))
	i := 11000 / z
	if phases == 3 {
		return i / math.Sqrt(3.0)
	} else if phases == 2 {
		return i / 2
	}
	return 0.0
}

func getResultText(mode string) string {
	var result string

	switch mode {
	case "2":
		result = "Підстанція не має аварійного режиму"
	case "0":
		result += "Струм трифазного КЗ в нормальному режимі: "
		result += fmt.Sprintf("%.2fА\n", calculateCurrency(10.65, 24.02, 3))
		result += "Струм двофазного КЗ в нормальному режимі: "
		result += fmt.Sprintf("%.2fА", calculateCurrency(10.65, 24.02, 2))
	case "1":
		result += "Струм трифазного КЗ в мінімальному режимі: "
		result += fmt.Sprintf("%.2fА\n", calculateCurrency(34.88, 65.68, 3))
		result += "Струм двофазного КЗ мінімальному режимі: "
		result += fmt.Sprintf("%.2fА", calculateCurrency(34.88, 65.68, 2))
	}
	return result
}

func handleTask3(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("static/task3.html")
		tmpl.Execute(w, Task3ViewData{})
	} else if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		mode := r.FormValue("mode")
		modeIndex, _ := strconv.ParseInt(mode, 10, 64)

		result := getResultText(mode)

		data := Task3ViewData{
			modeIndex,
			result,
		}

		tmpl, _ := template.ParseFiles("static/task3.html")
		tmpl.Execute(w, data)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")

}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleRoot)

	http.HandleFunc("/task1", handleTask1)

	http.HandleFunc("/task2", handleTask2)

	http.HandleFunc("/task3", handleTask3)

	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe("localhost:8000", nil)
}
