package main

import (
	"github.com/thesis-bkn/hfsd/internal/server"
)

func main() {
	server := server.NewServer()

	server.Logger.Fatal(server.Start(server.Server.Addr))
}
