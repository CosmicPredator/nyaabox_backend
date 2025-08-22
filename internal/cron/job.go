package cron

import (
	"context"
	"cosmic/nyaabox/internal/db"
	"cosmic/nyaabox/internal/storage"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type CronJob struct {
	interval time.Duration
	ctx context.Context
	dbHandle *db.DbHandler
	currTime time.Time
}

func NewCronJob(ctx context.Context, interval time.Duration, dbHandle *db.DbHandler) *CronJob {
	return &CronJob{
		interval: interval,
		ctx: ctx,
		dbHandle: dbHandle,
		currTime: time.Now(),
	}
}

func (c *CronJob) onCronTick() error {
	entires, err := c.dbHandle.GetEntriesByExpiredAt(c.currTime)
	if err != nil { return err }
	log.Infof("%v", entires)
	
	for _, entry := range entires {
		log.Infof("removing file entry: %s", entry.FilePath)
		err = storage.RemoveFile(entry.FilePath)
		if err != nil {
			return err
		}
	}
	err = c.dbHandle.DeleteEntriesByExpiredAt(c.currTime)
	return err
}

func (c *CronJob) RunCronJob() {
	timer := time.NewTicker(c.interval)
	for {
		select{
		case <- timer.C:
			c.currTime = time.Now()
			log.Debug("running cron job")
			err := c.onCronTick()
			if err != nil {
				log.Errorf("error occured while running cron job: %s", err.Error())
			}
			log.Debug("successfully ran crob job")
		case <- c.ctx.Done():
			return
		}
	}
}