package main

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	*dryRun = false
	*verbose = true
	err := generate("./testdata/cases/...")
	if err != nil {
		t.Fatal(err)
	}
}
