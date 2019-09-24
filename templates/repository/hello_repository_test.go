package repository

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	redcachekeeper "github.com/kumparan/cacher"
	"github.com/stretchr/testify/assert"
	"gitlab.kumparan.com/yowez/skeleton-service/repository/model"
)

var mock sqlmock.Sqlmock

func initializeHelloRepository() (hr HelloRepository, mm sqlmock.Sqlmock) {
	initializeConnection()
	k := redcachekeeper.NewKeeper()
	k.SetDisableCaching(true)
	db, mm := initializeCockroachMockConn()
	return NewHelloRepository(db, k), mm
}

func TestHelloRepo_FindByID(t *testing.T) {
	hr, sm := initializeHelloRepository()
	greeting := &model.Greeting{
		ID:   123,
		Name: "Skeleton service",
	}
	ctx := context.TODO()

	queryResult := sqlmock.NewRows([]string{"id", "name", "created_at"}).
		AddRow(greeting.ID, greeting.Name, greeting.CreatedAt)
	sm.ExpectQuery("^SELECT .+ FROM \"greetings\"").WillReturnRows(queryResult)
	sm.ExpectQuery("^SELECT .+ FROM \"greetings\"").WillReturnError(gorm.ErrRecordNotFound)

	res, err := hr.FindByID(ctx, greeting.ID)
	assert.NoError(t, err)
	assert.Equal(t, res.ID, greeting.ID)
	assert.Equal(t, res.Name, greeting.Name)

	_, err = hr.FindByID(ctx, 0)
	assert.NoError(t, err)
}

func TestRepo_Create(t *testing.T) {
	hr, sm := initializeHelloRepository()

	greeting := &model.Greeting{ID: 12345, Name: "Iwan Keren"}

	sm.ExpectBegin()
	queryResult := sm.NewRows([]string{"id"}).AddRow(greeting.ID)
	sm.ExpectQuery("INSERT INTO \"greetings\"").WithArgs(sqlmock.AnyArg(), greeting.Name, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(queryResult)
	sm.ExpectCommit()

	ctx := context.TODO()

	err := hr.Create(ctx, greeting)

	assert.NoError(t, err)
	assert.Equal(t, "Iwan Keren", greeting.Name)
}
