package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Episode struct {
	season    int
	formatted string
	title     string
	rating    float64
}

type View struct {
	Title string
	Year  string
	ID    string
	Error string
}

func homeHandler(rw http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("home.html")
	if err != nil {
		log.Println(err)
	}
	if newView.Error != "" {
		err = t.Execute(rw, newView)
		if err != nil {
			log.Println(err)
		}
	} else {
		err = t.Execute(rw, nil)
		if err != nil {
			log.Println(err)
		}
	}

}

var newView View

func viewHandler(rw http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("view.html")
	if err != nil {
		log.Println(err)
	}

	err = t.Execute(rw, newView)
	if err != nil {
		log.Println(err)
	}

}

func queryHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		if len(req.Form["query"]) < 1 {
			handleError(rw, req)
			return
		}
		title, year, imdbID, ratings := GetRatings(fmt.Sprintf("%s", req.Form["query"][0]))
		log.Println(title, year, imdbID, ratings)
		if title == "" {
			handleError(rw, req)
			return
		}
		newView = View{title, fmt.Sprintf("%d", int(year)), imdbID, ""}
		http.Redirect(rw, req, "/view/"+imdbID, http.StatusSeeOther)
	} else {
		handleError(rw, req)
	}
}

func handleError(rw http.ResponseWriter, req *http.Request) {
	newView = View{"", "", "", "Error: Could not find that show"}
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/query/", queryHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
