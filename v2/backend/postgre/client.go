package postgre

import (
	"context"
	"log"
	"time"

	pgxp "github.com/jackc/pgx/v5/pgxpool"
)

func Init(opt *options) (err error) {
	connStr, err := opt.Build()
	if err != nil {
		return err
	}
	log.Printf("ConnStr:{%s}", connStr)

	conf, err := pgxp.ParseConfig(connStr)
	if err != nil {
		c.lastInit.lError = err
		return
	}

	c.ctx = context.Background()
	c.conn, err = pgxp.NewWithConfig(c.ctx, conf)
	if err != nil {
		c.lastInit.lError = err
		return
	}
	c.isInit = true
	c.lastInit.lTime = time.Now()

	err = c.conn.Ping(c.ctx)
	if err != nil {
		c.lastInit.lPing = false
		c.lastInit.lError = err
		return
	}
	c.lastInit.lPing = true
	return
}

func GetConnect(ctx context.Context) (*pgxp.Conn, error) {
	conn, err := c.conn.Acquire(ctx)
	if err != nil {
		if !c.isInit || !c.lastInit.lPing || time.Since(c.lastInit.lTime).Minutes() > 5 {
			err = Init(&opt)
			if err != nil {
				return nil, err
			}
			return GetConnect(ctx)
		}
		return nil, err
	}
	return conn, nil
}
