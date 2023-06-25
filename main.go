package main

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"weavatar/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Start http server by facades.Route().
	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route run error: %v", err)
		}
	}()

	// Start queue server by facades.Queue().
	go func() {
		if err := facades.Queue().Worker(&queue.Args{
			Connection: "sync",
			Queue:      "process_avatar_check",
			Concurrent: 10,
		}).Run(); err != nil {
			facades.Log().Errorf("Queue run error: %v", err)
		}
	}()

	// Start schedule by facades.Schedule()
	go facades.Schedule().Run()

	select {}
}
