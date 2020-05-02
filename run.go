package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	app := application{}
	//初始化
	app.init()
	//启动
	app.start()

}

type application struct {
	config   *Config
	server   *HttpServer
	dispatch *DispatchManager
}

func (this *application) init() {
	//加载配置文件
	this.config = new(Config)
	this.config.LoadFromYaml("configs.yaml")
	//初始化api
	this.dispatch = new(DispatchManager)
	this.dispatch.Init()
	//加载apis目录下的api定义
	this.loadApis(string(os.PathSeparator), "apis")

}

func (this *application) loadApis(pathSeparator string, fileDir string) {
	files, _ := ioutil.ReadDir(fileDir)
	for _, onefile := range files {
		filename := fileDir + pathSeparator + onefile.Name()
		if onefile.IsDir() {
			//fmt.Println(tmpPrefix, onefile.Name(), "目录:")
			this.loadApis(pathSeparator, filename)
		} else {
			api := new(Api)
			api.LoadFromYaml(filename)
			fmt.Println(api)
			if api.Url != "" {
				this.dispatch.AddApi("", "datas", api)
			}
		}
	}
}

func (this *application) start() {
	//启动服务
	this.server = new(HttpServer)
	this.server.dispatch = this.dispatch
	this.server.port = this.config.Port
	this.server.Start()

}

func (this *application) shutdown() {

}
