package providers

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"weavatar/app/jobs"
)

type QueueServiceProvider struct {
}

func (receiver *QueueServiceProvider) Register() {
	facades.Queue.Register(receiver.Jobs())
}

func (receiver *QueueServiceProvider) Boot() {

}

func (receiver *QueueServiceProvider) Jobs() []queue.Job {
	return []queue.Job{
		&jobs.ProcessAvatarCheck{},
		&jobs.ProcessAvatarUpdate{},
	}
}
