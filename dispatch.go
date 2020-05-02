package main

import (
	"strings"
)

/**
  请求分发管理
*/
type DispatchManager struct {
	domainNode  map[string]*node    //特定域名下的node
	defaultNode *node               //默认
	apiMap      map[string]*Service //api列表
}

func (this *DispatchManager) Init() {
	this.domainNode = make(map[string]*node)
	this.apiMap = make(map[string]*Service)
	this.defaultNode = &node{}
}

//根据域名和url获取对应的API
func (this *DispatchManager) GetApi(domain, url string) *Service {
	node := this.domainNode[domain]
	if node == nil {
		node = this.defaultNode
	}

	if node != nil {
		paramIndex := strings.Index(url, "?")
		realuri := url
		if paramIndex != -1 {
			realuri = strings.TrimSpace((url[:paramIndex]))
		}

		h, _, _ := node.getValue(realuri)
		if h != nil {
			key := h.(string)
			return this.apiMap[key]
		}

	}
	return nil

}

func (this *DispatchManager) AddApi(domain string, root string, api *Api) {
	if api == nil {
		return
	}
	serv := NewService(root, api)
	this.apiMap[api.Name] = serv
	var apiNode *node
	if domain == "" {
		apiNode = this.defaultNode
	} else {
		//处理不同的域名的映射
		if apiNode == nil {
			apiNode = this.domainNode[domain]
			if apiNode == nil {
				apiNode = &node{}
				this.domainNode[domain] = apiNode
			}
		}

	}

	if api != nil {
		apiNode.addRoute(api.Url, api.Name)
	}
}
