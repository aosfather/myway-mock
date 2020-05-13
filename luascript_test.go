package main

import (
	luautils "github.com/aosfather/bingo_utils/lua"
	lua "github.com/yuin/gopher-lua"
	"testing"
	"time"
)

func TestLuaScript_Load(t *testing.T) {
	script := LuaScript{}
	script.pool = luautils.NewLuaPool(10, nil)
	s := `
a=2+1
print(a)
person = {
	  name = "Michel",
	  age  = "31", -- weakly input
	  work_place = "San Jose"
    }
return person
`
	script.Load("test", s)
	for i := 0; i < 100; i++ {
		go func() {
			l, v, err := script.Call()
			if err != nil {
				t.Log(err.Error())
				return
			}
			if v.Type() == lua.LTTable {
				t.Log(l.GetTable(v, lua.LString("name")).String())
				t.Log(l.GetField(v, "age").String())
			}
			t.Log(v.String())
		}()
		time.Sleep(10 * time.Microsecond)
	}
	time.Sleep(10 * time.Second)

}
