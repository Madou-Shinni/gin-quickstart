package snowflake

import (
	"fmt"
	"testing"
)

func TestGenerateID(t *testing.T) {
	// 1.初始化节点
	SnowflakeInit(1)
	// 2.生成id
	id := GenerateID()
	fmt.Println(id) // 1603230116786212864 1603230954631991296
}
