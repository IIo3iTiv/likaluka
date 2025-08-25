package engine

import (
	"context"
	"fmt"
	"miniogo/v2/postgre"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (t *Table) exists() error {
	ctx := context.Background()
	conn, err := postgre.GetConnect(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, fmt.Sprintf("SELECT 1 from %s LIMIT 1;", t.schemeTable))
	return err
}

func (t *Table) New(ctx context.Context, name, description string) (uid uuid.UUID, err error) {
	if len(name) == 0 {
		err = fmt.Errorf("parameter %s is empty", "name")
		return
	}
	if len(description) == 0 {
		err = fmt.Errorf("parameter %s is empty", "description")
		return
	}

	uid = uuid.New()
	ct := time.Now()
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(t.schemeTable).
		Columns(t.fields.UID, t.fields.CREATION_DATE, t.fields.NAME, t.fields.DESCRIPTION).
		Values(uid, ct, name, description)

	sql, args, err := query.ToSql()
	if err != nil {
		return
	}

	conn, err := postgre.GetConnect(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		return
	}
	return
}
