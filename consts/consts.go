package consts

import "runtime"

const ConfigPath = "config.toml"

var DirField = "/"

func init() {
	if runtime.GOOS == "windows" {
		DirField = "\\"
	}
}

var (
	ConfigTemplate = `[tinyPng]
apiKey = "pJ3FQxxxxxxxxxxxxxxxxxxxxxxzn8Tt"   #tinyPNG申请的key
inputDir = "imgs"     #需要压缩的图片所在的目录
#outputDir = "cpimgs"   #设置路径后会将压缩后的图片输出到该路径下
renameFormat = "20060102150405"   # 重命名图片的命名格式， 年月日时分秒的值必须是 6-1-3-4-5 （即：2006年01月02日 15:04:05）

[gitee]
token = "adf7dxxxxxxxxxxxxxxxxxxxxxxxxx5c3f"  # 码云的私人密钥， 不设置此项默认不推送到码云
owner = "Your Name"  # 码云的用户名
repo = "imgs"  # 仓库名
pathFormat = "20060102"  # 上传至仓库中文件夹名命名格式， 年月日时分秒的值必须是 6-1-3-4-5 （即：2006年01月02日 15:04:05）
#fixedPath = "img" #仓库中固定文件夹，设置此项后 pathFormat项无效
branch = "master"  # 分支名`

	HeadHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>图床图片列表</title>
</head>
<body>
<div>
`

	EndHtml = `
</div>
</body>
</html>
<style>
    img {
        width: 300px;
        height: 180px;
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

	DivTemple = `
	<div>
        <img src="%s" alt="%s">
        <p>%s</p>
	</div>
`
)
