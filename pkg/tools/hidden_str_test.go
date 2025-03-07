package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHidePhoneNumber(t *testing.T) {
	phone := "12345678901"
	result := HidePhoneNumber(phone)
	assert.Equal(t, "123****8901", result)
}

func TestHideBankCard(t *testing.T) {
	str := "6214764430576485192"
	result := HideBankCard(str)
	t.Log(result)
}

func TestHideAddr(t *testing.T) {
	str := "乌坦城萧家1号"
	result := HideAddr(str)
	t.Log(result)
}
