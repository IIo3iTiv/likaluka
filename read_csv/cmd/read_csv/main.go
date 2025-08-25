package readcsv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type Record struct {
	R float64 `csv:"current_R"`
	S float64 `csv:"current_S"`
	T float64 `csv:"current_T"`
}

func Read(file *os.File) (records []Record, err error) {
	dec, err := csvutil.NewDecoder(csv.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("csvutil.NewDecoder: %s", err.Error())
	}

	for {
		var r Record
		if err := dec.Decode(&r); err != nil {
			break
		}
		records = append(records, r)
	}

	return
}
