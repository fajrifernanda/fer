package worker

import (
	"gitlab.kumparan.com/yowez/skeleton-service/repository/model"

	"gitlab.kumparan.com/yowez/skeleton-service/config"
	"gitlab.kumparan.com/yowez/skeleton-service/repository"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/kumparan/tapao"
	log "github.com/sirupsen/logrus"
)

type caller string

const (
	contextCaller = caller("caller")

	// Define tasks here
	TaskUpdate = "skeleton.update"
)

// Worker :nodoc:
type Worker interface {
	SetConnectionPool(*redis.Pool)
	InitWorkers()

	RunWorkers()
	Stop()

	// Tasks
	Update(*model.Greeting) error
}

type (
	workerImpl struct {
		helloRepository  repository.HelloRepository
		redisHostBroker  string
		redisHostBackend string
		enqueuer         *work.Enqueuer
		pool             *work.WorkerPool
		redisPool        *redis.Pool
	}
)

// New :nodoc:
func New(helloRepository repository.HelloRepository) Worker {
	return &workerImpl{
		helloRepository: helloRepository,
	}
}

// SetConnectionPool(*redis.Pool)
func (w *workerImpl) SetConnectionPool(pool *redis.Pool) {
	w.redisPool = pool
}

// InitWorkers :nodoc:
func (w *workerImpl) InitWorkers() {
	w.enqueuer = work.NewEnqueuer(config.GocraftWorkerNamespace, w.redisPool)
}

// RunWorkers :nodoc:
func (w *workerImpl) RunWorkers() {
	w.pool = work.NewWorkerPool(taskContext{}, config.GocraftWorkerPoolConcurrency, config.GocraftWorkerNamespace, w.redisPool)

	w.pool.Middleware(func(c *taskContext, job *work.Job, next work.NextMiddlewareFunc) error {
		c.helloRepository = w.helloRepository

		if _, ok := job.Args["greeting"]; ok {
			d := job.ArgString("greeting")
			greeting := &model.Greeting{}
			err := tapao.Unmarshal([]byte(d), greeting, tapao.FallbackWith(tapao.JSON))
			if err != nil {
				log.WithFields(log.Fields{
					"greeting": d,
				}).Error(err)
				return err
			}
			c.greeting = greeting
		}
		return next()
	})

	// Scheduled jobs, if want using cron
	w.pool.PeriodicallyEnqueue("* * 24 * * * ", TaskUpdate) // This will enqueue a job every 12 am

	// Tasks
	w.pool.JobWithOptions(TaskUpdate, work.JobOptions{
		MaxFails: 3,
	}, (*taskContext).updateEventHandler)
	w.pool.Start()
}

// Stop Stop all pending jobs
func (w *workerImpl) Stop() {
	// TODO stop scheduled job
	if w.pool != nil {
		w.pool.Stop()
	}
}

// Update :nodoc:
func (w *workerImpl) Update(greeting *model.Greeting) (err error) {
	b, err := tapao.Marshal(greeting, tapao.With(tapao.JSON))
	if err != nil {
		return
	}
	_, err = w.enqueuer.Enqueue(TaskUpdate, work.Q{"document": string(b)})
	return
}
