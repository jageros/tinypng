package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var cfg *config

const configPath = "./config.toml"

type tinyPng struct {
	APIKey       string `toml:"apiKey"`
	InPutDir     string `toml:"inputDir"`
	OutputDir    string `toml:"outputDir"`
	RenameFormat string `toml:"renameFormat"`
}

type gitee struct {
	Token  string `toml:"token"`
	Owner  string `toml:"owner"`
	Repo   string `toml:"repo"`
	Path   string `toml:"path"`
	Branch string `toml:"branch"`
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
		fmt.Println(err)
	}
}
