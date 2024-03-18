package main

import (
	"github.com/thesis-bkn/hfsd/internal/server"
	"github.com/ztrue/tracerr"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		tracerr.PrintSourceColor(err)
        return
	}

	server.Logger.Fatal(server.Start(server.Server.Addr))
}
