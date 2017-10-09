package uaEngine

import (
	"testing"
)

func TestUaEngine(t *testing.T) {
	uaEngine := NewUaEngine("")
	t.Log(uaEngine.GetPcRandomeEngine())
	t.Log(uaEngine.GetMobileRandomeEngine())
}
