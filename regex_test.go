package main

import "testing"

func TestRegexpValidate(t *testing.T) {
	//整型
	t.Log(RegexpValidate("12", "", "[-]?[0-9]\\d+"))
	t.Log(RegexpValidate("-12", "", "[-]?[0-9]\\d+"))
	//日期
	t.Log(RegexpValidate("2019-12-20", "", "[1-9]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])"))
	//时间
	t.Log(RegexpValidate("10:23:00", "", "(20|21|22|23|[0-1]\\d):[0-5]\\d:[0-5]\\d"))
	t.Log(RegexpValidate("10:63:00", "", "(20|21|22|23|[0-1]\\d):[0-5]\\d:[0-5]\\d"))
	t.Log(RegexpValidate("24:43:00", "", "(20|21|22|23|[0-1]\\d):[0-5]\\d:[0-5]\\d"))
	//日期时间
	t.Log(RegexpValidate("2012-02-30 23:43:00", "", "[1-9]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\\s+(20|21|22|23|[0-1]\\d):[0-5]\\d:[0-5]\\d"))
}
