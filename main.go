package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/buger/jsonparser"
)

type LineReader interface {
	ReadBytes(delim byte) (line []byte, err error)
}

func main() {
	inputFile := flag.String("i", "", "/path/to/input.json (optional; default is stdin)")
	outputFile := flag.String("o", "", "/path/to/output.csv (optional; default is stdout)")
	outputDelim := flag.String("d", ",", "delimiter used for output values")
	showVersion := flag.Bool("version", false, "print version string")
	printHeader := flag.Bool("p", false, "prints header to output")
	keys := StringArray{}
	flag.Var(&keys, "k", "fields to output")
	flag.Parse()

	if *showVersion {
		fmt.Printf("json2csv v%s (built w/%s)\n", VERSION, runtime.Version())
		return
	}

	var reader *bufio.Reader
	var writer *csv.Writer
	if *inputFile != "" {
		file, err := os.OpenFile(*inputFile, os.O_RDONLY, 0600)
		if err != nil {
			log.Printf("Error %s opening input file %v", err, *inputFile)
			os.Exit(1)
		}
		reader = bufio.NewReader(file)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	if *outputFile != "" {
		file, err := os.OpenFile(*outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Printf("Error %s opening output file %v", err, *outputFile)
			os.Exit(1)
		}
		writer = csv.NewWriter(file)
	} else {
		writer = csv.NewWriter(os.Stdout)
	}

	delim, _ := utf8.DecodeRuneInString(*outputDelim)
	writer.Comma = delim

	json2csv(reader, writer, keys, *printHeader)
}

func json2csv(r LineReader, w *csv.Writer, keys []string, printHeader bool) {
	var line []byte
	var err error
	line_count := 0

	var expanded_keys [][]string
	for _, key := range keys {
		expanded_keys = append(expanded_keys, strings.Split(key, "."))
	}

	for {
		if err == io.EOF {
			return
		}
		line, err = r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Input ERROR: %s", err)
				break
			}
		}
		line_count++
		if len(line) == 0 {
			continue
		}

		if printHeader {
			w.Write(keys)
			w.Flush()
			printHeader = false
		}

		var record []string
		for _, expanded_key := range expanded_keys {
			//val, valuetype, offset, err
			val, vt, _, err := jsonparser.Get(line, expanded_key...)

			if err != nil {
				if errors.Is(err, jsonparser.KeyPathNotFoundError) {
					record = append(record, "")
					continue
				} else {
					log.Printf("ERROR Retrieving JSON key %s at line %d: %s\n%s",
						strings.Join(expanded_key[:], "."), line_count, err, line)
					record = append(record, "")
					continue
				}
			}

			value, err := convertJsonValue(val, vt)

			if err != nil {
				log.Printf("ERROR Decoding JSON at line %d: %s\n%s", line_count, err, line)
				continue
			}

			record = append(record, value)
		}

		w.Write(record)
		w.Flush()
	}
}

func convertJsonValue(v []byte, vt jsonparser.ValueType) (string, error) {
	switch vt {
	case jsonparser.String:
		return string(v[:]), nil
	case jsonparser.Boolean:
		return fmt.Sprintf("+%v", string(v[:])), nil
	case jsonparser.Number:
		f, _ := strconv.ParseFloat(string(v[:]), 64)
		if math.Mod(f, 1.0) == 0.0 {
			return fmt.Sprintf("%d", int(f)), nil
		} else {
			return fmt.Sprintf("%f", f), nil
		}
	case jsonparser.NotExist:
	case jsonparser.Null:
		return "", nil
	case jsonparser.Object:
		return "", fmt.Errorf("JSON value is an object: %s", v)
	case jsonparser.Array:
		return "", fmt.Errorf("JSON value is an array: %s", v)
	}

	return "", fmt.Errorf("JSON value is an unknown type: %s", v)
}
