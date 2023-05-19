package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"
	"weavatar/app/console/commands"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule.Command("avatar:update-expired").Hourly().SkipIfStillRunning(),
	}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.HashMake{},
		&commands.UpdateExpiredAvatar{},
	}
}
