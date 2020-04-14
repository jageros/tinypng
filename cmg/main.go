package main

import (
	"cimage/config"
	"cimage/genhtml"
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
	var inputDir string
	cfg := config.GetConfig()
	flag.BoolVar(&needDel, "d", false, "-d: 表示删除源文件")
	flag.StringVar(&inputDir, "dir", cfg.TinyPng.InPutDir, "-dir ./imgs  表示源图片路径")
	flag.Parse()
	if needDel {
		log.Printf("Delete src file end of compress!")
	}
	tinypng.SetKey(cfg.TinyPng.APIKey)
	compDir(inputDir, cfg.TinyPng.OutputDir, cfg.TinyPng.RenameFormat, needDel)
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
	if outDir != "" && !strings.HasSuffix(outDir, "/") {
		outDir = outDir + "/"
	}
	var imgUrls []string
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
		imUrl := compImage(inDir, outDir, inputFilename, outputFilename, needDel)
		if imUrl != "" {
			imgUrls = append(imgUrls, imUrl)
		}
	}
	genhtml.WriteUrlsToFile(imgUrls)
	gitee.BuildIndexHtmlToGitee()
}

func compImage(inPath, outPath, inputFilename, outputFilename string, needDel bool) string {
	start := time.Now()
	iPath := fmt.Sprintf("%s%s", inPath, inputFilename)
	log.Printf("Start Compress: %s ...", iPath)
	source, err := tinypng.FromFile(iPath)
	if err != nil {
		log.Print(err)
		return ""
	}
	content, err := source.ToBase64Str()
	if err != nil {
		log.Print(err)
		return ""
	}
	imgUrl := gitee.PushToGitee(content, outputFilename)
	if outPath != "" {
		oPath := fmt.Sprintf("%s%s", outPath, outputFilename)
		err = source.ToFile(oPath)
		if err != nil {
			log.Printf("source to file err: %v", err)
		}
	}
	if needDel {
		err = os.Remove(iPath)
		if err != nil {
			log.Printf("Delete file err: %v", err)
		}
	}
	takeTime := time.Now().Sub(start).Seconds()
	log.Printf("Compress successful: url: %s (takes %fs)", imgUrl, takeTime)
	return imgUrl
}
