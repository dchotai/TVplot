package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func writeCSV(episodes [][]Episode) {
	f, err := os.Create("episodes.csv")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	err := writer.Write([]string{"epNum", "season", "formatted", "title", "rating"})
	for _, season := range episodes {
		for _, ep := range season {
			data := []string{strconv.Itoa(ep.episodeNum), strconv.Itoa(ep.season), ep.formatted, ep.title, fmt.Sprintf("%.1f", ep.rating)}
			err = writer.Write(data)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
