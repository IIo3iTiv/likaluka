package engine

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Table struct {
	table       string
	scheme      string
	schemeTable string
	fields      fields
}

type fields struct {
	UID           string
	CREATION_DATE string
	NAME          string
	DESCRIPTION   string
}

type DataList []Data

type Data struct {
	Uid          uuid.UUID          `db:"UID" json:"UID"`
	CreationDate pgtype.Timestamptz `db:"CREATION_DATE" json:"creation_date"`
	Name         string             `db:"NAME" json:"name"`
	Description  string             `db:"DESCRIPTION" json:"description"`
}