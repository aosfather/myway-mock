package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strings"
)

/*
   基于fasthttp的实现
*/

type HttpServer struct {
	port     int
	server   *fasthttp.Server
	dispatch *DispatchManager
}

func (this *HttpServer) Start() {
	this.server = &fasthttp.Server{Handler: this.ServeHTTP}

	if this.port <= 0 {
		this.port = 80
	}

	addr := fmt.Sprintf("0.0.0.0:%d", this.port)
	this.server.ListenAndServe(addr)

}

func (this *HttpServer) ServeHTTP(ctx *fasthttp.RequestCtx) {
	//获取访问的url
	url := string(ctx.Request.URI().RequestURI())
	domain := string(ctx.Request.Header.Host())
	//处理favicon.ico
	if url == "/favicon.ico" {
		ico, _ := ioutil.ReadFile("favicon.ico")
		ctx.Response.Header.Set(CONTENT_TYPE, "image/x-icon")
		ctx.Response.SetBodyRaw(ico)
	} else {
		//通过dispatch，获取api的定义
		api := this.dispatch.GetApi(domain, url)
		if api == nil { //不存在的时候的处理
			ctx.Response.Header.Set(CONTENT_TYPE, "text/html;charset=utf-8")
			ctx.Response.SetBodyString("<b>the url not found!</b>")
			ctx.Response.SetStatusCode(404)

		} else {
			//检查http method 看是否支持该类型
			if api.IsSupportMethod(getMethod(ctx)) {
				this.call(api, ctx)
			} else {
				//不支持的 http method 处理
				ctx.Response.Header.Set(CONTENT_TYPE, "text/html;charset=utf-8")
				ctx.Response.SetBodyString("<b>the method not support !</b>")
				ctx.Response.SetStatusCode(405)
			}

		}

	}

	//设置服务器名称
	ctx.Response.Header.Set("Server", SERVER_NAME)
}

const CONTENT_TYPE = "Content-Type"
const SERVER_NAME = "nginx 1.3.0"

func (this *HttpServer) call(api *Service, ctx *fasthttp.RequestCtx) {
	//校验参数
	contentType := string(ctx.Request.Header.ContentType())
	fmt.Println(contentType)
	var input map[string]interface{} = make(map[string]interface{})
	if ctx.Request.IsBodyStream() {
		fmt.Println("body has content")
	}
	if strings.Contains(contentType, "application/json") {

	} else if strings.Contains(contentType, "text/xml") {

	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		ctx.Request.PostArgs().VisitAll(func(key, value []byte) {
			input[string(key)] = string(value)
		})

	} else {
		args := ctx.QueryArgs()
		if args != nil {
			args.VisitAll(func(key, value []byte) {
				input[string(key)] = string(value)
			})
		}
	}

	var st StyleType
	//校验参数
	st, err := api.ValidateInput(ctx.Response.BodyWriter(), input)
	if err == nil {
		//返回结果
		st = api.Select(ctx.Response.BodyWriter(), input)
	}

	switch st {
	case Json:
		ctx.Response.Header.Set(CONTENT_TYPE, "application/json; charset=utf-8")
	case Xml:
		ctx.Response.Header.Set(CONTENT_TYPE, "text/xml;charset=utf-8")
	default:
		ctx.Response.Header.Set(CONTENT_TYPE, "text/html;charset=utf-8")
	}
}

func getMethod(ctx *fasthttp.RequestCtx) HttpMethodType {
	method := string(ctx.Method())
	fmt.Println(method)
	switch method {
	case "GET":
		return Get
	case "POST":
		return Post
	case "PUT":
		return Put
	case "DELETE":
		return Del
	}
	return Get
}
