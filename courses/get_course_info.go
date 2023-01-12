package courses

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func GetCourseInfo(subject string, course, term int) (ClassInfo, error) {
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
			classType := trimspace(sel.FindNodes(childNodes[5]).Text())
			if classType == "Lab" {
				return
			}

			class.Time = trimspace(sel.FindNodes(childNodes[1]).Text())
			class.Weekday = trimspace(sel.FindNodes(childNodes[2]).Text())
			class.Building = trimspace(sel.FindNodes(childNodes[3]).Text())
			class.DateRange = trimspace(sel.FindNodes(childNodes[4]).Text())

			s.Location = append(s.Location, class)
			class = ClassLocation{}
		}
	})

	c.Visit(url(subject, course, term))

	return s, e
}

func trimspace(s string) string {
	return strings.TrimSpace(s)
}

func url(subject string, course, term int) string {
	/*
		https://www.uvic.ca/BAN1P/bwckctlg.p_disp_listcrse
		?term_in=202209
		&subj_in=STAT
		&crse_in=260
		&schd_in=
	*/
	baseUrl := "https://www.uvic.ca/BAN1P/bwckctlg.p_disp_listcrse"
	visitingUrl := fmt.Sprintf(
		"%s?term_in=%d&subj_in=%s&crse_in=%d&schd_in=",
		baseUrl,
		term,
		subject,
		course,
	)
	return visitingUrl
}
