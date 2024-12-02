package main

import "testing"

func TestNewBuilder(t *testing.T) {
	*dryRun = true
	*verbose = true
	err := generate("./testdata")
	if err != nil {
		t.Fatal(err)
	}
}
