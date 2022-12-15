package courses

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

func GetAllCourses() ([]Class, error) {
	allClasses := []Class{}

	res, err := http.Get("https://uvic.kuali.co/api/v1/catalog/courses/5d9ccc4eab7506001ae4c225")
	if err != nil {
		return allClasses, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return allClasses, err
	}

	var allCourses []struct {
		ID          string `json:"__catalogCourseId"`
		Start       string `json:"dateStart"`
		PID         string `json:"pid"`
		Title       string `json:"title"`
		SubjectCode struct {
			Description string `json:"description"`
		} `json:"subjectCode"`
	}

	err = json.Unmarshal([]byte(body), &allCourses)
	if err != nil {
		return allClasses, err
	}

	for _, v := range allCourses {
		/* COURSE NAME */
		courseName, err := regexp.Compile(`^[A-Z]+`)
		if err != nil {
			return allClasses, err
		}
		name := courseName.Find([]byte(v.ID))

		/* COURSE ID */
		courseID, err := regexp.Compile(`\d{3,4}\w{0,1}`)
		if err != nil {
			return allClasses, err
		}
		id := courseID.Find([]byte(v.ID))

		allClasses = append(allClasses, Class{string(name), string(id), v.PID, v.Title, v.SubjectCode.Description})
	}

	return allClasses, nil
}
