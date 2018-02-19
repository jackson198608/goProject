package tools

import (
	"testing"
)

func TestUserNameFromEmail(t *testing.T) {
	name := getUserNameFromEmail("zhoupengyuan@goumin.com")
	t.Log(name)
}
