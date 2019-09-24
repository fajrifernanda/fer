package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.kumparan.com/yowez/skeleton-service/config"

	log "github.com/sirupsen/logrus"
)

// DB :nodoc:
var DB *gorm.DB

// InitializeCockroachConn :nodoc:
func InitializeCockroachConn() {
	conn, err := gorm.Open("postgres", config.DatabaseDSN())
	if err != nil {
		log.WithField("databaseDSN", config.DatabaseDSN()).Fatal("Failed to connect cockroach database: ", err)
	}

	DB = conn
	log.Info("Connection to Cockroach Server success...")
}
