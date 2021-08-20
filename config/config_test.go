package config_test

import (
	"testing"

	zeroconfig "github.com/zerogo-hub/zero-helper/config"
)

func TestConfigJSON(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileJSON("./json_test.json")
	if err != nil {
		t.Errorf("TestConfigJSON failed, err: %s", err.Error())
	}

	framework, err := c.C("framework")
	if err != nil {
		t.Errorf("C[\"framework\"] failed, err: %s", err.Error())
		return
	}
	if framework != "gweb" {
		t.Errorf("framework: %s", framework)
	}

	j, err := c.CB("json")
	if err != nil {
		t.Errorf("C[\"json\"] failed, err: %s", err.Error())
		return
	}
	if j != true {
		t.Errorf("j: %v", j)
	}
}

func TestConfigTOML(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileTOML("./toml_test.toml")
	if err != nil {
		t.Errorf("TestConfigTOML failed, err: %s", err.Error())
	}
	if value, err := c.C("title"); err != nil || value != "test toml file" {
		if err != nil {
			t.Errorf("title err: %s", err.Error())
		}
		if value != "" {
			t.Errorf("title is invalid, title: %s", value)
		}
		return
	}

	user, err := c.Any("user")
	if err != nil {
		t.Errorf("user is invalid, err: %s", err.Error())
		return
	}
	userMap := user.(map[string]interface{})
	minAge, exist := userMap["minAge"]
	if !exist {
		t.Error("minAge does not exist")
		return
	}
	if minAge.(int64) != 18 {
		t.Error("minAge is not equal to 18")
		return
	}
}

func TestConfigYAML(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileYAML("./yaml_test.yaml")
	if err != nil {
		t.Errorf("TestConfigYAML failed, err: %s", err.Error())
	}

	title, err := c.C("title")
	if err != nil {
		t.Errorf("get title failed, err: %s", err.Error())
		return
	}
	if title != "hello title" {
		t.Errorf("title is invalid")
	}

	server, err := c.Any("server")
	if err != nil {
		t.Errorf("get server failed, err: %s", err.Error())
		return
	}
	serverMap := server.(map[interface{}]interface{})
	host, exist := serverMap["host"]
	if !exist {
		t.Error("server.host does not exist")
		return
	}
	if host.(string) != "127.0.0.1" {
		t.Errorf("server.host is invalid, host: %s", host.(string))
	}
}
