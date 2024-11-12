package xml

import "strings"

func GenerateEmptyMusicXML() string {
	return template
}

func GenerateProjectXML(title string) string {
	project := strings.Replace(template, "Untitled", title, 1)
	return project
}
