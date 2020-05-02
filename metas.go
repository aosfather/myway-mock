package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"text/template"
)

/*
  元信息定义

*/
//校验器类型列表
var validators map[string]*ValidatorType = make(map[string]*ValidatorType)
var types map[string]*Type = make(map[string]*Type)

//注册校验器
func init() {
	fmt.Println("init validators")
	//初始化
	validators["regex"] = &ValidatorType{"regex", RegexpValidate}
}

type Validate func(value string, option string) bool

//校验器类型
type ValidatorType struct {
	Name     string
	validate Validate //校验器
}

func (this *ValidatorType) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	t := validators[text]
	if t == nil {
		return fmt.Errorf("value is wrong! [ %s ]", text)
	}
	this.Name = t.Name
	this.validate = t.validate
	return nil
}

func (this *ValidatorType) MarshalYAML() (interface{}, error) {
	return this.Name, nil
}

//类型
type Type struct {
	Name      string        //类型名称
	Label     string        //类型描述
	Length    int           //长度限制
	Option    string        //校验参数
	Validator ValidatorType //校验器
}

func (this *Type) validate(v string) bool {
	//长度校验
	if this.Length > 0 {
		if len(v) > this.Length {
			return false
		}
	}
	//校验器校验
	if this.Validator.Name != "" {
		return this.Validator.validate(v, this.Name)
	}
	return true
}

//参数策略类型
type PolicyType byte

const (
	Must   PolicyType = 1
	Option PolicyType = 2
)

func (this *PolicyType) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	if text == "Must" {
		*this = Must
	} else if text == "Option" {
		*this = Option
	} else {
		*this = 0
		return fmt.Errorf("value is wrong! [ %s ]", text)
	}
	return nil
}

func (this PolicyType) MarshalYAML() (interface{}, error) {
	if this == Must {
		return "Must", nil
	} else if this == Option {
		return "Option", nil
	}
	return nil, fmt.Errorf("not surport %v", this)
}

//数据格式类型
type StyleType byte

const (
	Json    StyleType = 11
	Xml     StyleType = 12
	UrlForm StyleType = 13
)

func (this *StyleType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var text string
	unmarshal(&text)
	if text == "json" {
		*this = Json
	} else if text == "xml" {
		*this = Xml
	} else if text == "url-form" {
		*this = UrlForm
	} else {
		*this = 0
		return fmt.Errorf("value is wrong! [ %s ]", text)
	}
	return nil
}

func (this StyleType) MarshalYAML() (interface{}, error) {
	if this == Json {
		return "json", nil
	} else if this == Xml {
		return "xml", nil
	} else if this == UrlForm {
		return "url-form", nil
	}
	return nil, fmt.Errorf("not surport %v", this)
}

//http 访问方法类型
type HttpMethodType byte

const (
	Get  HttpMethodType = 20
	Post HttpMethodType = 21
	Put  HttpMethodType = 22
	Del  HttpMethodType = 23
	Head HttpMethodType = 24
)

func (this *HttpMethodType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var text string
	unmarshal(&text)
	if text == "GET" {
		*this = Get
	} else if text == "POST" {
		*this = Post
	} else if text == "PUT" {
		*this = Put
	} else if text == "DEL" {
		*this = Del
	} else if text == "HEAD" {
		*this = Head
	} else {
		*this = 0
		return fmt.Errorf("value is wrong! [ %s ]", text)
	}
	return nil
}

func (this HttpMethodType) MarshalYAML() (interface{}, error) {
	if this == Get {
		return "GET", nil
	} else if this == Post {
		return "POST", nil
	} else if this == Put {
		return "PUT", nil
	} else if this == Del {
		return "DEL", nil
	} else if this == Head {
		return "HEAD", nil
	}
	return nil, fmt.Errorf("not surport %v", this)
}

//参数
type Paramter struct {
	Name          string     `yaml:"paramter"` //参数名称
	TypeName      string     `yaml:"type"`     //参数类型
	Length        int        `yaml:"length"`   //长度限制
	Policy        PolicyType `yaml:"policy"`
	ValidatorName string     `yaml:"validator"`
	Expr          string     //表达式
}

func (this *Paramter) Validate(v string) bool {
	t := types[this.TypeName]
	if t == nil {
		return false
	}

	//类型校验
	if !t.validate(v) {
		return false
	}

	//长度校验
	if this.Length > 0 {
		if len(v) > this.Length {
			return false
		}
	}
	//额外校验器校验
	if this.ValidatorName != "" {
		vt := validators[this.ValidatorName]
		if vt != nil {
			return vt.validate(v, this.Expr)
		}
		return false
	}

	return true
}
func (this *Paramter) GetType() *Type {
	return types[this.TypeName]
}
func (this *Paramter) GetValidator() *ValidatorType {
	if this.ValidatorName != "" {
		return validators[this.ValidatorName]
	}
	return nil
}

//匹配条件
type MatchItem struct {
	Name  string `yaml:"paramter"` //参数名
	Value string //取值
}

type Request struct {
	Style  StyleType
	Error  string
	Fields []Paramter `yaml:"items"`
}

type Response struct {
	Default  string //默认返回值
	Triggers []ResponseTrigger
}

//返回值触发设置
type ResponseTrigger struct {
	Data  string      //返回数据id
	Match []MatchItem //匹配条件
}

//api 定义
type Api struct {
	Name        string
	Url         string
	Description string
	Methods     []HttpMethodType
	RequestSet  Request  `yaml:"request"`
	ResponseSet Response `yaml:"response"`
}

func (this *Api) LoadFromYaml(filename string) {
	loadfromYamlfile(filename, this)
}

//返回值定义
type ResponseData struct {
	Code        string             //编号
	Description string             //描述
	Style       StyleType          //格式
	Data        string             //数据类容
	t           *template.Template //模板
}

func (this *ResponseData) LoadFromYaml(filename string) {
	loadfromYamlfile(filename, this)
}
func (this *ResponseData) Render(writer io.Writer, p interface{}) {
	if this.t == nil {
		this.t = template.New(this.Code)
		this.t.Parse(this.Data)
	}
	this.t.Execute(writer, p)
}

//系统配置
type Config struct {
	Port    int
	Version string
	Types   []Type
}

func (this *Config) LoadFromYaml(filename string) {
	loadfromYamlfile(filename, this)
	//注册类型
	for _, t := range this.Types {
		types[t.Name] = &t
		if t.Validator.Name == "regex" {
			LoadPrepareExpex(t.Name, t.Option)
		}
	}
}

func loadfromYamlfile(filename string, target interface{}) {
	buffer, err := ioutil.ReadFile(filename)
	err = yaml.Unmarshal(buffer, target)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
