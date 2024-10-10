package utils

import (
	"regexp"
	"strings"
)

// StripFileName removes special characters from the filename and replaces "-" with "_".
func StripFileName(filename string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9._-]+`).ReplaceAllString(filename, "")
}

func OriginalFileName(url string) string {
	parts := strings.Split(url, ":")
	if parts[0] != "minio" {
		return url
	}
	uuidRegExp := "[a-fA-F0-9]{8}[-_][a-fA-F0-9]{4}[-_][a-fA-F0-9]{4}[-_][a-fA-F0-9]{4}[-_][a-fA-F0-9]{12}[-_]"
	return regexp.MustCompile(uuidRegExp).ReplaceAllString(parts[2], "")
}
