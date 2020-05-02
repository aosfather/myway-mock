package main

import (
	"fmt"
	"html/template"
	"io"
)

/*
  服务模拟
*/
const (
	ERR_PARAMETER = "_Parameter"
	ERR_MSG       = "_Msg"
)

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

func (this *Service) IsSupportMethod(m HttpMethodType) bool {
	for _, sm := range this.meta.Methods {
		if sm == m {
			return true
		}
	}
	return false
}

//校验输入
func (this *Service) ValidateInput(writer io.Writer, input map[string]interface{}) (StyleType, error) {
	for _, p := range this.meta.RequestSet.Fields {
		if !this.validateField(&p, input[p.Name]) {
			fmt.Println("not validate", p.Name)
			errMap := make(map[string]string)
			errMap[ERR_PARAMETER] = p.Name
			errMap[ERR_MSG] = fmt.Sprintf("%s 校验不通过", p.Name)
			this.outError(writer, errMap)
			return this.meta.RequestSet.Style, fmt.Errorf("%s")
		}
	}

	return this.meta.RequestSet.Style, nil
}

func (this *Service) outError(writer io.Writer, err map[string]string) StyleType {
	if this.errTemplate == nil {
		this.errTemplate = template.New(this.meta.Name + "_error")
		this.errTemplate.Parse(this.meta.RequestSet.Error)
	}
	fmt.Println(err)
	this.errTemplate.Execute(writer, err)
	return this.meta.RequestSet.Style
}

//校验字段
func (this *Service) validateField(p *Paramter, v interface{}) bool {
	//参数无值的情况
	if v == nil {
		if p.Policy == Must {
			//参数必须输入
			return false
		}
		return true
	}
	value := fmt.Sprintf("%s", v)
	return p.Validate(value)
}

//根据参数选择对应的数据结果集
func (this *Service) Select(writer io.Writer, input map[string]interface{}) StyleType {
	result := this.meta.ResponseSet.Default
	//根据触发器的条件进行匹配找的完全匹配的结果id

	return this.output(writer, result, input)
}

//输出数据
func (this *Service) output(writer io.Writer, id string, input interface{}) StyleType {
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
