package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)



type General struct {
	Nombre string
	Materia string
	Calificacion float64
}

type Materia struct {
	Alumno string
	calificacion float64
}

type Alumno struct {
	Materia string
	calificacion float64
}

var materias map [ string ][]Materia
var alumnos map [ string ][]Alumno


type Server struct{}

func ValidarNoRepetido(nombre string , materia string ) string {
	for mat , als := range materias {
		if mat == materia {
			for _ , al := range als {
				if al.Alumno == nombre {
					return "Ya existe esta calificacion"
				}
			}
		}
	}
	return ""
}

func NuevoAlumno (res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
		print("50")
	}
	nombre := req.FormValue("alumn")
	mat := req.FormValue("asignature")
	calificacion, err := strconv.ParseFloat(req.FormValue("calification"), 32)
	if err != nil {
		return
	}
	validar := ValidarNoRepetido(nombre, mat)
	if validar != ""{
		print("59")
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			iniciarHTML("confirmacion.html"),
			validar,
		)
	} else {
		alumnos[nombre] = append(alumnos[nombre], Alumno{mat, calificacion})
		materias[mat] = append(materias[mat], Materia{nombre, calificacion})
		c :=  nombre + " agregado a  " + mat 
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			iniciarHTML("confirmacion.html"),
			c,
		)
	}
}


func PromedioPorAlumno (res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
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
		cadenaFinal :=  nombre + " tiene: " + cal
		fmt.Fprintf(
			res,
			iniciarHTML("promedio_alumno.html"),
			cadenaFinal,
		)
	}
}

func PromedioPorMateria(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
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
		cadenaFinal := "EL promedio de " + mat + " es: " + strconv.FormatFloat(reply, 'f', 3, 64)
		fmt.Fprintf(
			res,
			iniciarHTML("materia.html"),
			cadenaFinal,
		)
	}
}


func PromedioGeneral ( res http.ResponseWriter, req *http.Request) {
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
		iniciarHTML("promedio_general.html"),
		cadenaFinal,
	)
}

func iniciarHTML(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func iniciar(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		iniciarHTML("general.html"),
	)
}

func main(){

	materias = make(map[string][]Materia)
	alumnos = make(map[string][]Alumno)
	http.HandleFunc("/", iniciar)
	http.HandleFunc("/alumno", NuevoAlumno)
	http.HandleFunc("/materia", PromedioPorMateria)
	http.HandleFunc("/promedio_alumno", PromedioPorAlumno)
	http.HandleFunc("/promedio_general", PromedioGeneral)
	http.ListenAndServe(":9000", nil)

	var input string
	fmt.Scanln(&input)
}
