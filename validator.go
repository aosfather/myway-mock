package main

import (
	"regexp"
)

/*
 校验实现

*/
//校验器
type Validator func(value string, option string) bool
type RegexpValidator struct {
	patterns map[string]*regexp.Regexp
}

//加载正则表达式到缓存中
func (this *RegexpValidator) Load(key string, p string) {
	if this.patterns == nil {
		this.patterns = make(map[string]*regexp.Regexp)
	}
	if key != "" && p != "" {
		pattern, err := regexp.Compile(p)
		if err != nil {
			//
		} else {
			this.patterns[key] = pattern
		}

	}
}

//校验参数
func (this *RegexpValidator) Validate(value string, option string) bool {
	pattern := this.patterns[option]
	if pattern == nil {
		pattern, _ = regexp.Compile(option)
	}
	if pattern != nil {
		return pattern.Match([]byte(value))
	}

	return false

}
