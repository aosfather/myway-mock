package main

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestPolicyType_MarshalYAML(t *testing.T) {
	t.Log(Must.MarshalYAML())
	t.Log(Option.MarshalYAML())
}

func TestPolicyType_UnmarshalYAML(t *testing.T) {
	var a PolicyType
	yaml.Unmarshal([]byte("Must"), &a)
	t.Log(a)

	yaml.Unmarshal([]byte("Option"), &a)
	t.Log(a)

	e := yaml.Unmarshal([]byte("yes"), &a)
	if e != nil {
		t.Errorf("ummarshal error %s", e.Error())
	}
	t.Log(a)

}

func TestApi_LoadFromYaml(t *testing.T) {
	var api Api
	api.LoadFromYaml("apis/hello.yaml")
	t.Log(api.Name, api.Methods)
	t.Log(api.RequestSet.Style)
	t.Log(api)
}

func TestValidatorType_UnmarshalYAML(t *testing.T) {
	var vt ValidatorType
	yaml.Unmarshal([]byte(""), &vt)
	t.Log(vt)
	yaml.Unmarshal([]byte("regex"), &vt)
	t.Log(vt)

}

func TestParamter_Validate(t *testing.T) {
	var api Api
	api.LoadFromYaml("apis/hello.yaml")
	p := api.RequestSet.Fields[0]
	c := new(Config)
	c.LoadFromYaml("configs.yaml")

	t.Log(p)
	t.Log(p.Validate("123"))
	t.Log(p.Validate("abc"))
	t.Log(p.Validate("123456789"))
}
