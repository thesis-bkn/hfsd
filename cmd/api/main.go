package main

import (
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/server"
	"github.com/thesis-bkn/hfsd/internal/worker"
)

func main() {
	w := worker.NewWorker()
	taskQueue := make(chan entity.Task)

	cfg := config.LoadConfig()
	client, err := database.NewClient(cfg)
	if err != nil {
		tracerr.PrintSourceColor(err)
		return
	}

	server, err := server.NewServer(
		taskQueue,
		cfg,
		client,
	)
	if err != nil {
		tracerr.PrintSourceColor(err)
		return
	}

	go w.Run(taskQueue, client)

	server.Logger.Fatal(server.Start(
		server.Server.Addr,
		// "./cert/thesis-hfsd.xyz-ssl-bundle/domain.cert.pem",
		// "./cert/thesis-hfsd.xyz-ssl-bundle/private.key.pem",
	))

	// server.Logger.Fatal(server.StartTLS(
	// 	server.Server.Addr,
	// 	"./cert/thesis-hfsd.xyz-ssl-bundle/domain.cert.pem",
	// 	"./cert/thesis-hfsd.xyz-ssl-bundle/private.key.pem",
	// ))
}
