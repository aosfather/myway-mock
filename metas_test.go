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

func TestStyleType_MarshalYAML(t *testing.T) {
	t.Log(Json.MarshalYAML())
	t.Log(Xml.MarshalYAML())
	t.Log(UrlForm.MarshalYAML())
}

func TestStyleType_UnmarshalYAML(t *testing.T) {
	var s StyleType
	yaml.Unmarshal([]byte("Json"), &s)
	t.Log(s)
	yaml.Unmarshal([]byte("Xml"), &s)
	t.Log(s)
	yaml.Unmarshal([]byte("UrlForm"), &s)
	t.Log(s)
}

func TestApi_LoadFromYaml(t *testing.T) {
	var api Api
	api.LoadFromYaml("apis/hello.yaml")
	t.Log(api.Name, api.Methods)
	t.Log(api.RequestSet.Style)
	t.Log(api)
}
