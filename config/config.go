package config

import (
	"cimage/consts"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var cfg *config

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
	if _, err := toml.DecodeFile(consts.ConfigPath, &cfg); err != nil {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("Create config file err=%v", err)
		}
		defer f.Close()
		f.WriteString(consts.ConfigTemplate)
		var path string
		path, err = os.Getwd()
		if err == nil {
			log.Fatalf("Modify Config.toml in file://%s%sconfig.toml to your info !", path, consts.DirField)
		}
		log.Fatalf("Modify Config.toml for your info !")
	}
}
