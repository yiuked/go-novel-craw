package storege

import "C"
import (
	"github.com/yiuked/go-novel/src/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DBClient *gorm.DB

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
	migrate()
}

func DB() *gorm.DB {
	if DBClient == nil {
		ConnectDB()
	}
	return DBClient
}
