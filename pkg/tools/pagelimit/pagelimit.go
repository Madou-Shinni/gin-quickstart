package pagelimit

// 限制每页的最大值100和页码的最小值1，转化为 offset和limit
// param: num页码 size每页显示数量
// return: offset偏移量l imit每页显示数量
func OffsetLimit(num int64, size int64) (offset int, limit int) {
	// 这里使用了map来实现三元表达式
	pageNum := map[bool]int{true: 1, false: int(num)}[num < 1]
	pageSize := map[bool]int{true: 10, false: int(size)}[size > 100 || size < 1]
	// 转 offset limit
	offset = (pageNum - 1) * pageSize
	limit = pageSize
	return
}
