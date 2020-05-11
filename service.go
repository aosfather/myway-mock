package main

import (
	"fmt"
	mvc "github.com/aosfather/bingo_mvc"
	"html/template"
	"io"
	"log"
	"math/rand"
	"time"
)

/*
  服务模拟
*/
const (
	ERR_PARAMETER = "_Parameter"
	ERR_MSG       = "_Msg"
)

var maxRandomDelay = 100

//随机种子
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Service struct {
	meta        *Api
	dataPath    string                   //数据目录
	dataMap     map[string]*ResponseData //模板缓存
	errTemplate *template.Template
}

//清除缓存
func (this *Service) ClearCache() {
	if this.dataMap != nil && len(this.dataMap) > 0 {
		this.dataMap = make(map[string]*ResponseData)
	}
}

func (this *Service) IsSupportMethod(m mvc.HttpMethodType) bool {
	for _, sm := range this.meta.Methods {
		if sm == m {
			return true
		}
	}
	return false
}

//校验输入
func (this *Service) validateInput(writer io.Writer, input map[string]interface{}) (mvc.StyleType, error) {
	for _, p := range this.meta.RequestSet.Fields {
		if !this.validateField(p, input[p.Name]) {
			log.Printf("parameter '%s' not validated!", p.Name)
			errMap := make(map[string]string)
			errMap[ERR_PARAMETER] = p.Name
			errMap[ERR_MSG] = fmt.Sprintf("%s 校验不通过", p.Name)
			this.outError(writer, errMap)
			return this.meta.RequestSet.Style, fmt.Errorf("%s", errMap[ERR_MSG])
		}
	}

	return this.meta.RequestSet.Style, nil
}

func (this *Service) outError(writer io.Writer, err map[string]string) mvc.StyleType {
	if this.errTemplate == nil {
		this.errTemplate = template.New(this.meta.Name + "_error")
		this.errTemplate.Parse(this.meta.RequestSet.Error)
	}
	this.errTemplate.Execute(writer, err)
	return this.meta.RequestSet.Style
}

//校验字段
func (this *Service) validateField(p Paramter, v interface{}) bool {
	//参数无值的情况
	if v == nil {
		if p.Policy == Must {
			//参数必须输入
			return false
		}
		return true
	}
	value := fmt.Sprintf("%v", v)
	return p.Validate(value)
}

//根据参数选择对应的数据结果集
func (this *Service) Select(writer io.Writer, inputfunc func(interface{}) error) mvc.StyleType {
	var input map[string]interface{} = make(map[string]interface{})
	inputfunc(input)
	//随机延时100毫秒
	time.Sleep(time.Millisecond * time.Duration(r.Intn(maxRandomDelay)))
	//校验参数
	st, err := this.validateInput(writer, input)
	if err != nil {
		//返回结果
		return st
	}

	//根据设置的延时时间，随机选取设置的多个值中的一个进行延时处理
	if this.meta.Delay != nil && len(this.meta.Delay) > 0 {
		index := r.Intn(len(this.meta.Delay))
		log.Println("mock service uesed time ", this.meta.Delay[index], " ms")
		time.Sleep(time.Millisecond * time.Duration(this.meta.Delay[index]))
	}

	//根据触发器的条件进行匹配找的完全匹配的结果id
	result := this.match(input)
	log.Printf("match '%s' data", result)
	return this.output(writer, result, input)
}

//检查匹配的条件，获取到与输入的条件最匹配的结果集编号
func (this *Service) match(input map[string]interface{}) string {
	for _, trigger := range this.meta.ResponseSet.Triggers {
		if trigger.IsMatch(input) {
			return trigger.Data
		}
	}

	return this.meta.ResponseSet.Default
}

//输出数据
func (this *Service) output(writer io.Writer, id string, input interface{}) mvc.StyleType {
	data := this.dataMap[id]
	if data == nil {
		filename := this.dataPath + "/" + id + ".yaml"
		data = new(ResponseData)
		data.LoadFromYaml(filename)
		this.dataMap[id] = data
	}
	//格式化模板
	data.Render(writer, input)

	return data.Style
}

//新建服务
func NewService(dataroot string, api *Api) *Service {
	if dataroot != "" && api != nil {
		var s Service = Service{api, dataroot + "/" + api.Name, make(map[string]*ResponseData), nil}
		return &s
	}

	return nil
}
