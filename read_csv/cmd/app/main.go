package main

import (
	"it_camp_case/internal/util"
	readcsv "it_camp_case/read_csv"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	execPath := util.GetExecuteFilePath()
	datasetPath := filepath.Join(execPath, "..", "..", "v3", "datasets", "current_1.csv")
	dataset1, err := os.Open(datasetPath)
	if err != nil {
		log.Printf("ERROR: Open dataset: %s", err.Error())
		return
	}

	log.Println("Start ReadCSV")
	t1 := time.Now()
	records, err := readcsv.Read(dataset1)
	t2 := time.Now()
	if err != nil {
		log.Printf("ERROR: readcsv.Read: %s", err.Error())
		return
	}
	log.Println("Success")
	log.Printf("TimeStart: %s", t1.Format(time.RFC3339Nano))
	log.Printf("TimeEnd: %s", t2.Format(time.RFC3339Nano))
	log.Printf("TimeDif: %s", t2.Sub(t1).String())
	log.Printf("RecordCount: %d", len(records))
	for i, r := range records {
		if i >= 100 {
			break
		}
		log.Printf("%e\t%e\t%e", r.R, r.S, r.T)
	}
}
