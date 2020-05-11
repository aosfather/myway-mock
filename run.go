package main

import (
	"github.com/aosfather/bingo_mvc"
	fsd "github.com/aosfather/bingo_mvc/fasthttp"
	"io/ioutil"
	"log"
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
	server   []*fsd.FastHTTPDispatcher
	dispatch *bingo_mvc.DispatchManager
}

func (this *application) init() {
	log.Println("myway-mock service begin start")
	//加载配置文件
	this.config = new(Config)
	this.config.LoadFromYaml("configs.yaml")
	//初始化api
	this.dispatch = new(bingo_mvc.DispatchManager)
	this.dispatch.Init()
	//加载apis目录下的api定义
	this.loadApis(string(os.PathSeparator), "apis")

}

func (this *application) loadApis(pathSeparator string, fileDir string) {
	files, _ := ioutil.ReadDir(fileDir)
	for _, onefile := range files {
		filename := fileDir + pathSeparator + onefile.Name()
		if onefile.IsDir() {
			this.loadApis(pathSeparator, filename)
		} else {
			api := new(Api)
			api.LoadFromYaml(filename)
			if api.Url != "" {
				log.Println("load the define of api:", api.Name)
				serv := NewService("datas", api)
				log.Println(api.Url)
				this.dispatch.AddApi("", api.Name, api.Url, serv)
			}
		}
	}
}

/**
  启动应用
  默认端口 80
  持多个端口监听
*/
func (this *application) start() {

	default_port := 80
	if this.config.Port != nil && len(this.config.Port) > 0 {
		default_port = this.config.Port[0]
	}
	if this.config.Maxdelay != 0 {
		maxRandomDelay = this.config.Maxdelay
	}
	//多个端口时候，对第一个以为的端口使用协程启动
	if len(this.config.Port) > 1 {
		//启动服务
		for index, port := range this.config.Port {
			if index == 0 {
				continue
			}
			log.Println("start server in port: ", port)
			go this.runServer(port)
		}
		log.Println("start server in port: ", default_port)
		this.runServer(default_port)
	}

}

//启动端口监听服务，提供http服务
func (this *application) runServer(port int) {
	server := new(fsd.FastHTTPDispatcher)
	server.Port = port
	server.SetDispatchManager(this.dispatch)
	//server.port = port
	this.server = append(this.server, server)
	server.Run()
}

func (this *application) shutdown() {

}
