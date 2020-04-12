package genhtml

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	head = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>图床图片列表</title>
</head>
<body>
<div>
`

	end = `
</div>
</body>
</html>
<style>
    img {
        width: 300px;
        height: auto;
    }

    div {
        float:left;
        margin: 20px;
    }
    p {
        width: 300px;
        word-wrap: break-word;
        word-break: break-all;
        overflow: hidden;
    }
</style>
`

	temple = `
	<div>
        <img src="%s" alt="%s">
        <p>%s</p>
	</div>
`
)

func genByTemple(url string) string {
	strs := strings.Split(url, "/")
	filename := strs[len(strs)-1]
	return fmt.Sprintf(temple, url, filename, url)
}

func genAllContent(urls []string) string {
	result := head
	for _, url := range urls {
		result += genByTemple(url)
	}
	result += end
	return result
}

func writeToHtml(urls []string) {
	f, err := os.Create("index.html")
	if err != nil {
		log.Printf("create file err: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(genAllContent(urls))
	if err != nil {
		log.Printf("write string to html file err: %v", err)
	}
	dir, err := os.Getwd()
	if err == nil {
		log.Printf("You can look through at browser: file://%s/index.html", dir)
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
		if line != "" {
			lines = append(lines, line)
		}
	}
	writeToHtml(lines)
}
