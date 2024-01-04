package excel

import "testing"

// 导出excel表，复杂类型给定string，业务上自己转换
type Data struct {
	Name string `excel:"姓名"`
	Age  int    `excel:"年龄"`
}

func TestNewExcelTool(t *testing.T) {
	tool := NewExcelTool("Sheet1")
	if tool == nil {
		t.Error("tool is nil")
		return
	}
	tool.
		WriteHead(&Data{}).
		WriteBody([]*Data{
			{Name: "张三", Age: 18},
		}).
		Flush()
	bytesBuf, err := tool.WriteToBuffer()
	if err != nil {
		t.Error(err)
		return
	}
	// bytesBuf.Bytes() 作为流传递给前端
	t.Log("success", bytesBuf.Bytes())
}
