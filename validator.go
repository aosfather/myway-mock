package main

import (
	"fmt"
	"regexp"
)

/*
 校验实现

*/
//注册校验器
func init() {
	//初始化
	validators["regex"] = &ValidatorType{"regex", RegexpValidate}
	validators["dict"] = &ValidatorType{"dict", DictValidate}
}

// 正则校验器
var patterns map[string]*regexp.Regexp = make(map[string]*regexp.Regexp)

//加载正则表达式到缓存中
func LoadPrepareExpex(key string, p string) {
	if key != "" && p != "" {
		pattern, err := regexp.Compile(p)
		if err != nil {
			//
		} else {
			patterns[key] = pattern
		}

	}
}

//字典校验
func DictValidate(value string, name string, dict string) bool {
	fmt.Println(dictionary)
	catalog := dictionary[dict]
	//轮询code看是否属于取值范围内的值
	if catalog.Code != "" {
		for _, v := range catalog.Values {
			if v.Key == value {
				return true
			}
		}
	}

	return false
}

//正则表达式校验参数
func RegexpValidate(value string, name string, option string) bool {
	pattern := patterns[name]
	if pattern == nil {
		pattern, _ = regexp.Compile(option)
	}
	if pattern != nil {
		return pattern.Match([]byte(value))
	}

	return false

}

//-----------------------------
