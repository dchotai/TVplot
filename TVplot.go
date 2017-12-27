package main

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var TRAKT_API_KEY = "XXX"

func ShowQuery(q string) (t string, y float64, i string, n int) {
	client := &http.Client{}
	q = strings.TrimSpace(q)
	query := "https://api.trakt.tv/search/show?query=" + strings.Replace(q, " ", "%20", -1)
	req, _ := http.NewRequest("GET", query, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", TRAKT_API_KEY)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.Status, "Show info")

	var data []interface{}
	json.Unmarshal(resp_body, &data)

	if len(data) < 1 {
		log.Fatal("No shows found with given query")
	}

	title := data[0].(map[string]interface{})["show"].(map[string]interface{})["title"].(string)
	year := data[0].(map[string]interface{})["show"].(map[string]interface{})["year"].(float64)
	imdbID := data[0].(map[string]interface{})["show"].(map[string]interface{})["ids"].(map[string]interface{})["imdb"].(string)
	traktSlug := data[0].(map[string]interface{})["show"].(map[string]interface{})["ids"].(map[string]interface{})["slug"].(string)

	query = fmt.Sprintf("https://api.trakt.tv/shows/%s/seasons", traktSlug)
	req, _ = http.NewRequest("GET", query, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", TRAKT_API_KEY)
	resp, err = client.Do(req)

	if err != nil {
		log.Fatal("Errored when sending request to the server")
	}

	defer resp.Body.Close()
	resp_body, _ = ioutil.ReadAll(resp.Body)

	log.Println(resp.Status, "Number of seasons")
	json.Unmarshal(resp_body, &data)

	if len(data) < 1 {
		log.Fatal("No shows found with given slug")
	}

	numSeasons := len(data)

	if data[0].(map[string]interface{})["number"].(float64) == 0 {
		numSeasons -= 1
	}

	return title, year, imdbID, numSeasons
}

func DownloadPage(id string, year int) *html.Node {
	url := fmt.Sprintf("http://www.imdb.com/title/%s/episodes?season=%d", id, year)

	tempFile, err := os.Create("temp.html")
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()

	resp, err := http.Get(url)
	defer resp.Body.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("temp.html")
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(f)
	if err != nil {
		panic(err)
	}

	return doc
}

type Episode struct {
	season    int
	formatted string
	title     string
	rating    float64
}

func GetRatings() [][]Episode {
	title, year, imdbID, numSeasons := ShowQuery("BREAKING BAD")
	log.Println(title, year, imdbID, numSeasons)
	rv := make([][]Episode, numSeasons)

	for i := 0; i < numSeasons; i += 1 {
		doc := DownloadPage(imdbID, i+1)

		episodes := htmlquery.Find(doc, "//meta[@itemprop = 'episodeNumber']")

		rv[i] = make([]Episode, len(episodes))

		for j, n := range episodes {
			rating := htmlquery.Find(doc, "//div[@class = 'ipl-rating-star ']/span[@class='ipl-rating-star__rating']/text()")[j].Data
			episodeTitle := htmlquery.Find(doc, "//a[@itemprop='name']/text()")[j].Data
			episodeNum, _ := strconv.Atoi(htmlquery.SelectAttr(n, "content"))
			ratingFloat, _ := strconv.ParseFloat(rating, 64)
			rv[i][j] = Episode{i + 1, fmt.Sprintf("S%02dE%02d", i+1, episodeNum), episodeTitle, ratingFloat}

			// log.Printf("S%02dE%02d - %s - %s\n", i+1, episodeNum, episodeTitle, rating)
		}
	}
	return rv
	// strings.TrimSpace()
}

func main() {
	ratings := GetRatings()
	log.Println(ratings)
	// http.HandleFunc(pattern, handler)
}
