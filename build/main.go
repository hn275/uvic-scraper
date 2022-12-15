package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/hn275/uvic-scraper/courses"
)

type course struct {
	subject string
	course  string
	term    int
}

var ch = make(chan courses.ClassInfo)
var data []courses.ClassInfo
var wg sync.WaitGroup

const TERM = 202301

/* NOTE: Don't run yet, run the shape by the team to make sure it's good
func main() {
	allCourses, err := courses.GetAllCourses()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range allCourses {
		wg.Add(1)
		go func(ch chan<- courses.ClassInfo, a courses.Class) {
			fmt.Printf("Scraping %s%s...\n", a.Subject, a.Course)

			s, err := courses.GetCourseInfo(a.Subject, a.Course, TERM)
			if err != nil {
				log.Printf("failed to parsed %s%s:%v\n", a.Subject, a.Course, err)
			}
			ch <- s

		}(ch, v)

		go func(ch <-chan courses.ClassInfo) {
			defer wg.Done()

			courseInfo := <-ch
			if len(courseInfo.Location) != 0 {
				data = append(data, courseInfo)
			}
		}(ch)
	}

	wg.Wait()
	close(ch)

	s, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("./data.json", s, 0666); err != nil {
		log.Fatal(err)
	}
}
*/
