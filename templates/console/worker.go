package console

import (
	"errors"
	"os"
	"os/signal"

	"gitlab.kumparan.com/yowez/skeleton-service/db"

	"gitlab.kumparan.com/yowez/skeleton-service/config"
	"gitlab.kumparan.com/yowez/skeleton-service/repository"
	"gitlab.kumparan.com/yowez/skeleton-service/worker"

	"github.com/kumparan/cacher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "worker",
	Long:  `This subcommand used to run worker`,
	Run:   runWorker,
}

func init() {
	RootCmd.AddCommand(workerCmd)
}

func runWorker(_ *cobra.Command, _ []string) {
	// Initiate DB connection
	db.InitializeCockroachConn()
	// Initialize CacheKeeper
	cacheKeeper := cacher.NewKeeper()
	redisWorkerConn := db.NewRedisConnPool(config.RedisWorkerBrokerHost())
	if !config.DisableCaching() {
		redisConn := db.NewRedisConnPool(config.RedisCacheHost())
		redisLockConn := db.NewRedisConnPool(config.RedisLockHost())

		cacheKeeper.SetConnectionPool(redisConn)
		cacheKeeper.SetLockConnectionPool(redisLockConn)
		cacheKeeper.SetDefaultTTL(config.CacheTTL())
	}
	cacheKeeper.SetDisableCaching(config.DisableCaching())

	helloRepo := repository.NewHelloRepository(db.DB, cacheKeeper)

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		select {
		case <-sigCh:
			errCh <- errors.New("received an interrupt")
		}
	}()

	wrk := worker.New(helloRepo)
	wrk.SetConnectionPool(redisWorkerConn)
	wrk.InitWorkers()

	go wrk.RunWorkers()

	log.Error(<-errCh)

	wrk.Stop()
}
