package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	columns := flag.Int("columns", 2, "number of columns")
	date := flag.String("date", "", "date in text format")
	flag.Parse()

	reader := csv.NewReader(os.Stdin)

	var headers []string
	var data []map[string]string

	for i := 0; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read csv: %s", err)
		}

		if i == 0 {
			// first line. gather headers
			headers = record
		} else {
			vals := map[string]string{}
			for j := 0; j < len(headers); j++ {
				vals[headers[j]] = record[j]
			}
			data = append(data, vals)
		}
	}

	f, err := os.Create("output.csv")
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	var outputLine []string
	for i, line := range data {
		fmt.Println("processing line ", i, " outputLine len: ", len(outputLine), " orientation: ", line["Orientation"])
		if len(outputLine) == *columns {
			fmt.Println("writing output line")
			if err := w.Write(outputLine); err != nil {
				log.Fatalln("error writing record to file", err)
			}
			outputLine = nil
		}
		if len(line["Orientation"]) > 0 {
			outputLine = append(outputLine, buildLabelRegular(line, *date))
		}

	}

}

func buildLabelRegular(line map[string]string, date string) string {
	fmt.Print("symbol=", line["Symbol"], "cat: ", line["Cat"], "\n")
	outputLineOne := line["Symbol"] + " - " + line["Cat"][len(line["Cat"])-3:]

	tmiNumber := line["TMI"][2:]
	tmiShortVal := 0
	tmiShortVal, err := strconv.Atoi(tmiNumber)
	if err != nil {
		fmt.Println("s=", tmiNumber)
	}

	outputLineTwo := ""
	if line["Orientation"] == "CAP" || line["Orientation"] == "DUAL" {
		outputLineTwo = fmt.Sprint("302/", tmiShortVal)
	} else {
		outputLineTwo = fmt.Sprint(tmiShortVal, "L/148")
	}

	outputLineThree := "200nM"
	outputLineFour := date

	multiValue := []string{outputLineOne, outputLineTwo, outputLineThree, outputLineFour}
	fmt.Println("output line: ", fmt.Sprint(multiValue))
	return strings.Join(multiValue, "\n")
}
