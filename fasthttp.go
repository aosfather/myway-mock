package main

import (
	"encoding/json"
	"encoding/xml"
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
			if api.IsSupportMethod(ParseHttpMethodType(string(ctx.Method()))) {
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
	var input map[string]interface{} = make(map[string]interface{})
	//不支持文件流
	if ctx.Request.IsBodyStream() {
		ctx.Response.SetBodyString("<b>not surpport stream!</b>")
		ctx.Response.SetStatusCode(400)
		return
	}

	//json格式请求处理
	if strings.Contains(contentType, "application/json") {
		body := ctx.Request.Body()
		err := json.Unmarshal(body, &input)
		if err != nil {

		}
		//xml格式请求处理
	} else if strings.Contains(contentType, "text/xml") {
		body := ctx.Request.Body()
		err := xml.Unmarshal(body, &input)
		if err != nil {

		}
		//form方式
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		ctx.Request.PostArgs().VisitAll(func(key, value []byte) {
			input[string(key)] = string(value)
		})
		//普通查询
	} else {
		args := ctx.QueryArgs()
		if args != nil {
			args.VisitAll(func(key, value []byte) {
				input[string(key)] = string(value)
			})
		}
	}

	st := api.Select(ctx.Response.BodyWriter(), input)
	ctx.Response.Header.Set(CONTENT_TYPE, st.GetContentType())
}
