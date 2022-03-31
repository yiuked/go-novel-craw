package storege

import "C"
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func init() {
	var dir string
	var err error
	if len(os.Getenv("CRAW_DATA_DIR")) > 0 {
		dir = os.Getenv("CRAW_DATA_DIR")
	} else {
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	dataDir := dir + string(os.PathSeparator) + "sqlite"
	_, err = os.Stat(dataDir)
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
