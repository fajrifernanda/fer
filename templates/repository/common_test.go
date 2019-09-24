package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"gitlab.kumparan.com/yowez/skeleton-service/db"
	"os"
	"strconv"

	"github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
	"gitlab.kumparan.com/yowez/skeleton-service/config"
)

func initializeConnection() {
	config.GetConf()

	mockDB, newMock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	mock = newMock

	db.DB, err = gorm.Open("postgres", mockDB)
	if err != nil {
		panic(err)
	}

	db.DB.LogMode(true)
	setupLogger()
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		},
		Line: true,
		File: true,
	}

	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	verbose, _ := strconv.ParseBool(os.Getenv("VERBOSE"))
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func initializeCockroachMockConn() (db *gorm.DB, mock sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	db, err = gorm.Open("postgres", mockDB)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	return
}
