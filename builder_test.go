package main

import (
	"os"
	"testing"

	"github.com/ycl2018/go-conv/internal"
)

func TestNewBuilder(t *testing.T) {
	*dryRun = false
	*verbose = true
	*quiet = true
	internal.DefaultLogger = internal.NewLogger(os.Stdout, *verbose)
	err := generate("./testdata/cases/apply")
	if err != nil {
		t.Fatal(err)
	}
}
