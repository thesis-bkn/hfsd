package main

import (
	"github.com/thesis-bkn/hfsd/internal/server"
	"github.com/thesis-bkn/hfsd/internal/worker"
	"github.com/ztrue/tracerr"
)

func main() {
	w := worker.NewWorker()
	server, err := server.NewServer(w)
	if err != nil {
		tracerr.PrintSourceColor(err)
		return
	}

	go w.Run()

	server.Logger.Fatal(server.Start(server.Server.Addr))
}
