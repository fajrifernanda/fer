package client

import (
	"context"
	pb "gitlab.kumparan.com/yowez/skeleton-service/pb/skeleton"
	"google.golang.org/grpc"
)

func (c *client) SayHello(ctx context.Context, req *pb.HelloRequest, opts ...grpc.CallOption) (*pb.HelloResponse, error) {
	conn, err := c.Conn.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pb.NewHelloServiceClient(conn.ClientConn)
	return cli.SayHello(ctx, req, opts...)
}

func (c *client) FindByID(ctx context.Context, req *pb.FindByIDRequest, opts ...grpc.CallOption) (*pb.Greeting, error) {
	conn, err := c.Conn.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pb.NewHelloServiceClient(conn.ClientConn)
	return cli.FindByID(ctx, req, opts...)
}

// Create :nodoc:
func (c *client) Create(ctx context.Context, req *pb.Greeting, opts ...grpc.CallOption) (*pb.Greeting, error) {
	conn, err := c.Conn.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := pb.NewHelloServiceClient(conn.ClientConn)
	return cli.Create(ctx, req, opts...)
}
