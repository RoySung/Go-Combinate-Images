// +build !release

package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type settingsType struct {
	Folders []string `json:"folders"`
}

var Config = settingsType{}

func init() {
	// settingsPath, _ := filepath.Abs("./settings.json")
	rootPath := GetRootPath()
	settingsPath := fmt.Sprint(rootPath, "/settings.json")
	fmt.Println(settingsPath)
	file, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetRootPath() string {
	root, _ := os.Getwd()

	return root
}
