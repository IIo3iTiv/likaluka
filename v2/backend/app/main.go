package main

import (
	"context"
	"log"
	"miniogo/v2/minio"
	"miniogo/v2/postgre"
	"miniogo/v2/server"
)

func main() {
	var err error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Config
	LoadConfig()

	// Minio
	err = minio.Init(context.Background(), conf.MinioEndpoint, conf.MinioU, conf.MinioP, conf.MinioBucketName, conf.MinioUseSSL)
	if err != nil {
		log.Printf("ERROR: Minio: %s", err.Error())
		return
	}

	// Postgre
	err = postgre.Init(postgre.NewOption().Set_Full(conf.PgHost, conf.PgPort, conf.PgDB, conf.PgU, conf.PgP, conf.PgMOC, conf.PgMCLT, conf.PgMILT, conf.PgSSL))
	if err != nil {
		log.Printf("ERROR: Postgre: %s", err.Error())
		return
	}

	// Server
	s := server.Init(conf.Mode)
	s.SetMainHandlers()
	s.Run(conf.ServerPort)
}
