package courses

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"regexp"
)

type course struct {
	ID string `json:"__catalogCourseId"`
}

type Class struct {
	Subject string `json:"subject"`
	Course  string `json:"course"`
}

func GetAllCourses(data *[]Class) error {

	res, err := http.Get("https://uvic.kuali.co/api/v1/catalog/courses/5d9ccc4eab7506001ae4c225")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var allCourses []course
	err = json.Unmarshal([]byte(body), &allCourses)
	if err != nil {
		return err
	}

	for _, v := range allCourses {
		/* COURSE NAME */
		courseName, err := regexp.Compile(`^[A-Z]+`)
		handleError("Failed to parse course name", err)
		name := courseName.Find([]byte(v.ID))

		/* COURSE ID */
		courseID, err := regexp.Compile(`\d{3,4}\w{0,1}`)
		handleError("Failed to parse course ID", err)
		id := courseID.Find([]byte(v.ID))

		newEntry := Class{string(name), string(id)}
		fmt.Println(newEntry)
		*data = append(*data, Class{string(name), string(id)})
	}

	return nil
}

func handleError(message string, err error) {
	if err != nil {
		log.Fatalf("%s:\n%v", message, err)
	}
}
