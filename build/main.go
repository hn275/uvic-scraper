package main

import (
	"encoding/json"
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

var ch = make(chan courses.ClassInfo)
var data []courses.ClassInfo
var wg sync.WaitGroup

const TERM = 202301

func main() {
	log.Printf("%s\n", util.LogColor("GREEN", "Fetching all courses"))
	allCourses, err := courses.GetAllCourses()
	if err != nil {
		log.Printf("%s: %v\n", util.LogColor("RED", "FAILED"), err)
	}

	total := len(allCourses)

	for i, v := range allCourses {
		if i > 10 {
			break
		}
		wg.Add(1)
		time.Sleep(time.Millisecond * 25) // otherwise uvic server can't handle all the process at the same time :/
		log.Printf("%s %s%s\n", util.LogColor("YELLOW", "Fetching"), v.Subject, v.Course)
		go func(ch chan<- courses.ClassInfo, a courses.Class, i int) {
			s, err := courses.GetCourseInfo(a.Subject, a.Course, TERM)
			if err != nil {
				log.Printf("%s: %s%s:\n%v\n", util.LogColor("RED", "FAILED"), a.Subject, a.Course, err)
			} else {
				log.Printf("%s: %s%s\t(%d/%d)\n", util.LogColor("GREEN", "SUCCESS"), v.Subject, v.Course, i, total)
			}
			ch <- s

		}(ch, v, i)

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

	log.Printf("%s, generating json\n", util.LogColor("GREEN", "Fetch success"))

	s, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("./data.json", s, 0666); err != nil {
		log.Printf("%s: %v\n", util.LogColor("RED", "FAILED"), err)
	}

	log.Printf("%s\n", util.LogColor("GREEN", "DONE"))
}
