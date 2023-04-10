package str

import (
	"fmt"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	code1 := ""
	code2 := ""

	code1 = GenerateCode(20)
	code2 = GenerateCode(20)

	fmt.Println(code1)
	fmt.Println(code2)
}
