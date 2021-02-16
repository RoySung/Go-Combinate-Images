package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/RoySung/Go-Combinate-Images/settings"
)

func main() {
	log.Print("Main")

	assetsFiles := getAssetsFiles(settings.Config.Folders)

	// log.Print(assetsFiles)
	filesSets := getCombination(assetsFiles)
	log.Print(filesSets)
}

func getCombination(data [][]string) [][]string {
	getLength := func(list [][]string) int {
		var length = 1
		for _, items := range list {
			length *= len(items)
		}
		return length
	}
	combinate := func(acc [][]string, list []string) [][]string {
		result := make([][]string, 0)
		if len(acc) == 0 {
			acc = make([][]string, 1)
		}
		for _, preSet := range acc {
			for _, item := range list {
				set := append(preSet, item)
				result = append(result, set)
			}
		}

		return result
	}
	length := getLength(data)

	result := make([][]string, 0, length)
	for _, items := range data {
		result = combinate(result, items)
	}
	return result
}

func getAssetsFiles(folderNames []string) [][]string {
	result := make([][]string, len(folderNames))

	for index, foldeName := range folderNames {
		dirPath := fmt.Sprintf("assets/%s", foldeName)
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
		}

		row := make([]string, len(files))
		for index, file := range files {
			row[index] = fmt.Sprintf("%s/%s", dirPath, file.Name())
		}
		result[index] = row
	}
	return result
}
