package string_utils

import "regexp"

func AppendStringToFileName(fileName string, stringToAppend string) string {
	m1 := regexp.MustCompile(`\.([^.]*)$`)
	fileExtWithDot := m1.FindString(fileName)
	withSuffix := m1.ReplaceAllString(fileName, "_"+stringToAppend+fileExtWithDot)

	return withSuffix
}
