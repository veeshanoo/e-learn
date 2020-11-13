package main

import (
	"e-learn/dbg"
	"e-learn/server"
)

func main() {
	dbg.SetDebug(true)
	srv := &server.Server{}
	srv.Run()
}
