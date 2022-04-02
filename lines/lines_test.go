package lines

import (
	"bytes"
	"testing"
)

type MyInt uint32

func TestEachNumeric(t *testing.T) {
	data := bytes.NewBuffer([]byte("99"))
	i := MyInt(0)
	err := Each(data, func(x MyInt) { i = x })
	if err != nil {
		t.Error("unexpected error", err)
	}
	if i != 99 {
		t.Error("unexpected value", i)
	}

	data = bytes.NewBuffer([]byte("-99"))
	err = Each(data, func(x MyInt) { i = x })
	if err == nil {
		t.Error("should have had error")
	}
}
