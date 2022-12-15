package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hn275/uvic-scraper/courses"
)

const (
	SUBJECT = "CSC"
	COURSE  = 116
	TERM    = 202209
)

func main() {
	classes, err := courses.GetCourseInfo(SUBJECT, COURSE, TERM)

	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.MarshalIndent(&classes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
