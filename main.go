package main

import (
	"cimage/config"
	"cimage/gitee"
	"cimage/tinypng"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	needDel := false
	flag.BoolVar(&needDel, "d", false, "-d: 表示删除源文件")
	flag.Parse()
	if needDel {
		log.Printf("Delete src file end of compress!")
	}
	cfg := config.GetConfig()
	tinypng.SetKey(cfg.TinyPng.APIKey)
	compDir(cfg.TinyPng.InPutDir, cfg.TinyPng.OutputDir, cfg.TinyPng.RenameFormat, needDel)
}

func compDir(inDir, outDir, outputFilenameFormat string, needDel bool) {
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
	for i, f := range files {
		inputFilename := f.Name()
		var fileType string
		switch {
		case strings.HasSuffix(inputFilename, ".jpg"):
			fileType = ".jpg"
		case strings.HasSuffix(inputFilename, ".jpeg"):
			fileType = ".jpeg"
		case strings.HasSuffix(inputFilename, ".png"):
			fileType = ".png"
		default:
			continue
		}
		outputFilename := time.Now().Format(outputFilenameFormat) + strconv.Itoa(i) + fileType
		compImage(inDir, outDir, inputFilename, outputFilename, needDel)

	}
}

func compImage(inPath, outPath, inputFilename, outputFilename string, needDel bool) {
	start := time.Now()
	iPath := fmt.Sprintf("%s%s", inPath, inputFilename)
	log.Printf("Start Compress: %s ...", iPath)
	source, err := tinypng.FromFile(iPath)
	if err != nil {
		log.Print(err)
		return
	}
	//oPath := fmt.Sprintf("%s%s", outPath, outputFilename)
	content, err := source.ToBase64Str()
	if err != nil {
		log.Print(err)
		return
	}
	imgUrl := gitee.PushToGitee(content, outputFilename)
	if needDel {
		err = os.Remove(iPath)
		if err != nil {
			log.Printf("Delete file err: %v", err)
		}
	}
	takeTime := time.Now().Sub(start).Seconds()
	log.Printf("Compress successful: url: %s (takes %fs)", imgUrl, takeTime)
}

func pushToGitee(content string) {

}
