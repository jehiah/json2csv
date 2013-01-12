package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type LineReader interface {
	ReadBytes(delim byte) (line []byte, err error)
}

var (
	inputFilename  = flag.String("i", "", "/path/to/input.json (optional; default is stdin)")
	outputFilename = flag.String("o", "", "/path/to/output.json (optional; default is stdin)")
	verbose        = flag.Bool("v", false, "verbose output (to stderr)")
	showVersion    = flag.Bool("version", false, "print version string")
	keys           = StringArray{}
)

func init() {
	flag.Var(&keys, "k", "fields to output")
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("json2csv v1.1\n")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	writer := csv.NewWriter(os.Stdout)
	json2csv(reader, writer, keys)
}

func json2csv(r LineReader, w *csv.Writer, keys []string) {
	line, err := r.ReadBytes('\n')
	for {
		if err != nil {
			if err != io.EOF {
				log.Printf("Input ERROR: %s", err)
				break
			}
		}

		var data map[string]interface{}
		err = json.Unmarshal(line, &data)
		if err != nil {
			log.Printf("ERROR Json Decoding: %s - %v", err, line)
			continue
		}
		var record []string
		for _, key := range keys {
			if v, ok := data[key]; ok {
				switch v.(type) {
				case nil:
					record = append(record, "")
				default:
					record = append(record, fmt.Sprintf("%+v", v))
				}
			} else {
				record = append(record, "")
			}
		}
		w.Write(record)
		w.Flush()
		line, err = r.ReadBytes('\n')
	}
}
