package excel

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
	"testing"
)

// 导出excel表，复杂类型给定string，业务上自己转换
type Data struct {
	ID      uint   `excel:"id"`
	Name    string `excel:"姓名"`
	Age     int    `excel:"年龄"`
	IsAdult bool   `excel:"是否成年"`
	IsOk    *bool  `excel:"是否"`
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
	// 创建一个文件流
	ex := excelize.NewFile()
	streamWriter, _ := ex.NewStreamWriter("Sheet1")

	// 创建一个带有下拉列表的列
	// B 列添加下拉列表，显示文本，值为数值
	//dropDownOptions := map[string]string{
	//	"Apple":  "1",
	//	"Orange": "2",
	//	"Banana": "3",
	//}

	//validation := &excelize.DataValidation{
	//	//Type:             "list",
	//	Sqref:            "A2:A65535",
	//	ShowInputMessage: true,
	//	ShowErrorMessage: true,
	//}
	//validation.SetDropList([]string{"1", "2", "3"})
	//err := ex.AddDataValidation("Sheet1", validation)

	validation := excelize.NewDataValidation(true)
	validation.SetSqref("A2:A65535")
	validation.SetDropList([]string{"1", "2", "3"})
	err := ex.AddDataValidation("Sheet1", validation)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 将数据流式写入文件

	// 写入表头
	cell, _ := excelize.CoordinatesToCellName(1, 1)
	streamWriter.SetRow(cell, []interface{}{"选项"})

	// 以流式方式写入数据
	//for i := 2; i <= 100; i++ {
	//	rowData := []interface{}{"", dropDownOptions["Apple"]} // 下拉列初始为 Apple
	//	cell, _ = excelize.CoordinatesToCellName(1, i)
	//	streamWriter.SetRow(cell, rowData)
	//}

	// 保存文件到流
	if err := streamWriter.Flush(); err != nil {
		fmt.Println(err)
		return
	}

	// 将文件流保存到磁盘
	if err := ex.SaveAs("test.xlsx"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Excel 创建完成")
}

func TestExcelTool_SetDropListPro(t *testing.T) {
	// 创建一个文件流
	tool := NewExcelTool("Sheet1")
	if tool == nil {
		t.Error("tool is nil")
		return
	}

	tool.Model(&Data{})

	tool.SetDropList(map[string][]string{
		"姓名": {"张三", "李四", "王五"},
	})

	// 设置选项
	var dropList []string
	for i := 0; i < 2000; i++ {
		dropList = append(dropList, fmt.Sprintf("%d", i+1))
	}

	tool.SetDropListPro(map[string][]string{
		"年龄": dropList,
	})

	// Flush
	err := tool.Flush()
	assert.Equal(t, nil, err)

	// 将文件流保存到磁盘
	if err := tool.SaveAs("test.xlsx"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Excel 创建完成")
}

func TestExcelTool_Remark(t *testing.T) {
	// 创建一个文件流
	tool := NewExcelTool("Sheet1")
	if tool == nil {
		t.Error("tool is nil")
		return
	}

	err := tool.Model(&Data{}).Remark(`填写说明:
1.请仔细阅读本填写说明，若填写不符合规则，将导致数据导入失败
`).Flush()
	assert.Equal(t, nil, err)

	// 将文件流保存到磁盘
	if err := tool.SaveAs("test.xlsx"); err != nil {
		t.Log(err)
		return
	}
}

func TestExcelTool_FormatBool(t *testing.T) {
	// 创建一个文件流
	tool := NewExcelTool("Sheet1")
	if tool == nil {
		t.Error("tool is nil")
		return
	}

	err := tool.Model(&Data{}).WriteBody([]*Data{
		{ID: 1, Name: "张三", Age: 17},
		{ID: 2, Name: "李四", Age: 27, IsAdult: true, IsOk: new(bool)},
	}).FormatBool(map[bool]string{
		true:  "是",
		false: "否",
	}).Flush()
	assert.Equal(t, nil, err)

	// 将文件流保存到磁盘
	if err := tool.SaveAs("test.xlsx"); err != nil {
		t.Log(err)
		return
	}
}
