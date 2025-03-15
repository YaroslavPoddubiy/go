package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Task1ViewData struct {
	Length string
	Result string
}

type Task2ViewData struct {
	RefuseFrequency string
	RepairTime      string
	PlanDownTime    string
	EmergencyLoss   string
	PlanLoss        string
	Result          string
}

func calculateReliability(lineLength float64) float64 {
	lineReliability := 0.007
	lineRepairTime := 10.0
	lineComponentsReliability := []float64{0.01, 0.015, 0.02, 0.03}
	lineComponentsQuantity := []int{1, 1, 1, 6}
	lineComponentsRepairTime := []float64{30, 100, 15, 2}

	refuseFrequency := lineReliability * lineLength
	averageRepairTime := refuseFrequency * lineRepairTime

	for i := range lineComponentsReliability {
		refuseFrequency += lineComponentsReliability[i] * float64(lineComponentsQuantity[i])
		averageRepairTime += lineComponentsReliability[i] * float64(lineComponentsQuantity[i]) * lineComponentsRepairTime[i]
	}

	averageRepairTime /= refuseFrequency
	emergencyDownTimeCoefficient := refuseFrequency * averageRepairTime / 8760
	planDownTimeCoefficient := 1.2 * 43 / 8760
	doubleCircleRefuseFrequency := 2*refuseFrequency*(emergencyDownTimeCoefficient+planDownTimeCoefficient) + 0.02

	return refuseFrequency / doubleCircleRefuseFrequency
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

		lengthToShow := r.FormValue("length")

		length, _ := strconv.ParseFloat(lengthToShow, 64)

		result := calculateReliability(length)

		data := Task1ViewData{
			lengthToShow,
			"Двоколова система надійніша за одноколову в " + fmt.Sprintf("%.2f", result) + "разів",
		}

		tmpl, _ := template.ParseFiles("static/task1.html")
		tmpl.Execute(w, data)
	}
}

func calculateLoss(refuseFrequency, repairTime, planDownTime, emergencyLoss, planLoss float64) float64 {
	mEmergency := refuseFrequency * repairTime * 5120 * 6451
	mPlan := planDownTime * 5120 * 6451
	m := mEmergency*emergencyLoss + mPlan*planLoss
	return m
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

		refuseFrequencyToShow := r.FormValue("refuseFrequency")
		repairTimeToShow := r.FormValue("repairTime")
		planDownTimeToShow := r.FormValue("planDownTime")
		emergencyLossToShow := r.FormValue("emergencyLoss")
		planLossToShow := r.FormValue("planLoss")

		refuseFrequency, _ := strconv.ParseFloat(refuseFrequencyToShow, 64)
		repairTime, _ := strconv.ParseFloat(repairTimeToShow, 64)
		planDownTime, _ := strconv.ParseFloat(planDownTimeToShow, 64)
		emergencyLoss, _ := strconv.ParseFloat(emergencyLossToShow, 64)
		planLoss, _ := strconv.ParseFloat(planLossToShow, 64)

		result := calculateLoss(refuseFrequency, repairTime, planDownTime, emergencyLoss, planLoss)

		data := Task2ViewData{
			refuseFrequencyToShow,
			repairTimeToShow,
			planDownTimeToShow,
			emergencyLossToShow,
			planLossToShow,
			"Збитки від переривання електропостачання становитимуть: ₴" + fmt.Sprintf("%.2f", result),
		}

		tmpl, _ := template.ParseFiles("static/task2.html")
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

	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe("localhost:8000", nil)
}
