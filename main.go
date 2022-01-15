package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

// global variables and struct for data from form
var TextWeb, FontWeb, errResult string
var err int

type Errors struct {
	Number  int
	Message string
}

// func for getting data from form
func get(w http.ResponseWriter, r *http.Request) {
	FontWeb = r.FormValue("font")
	TextWeb = r.FormValue("text")
	http.Redirect(w, r, "http://localhost:4000/", http.StatusSeeOther)
}

// index page, if address != index, you are redirect to 404err func
func index(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	var result string
	result, err = textCreate()
	if tmplErr != nil {
		err = 404
		errResult = "This page is not exist"
		w.WriteHeader(http.StatusNotFound)
	}
	if r.URL.Path != "/" {
		err = 404
		errResult = "This page is not exist"
		err404(w, r)
		return
	} else if err != 200 {
		errResult = result
		err404(w, r)
		return
	} else {
		tmpl.ExecuteTemplate(w, "index", result)
	}
}

func err404(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("templates/404.html", "templates/header.html", "templates/footer.html")
	dataErr := Errors{err, errResult}

	if tmplErr != nil {
		err = 404
		errResult = "This page is not exist"
		w.WriteHeader(http.StatusNotFound)
	}
	if err == 404 {
		w.WriteHeader(http.StatusNotFound)
	} else if err == 400 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	tmpl.ExecuteTemplate(w, "404", dataErr)
	TextWeb = ""
	FontWeb = "standard.txt"

}

// routing func
func handleRequest() {
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/styles/"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./templates/static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/get", get)
	http.HandleFunc("/404", err404)

	log.Println("Server running on: http://localhost:4000")
	http.ListenAndServe(":4000", nil)
}

func main() {
	handleRequest()
}

func textCreate() (string, int) {
	var result string
	fontFile := "font/standard.txt"
	if len(FontWeb) != 0 {
		fontFile = "font/" + FontWeb
	}
	fContent, err := ioutil.ReadFile(fontFile)
	if err != nil {
		fmt.Println(err)
		return "Couldn't open file", 400
	}
	fontData := string(fContent)
	font := strings.Split(fontData, "\n")

	if len(TextWeb) == 0 {
		TextWeb = "Hello world!"
	}

	text := strings.Split(TextWeb, "\r\n")
	for n := 0; n < len(text); n++ { // words
		if len(text[n]) > 0 {
			for i := 1; i < 9; i++ { // lines in font file
				for j := 0; j < len(text[n]); j++ { // symbol
					if text[n][j] >= 32 && text[n][j] <= 126 {
						pos := int(byte(text[n][j])-32)*9 + i
						result += thinkertoy(font[pos])
					} else if text[n][j] == 13 || text[n][j] == 10 {
						continue
					} else {
						errResult = "Illegal symbol"
						return "Illegal symbol", 400
					}
				}

				result += "\n"
			}
		} else {
			result += "\n"
		}
	}
	return result, 200
}

// trim for thinkertoy font
func thinkertoy(line string) string {
	line = strings.TrimSuffix(line, "\a")
	line = strings.TrimSuffix(line, "\b")
	line = strings.TrimSuffix(line, "\t")
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\v")
	line = strings.TrimSuffix(line, "\f")
	line = strings.TrimSuffix(line, "\r")
	return line
}
