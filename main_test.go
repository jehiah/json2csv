package main

import (
	"bytes"
	"encoding/csv"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGetTopic(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	reader := bytes.NewBufferString(`{"a": 1, "b": "asdf\n"}
{"a" : null}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a", "c"})

	output := buf.String()
	assert.Equal(t, output, "1,\"\"\n\"\",\"\"\n")
}
