package config_test

import (
	"testing"

	zeroconfig "github.com/zerogo-hub/zero-helper/config"
)

func TestLoadJSON_NilBytes(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.LoadJSON(nil)
	if err == nil || err.Error() != "bytes cannot be empty" {
		t.Errorf("Expected error 'bytes cannot be empty', but got %v", err)
	}
}

func TestLoadJSON_ValidJSON(t *testing.T) {
	jsonData := []byte(`{"framework": "Go", "version": "3"}`)
	c := zeroconfig.NewConfig()
	err := c.LoadJSON(jsonData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	framework, _ := c.C("framework")
	version, _ := c.CI("version")

	if framework != "Go" || version != 3 {
		t.Errorf("Expected framework='Go' and version=1.15, but got framework='%s' and version=%d", framework, version)
	}
}

func TestLoadJSON_InvalidJSON(t *testing.T) {
	jsonData := []byte(`{"framework": "Go", "version": "1.15",}`)
	c := zeroconfig.NewConfig()
	err := c.LoadJSON(jsonData)
	if err == nil {
		t.Error("Excepted got err")
	}
}

func TestLoadJSON_EmptyJSON(t *testing.T) {
	jsonData := []byte(`{}`)
	c := zeroconfig.NewConfig()
	err := c.LoadJSON(jsonData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = c.C("framework")
	if err == nil || err.Error() != "configuration does not exist" {
		t.Errorf("Expected error 'configuration does not exist', but got %v", err)
	}
}

func TestLoadTOML_NilBytes(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.LoadTOML(nil)
	if err == nil || err.Error() != "bytes cannot be empty" {
		t.Errorf("Expected error for nil bytes, but got: %v", err)
	}
}

func TestLoadTOML_DecodeError(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.LoadTOML([]byte("invalid toml data"))
	if err == nil {
		t.Error("Expected error for invalid TOML data, but got nil")
	}
}

func TestLoadTOML_Successful(t *testing.T) {
	c := zeroconfig.NewConfig()
	tomlData := []byte("key = \"value\"")
	err := c.LoadTOML(tomlData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestLoadTOML_ValidData(t *testing.T) {
	c := zeroconfig.NewConfig()
	tomlData := []byte("key = \"value\"")
	err := c.LoadTOML(tomlData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	value, _ := c.C("key")
	if value != "value" {
		t.Errorf("Expected value for key 'key' to be 'value', but got: %s", value)
	}
}

func TestLoadTOML_ValidDataAny(t *testing.T) {
	c := zeroconfig.NewConfig()
	tomlData := []byte("key = \"value\"")
	err := c.LoadTOML(tomlData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	value, _ := c.Any("key")
	if value != "value" {
		t.Errorf("Expected value for key 'key' to be 'value', but got: %s", value)
	}
}

func TestLoadYAML_NilBytes(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.LoadYAML(nil)
	if err == nil || err.Error() != "bytes cannot be empty" {
		t.Errorf("Expected error for nil bytes, but got: %v", err)
	}
}

func TestLoadYAML_UnmarshalError(t *testing.T) {
	c := zeroconfig.NewConfig()
	bytes := []byte("invalid yaml")
	err := c.LoadYAML(bytes)
	if err == nil {
		t.Error("Expected error for unmarshal failure, but got nil")
	}
}

func TestLoadYAML_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	bytes := []byte("key: value")
	err := c.LoadYAML(bytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFileJSON_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileJSON("json_test.json")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
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

func TestFileJSON_LoadFileError(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileJSON("nonexistent.json")
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestFileJSON_EmptyPath(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileJSON("")
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestConfig_FileTOML_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileTOML("toml_test.toml")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
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

func TestConfig_FileTOML_EmptyPath(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileTOML("")
	if err == nil {
		t.Error("Expected an error for empty path, but got nil")
	}
}

func TestConfig_FileTOML_FileReadError(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileTOML("nonexistent.toml")
	if err == nil {
		t.Error("Expected an error for non-existent file, but got nil")
	}
}

func TestConfig_FileTOML_NilBytes(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.LoadTOML(nil)
	if err == nil {
		t.Error("Expected an error for nil bytes, but got nil")
	}
}

func TestFileYAML_LoadsYAMLFileSuccessfully(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileYAML("yaml_test.yaml")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
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

func TestConfig_FileYAML_FileReadError(t *testing.T) {
	c := zeroconfig.NewConfig()
	err := c.FileYAML("nonexistent.toml")
	if err == nil {
		t.Error("Expected an error for non-existent file, but got nil")
	}
}

func TestCB_TrueValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": "true"}`))
	result, err := c.CB("key")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if result != true {
		t.Errorf("Expected true, but got %v", result)
	}
}

func TestCB_FalseValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": "false"}`))
	result, err := c.CB("key")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if result != false {
		t.Errorf("Expected false, but got %v", result)
	}
}

func TestCB_EmptyValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": ""}`))
	_, err := c.CB("key")
	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestCB_InvalidValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": "invalid"}`))
	_, err := c.CB("key")
	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestConfig_CI_ConfigNil(t *testing.T) {
	c := zeroconfig.NewConfig()
	_, err := c.CI("key")
	if err == nil {
		t.Error("Expected an error when config is nil, but got nil")
	}
}

func TestCI32_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": "42"}`))
	result, err := c.CI32("key")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if result != 42 {
		t.Errorf("Expected result to be 42, but got: %v", result)
	}
}

func TestCI32_KeyNotFound(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"anotherKey": "42"}`))
	_, err := c.CI32("key")
	if err == nil {
		t.Error("Expected an error for key not found, but got nil")
	}
}

func TestCI32_ParseError(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": "not_an_integer"}`))
	_, err := c.CI32("key")
	if err == nil {
		t.Error("Expected an error for parse failure, but got nil")
	}
}

func TestConfig_CI64_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": "123"}`))
	result, err := c.CI64("key")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if result != 123 {
		t.Errorf("Expected result to be 123, but got %v", result)
	}
}

func TestConfig_CI64_KeyNotFound(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"anotherKey": "123"}`))
	_, err := c.CI64("key")
	if err == nil {
		t.Error("Expected an error for key not found, but got nil")
	}
}

func TestConfig_CI64_ParseError(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key": "not_an_integer"}`))
	_, err := c.CI64("key")
	if err == nil {
		t.Error("Expected an error for parse failure, but got nil")
	}
}

func TestCF32_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "3.14"}`))
	f32, err := c.CF32("key1")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expected := float32(3.14)
	if f32 != expected {
		t.Errorf("Expected %v, but got: %v", expected, f32)
	}
}

func TestCF32_KeyNotFound(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "3.14"}`))
	_, err := c.CF32("key2")
	if err == nil {
		t.Error("Expected an error for key not found, but got nil")
	}
}

func TestCF32_InvalidValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "invalid"}`))
	_, err := c.CF32("key1")
	if err == nil {
		t.Error("Expected an error for invalid value, but got nil")
	}
}

func TestConfig_D_ExistingKey(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "value1", "key2": "value2"}`))
	expected := "value1"
	result := c.D("key1", "default")
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestConfig_D_NonExistingKey(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "value1", "key2": "value2"}`))
	expected := "default"
	result := c.D("key3", "default")
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestConfig_D_EmptyConfig(t *testing.T) {
	c := zeroconfig.NewConfig()
	expected := "default"
	result := c.D("key1", "default")
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestConfig_D_EmptyDefaultValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "value1", "key2": "value2"}`))
	expected := "value1"
	result := c.D("key1", "")
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestConfig_D_EmptyKey(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "value1", "key2": "value2"}`))
	expected := "default"
	result := c.D("", "default")
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestConfig_DB_KeyExists(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "true"}`))
	result := c.DB("key1", false)
	if result != true {
		t.Errorf("Expected true, but got %v", result)
	}
}

func TestConfig_DB_KeyNotExists(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "true"}`))
	result := c.DB("key2", false)
	if result != false {
		t.Errorf("Expected false, but got %v", result)
	}
}

func TestConfig_DB_KeyExistsWithDefaultValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "true"}`))
	result := c.DB("key1", true)
	if result != true {
		t.Errorf("Expected true, but got %v", result)
	}
}

func TestConfig_DB_KeyNotExistsWithDefaultValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "true"}`))
	result := c.DB("key2", true)
	if result != true {
		t.Errorf("Expected true, but got %v", result)
	}
}

func TestConfig_DB_EmptyConfig(t *testing.T) {
	c := zeroconfig.NewConfig()
	result := c.DB("key1", true)
	if result != true {
		t.Errorf("Expected true, but got %v", result)
	}
}

func TestConfig_DI_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "10"}`))
	result := c.DI("key1", 5)
	if result != 10 {
		t.Errorf("Expected DI to return 10, but got %d", result)
	}
}

func TestConfig_DI_DefaultValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "abc"}`))
	result := c.DI("key2", 5)
	if result != 5 {
		t.Errorf("Expected DI to return default value 5, but got %d", result)
	}
}

func TestConfig_DI_EmptyConfig(t *testing.T) {
	c := zeroconfig.NewConfig()
	result := c.DI("key1", 5)
	if result != 5 {
		t.Errorf("Expected DI to return default value 5, but got %d", result)
	}
}

func TestConfig_DI_EmptyDefaultValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "10"}`))
	result := c.DI("key2", 0)
	if result != 0 {
		t.Errorf("Expected DI to return default value 0, but got %d", result)
	}
}

func TestConfig_DI_InvalidValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "abc"}`))
	result := c.DI("key1", 5)
	if result != 5 {
		t.Errorf("Expected DI to return default value 5, but got %d", result)
	}
}

func TestDI32_ConfigValueExists(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "42"}`))
	result := c.DI32("key1", 0)
	if result != 42 {
		t.Errorf("Expected 42, but got %d", result)
	}
}

func TestDI32_ConfigValueDoesNotExist(t *testing.T) {
	c := zeroconfig.NewConfig()
	result := c.DI32("nonexistent", 42)
	if result != 42 {
		t.Errorf("Expected 42, but got %d", result)
	}
}

func TestDI32_ConfigValueNotParsable(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "not_an_integer"}`))
	result := c.DI32("key1", 42)
	if result != 42 {
		t.Errorf("Expected 42, but got %d", result)
	}
}

func TestConfig_DI64_Success(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "123"}`))
	result := c.DI64("key1", 0)
	if result != 123 {
		t.Errorf("Expected result to be 123, but got %v", result)
	}
}

func TestConfig_DI64_DefaultValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "abc"}`))
	result := c.DI64("key2", 456)
	if result != 456 {
		t.Errorf("Expected result to be 456, but got %v", result)
	}
}

func TestConfig_DI64_EmptyConfig(t *testing.T) {
	c := zeroconfig.NewConfig()
	result := c.DI64("key1", 789)
	if result != 789 {
		t.Errorf("Expected result to be 789, but got %v", result)
	}
}

func TestConfig_DI64_EmptyKey(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "123"}`))
	result := c.DI64("", 0)
	if result != 0 {
		t.Errorf("Expected result to be 0, but got %v", result)
	}
}

func TestConfig_DI64_InvalidValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "abc"}`))
	result := c.DI64("key1", 0)
	if result != 0 {
		t.Errorf("Expected result to be 0, but got %v", result)
	}
}

func TestDF32_ConfigValueExists(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "3.14"}`))
	expected := float32(3.14)
	if result := c.DF32("key1", 0); result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestDF32_ConfigValueDoesNotExist(t *testing.T) {
	c := zeroconfig.NewConfig()
	expected := float32(0)
	if result := c.DF32("nonexistent", 0); result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestDF32_ConfigValueConversionError(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "not_a_number"}`))
	expected := float32(0)
	if result := c.DF32("key1", 0); result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestDF32_ConfigDefaultValue(t *testing.T) {
	c := zeroconfig.NewConfig()
	expected := float32(5.5)
	if result := c.DF32("nonexistent", 5.5); result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestDF32_ConfigDefaultValueZero(t *testing.T) {
	c := zeroconfig.NewConfig()
	expected := float32(0)
	if result := c.DF32("nonexistent", 0); result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestDF64_ConfigValueExists(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "3.14"}`))
	result := c.DF64("key1", 0.0)
	expected := 3.14
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestDF64_ConfigValueDoesNotExist(t *testing.T) {
	c := zeroconfig.NewConfig()
	result := c.DF64("nonexistent", 3.14)
	expected := 3.14
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestDF64_ConfigValueConversionError(t *testing.T) {
	c := zeroconfig.NewConfig()
	c.LoadJSON([]byte(`{"key1": "not_a_number"}`))
	result := c.DF64("key1", 3.14)
	expected := 3.14
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}
