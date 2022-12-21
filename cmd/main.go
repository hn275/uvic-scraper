package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hn275/uvic-scraper/courses"
	"github.com/hn275/uvic-scraper/util"
)

type course struct {
	subject string
	course  string
	term    int
}

var data []courses.ClassInfo
var wg sync.WaitGroup
var count int = 0

const TERM = 202301

func main() {
	log.Printf("%s\n", util.LogColor("GREEN", "Fetching all courses"))
	allCourses, err := courses.GetAllCourses()
	if err != nil {
		log.Printf("%s: %v\n", util.LogColor("RED", "FAILED"), err)
	}

	// fetchCount := len(allCourses)
	fetchCount := 10

	for index, course := range allCourses {
		if index >= fetchCount {
			break
		}

		wg.Add(1)
		time.Sleep(time.Millisecond * 25) // otherwise requests get aired out

		log.Printf("%s %s%s\n", util.LogColor("YELLOW", "Fetching"), course.Subject, course.Course)

		/* Make request then send fetched data */
		go func(a courses.Class, i int) {
			defer wg.Done()

			courseInfo, err := courses.GetCourseInfo(a.Subject, a.Course, TERM)
			if err != nil {
				log.Printf("%s: %s%s:\n%v\n", util.LogColor("RED", "FAILED"), a.Subject, a.Course, err)
				return
			}

			log.Printf("%s: %s%s\t(%d/%d)\n", util.LogColor("GREEN", "PARSED"), a.Subject, a.Course, i+1, fetchCount)
			if len(courseInfo.Location) != 0 {
				data = append(data, courseInfo)
				count++
			}

		}(course, index)

	}

	wg.Wait()

	fmt.Println(util.LogColor("GREEN", fmt.Sprintf("PARSED %d RECORDS", count)))
	fmt.Printf(util.LogColor("RED", fmt.Sprintf("SKIPPED %d RECORDS\n", fetchCount-count)))

	s, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("./data.json", s, 0666); err != nil {
		fmt.Printf("%s: %v\n", util.LogColor("RED", "FAILED"), err)
	}

	fmt.Println("Done")
}
