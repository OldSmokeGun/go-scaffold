package vfs

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _filesystemea5977906e2e7c4ab2e8f1ecb13cb0b0cb1c212f = "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title></title>\r\n</head>\r\n<body>\r\n<h1>Hello</h1>\r\n</body>\r\n<script src=\"https://cdnjs.cloudflare.com/ajax/libs/axios/0.21.0/axios.min.js\"></script>\r\n<script>\r\n    const search = new URLSearchParams(window.location.search)\r\n    const code = search.get(\"code\")\r\n    const pc = search.get(\"pc\")\r\n\r\n    if (!code) {\r\n        const redirectUrl = \"https://test.business.zhzf.zsyh99.com?pc=\" + pc\r\n        window.location.href = 'https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx0c27bccd855983fe&redirect_uri=' + redirectUrl + '=&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect'\r\n\r\n    }\r\n\r\n    axios.get(\"/api/v1/payment/index?code=\" + code + \"&pc=\" + pc).then((res) => {\r\n        console.log(res)\r\n    })\r\n\r\n\r\n    console.log(search.get(\"code\"))\r\n</script>\r\n</html>"

// filesystem returns go-assets FileSystem
var filesystem = assets.NewFileSystem(map[string][]string{"/": []string{"app"}, "/app": []string{"templates"}, "/app/templates": []string{"index.html"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ff,
		Mtime:    time.Unix(1607506272, 1607506272516804100),
		Data:     nil,
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ff,
		Mtime:    time.Unix(1607501897, 1607501897541716700),
		Data:     nil,
	}, "/app/templates": &assets.File{
		Path:     "/app/templates",
		FileMode: 0x800001ff,
		Mtime:    time.Unix(1607503493, 1607503493503841600),
		Data:     nil,
	}, "/app/templates/index.html": &assets.File{
		Path:     "/app/templates/index.html",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1607485251, 1607485251192307100),
		Data:     []byte(_filesystemea5977906e2e7c4ab2e8f1ecb13cb0b0cb1c212f),
	}}, "")
