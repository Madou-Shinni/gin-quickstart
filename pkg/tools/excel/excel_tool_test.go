package excel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 导出excel表，复杂类型给定string，业务上自己转换
type Data struct {
	ID   uint   `excel:"id"`
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

func TestExcelTool_SaveAs(t *testing.T) {
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
	err := tool.SaveAs("test.xlsx")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestExcelTool_StreamWriteBodyWithMerge(t *testing.T) {
	tool := NewExcelTool("Sheet1")
	if tool == nil {
		t.Error("tool is nil")
		return
	}
	tool.
		WriteHead(&Data{}).
		WriteBody([]*Data{
			{ID: 1, Name: "张三", Age: 18},
			{ID: 1, Name: "张三", Age: 27},
		}).
		MergeCols("id", "姓名").
		Flush()
	err := tool.SaveAs("test.xlsx")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestExcelTool_SetDropList(t *testing.T) {
	tool := NewExcelTool("Sheet1")
	if tool == nil {
		t.Error("tool is nil")
		return
	}
	// 添加下拉选项 map key必须跟 excel表头一致
	var dropList = make(map[string][]string)
	dropList["年龄"] = []string{"18", "21"}
	err := tool.Model(&Data{}).SetDropList(dropList)
	assert.Equal(t, nil, err)
	err = tool.Flush()
	assert.Equal(t, nil, err)
	tool.SaveAs("test.xlsx")
}
