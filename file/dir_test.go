package file_test

import (
	"testing"

	zerofile "github.com/zerogo-hub/zero-helper/file"
)

func TestListDirs(t *testing.T) {
	dirs, err := zerofile.ListDirs("../")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(len(dirs))
}

func TestDirContains(t *testing.T) {
	if !zerofile.DirContains("..", ".git") {
		t.Fatal("contains .git")
	}

	if zerofile.DirContains("..", ".unknwon") {
		t.Fatal("unknon exist")
	}
}
