package tools

import (
	"strings"
)

func getUserNameFromEmail(email string) string {
	//check params
	if email == "" {
		return ""
	}

	index := strings.Index(email, "@")
	if index == -1 {
		return ""
	}
	return email[0:index]

}
