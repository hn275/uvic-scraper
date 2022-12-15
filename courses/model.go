package courses

type Class struct {
	Subject     string `json:"subject"`
	Course      string `json:"course"`
	PID         string `json:"pid"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ClassLocation struct {
	Weekday   string `json:"weekday"`
	Time      string `json:"time"`
	Building  string `json:"building"`
	DateRange string `json:"date_range"`
}

type ClassInfo struct {
	Subject  string          `json:"subject"`
	Course   int             `json:"course"`
	Term     int             `json:"term"`
	Location []ClassLocation `json:"location"`
}
