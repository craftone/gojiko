package stats

import (
	"os"
	"testing"

	"github.com/craftone/gojiko/config"
)

func TestMain(m *testing.M) {
	config.Init()
	Init()
	code := m.Run()
	os.Exit(code)
}
