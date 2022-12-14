package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// var mockdata []courses.Class

func init() {
	// In case I need this later.
	// courses, err := ioutil.ReadFile("./all_uvic_courses.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := json.Unmarshal(courses, &mockdata); err != nil {
	// 	log.Fatal(err)
	// }
}

const (
	SUBJECT = "CSC"
	COURSE  = "116"
	TERM    = 202209
)

type ClassInfo struct {
	Weekday   string `json:"weekday"`
	Time      string `json:"time"`
	Building  string `json:"building"`
	DateRange string `json:"date_range"`
}

var class ClassInfo

var allClasses []ClassInfo

func main() {
	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("error fetching url %v:\n%s", r.Request.URL, err.Error())
	})

	c.OnHTML("table tbody tr:nth-child(2), table tbody tr th.ddtitle", func(h *colly.HTMLElement) {

		sel := h.DOM
		childNodes := sel.Children().Nodes

		if len(childNodes) == 7 {
			class.Time = trimspace(sel.FindNodes(childNodes[1]).Text())
			class.Weekday = trimspace(sel.FindNodes(childNodes[2]).Text())
			class.Building = trimspace(sel.FindNodes(childNodes[3]).Text())
			class.DateRange = trimspace(sel.FindNodes(childNodes[4]).Text())

			allClasses = append(allClasses, class)
			class = ClassInfo{}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(allClasses)
		output, err := json.MarshalIndent(&allClasses, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(output))
	})

	c.Visit(url(SUBJECT, COURSE, TERM))
}

func trimspace(s string) string {
	return strings.TrimSpace(s)
}

func url(subject, course string, term int) string {
	/*
		https://www.uvic.ca/BAN1P/bwckctlg.p_disp_listcrse
		?term_in=202209
		&subj_in=STAT
		&crse_in=260
		&schd_in=
	*/
	baseUrl := "https://www.uvic.ca/BAN1P/bwckctlg.p_disp_listcrse"
	visitingUrl := fmt.Sprintf(
		"%s?term_in=%d&subj_in=%s&crse_in=%s&schd_in=",
		baseUrl,
		term,
		subject,
		course,
	)
	return visitingUrl
}
