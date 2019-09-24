package service

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.kumparan.com/yowez/skeleton-service/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kumparan/go-lib/utils"
	"github.com/kumparan/tapao"
	pb "gitlab.kumparan.com/yowez/skeleton-service/pb/skeleton"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/model"
)

// SayHello :nodoc:
func (s *Service) SayHello(ctx context.Context, req *pb.HelloRequest) (res *pb.HelloResponse, err error) {
	return &pb.HelloResponse{
		Greeting: "Hello " + req.GetName(),
	}, nil
}

// FindByID :nodoc:
func (s *Service) FindByID(ctx context.Context, req *pb.FindByIDRequest) (res *pb.Greeting, err error) {
	greeting, err := s.helloRepo.FindByID(ctx, req.GetId())
	if err != nil {
		log.WithFields(log.Fields{
			"context": utils.Dump(ctx),
			"req":     utils.Dump(req)}).
			Error(err)
		return
	}

	res = greeting.ToProto()
	return
}

// Create :nodoc:
func (s *Service) Create(ctx context.Context, req *pb.Greeting) (result *pb.Greeting, err error) {
	greet := model.NewFromProto(req)

	err = s.helloRepo.Create(ctx, greet)
	if err != nil {
		log.WithFields(log.Fields{
			"context": utils.DumpIncomingContext(ctx),
			"req":     utils.Dump(req)}).
			Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	go func() {
		b, err := tapao.Marshal(event.NatsHelloMessage{
			ID:   greet.ID,
			Body: utils.Dump(greet),
			Type: event.TypeHello,
			Time: time.Now().UTC().Format(time.RFC3339),
		})
		if err != nil {
			log.WithFields(log.Fields{
				"context":  utils.DumpIncomingContext(ctx),
				"greeting": utils.Dump(greet)}).
				Error(err)
			return
		}
		_ = s.nats.Publish(event.NatsHelloChannel, b)
	}()

	err = s.worker.Update(greet)
	if err != nil {
		log.WithFields(log.Fields{
			"context": utils.DumpIncomingContext(ctx),
			"req":     utils.Dump(*req)}).
			Error(err)
		return
	}

	result = greet.ToProto()
	return

}
