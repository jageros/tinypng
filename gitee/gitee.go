package gitee

import (
	"cimage/config"
	"cimage/genhtml"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	baseUrl     = "https://gitee.com/api/v5/repos"
	contentType = "application/json;charset=UTF-8"
)

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
	if err != nil {
		log.Printf("post err: %v", err)
		return "--"
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 {
		log.Printf("push %s to gitee successfuly!", filename)
	} else {
		log.Printf("push %s to gitee faild StatusCode=%d", filename, resp.StatusCode)
	}
	imgUrl := fmt.Sprintf("https://gitee.com/%s/%s/raw/master/%s/%s", cfg.Owner, cfg.Repo, path, filename)
	return imgUrl
}

func createFileToGitee(content, filePath string) {

}

func UpdateFile(filePath, content string) error {
	cfg := config.GetConfig().Gitee
	if cfg.Token == "" {
		return errors.New("not set token")
	}
	sha := getFileSha(filePath)
	if sha == "" {
		//log.Printf("Error: get file sha = nil !")
		return errors.New("get file sha failed")
	}
	method := "PUT"
	if sha == "404" {
		method = "POST"
		sha = ""
	}
	// https://gitee.com/api/v5/repos/jayos/imgs/contents/index.html
	url := fmt.Sprintf("%s/%s/%s/contents/%s", baseUrl, cfg.Owner, cfg.Repo, filePath)
	msg := "update " + filePath + " !"
	body := map[string]interface{}{
		"access_token": cfg.Token,
		"content":      content,
		"sha":          sha,
		"message":      msg,
		"branch":       cfg.Branch,
	}
	byt, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyReader := strings.NewReader(string(byt))
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("update %s successfuly!", filePath)
	} else {
		errMsg := fmt.Sprintf("Update %s faild StatusCode=%d", filePath, resp.StatusCode)
		return errors.New(errMsg)
	}
	return nil
}

func getFileSha(filePath string) string {
	cfg := config.GetConfig().Gitee
	if cfg.Token == "" {
		return ""
	}
	// https://gitee.com/api/v5/repos/jayos/imgs/contents/index.html?access_token=adf7d13db9231626f5268292bc645c3f&ref=master
	url := fmt.Sprintf("%s/%s/%s/contents/%s?access_token=%s&ref=%s", baseUrl, cfg.Owner, cfg.Repo, filePath, cfg.Token, cfg.Branch)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Get request err: %v", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return "404"
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}
	type shaSt struct {
		Sha string `json:"sha"`
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("get reply read body err: %v", err)
		return ""
	}
	var sha = shaSt{}
	err = json.Unmarshal(body, &sha)
	if err != nil {
		log.Printf("Unmarshal sha err: %v", err)
		return ""
	}
	return sha.Sha
}

func GetAllPicUrl() []string {
	var urls []string
	cfg := config.GetConfig().Gitee
	if cfg.Token == "" {
		return urls
	}
	type treeSt struct {
		PATH string `json:"path"`
	}
	type pathData struct {
		Trees []*treeSt `json:"tree"`
	}
	var ts = pathData{}
	// https://gitee.com/api/v5/repos/jayos/imgs/git/trees/master?access_token=adf7d13db9231626f5268292bc645c3f&recursive=1
	url := fmt.Sprintf("%s/%s/%s/git/trees/%s?access_token=%s&recursive=1", baseUrl, cfg.Owner, cfg.Repo, cfg.Branch, cfg.Token)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Get request err: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Get all Pic url read body data err: %v", err)
	}

	err = json.Unmarshal(body, &ts)
	if err != nil {
		log.Printf("json unmarshal err: %v", err)
	}

	for _, t := range ts.Trees {
		if strings.HasSuffix(t.PATH, ".jpg") || strings.HasSuffix(t.PATH, ".png") || strings.HasSuffix(t.PATH, ".jpeg") {
			imgUrl := fmt.Sprintf("https://gitee.com/%s/%s/raw/master/%s", cfg.Owner, cfg.Repo, t.PATH)
			urls = append(urls, imgUrl)
		}
	}
	return urls
}

func BuildGiteePage() {
	cfg := config.GetConfig().Gitee
	if cfg.Token == "" {
		return
	}
	url := fmt.Sprintf("%s/%s/%s/pages/builds", baseUrl, cfg.Owner, cfg.Repo)
	body := map[string]interface{}{
		"access_token": cfg.Token,
	}
	byt, err := json.Marshal(body)
	if err != nil {
		log.Printf("json.Marshal err: %v", err)
	}
	bodyReader := strings.NewReader(string(byt))
	resp, err := http.Post(url, contentType, bodyReader)
	if err != nil {
		log.Printf("post err: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		log.Printf("Build giteePage successfuly, website: http://%s.gitee.io/%s", cfg.Owner, cfg.Repo)
	} else {
		log.Printf("Build giteePage faild StatusCode=%d", resp.StatusCode)
	}
}

func BuildIndexHtmlToGitee() {
	urls := GetAllPicUrl()
	if len(urls) <= 0 {
		log.Printf("not has pic url !")
		return
	}
	content := genhtml.GenAllContent(urls)
	contentBase64 := base64.StdEncoding.EncodeToString([]byte(content))
	err := UpdateFile("index.html", contentBase64)
	if err != nil {
		log.Printf("update gitee file err: %v", err)
	}
	BuildGiteePage()
}