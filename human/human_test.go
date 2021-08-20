package human_test

import (
	"testing"

	zerohuman "github.com/zerogo-hub/zero-helper/human"
)

// 以下测试用的身份证号码均为随机产生，如有雷同，实属意外

func TestIDCheck(t *testing.T) {
	codes := []string{"110101199003070011", "11010119900307993X", "120101199003073900"}
	for _, code := range codes {
		if !zerohuman.IDCheck(code) {
			t.Error("IDCheck error")
		}
	}
}

func TestIDInfo(t *testing.T) {
	code := "440106199910017896"
	idCard, err := zerohuman.IDInfo(code)
	if err != nil {
		t.Errorf("IDInfo error: %s", err.Error())
	} else if idCard.Province != "广东省" || idCard.City != "广州市" || idCard.County != "天河区" {
		t.Error("IDInfo parse failed")
	}
}

func TestIDGenerate(t *testing.T) {
	_, err := zerohuman.IDGenerate(1999, 10, 1, 1, "440106", 5)
	if err != nil {
		t.Errorf("IDGenerate error: %s", err.Error())
	}
}
