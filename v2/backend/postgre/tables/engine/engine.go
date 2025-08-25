package engine

import "fmt"

func New(scheme, table string) (t Table, err error) {
	if len(scheme) == 0 {
		err = fmt.Errorf("parameter scheme is empty")
		return
	}
	if len(table) == 0 {
		err = fmt.Errorf("parameter table is empty")
		return
	}
	t.scheme = scheme
	t.table = table
	t.schemeTable = fmt.Sprintf("%s.%s", scheme, table)
	t.fields = fields{
		UID:           "UID",
		CREATION_DATE: "CREATION_DATE",
		NAME:          "NAME",
		DESCRIPTION:   "DESCRIPTION",
	}
	err = t.exists()
	return
}
