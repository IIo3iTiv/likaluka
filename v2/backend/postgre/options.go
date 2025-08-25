package postgre

import (
	"fmt"
	"strings"
)

func NewOption() *options {
	return opt.clear()
}

func (o *options) clear() *options {
	opt = options{}
	return &opt
}

func (opt *options) Set_Server(server string) *options {
	opt.Server = server
	return opt
}

func (opt *options) Set_Port(port string) *options {
	opt.Port = port
	return opt
}

func (opt *options) Set_Db(db string) *options {
	opt.Db = db
	return opt
}

func (opt *options) Set_User(user string) *options {
	opt.User = user
	return opt
}

func (opt *options) Set_Pass(pass string) *options {
	opt.Pass = pass
	return opt
}

func (opt *options) Set_MaxOpenConn(maxOpenConn string) *options {
	opt.MaxOpenConn = maxOpenConn
	return opt
}

func (opt *options) Set_MaxConnLifeTime(maxConnLifeTime string) *options {
	opt.MaxConnLifeTime = maxConnLifeTime
	return opt
}

func (opt *options) Set_MaxIdleLifeTime(maxIdleLifeTime string) *options {
	opt.MaxIdleLifeTime = maxIdleLifeTime
	return opt
}

func (opt *options) Set_SslMode(sslMode string) *options {
	opt.SslMode = sslMode
	return opt
}

func (opt *options) Set_Full(server, port, db, user, pass, maxOpenConn, maxConnLifeTime, maxIdleLifeTime, sslMode string) *options {
	return opt.Set_Server(server).
		Set_Port(port).
		Set_Db(db).
		Set_User(user).
		Set_Pass(pass).
		Set_MaxOpenConn(maxOpenConn).
		Set_MaxConnLifeTime(maxConnLifeTime).
		Set_MaxIdleLifeTime(maxIdleLifeTime).
		Set_SslMode(sslMode)
}

func (opt *options) check() (err error) {
	if len(opt.Server) == 0 {
		err = fmt.Errorf("parameter %s is empty.", keyConnectHost[0:len(keyConnectHost)-1])
		return
	}
	if len(opt.Port) == 0 {
		err = fmt.Errorf("parameter %s is empty.", keyConnectPort[0:len(keyConnectPort)-1])
		return
	}
	if len(opt.Db) == 0 {
		err = fmt.Errorf("parameter %s is empty.", keyConnectDbName[0:len(keyConnectDbName)-1])
		return
	}
	if len(opt.User) == 0 {
		err = fmt.Errorf("parameter %s is empty.", keyConnectUser[0:len(keyConnectUser)-1])
		return
	}
	if len(opt.Pass) == 0 {
		err = fmt.Errorf("parameter %s is empty.", keyConnectPass[0:len(keyConnectPass)-1])
		return
	}
	return
}

func (opt *options) Build() (conf string, err error) {
	err = opt.check()
	if err != nil {
		return
	}
	conf = strings.Join([]string{
		fmt.Sprintf("%s%s", keyConnectHost, opt.Server),
		fmt.Sprintf("%s%s", keyConnectPort, opt.Port),
		fmt.Sprintf("%s%s", keyConnectDbName, opt.Db),
		fmt.Sprintf("%s%s", keyConnectUser, opt.User),
		fmt.Sprintf("%s%s", keyConnectPass, opt.Pass),
		fmt.Sprintf("%s%s", keyConnectSSL, opt.SslMode),
		fmt.Sprintf("%s%s", keyConnectPMC, opt.MaxOpenConn),
		fmt.Sprintf("%s%s", keyConnectPMCLT, opt.MaxConnLifeTime),
		fmt.Sprintf("%s%s", keyConnectPMILT, opt.MaxIdleLifeTime),
	}, " ")
	return
}
