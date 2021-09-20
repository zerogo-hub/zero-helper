package words_test

import (
	"testing"

	zerowords "github.com/zerogo-hub/zero-helper/words"
)

func TestWordsCount(t *testing.T) {
	wc1 := zerowords.Count("abcd")
	if wc1.Total() != 1 || wc1.Pure() != 1 {
		t.Fatal("wc1 test failed")
	}

	wc2 := zerowords.Count("ab,cd")
	if wc2.Total() != 3 || wc2.Pure() != 2 {
		t.Fatal("wc2 test failed")
	}

	// 空格忽略
	wc3 := zerowords.Count("ab  ,  cd")
	if wc3.Total() != 3 || wc3.Pure() != 2 {
		t.Fatal("wc3 test failed")
	}

	// 纯中文
	wc4 := zerowords.Count("你好啊")
	if wc4.Total() != 3 || wc4.Pure() != 3 {
		t.Fatal("wc4 test failed")
	}

	// 中文 + 中文逗号
	wc5 := zerowords.Count("你好，啊")
	if wc5.Total() != 4 || wc5.Pure() != 3 {
		t.Fatal("wc5 test failed")
	}

	// 中文 + 英文逗号
	wc6 := zerowords.Count("你好,啊")
	if wc6.Total() != 4 || wc6.Pure() != 3 {
		t.Fatal("wc6 test failed")
	}

	// 中文 + 空格，空格忽略
	wc7 := zerowords.Count("你好 啊")
	if wc7.Total() != 3 || wc7.Pure() != 3 {
		t.Fatal("wc7 test failed")
	}

	// 大杂烩
	wc8 := zerowords.Count("你好啊， Bill，hi")
	if wc8.Total() != 7 || wc8.Pure() != 5 {
		t.Fatal("wc8 test failed")
	}
}
