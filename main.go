package main

import (
	"github.com/goravel/framework/facades"

	"weavatar/bootstrap"
)

func main() {
	// WeAvatar，启动！
	bootstrap.Boot()

	// 路由，启动！
	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route run error: %v", err)
		}
	}()

	// 计划任务，启动！
	go facades.Schedule().Run()

	// 队列，启动！
	go func() {
		if err := facades.Queue().Worker(nil).Run(); err != nil {
			facades.Log().Errorf("Queue run error: %v", err)
		}
	}()

	select {}
}
