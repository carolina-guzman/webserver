package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var materias map[string][]materia
var alumnos map[string][]alumno

type materia struct {
	Alumno       string
	calificacion float64
}

type alumno struct {
	Materia      string
	calificacion float64
}

type Info struct {
	Nombre       string
	Materia      string
	Calificacion float64
}

func agregaAlumno(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
		// return "Error al obtener los datos"
	}
	nombre := req.FormValue("alumn")
	mat := req.FormValue("asignature")
	calificacion, err := strconv.ParseFloat(req.FormValue("calification"), 32)
	if err != nil {
		return
	}
	err2 := detectCal(nombre, mat)
	if err2 != "" {
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("response.html"),
			err2,
		)
	} else {
		alumnos[nombre] = append(alumnos[nombre], alumno{mat, calificacion})
		materias[mat] = append(materias[mat], materia{nombre, calificacion})

		// falta retornar el HTML con el mensaje:
		c := "El alumno " + nombre + " Ha sido asignado a  " + mat + " con éxito."
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHTML("response.html"),
			c,
		)
	}
}

func detectCal(nombre string, materia string) string {
	for mat, als := range materias {
		if mat == materia {
			for _, al := range als {
				if al.Alumno == nombre {
					return "La calificacion ya fue dada de alta"
				}
			}
		}
	}
	return ""
}

func promedioAlumno(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
		// return "Error al obtener los datos"
	}
	nombre := req.FormValue("calification_prom")
	reply := 0.0
	var counter float64
	counter = 0.0
	_, verification := alumnos[nombre]
	if verification {
		for _, mat := range alumnos[nombre] {
			reply += mat.calificacion
			counter++
		}
		reply = reply / counter
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		cal := strconv.FormatFloat(reply, 'f', 3, 32)
		cadenaFinal := "El promedio de: " + nombre + " es: " + cal
		fmt.Fprintf(
			res,
			cargarHTML("promedio.html"),
			cadenaFinal,
		)
	}
}

func promedioGeneral(res http.ResponseWriter, req *http.Request) {
	reply := 0.0
	var total float64
	var counter float64
	var mat float64
	for name := range alumnos {
		counter = 0
		mat = 0
		for _, alumno := range alumnos[name] {
			counter += alumno.calificacion
			mat++
		}
		total += counter / mat
	}
	reply = total / float64(len(alumnos))
	replyS := strconv.FormatFloat(reply, 'f', 3, 64)
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	cadenaFinal := "El promedio general es: " + replyS
	fmt.Fprintf(
		res,
		cargarHTML("promedio.html"),
		cadenaFinal,
	)
}

func promedioMateria(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
		// return "Error al obtener los datos"
	}
	mat := req.FormValue("asignature_prom")
	reply := 0.0
	var counter float64
	counter = 0.0
	_, verification := materias[mat]
	if verification {
		for _, mat := range materias[mat] {
			reply += mat.calificacion
			counter++
		}
		reply = reply / counter
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		cadenaFinal := "EL promedio de: " + mat + " es: " + strconv.FormatFloat(reply, 'f', 3, 64)
		fmt.Fprintf(
			res,
			cargarHTML("materia.html"),
			//cadena a enviar
			cadenaFinal,
		)
	}
}

func cargarHTML(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func general(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHTML("general.html"),
	)
}

func main() {
	materias = make(map[string][]materia)
	alumnos = make(map[string][]alumno)
	fmt.Println("Corriendo servirdor de Alumnos...")
	http.HandleFunc("/", general)
	http.HandleFunc("/alumno", agregaAlumno)
	http.HandleFunc("/materia", promedioMateria)
	http.HandleFunc("/promedio", promedioAlumno)
	http.HandleFunc("/promedio_general", promedioGeneral)
	http.ListenAndServe(":9000", nil)

}
