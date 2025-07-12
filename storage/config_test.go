package storage

import (
	"testing"
)

func TestNoConfigFilePassed(t *testing.T) {
	cfgPath, _ := ParseFlags()
	want := ""

	if cfgPath != want {
		t.Errorf("got %q want %q", cfgPath, want)
	}
}

//func TestExists(t *testing.T) {
//	_, err := os.Create("./config.yml")
//	if err != nil {
//		return
//	}
//	cfgPath, _ := ParseFlags()
//	want := "./config.yml"
//
//	if cfgPath != want {
//		t.Errorf("got %q want %q", cfgPath, want)
//	}
//}
