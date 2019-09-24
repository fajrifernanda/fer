package console

import (
	"context"
	"fmt"
	"net"

	// "time"

	// redigo "github.com/gomodule/redigo/redis"
	redcachekeeper "github.com/kumparan/cacher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.kumparan.com/yowez/skeleton-service/config"
	"gitlab.kumparan.com/yowez/skeleton-service/db"
	pb "gitlab.kumparan.com/yowez/skeleton-service/pb/skeleton"
	"gitlab.kumparan.com/yowez/skeleton-service/repository"
	"gitlab.kumparan.com/yowez/skeleton-service/service"
	"google.golang.org/grpc"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	// Initiate all connection
	db.InitializeCockroachConn()
	cacheKeeper := redcachekeeper.NewKeeper()

	if !config.DisableCaching() {
		redisConn := db.NewRedisConnPool(config.RedisCacheHost())
		redisLockConn := db.NewRedisConnPool(config.RedisLockHost())

		cacheKeeper.SetConnectionPool(redisConn)
		cacheKeeper.SetLockConnectionPool(redisLockConn)
		cacheKeeper.SetDisableCaching(config.DisableCaching())
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port()))
	if err != nil {
		log.WithField("port", config.Port()).Fatalf("failed to listen: %v", err)
	}

	log.Info("Listening on ", config.Port())

	// Service definition
	svc := service.NewHelloService()
	helloRepo := repository.NewHelloRepository(db.DB, cacheKeeper)
	svc.RegisterHelloRepository(helloRepo)

	server := grpc.NewServer()
	pb.RegisterHelloServiceServer(server, svc)

	if err := server.Serve(lis); err != nil {
		log.WithField("lis", lis).Fatalf("failed to serve: %v", err)
	}
}

func serverInterceptor(ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, config.RPCServerTimeout())
	defer cancel()
	return handler(ctx, req)
}
