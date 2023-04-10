package str

import (
	"fmt"
	"testing"
)

func TestCamelToSnake(t *testing.T) {
	str := "camelCaseString"
	fmt.Println(CamelToSnake(str))
}
