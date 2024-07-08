package task

import (
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func init() {
	// Seconds field, optional
	c = cron.New(cron.WithSeconds())

	c.Start()
}
