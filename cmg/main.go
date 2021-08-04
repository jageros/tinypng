package main

import (
	"cimage/config"
	"cimage/consts"
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
	"sync"
	"time"
)

func main() {
	needDel := false
	var inputDir string
	cfg := config.GetConfig()
	flag.BoolVar(&needDel, "d", false, "删除源文件")
	flag.StringVar(&inputDir, "dir", cfg.TinyPng.InPutDir, "源图片路径")
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
		log.Fatalf("ReadDir error: %v", err)
	}
	if len(files) <= 0 {
		log.Fatalf("Dir %s not has file", inDir)
	}
	if !strings.HasSuffix(inDir, consts.DirField) {
		inDir = inDir + consts.DirField
	}
	if outDir != "" && !strings.HasSuffix(outDir, consts.DirField) {
		outDir = outDir + consts.DirField
	}

	imgUrlCh := make(chan string, 10)
	g := &sync.WaitGroup{}
	for i, f := range files {
		inputFilename := f.Name()
		var fileType string
		switch {
		case strings.HasSuffix(inputFilename, ".jpg"), strings.HasSuffix(inputFilename, ".JPG"):
			fileType = ".jpg"
		case strings.HasSuffix(inputFilename, ".jpeg"), strings.HasSuffix(inputFilename, ".JPEG"):
			fileType = ".jpeg"
		case strings.HasSuffix(inputFilename, ".png"), strings.HasSuffix(inputFilename, ".PNG"):
			fileType = ".png"
		default:
			continue
		}
		var outputFilename string
		if outputFilenameFormat != "" {
			outputFilename = time.Now().Format(outputFilenameFormat) + strconv.Itoa(i) + fileType
		} else {
			outputFilename = inputFilename
		}

		g.Add(1)
		go compImage(inDir, outDir, inputFilename, outputFilename, needDel, imgUrlCh, g)
	}
	g2 := &sync.WaitGroup{}
	g2.Add(1)
	go func() {
		defer g2.Done()
		var imgUrls []string
		for imgUrl := range imgUrlCh {
			imgUrls = append(imgUrls, imgUrl)
		}
		if len(imgUrls) > 0 {
			genhtml.WriteUrlsToFile(imgUrls)
			gitee.BuildIndexHtmlToGitee()
		}
	}()

	g.Wait()
	close(imgUrlCh)
	g2.Wait()
	log.Printf("All Compress End!")
}

func compImage(inPath, outPath, inputFilename, outputFilename string, needDel bool, imgUrlCh chan string, g *sync.WaitGroup) {
	defer g.Done()
	start := time.Now()
	iPath := fmt.Sprintf("%s%s", inPath, inputFilename)
	log.Printf("Start Compress: %s ...", iPath)
	source, err := tinypng.FromFile(iPath)
	if err != nil {
		log.Print(err)
		return
	}
	content, err := source.ToBase64Str()
	if err != nil {
		log.Print(err)
		return
	}
	imgUrl := gitee.PushToGitee(content, outputFilename)
	if outPath != "" {
		_, err := os.Stat(outPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(outPath, os.ModePerm)
			if err != nil {
				log.Printf("create path err: %v", err)
			}
		}
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
	log.Printf("Compress successful! (takes %fs)", takeTime)
	if imgUrl != "" {
		log.Printf("imgUrl=%s", imgUrl)
		imgUrlCh <- imgUrl
	}
}
