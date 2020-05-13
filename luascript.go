package main

import (
	"fmt"
	luautils "github.com/aosfather/bingo_utils/lua"
	lua "github.com/yuin/gopher-lua"
	"log"
)

type LuaScript struct {
	pool     *luautils.LuaPool
	function *lua.FunctionProto
}

//加载脚本文件
func (this *LuaScript) Loadfile(filename string) {
	var err error
	this.function, err = luautils.CompileByfile(filename)
	if err != nil {
		log.Println(err.Error())
	}
}

//加载脚本
func (this *LuaScript) Load(name, content string) {
	var err error
	this.function, err = luautils.CompileByString(name, content)
	if err != nil {
		log.Println(err.Error())
	}
}

func (this *LuaScript) Call() (*lua.LState, lua.LValue, error) {
	var l *lua.LState
	if this.pool != nil {
		l = this.pool.Get()
		defer func() {
			this.pool.Put(l)
		}()
	}
	if l == nil {
		return nil, nil, fmt.Errorf("no lua vm!")
	}
	lfunc := l.NewFunctionFromProto(this.function)
	l.Push(lfunc)
	err := l.PCall(0, 1, l.NewFunction(this.errHandle))
	return l, l.Get(-1), err
}

//错误捕获器
func (this *LuaScript) errHandle(l *lua.LState) int {
	log.Println(l.Get(-1).String())
	return 1
}
