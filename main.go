package main

import (
	"cimage/config"
	"fmt"
	"github.com/gwpp/tinify-go/tinify"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	cfg := config.GetConfig()
	Tinify.SetKey(cfg.TinyPng.APIKey)
	compDir(cfg.TinyPng.InPutDir, cfg.TinyPng.OutputDir)
}

func compDir(inDir, outDir string) {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Panicf("ReadDir error: %v", err)
	}
	if len(files) <= 0 {
		log.Panicf("Dir %s not has file", inDir)
	}
	if !strings.HasSuffix(inDir, "/") {
		inDir = inDir + "/"
	}
	if !strings.HasSuffix(outDir, "/") {
		outDir = outDir + "/"
	}
	for _, f := range files {
		fileName := f.Name()
		if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") || strings.HasSuffix(fileName, ".png") {
			compImage(inDir, outDir, fileName)
		}
	}
}

func compImage(inPath, outPath, fileName string) {
	start := time.Now()
	iPath := fmt.Sprintf("%s%s", inPath, fileName)
	log.Printf("Start Compress: %s ...", iPath)
	source, err := Tinify.FromFile(iPath)
	if err != nil {
		log.Print(err)
		return
	}
	oPath := fmt.Sprintf("%s%s", outPath, fileName)
	err = source.ToFile(oPath)
	if err != nil {
		log.Print(err)
		return
	}
	err = os.Remove(iPath)
	if err != nil {
		log.Printf("Delete file err: %v", err)
	}
	takeTime := time.Now().Sub(start).Seconds()
	log.Printf("Compress successful: %s (takes %fs)", oPath, takeTime)
}
