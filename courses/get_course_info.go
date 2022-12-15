package courses

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func GetCourseInfo(subject, course string, term int) (ClassInfo, error) {
	log.Printf("Fetching %s%s\n", subject, course)
	s := ClassInfo{
		Subject:  subject,
		Course:   course,
		Term:     term,
		Location: []ClassLocation{},
	}
	var e error = nil

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		e = err
	})

	c.OnHTML("table tbody tr:nth-child(2), table tbody tr th.ddtitle", func(h *colly.HTMLElement) {
		var class ClassLocation

		sel := h.DOM
		childNodes := sel.Children().Nodes

		if len(childNodes) == 7 {
			class.Time = trimspace(sel.FindNodes(childNodes[1]).Text())
			class.Weekday = trimspace(sel.FindNodes(childNodes[2]).Text())
			class.DateRange = trimspace(sel.FindNodes(childNodes[4]).Text())

			/* PARSING BUILDING AND ROOM */
			b := trimspace(sel.FindNodes(childNodes[3]).Text())
			bSplit := strings.Split(b, " ")

			room := bSplit[len(bSplit)-1:]
			building := bSplit[:len(bSplit)-1]

			class.Room = room[0]
			class.Building = strings.Join(building, " ")

			s.Location = append(s.Location, class)
			class = ClassLocation{} // reset class for the next table
		}
	})

	c.OnScraped(func(_ *colly.Response) {
		log.Printf("Sucessfully scraped %s%s\n", subject, course)
	})

	c.Visit(url(subject, course, term))

	return s, e
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
