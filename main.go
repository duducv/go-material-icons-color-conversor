package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	root := "./"
	index := 0

	var wg sync.WaitGroup

	fmt.Println("Executando, aguarde...")

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == "materialicons" {
			wg.Add(3)
			finalPath := path + `\` + "24px.svg"
			fileName := fmt.Sprintf("%d", index)

			go func(finalPath string, newName string, color string, outputDestination string, w *sync.WaitGroup) {
				defer wg.Done()
				modifySVG(finalPath, newName, color, outputDestination)
			}(finalPath, fileName, "black", "output-black", &wg)

			go func(finalPath string, newName string, color string, outputDestination string, w *sync.WaitGroup) {
				defer wg.Done()
				modifySVG(finalPath, newName, color, outputDestination)
			}(finalPath, fileName, "white", "output-white", &wg)

			go func(finalPath string, newName string, color string, outputDestination string, w *sync.WaitGroup) {
				defer wg.Done()
				modifySVG(finalPath, newName, color, outputDestination)
			}(finalPath, fileName, "#9AA0A6", "output-gray", &wg)

			index++
		}
		return nil
	})

	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("finalizado em %f segundos", duration.Seconds())
}

func modifySVG(path string, newName string, color string, outputDestination string) {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Não foi possível importar o arquivo")
	}

	lines := strings.Split(string(file), ">")
	for i, line := range lines {
		if strings.Contains(line, "xmlns") {

			lines[i] = lines[i] + " fill=" + `"` + color + `"` + ">"

		} else {
			lines[i] = lines[i] + "> "
		}
		lines[len(lines)-1] = " "
	}

	output := strings.Join(lines, " ")

	if !verifyIfFolderAlreadyExistsBeforeCreate(outputDestination) {
		if err := os.Mkdir(outputDestination, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}

	err = ioutil.WriteFile("./"+outputDestination+"/"+newName+".svg", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}

}

func verifyIfFolderAlreadyExistsBeforeCreate(outputDestination string) bool {
	dir, err := ioutil.ReadDir("./")

	if err != nil {
		log.Fatal("Erro ao verificar se a pasta já foi criada")
	}

	for _, file := range dir {
		if file.Name() == outputDestination {
			return true
		}
	}
	return false
}
