package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var cfg *config

const configPath = "./config.toml"

var configTemplate = `[tinyPng]
apiKey = "pJ3FQxxxxxxxxxxxxxxxxxxxxxxzn8Tt"   #tinyPNG申请的key
inputDir = "./imgs/"     #需要压缩的图片所在的目录
#outputDir = "./cpimgs/"   #设置路径后会将压缩后的图片输出到该路径下
renameFormat = "20060102150405"   # 重命名图片的命名格式， 年月日时分秒的值必须是 6-1-3-4-5 （即：2006年01月02日 15:04:05）

[gitee]
token = "adf7dxxxxxxxxxxxxxxxxxxxxxxxxx5c3f"  # 码云的私人密钥， 不设置此项默认不推送到码云
owner = "Your Name"  # 码云的用户名
repo = "imgs"  # 仓库名
pathFormat = "20060102"  # 上传至仓库中文件夹名命名格式， 年月日时分秒的值必须是 6-1-3-4-5 （即：2006年01月02日 15:04:05）
#fixedPath = "img" #仓库中固定文件夹，设置此项后 pathFormat项无效
branch = "master"  # 分支名`

type tinyPng struct {
	APIKey       string `toml:"apiKey"`
	InPutDir     string `toml:"inputDir"`
	OutputDir    string `toml:"outputDir"`
	RenameFormat string `toml:"renameFormat"`
}

type gitee struct {
	Token      string `toml:"token"`
	Owner      string `toml:"owner"`
	Repo       string `toml:"repo"`
	PathFormat string `toml:"pathFormat"`
	FixedPath  string `toml:"fixedPath"`
	Branch     string `toml:"branch"`
}

type config struct {
	TinyPng *tinyPng `toml:"tinyPng"`
	Gitee   *gitee   `toml:"gitee"`
}

func GetConfig() *config {
	return cfg
}

func init() {
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("Create config file err=%v", err)
		}
		defer f.Close()
		f.WriteString(configTemplate)
		var path string
		path, err = os.Getwd()
		if err == nil {
			log.Fatalf("Modify Config.toml in file://%s/config.toml to your info !", path)
		}
		log.Fatalf("Modify Config.toml for your info !")
	}
}
