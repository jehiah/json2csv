package main

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTopic(t *testing.T) {
	reader := bytes.NewBufferString(`{"a": 1, "b": "asdf\n"}
{"a" : null}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a", "c"}, false)

	output := buf.String()
	assert.Equal(t, output, "1,\n,\n")
}

func TestGetLargeInt(t *testing.T) {
	reader := bytes.NewBufferString(`{"a": 1356998399}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a"}, false)

	output := buf.String()
	assert.Equal(t, output, "1356998399\n")
}

func TestGetExactIntegers(t *testing.T) {
	for _, value := range []string{
		"0",
		"1",
		"-1",
		"2147483648",
		"9007199254740993",
		"-9007199254740993",
		"9223372036854775807",
		"-9223372036854775808",
		"9223372036854775808",
		"-9223372036854775809",
		"18446744073709551615",
		"110750513297351875876238378613844499030",
		"-110750513297351875876238378613844499030",
	} {
		t.Run(value, func(t *testing.T) {
			reader := bytes.NewBufferString(`{"a": ` + value + `, "nested": {"a": ` + value + `}}`)
			buf := new(bytes.Buffer)
			writer := csv.NewWriter(buf)

			json2csv(reader, writer, []string{"a", "nested.a"}, false)

			assert.Equal(t, value+","+value+"\n", buf.String())
		})
	}
}

func TestGetNumberFormatting(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  string
	}{
		{"1.0", "1"},
		{"-1.0", "-1"},
		{"-0", "0"},
		{"-0.0", "0"},
		{"-0e2", "0"},
		{"0.32", "0.320000"},
		{"-0.32", "-0.320000"},
		{"1e3", "1000"},
		{"1E-3", "0.001000"},
		{"1.25e1", "12.500000"},
		{"1e20", "100000000000000000000"},
	} {
		t.Run(tc.value, func(t *testing.T) {
			reader := bytes.NewBufferString(`{"a": ` + tc.value + `}`)
			buf := new(bytes.Buffer)
			writer := csv.NewWriter(buf)

			json2csv(reader, writer, []string{"a"}, false)

			assert.Equal(t, tc.want+"\n", buf.String())
		})
	}
}

func TestGetCompositeNumbers(t *testing.T) {
	reader := bytes.NewBufferString(`{"a": [9007199254740993, 1e3, 1.20], "b": {"n": 9223372036854775807}}`)
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a", "b"}, false)

	assert.Equal(t, "[9007199254740993 1e3 1.20],map[n:9223372036854775807]\n", buf.String())
}

func TestRejectTrailingJSON(t *testing.T) {
	for _, suffix := range []string{
		`{"a": 2}`,
		`true`,
		`123`,
		`invalid`,
		`[`,
	} {
		t.Run(suffix, func(t *testing.T) {
			reader := bytes.NewBufferString(`{"a": 1}` + suffix + "\n" + `{"a": 2}`)
			buf := new(bytes.Buffer)
			writer := csv.NewWriter(buf)

			json2csv(reader, writer, []string{"a"}, false)

			assert.Equal(t, "2\n", buf.String())
		})
	}
}

func TestSkipEmptyLines(t *testing.T) {
	for _, padding := range []string{"\n", " \t\r\n", "\n\n"} {
		t.Run(padding, func(t *testing.T) {
			reader := bytes.NewBufferString(padding + `{"a": 1}` + "\n" + padding + `{"a": 2}`)
			buf := new(bytes.Buffer)
			writer := csv.NewWriter(buf)

			json2csv(reader, writer, []string{"a"}, false)

			assert.Equal(t, "1\n2\n", buf.String())
		})
	}
}

func TestGetFloat(t *testing.T) {

	reader := bytes.NewBufferString(`{"a": 1356998399.32}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a"}, false)

	output := buf.String()
	assert.Equal(t, output, "1356998399.320000\n")
}

func TestGetNested(t *testing.T) {
	reader := bytes.NewBufferString(`{"a": {"b": "asdf"}}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a.b"}, false)

	output := buf.String()
	assert.Equal(t, output, "asdf\n")
}

func TestHeader(t *testing.T) {
	reader := bytes.NewBufferString(`{"a": "b"}
{"a": "c"}`)
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	json2csv(reader, writer, []string{"a"}, true)

	output := buf.String()
	assert.Equal(t, output, "a\nb\nc\n")
}
