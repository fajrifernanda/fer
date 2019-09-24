package client

import (
	"context"
	"time"

	grpcpool "github.com/processout/grpc-go-pool"
	log "github.com/sirupsen/logrus"
	pb "gitlab.kumparan.com/yowez/skeleton-service/pb/skeleton"
	"google.golang.org/grpc"
)

type client struct {
	Conn *grpcpool.Pool
}

//NewClient is a func to create Client
func NewClient(target string, timeout time.Duration, idleConnPool, maxConnPool int) (pb.HelloServiceClient, error) {
	factory := newFactory(target, timeout)

	pool, err := grpcpool.New(factory, idleConnPool, maxConnPool, time.Second)
	if err != nil {
		log.Errorf("Error : %v", err)
		return nil, err
	}

	return &client{
		Conn: pool,
	}, nil
}

func newFactory(target string, timeout time.Duration) grpcpool.Factory {
	return func() (*grpc.ClientConn, error) {
		conn, err := grpc.Dial(target, grpc.WithInsecure(), withClientUnaryInterceptor(timeout))
		if err != nil {
			log.Errorf("Error : %v", err)
			return nil, err
		}

		return conn, err
	}
}

func withClientUnaryInterceptor(timeout time.Duration) grpc.DialOption {
	return grpc.WithUnaryInterceptor(func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	})
}
