package worker

import (
	"context"
	"errors"

	"github.com/gocraft/work"
	"github.com/kumparan/go-lib/utils"
	log "github.com/sirupsen/logrus"
	"gitlab.kumparan.com/yowez/skeleton-service/repository"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/model"
)

type taskContext struct {
	helloRepository repository.HelloRepository

	greeting *model.Greeting
}

// updateEventHandler :nodoc:
func (t *taskContext) updateEventHandler(job *work.Job) (err error) {
	if t.greeting == nil {
		return errors.New("Invalid data")
	}

	err = t.helloRepository.Update(context.TODO(), t.greeting.ID, t.greeting)
	if err != nil {
		log.WithFields(log.Fields{
			"greeting": utils.Dump(t.greeting),
		}).Error(err)
		return err
	}
	return nil
}
