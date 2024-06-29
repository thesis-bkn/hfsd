package main

import (
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/server"
	"github.com/thesis-bkn/hfsd/internal/worker"
)

func main() {
	w := worker.NewWorker()
    taskQueue := make(chan interface{})
	server, err := server.NewServer(taskQueue)
	if err != nil {
		tracerr.PrintSourceColor(err)
		return
	}

	go w.Run(taskQueue)

	server.Logger.Fatal(server.Start(server.Server.Addr))
}
