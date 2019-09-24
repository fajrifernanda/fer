package service

import (
	"context"
	"testing"

	pb "gitlab.kumparan.com/yowez/skeleton-service/pb/skeleton"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/mock"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/model"

	"github.com/golang/mock/gomock"
	redcachekeeper "github.com/kumparan/cacher"
	"github.com/stretchr/testify/assert"
	"gitlab.kumparan.com/yowez/skeleton-service/config"
	"gitlab.kumparan.com/yowez/skeleton-service/db"
	"gitlab.kumparan.com/yowez/skeleton-service/repository"
	workerMock "gitlab.kumparan.com/yowez/skeleton-service/worker/mock"
)

func initializeConnection() {
	config.GetConf()
}

func initializeHelloService() *Service {
	initializeConnection()

	s := NewHelloService()
	k := redcachekeeper.NewKeeper()
	k.SetDisableCaching(true)

	helloRepo := repository.NewHelloRepository(db.DB, k)
	s.RegisterHelloRepository(helloRepo)

	return s
}

func TestSayHello(t *testing.T) {
	svc := initializeHelloService()

	ctx := context.TODO()
	res, err := svc.SayHello(ctx, &pb.HelloRequest{
		Name: "Koji Keren",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Hello Koji Keren", res.GetGreeting())
}

func TestFindByID(t *testing.T) {
	svc := initializeHelloService()
	ctx := context.TODO()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelloRepo := mock.NewMockHelloRepository(ctrl)
	mockHelloRepo.EXPECT().FindByID(ctx, gomock.Any()).MinTimes(1).MaxTimes(1).Return(&model.Greeting{ID: 1, Name: "Koji Keren"}, nil)

	svc.RegisterHelloRepository(mockHelloRepo)

	res, err := svc.FindByID(ctx, &pb.FindByIDRequest{
		Id: 1,
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(1), res.GetId())
	assert.Equal(t, "Koji Keren", res.GetName())
}

func TestService_Create(t *testing.T) {
	svc := initializeHelloService()
	ctx := context.TODO()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelloRepo := mock.NewMockHelloRepository(ctrl)
	mockHelloRepo.EXPECT().Create(ctx, gomock.Any()).
		MinTimes(1).
		MaxTimes(1).
		Return(nil)

	wrk := workerMock.NewMockWorker(ctrl)
	svc.RegisterWorker(wrk)

	svc.RegisterHelloRepository(mockHelloRepo)

	wrk.EXPECT().Update(gomock.Any()).Times(1).
		DoAndReturn(func(g *model.Greeting) error {
			return nil
		})

	res, err := svc.Create(ctx, &pb.Greeting{})

	assert.NoError(t, err)
	assert.NotEqual(t, 0, res.Id)
}
