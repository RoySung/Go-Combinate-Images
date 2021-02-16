package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/RoySung/Go-Combinate-Images/settings"
)

var wg sync.WaitGroup

func main() {
	log.Print("Main")

	assetsFiles := getAssetsFiles(settings.Config.Folders)

	// log.Print(assetsFiles)
	filesSets := getCombination(assetsFiles)
	// log.Print(filesSets)
	wg.Add(len(filesSets))
	for _, set := range filesSets {
		go mergeImagesSet(set)
	}
	wg.Wait()
	log.Print("Tasks is Done !")
}

func mergeImagesSet(set []string) {
	defer wg.Done()
	var firstImg image.Image
	var canvasRange image.Rectangle
	var result draw.Image
	var outputFilePath string
	fileRegex := regexp.MustCompile(`([^\/\.]+)\.[^\/]+`)

	for _, path := range set {
		img, err := getImageFromFilePath(path)
		if err != nil {
			log.Fatal(err)
		}
		if firstImg == nil {
			firstImg = img
			canvasRange = img.Bounds()
			result = image.NewRGBA(image.Rect(0, 0, canvasRange.Max.X, canvasRange.Max.Y))
		}

		draw.Draw(result, canvasRange, img, image.ZP, draw.Over)

		fileName := fileRegex.FindStringSubmatch(path)[1]
		if outputFilePath != "" {
			outputFilePath += "-"
		}
		outputFilePath += fileName
	}

	// open file to save
	outputFilePath += ".png"
	outputFilePath = fmt.Sprintf("./assets/output/%s", outputFilePath)
	log.Print(outputFilePath)
	dstFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	// encode as .png to the file
	err = png.Encode(dstFile, result)
	// close the file
	dstFile.Close()
	if err != nil {
		log.Fatal(err)
	}

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

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(f)
	return image, err
}
