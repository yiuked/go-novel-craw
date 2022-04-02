package storege

import "C"
import (
	"github.com/yiuked/go-novel/src/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"sync"
)

var DBClient *gorm.DB

var mu sync.Mutex

func ConnectDB() {
	dataDir := utils.GetDataDir() + string(os.PathSeparator) + "sqlite"
	_, err := os.Stat(dataDir)
	if err != nil {
		err := os.MkdirAll(dataDir, 0777)
		if err != nil {
			panic(err)
			return
		}
	}
	dataFile := dataDir + string(os.PathSeparator) + "gorm.db"
	DBClient, err = gorm.Open(sqlite.Open(dataFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := DBClient.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(20)
	DBClient = DBClient.Debug()
	migrate()
}

func DB() *gorm.DB {
	mu.Lock()
	defer mu.Unlock()
	if DBClient == nil {
		ConnectDB()
	}
	return DBClient
}
