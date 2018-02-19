package tools

import (
	"testing"
)

func TestUserNameFromEmail(t *testing.T) {
	name := GetUserNameFromEmail("zhoupengyuan@goumin.com")
	t.Log(name)
}
