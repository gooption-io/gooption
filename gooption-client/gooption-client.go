//go:generate gooption-gen client ImpliedVol
package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	r "gopkg.in/gorethink/gorethink.v3"
)

var (
	samplePath = "./Sample_2015_October/"
)

func main() {
	filepath.Walk(samplePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if data, err := readCSV(samplePath + info.Name()); err == nil {
				println(data[0])
			} else {
				panic(err)
			}
		}
		return nil
	})
}

func readCSV(filepath string) ([][]string, error) {
	csvfile, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	fields, err := reader.ReadAll()

	return fields, nil
}

func insertCSV(csvFile string[][]) {
	for index := 1; index < len(csvFile); index++ {
	}

	session, err := r.Connect(r.ConnectOpts{Address: ""})
	if err != nil {
		log.Fatalln(err)
	}
	r.Table("sample_2015_October").Insert(objCSV, r.InsertOpts{
		Conflict: func(id, oldCSV, newCSV r.Term) interface{} {
			return newCSV.Merge(map[string]interface{}{
				"count": oldCSV.Add(newCSV.Field("count")),
			})
		},
	})
}
