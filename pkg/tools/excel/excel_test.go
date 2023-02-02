package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"testing"
)

type Excel struct {
	ID    int    `ex:"head:序号;type:int;required;color:#0070C0"`
	Name  string `ex:"head:姓名;type:string;required;color:#0070C0"`
	Phone string `ex:"head:手机号;type:string;required;color:#0070C0"`
}

func TestParseExcelTag(t *testing.T) {
	settingList := ParseExcelTag(&Excel{})
	fmt.Printf("settingList:%v", settingList)
}

// 流式导出excel
func TestStreamWrite(t *testing.T) {
	ef := excelize.NewFile()
	sw, _ := ef.NewStreamWriter("Sheet1")
	StreamWriteHead(sw, &Excel{})
	list := []Excel{
		{ID: 1, Name: "张三", Phone: "123"},
		{ID: 2, Name: "李四", Phone: "123"},
	}
	StreamWriteBody(sw, list)
	// 流式必须在结束使用Flush()
	sw.Flush()
	ef.SaveAs("stream_info.xlsx")
}
