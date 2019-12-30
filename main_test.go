package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTopic(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	reader := bytes.NewBufferString(`{"a": 1, "b": "asdf\n"}
{"a" : null}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a", "c"}, false)

	output := buf.String()
	assert.Equal(t, output, "1,\n,\n")
}

func TestGetLargeInt(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	reader := bytes.NewBufferString(`{"a": 1356998399}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a"}, false)

	output := buf.String()
	assert.Equal(t, output, "1356998399\n")
}

func TestGetFloat(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	reader := bytes.NewBufferString(`{"a": 1356998399.32}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a"}, false)

	output := buf.String()
	assert.Equal(t, output, "1356998399.320000\n")
}

func TestGetNested(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	reader := bytes.NewBufferString(`{"a": {"b": "asdf"}}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a.b"}, false)

	output := buf.String()
	assert.Equal(t, output, "asdf\n")
}

func TestHeader(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	reader := bytes.NewBufferString(`{"a": "b"}
{"a": "c"}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a"}, true)

	output := buf.String()
	assert.Equal(t, output, "a\nb\nc\n")
}
