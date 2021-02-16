package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type settingsType struct {
	Folders []string `json:"folders"`
}

var Config = settingsType{}

func init() {
	file, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
