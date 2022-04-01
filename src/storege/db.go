package storege

import "C"
import (
	"github.com/yiuked/go-novel/src/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func init() {
	dataDir := utils.DataDir() + string(os.PathSeparator) + "sqlite"
	_, err := os.Stat(dataDir)
	if err != nil {
		err := os.MkdirAll(dataDir, 0777)
		if err != nil {
			panic(err)
			return
		}
	}
	dataFile := dataDir + string(os.PathSeparator) + "gorm.db"
	DB, err = gorm.Open(sqlite.Open(dataFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migrate()
}
