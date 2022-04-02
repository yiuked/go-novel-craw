package utils

import "os"

var DataDir string

func GetDataDir() string {
	if len(DataDir) <= 0 {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		DataDir = pwd + string(os.PathSeparator) + "storage"
	}
	return DataDir
}
