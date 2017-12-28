package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Episode struct {
	EpisodeNum int
	Season     int
	Formatted  string
	Title      string
	Rating     float64
}

type View struct {
	Title    string
	Year     string
	ID       string
	Error    string
	Episodes [][]Episode
}

func homeHandler(rw http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/"):]
	idExists, err := regexp.MatchString("^tt\\d+", id)
	if idExists {
		viewHandler(rw, req, id)
		return
	}

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

func viewHandler(rw http.ResponseWriter, req *http.Request, id string) {
	if newView.ID != id {
		title, year, imdbID, episodes := GetRatings(id)
		log.Println(title, year, imdbID, episodes)
		if title == "" {
			handleError(rw, req)
			return
		}
		newView = View{title, fmt.Sprintf("%d", int(year)), imdbID, "", episodes}
		http.Redirect(rw, req, "/"+imdbID, http.StatusSeeOther)
	}

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
		title, year, imdbID, episodes := GetRatings(fmt.Sprintf("%s", req.Form["query"][0]))
		log.Println(title, year, imdbID, episodes)
		if title == "" {
			handleError(rw, req)
			return
		}
		newView = View{title, fmt.Sprintf("%d", int(year)), imdbID, "", episodes}
		http.Redirect(rw, req, "/"+imdbID, http.StatusSeeOther)
	} else {
		handleError(rw, req)
	}
}

func handleError(rw http.ResponseWriter, req *http.Request) {
	newView = View{"", "", "", "Error: Could not find that show", nil}
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/query/", queryHandler)
	port := os.Getenv("PORT")
	if port != "" {
		port = ":8000"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
