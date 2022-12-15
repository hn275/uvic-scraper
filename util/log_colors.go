package util

func LogColor(c, s string) string {
	var color string

	switch c {
	case "GREEN":
		color = "\033[0;32m"
		break
	case "RED":
		color = "\033[1;90m"
	case "YELLOW":
		color = "\033[0;33m"
	default:
		return ""
	}

	return color + s + "\033[0m"
}
