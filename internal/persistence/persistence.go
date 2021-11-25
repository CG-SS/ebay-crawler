package persistence

import (
	"ebay-crawler/internal/config"
	"ebay-crawler/internal/model"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var app *config.AppConfig

func Init(config *config.AppConfig) {
	app = config
}

func Manager(persistenceChan chan model.ItemModel) {
	jsonFolderFilename := "." + string(filepath.Separator) + "json"

	if !exists(jsonFolderFilename) {
		app.InfoLog.Println("Creating folder ", jsonFolderFilename)
		err := os.Mkdir(jsonFolderFilename, 0777)
		if err != nil {
			panic(err)
		}
	}

	for item := range persistenceChan {
		app.InfoLog.Println("Saving item ", item)
		file, err := json.Marshal(item)
		if err != nil {
			app.ErrorLog.Println("Failed to marshal item ", item)
			continue
		}

		jsonFilename := jsonFolderFilename + string(filepath.Separator) + item.Id + ".json"

		if !exists(jsonFilename) {
			err = ioutil.WriteFile(jsonFilename, file, 0644)
			if err != nil {
				app.ErrorLog.Println("Failed to write item ", item)
				continue
			}
		}
	}
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
