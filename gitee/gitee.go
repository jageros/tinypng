package gitee

import (
	"cimage/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	baseUrl     = "https://gitee.com/api/v5/repos"
	contentType = "application/json;charset=UTF-8"
)

//https://gitee.com/api/v5/repos/jayos/imgs/contents/img/README.md

func PushToGitee(fileContent, filename string) string {
	cfg := config.GetConfig().Gitee
	if cfg.Token == "" {
		return "push to gitee must set the gitee token!"
	}
	path := time.Now().Format(cfg.PathFormat)
	if cfg.FixedPath != "" {
		path = cfg.FixedPath
	}
	url := fmt.Sprintf("%s/%s/%s/contents/%s/%s", baseUrl, cfg.Owner, cfg.Repo, path, filename)
	body := map[string]interface{}{
		"access_token": cfg.Token,
		"content":      fileContent,
		"message":      "上传图片!",
		"branch":       cfg.Branch,
	}
	byt, err := json.Marshal(body)
	if err != nil {
		log.Printf("json.Marshal err: %v", err)
	}
	bodyReader := strings.NewReader(string(byt))
	resp, err := http.Post(url, contentType, bodyReader)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("post err: %v", err)
		return "--"
	}
	if resp.StatusCode == 201 {
		log.Printf("push %s to gitee successfuly!", filename)
	} else {
		log.Printf("push %s to gitee faild StatusCode=%d", filename, resp.StatusCode)
	}
	imgUrl := fmt.Sprintf("https://gitee.com/%s/%s/raw/master/%s/%s", cfg.Owner, cfg.Repo, path, filename)
	return imgUrl
}
