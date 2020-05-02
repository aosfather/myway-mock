package main

import (
	"io"
)

/*
  服务模拟
*/
type Service struct {
	meta     *Api
	dataPath string                   //数据目录
	dataMap  map[string]*ResponseData //模板缓存
}

//清除缓存
func (this *Service) ClearCache() {
	if this.dataMap != nil && len(this.dataMap) > 0 {
		this.dataMap = make(map[string]*ResponseData)
	}
}

//校验输入
func (this *Service) validateInput(input interface{}) {

}

//校验字段
func (this *Service) validateField(p *Paramter, v string) bool {
	//参数无值的情况
	if v == "" {
		if p.Policy == Must {
			//参数必须输入
			return false
		}
		return true
	}

	return p.Validate(v)
}

//输出数据
func (this *Service) output(writer io.Writer, id string, input interface{}) StyleType {
	data := this.dataMap[id]
	if data == nil {
		filename := this.dataPath + "/" + id + ".yaml"
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
		var s Service = Service{api, dataroot + "/" + api.Name, make(map[string]*ResponseData)}
		return &s
	}

	return nil
}
