package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"testing"
)

type Excel struct {
	ID     int    `ex:"head:序号;type:int;required;color:#0070C0"`
	Name   string `ex:"head:姓名;type:string;required;color:#0070C0"`
	Phone  string `ex:"head:手机号;type:string;required;color:#0070C0"`
	Gender string `ex:"head:性别;type:string;required;color:#0070C0;select:男,女"`
}

func TestParseExcelTag(t *testing.T) {
	settingList := ParseExcelTag(&Excel{})
	fmt.Printf("settingList:%v", settingList)
}

// 流式导出excel
func TestStreamWrite(t *testing.T) {
	ef := excelize.NewFile()

	defer ef.Close()
	sw, _ := ef.NewStreamWriter("Sheet1")
	err2 := StreamWriteHead(sw, &Excel{})
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	list := []*Excel{
		{ID: 1, Name: "张三", Phone: "123", Gender: "男"},
		{ID: 2, Name: "李四", Phone: "123", Gender: "女"},
	}
	StreamWriteBody(sw, list)
	// 流式必须在结束使用Flush()
	sw.Flush()

	err := ef.SaveAs("stream_info.xlsx")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Test2(t *testing.T) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := f.SetCellValue("Sheet1", "A6", 42920.5); err != nil {
		fmt.Println(err)
		return
	}
	exp := "[$-380A]dddd\\,\\ dd\" de \"mmmm\" de \"yyyy;@"
	style, err := f.NewStyle(&excelize.Style{CustomNumFmt: &exp})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = f.SetCellStyle("Sheet1", "A6", "A6", style)
	err = f.SaveAs("stream_info.xlsx")
	if err != nil {
		fmt.Println(err.Error())
	}
}
