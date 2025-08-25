package postgre

import (
	"context"
	"time"

	pgxp "github.com/jackc/pgx/v5/pgxpool"
)

type client struct {
	conn     *pgxp.Pool
	isInit   bool
	ctx      context.Context
	lastInit struct {
		lTime  time.Time
		lPing  bool
		lError error
	}
}

type options struct {
	Server          string
	Port            string
	Db              string
	User            string
	Pass            string
	MaxOpenConn     string
	MaxConnLifeTime string
	MaxIdleLifeTime string
	SslMode         string
}

const (
	// Keys connection
	keyConnectHost   string = "host="
	keyConnectPort   string = "port="
	keyConnectDbName string = "dbname="
	keyConnectUser   string = "user="
	keyConnectPass   string = "password="
	keyConnectSSL    string = "sslmode="
	keyConnectPMC    string = "pool_max_conns="
	keyConnectPMCLT  string = "pool_max_conn_lifetime="
	keyConnectPMILT  string = "pool_max_conn_idle_time="
)

var (
	c   = client{}
	opt = options{}
)
