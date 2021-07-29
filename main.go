package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"

	"github.com/RoySung/Go-Combinate-Images/settings"
	"github.com/thoas/go-funk"
)

var allowFileTypes = [...]string{".png", ".jpg", ".jpeg"}

func main() {
	log.Print("Main")
	assetsFiles := getAssetsFiles(settings.Config.Folders)
	filesSets := getCombination(assetsFiles)

	// wg.Add(len(filesSets))
	// for _, set := range filesSets {
	// 	go mergeImagesSet(set)
	// }
	// wg.Wait()

	jobsChan := make(chan []string, 1024)
	resultChan := make(chan string, len(filesSets))
	// finishChan := make(chan bool)
	workerCount := runtime.NumCPU()
	fmt.Println("worker count: ", workerCount)

	for i := 0; i < workerCount; i++ {
		go worker(mergeImagesSet, jobsChan, resultChan)
	}

	for _, set := range filesSets {
		jobsChan <- set
	}
	close(jobsChan)

	// <-finishChan

	for i := 0; i < len(filesSets); i++ {
		log.Println("Path: ", <-resultChan)
	}

	log.Print("Tasks is Done !")
}

func worker(process func([]string) string, jobs <-chan []string, result chan<- string) {
	for job := range jobs {
		result <- process(job)
	}
}

func mergeImagesSet(set []string) string {
	var firstImg image.Image
	var canvasRange image.Rectangle
	var result draw.Image
	var outputFilePath string
	fileRegex := regexp.MustCompile(`([^\/\.]+)\.[^\/]+$`)

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
	rootPath := settings.GetRootPath()
	outputFilePath = fmt.Sprintf("%s/assets/output/%s", rootPath, outputFilePath)
	// log.Print(outputFilePath)
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

	return outputFilePath
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
				set := make([]string, 0, len(preSet)+1)
				set = append(set, preSet...)
				set = append(set, item)
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
		rootPath := settings.GetRootPath()
		dirPath := fmt.Sprintf("%s/assets/%s", rootPath, foldeName)
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
		}
		files = funk.Filter(files, func(fileInfo fs.FileInfo) bool {
			return funk.Contains(allowFileTypes, func(fileType string) bool {
				return funk.Contains(fileInfo.Name(), fileType)
			})
		}).([]fs.FileInfo)

		row := make([]string, len(files))
		for index, file := range files {
			row[index] = fmt.Sprintf("%s/%s", dirPath, file.Name())
		}
		// row = funk.Filter(row, func(path string) bool {
		// 	return funk.Contains(allowFileTypes, func(fileType string) bool {
		// 		return funk.Contains(path, fileType)
		// 	})
		// }).([]string)
		result[index] = row
	}

	return result
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(f)
	return image, err
}
