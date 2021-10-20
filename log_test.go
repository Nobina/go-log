package log

import "testing"

func TestLog(t *testing.T) {
	c := New()

	c.Infof("%d Test on line 8", 1337)
	c.Debug("Debug on line 9")
	c.Error("Error on line 10")
	c.Critical("Critical on line 11")
	c.Warning("Warning on line 12")
}
