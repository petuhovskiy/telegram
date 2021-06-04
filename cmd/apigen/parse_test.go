package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/petuhovskiy/telegram/tools/apigen"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("api.html")
	assert.Nil(t, err)
	defer f.Close()

	p, err := apigen.Parse(f, apigen.DefaultParseOpts)
	spew.Dump(p, err)

	assert.Nil(t, err)
}
