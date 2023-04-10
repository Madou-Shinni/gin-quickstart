package str

import (
	"math/rand"
	"time"
)

var randSeedMap = make(map[string]bool)

// 生成随机的 num 位大小写字母
// param: num生成code的位数
func GenerateCode(num int) string {
	// 定义一个字符串
	code := ""
	// 定义一个包含所有字母的字符串
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 置随机数种子
	if _, ok := randSeedMap["randSeed"]; !ok {
		rand.Seed(time.Now().UnixNano())
		randSeedMap["randSeed"] = true
	}

	// 生成 num 个随机整数
	nums := rand.Perm(num)

	// 根据随机整数的值，从 letters 中取出相应的字母
	for _, num := range nums {
		code += string(letters[num])
	}

	return code
}
