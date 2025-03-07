package tools

import (
	"fmt"
	"regexp"
	"strings"
)

// HidePhoneNumber 隐藏手机号
func HidePhoneNumber(phone string) string {
	if len(phone) < 7 { // 确保字符串长度足够
		return phone // 返回原始字符串
	}
	return fmt.Sprintf("%s****%s", phone[:3], phone[len(phone)-4:])
}

// HideBankCard 隐藏银行卡
func HideBankCard(str string) string {
	// 获取字符串的长度
	strLen := len(str)

	// 如果字符串长度小于等于4，直接返回原始字符串
	if strLen <= 4 {
		return str
	}

	// 计算需要隐藏的字符数量
	hiddenLength := strLen - 4

	// 创建替换部分：将前面的字符替换为 *
	hiddenPart := strings.Repeat("*", hiddenLength)

	// 获取最后四位字符
	lastFour := str[strLen-4:]

	// 拼接结果
	return hiddenPart + lastFour
}

// HideAddr 隐藏地址
func HideAddr(address string) string {
	// 正则表达式：匹配纯数字（可能是多位数字）
	re := regexp.MustCompile(`\d+`)

	// 如果地址包含纯数字
	if re.MatchString(address) {
		// 替换数字部分为"***"
		return re.ReplaceAllString(address, "***")
	}

	// 如果地址不包含纯数字，替换最后6个字符为"***"
	if len(address) > 6 {
		return address[:len(address)-6] + "***"
	}

	// 如果地址长度小于6，直接返回原始字符串
	return address
}
