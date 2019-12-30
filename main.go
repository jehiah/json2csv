package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"unicode/utf8"
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

func get_value(data map[string]interface{}, keyparts []string) string {
	if len(keyparts) > 1 {
		subdata, _ := data[keyparts[0]].(map[string]interface{})
		return get_value(subdata, keyparts[1:])
	} else if v, ok := data[keyparts[0]]; ok {
		switch v.(type) {
		case nil:
			return ""
		case float64:
			f, _ := v.(float64)
			if math.Mod(f, 1.0) == 0.0 {
				return fmt.Sprintf("%d", int(f))
			} else {
				return fmt.Sprintf("%f", f)
			}
		default:
			return fmt.Sprintf("%+v", v)
		}
	}

	return ""
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

		var data map[string]interface{}
		err = json.Unmarshal(line, &data)
		if err != nil {
			log.Printf("ERROR Decoding JSON at line %d: %s\n%s", line_count, err, line)
			continue
		}

		var record []string
		for _, expanded_key := range expanded_keys {
			record = append(record, get_value(data, expanded_key))
		}

		w.Write(record)
		w.Flush()
	}
}
