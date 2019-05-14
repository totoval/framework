package failed_queue

import (
	"strconv"

	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/queue"
)

func init() {
	cmd.Add(&Retry{})
}

type Retry struct {
}

func (r *Retry) Command() string {
	return "failed_queue:retry {queue_id}"
}

func (r *Retry) Aliases() []string {
	return []string{"c"}
}

func (r *Retry) Description() string {
	return "Retry a failed queue from database"
}

func (r *Retry) Handler(arg *cmd.Arg) error {
	queueIdStr, err := arg.Get("queue_id")
	if err != nil {
		return err
	}

	queueId, err := strconv.ParseUint(queueIdStr, 10, 32)
	if err != nil {
		return err
	}

	if err := queue.Retry(uint(queueId)); err != nil {
		return err
	}

	return nil
}
