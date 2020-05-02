package main

import (
	"regexp"
)

/*
 校验实现

*/

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

//正则表达式校验参数
func RegexpValidate(value string, option string) bool {
	pattern := patterns[option]
	if pattern == nil {
		pattern, _ = regexp.Compile(option)
	}
	if pattern != nil {
		return pattern.Match([]byte(value))
	}

	return false

}

//-----------------------------
