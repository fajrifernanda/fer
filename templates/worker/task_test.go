package worker

import (
	"testing"

	"github.com/gocraft/work"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/mock"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/model"
)

func TestWorkerImpl_UpdateEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHelloRepo := mock.NewMockHelloRepository(ctrl)

	ctx := new(taskContext)
	ctx.helloRepository = mockHelloRepo

	greeting := model.Greeting{
		ID:   1234,
		Name: "task test",
	}
	ctx.greeting = &greeting

	mockHelloRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)

	err := ctx.updateEventHandler(&work.Job{})
	assert.NoError(t, err)
}
