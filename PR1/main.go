package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Task1ViewData struct {
	DryComponents     string
	AshFreeComponents string
	WorkWormth        string
	DryWormth         string
	AshFreeWormth     string
}

type Task2ViewData struct {
	WorkComponents string
	WorkWormth     string
}

func getText(title string, components []float64, components_names []string) string {
	var text string = title + "\n"
	for i := 0; i < len(components); i++ {
		text = text + components_names[i] + ":" + fmt.Sprintf("%.3f", components[i])
		if components_names[i] == "Ванадій" {
			text += "мг/кг\n"
		} else {
			text += "%\n"
		}
	}
	return text
}

func getRdCoef(w float64) float64 {
	return 100 / (100 - w)
}

func getRafCoef(w, a float64) float64 {
	return 100 / (100 - w - a)
}

func getDryComponents(workComponents []float64) []float64 {
	coef := getRdCoef(workComponents[5])
	dryComponents := make([]float64, len(workComponents))
	copy(dryComponents, workComponents)
	dryComponents[5] = 0.0
	for i := range dryComponents {
		dryComponents[i] *= coef
	}
	return dryComponents
}

func getAshFreeComponents(workComponents []float64) []float64 {
	coef := getRafCoef(workComponents[5], workComponents[6])
	ashFreeComponents := make([]float64, len(workComponents))
	copy(ashFreeComponents, workComponents)

	ashFreeComponents[5] = 0.0
	ashFreeComponents[6] = 0.0
	for i := range ashFreeComponents {
		ashFreeComponents[i] *= coef
	}
	return ashFreeComponents
}

func getWormth(components []float64) float64 {
	return 339*components[0] + 1030*components[1] - 108.8*(components[4]-components[2]) -
		25*components[5]
}

func getDryWormth(wormth, w float64) float64 {
	return (wormth + 0.025*w) * (100 / (100 - w))
}

func getAfWormth(wormth, w, a float64) float64 {
	return (wormth + 0.025*w) * (100 / (100 - w - a))
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

		values := [7]string{
			r.FormValue("carbon"),
			r.FormValue("hydrogen"),
			r.FormValue("sulfur"),
			r.FormValue("nitrogen"),
			r.FormValue("oxygen"),
			r.FormValue("moisture"),
			r.FormValue("ash"),
		}

		var components [7]float64
		for i, v := range values {
			num, err := strconv.ParseFloat(v, 64)
			if err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			components[i] = num
		}

		dryComponents := getDryComponents(components[:])
		ashFreeComponents := getAshFreeComponents(components[:])
		workWormth := getWormth(components[:]) / 1000
		dryWormth := getDryWormth(workWormth, components[5])
		ashFreeWormth := getAfWormth(workWormth, components[5], components[6])

		data := Task1ViewData{
			DryComponents:     getText("Склад сухого палива", dryComponents, []string{"Вуглець", "Водень", "Сірка", "Азот", "Кисень", "Волога", "Зола"}),
			AshFreeComponents: getText("Склад горючого палива", ashFreeComponents, []string{"Вуглець", "Водень", "Сірка", "Азот", "Кисень", "Волога", "Зола"}),
			WorkWormth:        "Нижча теплота згоряння робочого складу палива: " + fmt.Sprintf("%.3f", workWormth) + "МДж/кг",
			DryWormth:         "Нижча теплота згоряння сухого складу палива: " + fmt.Sprintf("%.3f", dryWormth) + "МДж/кг",
			AshFreeWormth:     "Нижча теплота згоряння горючого складу палива: " + fmt.Sprintf("%.3f", ashFreeWormth) + "МДж/кг",
		}

		tmpl, _ := template.ParseFiles("static/task1.html")
		tmpl.Execute(w, data)
	}
}

func getAfrCoef(w, a float64) float64 {
	return (100 - w - a) / 100
}

func getWorkComponents(components []float64) []float64 {
	coef := getAfrCoef(components[6], components[5])
	workComponents := make([]float64, len(components))
	copy(workComponents, components)

	for i := range workComponents {
		workComponents[i] = coef * workComponents[i]
	}

	workComponents[5] = components[5]
	workComponents[6] = ((100 - components[5]) / 100) * components[6]
	return workComponents
}

func getWorkWormth(wormth, w, a float64) float64 {
	return (wormth * (100 - w - a) / 100) - (0.025 * w)
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

		values := [8]string{
			r.FormValue("carbon"),
			r.FormValue("hydrogen"),
			r.FormValue("sulfur"),
			r.FormValue("vanadium"),
			r.FormValue("oxygen"),
			r.FormValue("moisture"),
			r.FormValue("ash"),
			r.FormValue("wormth"),
		}

		var components [8]float64
		for i, v := range values {
			num, err := strconv.ParseFloat(v, 64)
			if err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			components[i] = num
		}

		workComponents := getWorkComponents(components[:])
		workWormth := getWorkWormth(components[7], components[5], components[6])

		data := Task2ViewData{
			WorkComponents: getText("Склад сухого палива", workComponents[:len(workComponents)-1], []string{"Вуглець", "Водень", "Сірка", "Ванадій", "Кисень", "Волога", "Зола"}),
			WorkWormth:     "Нижча теплота згоряння робочого складу палива: " + fmt.Sprintf("%.3f", workWormth) + "МДж/кг",
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
