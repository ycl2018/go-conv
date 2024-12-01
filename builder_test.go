package main

import "testing"

func TestNewBuilder(t *testing.T) {
	*dryRun = true
	err := generate("./testdata")
	if err != nil {
		t.Fatal(err)
	}
}
