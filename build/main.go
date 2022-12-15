package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

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

func main() {
	allCourses, err := courses.GetAllCourses()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range allCourses {
		wg.Add(1)
		time.Sleep(time.Millisecond * 50) // otherwise uvic server can't handle all the process at the same time :/
		go func(ch chan<- courses.ClassInfo, a courses.Class) {
			s, err := courses.GetCourseInfo(a.Subject, a.Course, TERM)
			if err != nil {
				log.Printf("failed to parsed %s%s:\n%v\n", a.Subject, a.Course, err)
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
