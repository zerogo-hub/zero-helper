package file_test

import (
	"testing"

	zerofile "github.com/zerogo-hub/zero-helper/file"
)

func TestFile(t *testing.T) {
	filePath := "./file_test.go"
	dirPath := "."

	if !zerofile.IsFile(filePath) {
		t.Error("IsFile error")
	}

	if zerofile.IsFile(dirPath) {
		t.Error("IsFile error, dirPath")
	}

	if !zerofile.IsDir(dirPath) {
		t.Error("IsDir error")
	}

	if zerofile.IsExist("/hi/ROS9cRYp") {
		t.Error("IsExist error")
	}

	if !zerofile.IsExist(dirPath) {
		t.Error("IsExist error, dirPath")
	}
}

func TestFileName(t *testing.T) {
	filePath := "./file_test.go"
	name := zerofile.Name(filePath)
	if name != "file_test.go" {
		t.Errorf("Name error, name: %s", name)
	}
}

func TestFileNameRandom(t *testing.T) {
	name := zerofile.NameRand("score.txt", 8)
	if name == "" {
		t.Error("name is empty")
	} else {
		t.Logf("name: %s", name)
	}
}

func TestBaseName(t *testing.T) {
	filePath := "./file_test.go"
	name := zerofile.BaseName(filePath)
	if name != "file_test" {
		t.Errorf("BaseName error, name: %s", name)
	}
}

func TestExtensionName(t *testing.T) {
	filePath := "./file_test.go"
	name := zerofile.ExtensionName(filePath)
	if name != ".go" {
		t.Errorf("ExtensionName error, name: %s", name)
	}
}
