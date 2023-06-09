package genhtml

import (
	"bufio"
	"cimage/consts"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func genByTemple(url string) string {
	strs := strings.Split(url, "/")
	filename := strs[len(strs)-1]
	return fmt.Sprintf(consts.DivTemple, url, filename, url)
}

func GenAllContent(urls []string) string {
	result := consts.HeadHtml
	for i := len(urls) - 1; i >= 0; i-- {
		result += genByTemple(urls[i])
	}
	result += consts.EndHtml
	return result
}

func writeToHtml(content string) {
	f, err := os.Create("index.html")
	if err != nil {
		log.Printf("create file err: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		log.Printf("write string to html file err: %v", err)
	}
	dir, err := os.Getwd()
	if err == nil {
		log.Printf("You can look through at browser: file://%s%sindex.html", dir, consts.DirField)
	}
}

func WriteUrlsToFile(urls []string) {
	f, err := os.OpenFile("img_urls.txt", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		f, err = os.Create("img_urls.txt")
		if err != nil {
			log.Printf("create file err: %v", err)
		}
	}
	defer f.Close()
	for _, url := range urls {
		content := url + "\n"
		_, err = f.WriteString(content)
		if err != nil {
			log.Printf("write url to txt err: %v", err)
		}
	}
	f.Seek(0, 0)
	rd := bufio.NewReader(f)
	var lines []string
	for {
		line, err := rd.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		strs := strings.Split(line, "\n")
		line = strs[0]
		if line != "" {
			lines = append(lines, line)
		}
	}
	content := GenAllContent(lines)
	writeToHtml(content)
}
