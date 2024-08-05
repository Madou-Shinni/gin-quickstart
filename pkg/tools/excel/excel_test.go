package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"testing"
	"time"
)

type Excel struct {
	ID     int    `excel:"序号"`
	Name   string `excel:"姓名"`
	Phone  string `excel:"手机号"`
	Gender string `excel:"性别"`
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

func TestStreamWriterAllRows(t *testing.T) {
	var err error
	ef := excelize.NewFile()

	defer ef.Close()
	sw, _ := ef.NewStreamWriter("Sheet1")

	// 业务处理逻辑

	// 最终数据集
	rows := make([][]interface{}, 0)

	// 边框
	borders := []excelize.Border{
		excelize.Border{Type: "bottom", Color: "#000000", Style: 1},
		excelize.Border{Type: "right", Color: "#000000", Style: 1},
	}
	// 对齐方式
	alignment := excelize.Alignment{
		Horizontal: "center",
	}
	// 设置样式
	styleTitleID, err := ef.NewStyle(&excelize.Style{Border: borders, Alignment: &alignment, Font: &excelize.Font{Family: "宋体", Size: 14}})
	styleContentCenterID, err := ef.NewStyle(&excelize.Style{Border: borders, Alignment: &alignment, Font: &excelize.Font{Family: "宋体", Size: 10.5}})
	styleContentID, err := ef.NewStyle(&excelize.Style{Border: borders, Alignment: &excelize.Alignment{Horizontal: "left"}, Font: &excelize.Font{Family: "宋体", Size: 10.5}})

	// 表格数据
	for i := 0; i < 7; i++ {
		// 前7行
		row := make([]interface{}, 0)

		// 每一行的需要写入几列数据，不需要写入数据的使用
		//	excelize.Cell{StyleID: styleTitleID} 置空，否则会丢失样式
		if i == 0 {
			// 第一行
			row = append(row, excelize.Cell{
				StyleID: styleTitleID,
				Value:   "标题",
			}, excelize.Cell{
				StyleID: styleTitleID,
			}, excelize.Cell{
				StyleID: styleTitleID,
			}, excelize.Cell{
				StyleID: styleTitleID,
			})
		} else if i == 6 {
			// 第七行
			row = append(row, excelize.Cell{
				StyleID: styleContentID,
				Value:   "内容-标题",
			}, excelize.Cell{
				StyleID: styleContentCenterID,
				Value:   "内容-标题",
			}, excelize.Cell{
				StyleID: styleContentCenterID,
				Value:   "内容-标题",
			}, excelize.Cell{
				StyleID: styleContentCenterID,
				Value:   "内容-标题",
			})
		} else {
			// 第2-6行
			row = append(row, excelize.Cell{
				StyleID: styleContentID,
				Value:   "内容-标题",
			}, excelize.Cell{
				StyleID: styleContentCenterID,
			}, excelize.Cell{
				StyleID: styleContentCenterID,
			}, excelize.Cell{
				StyleID: styleContentCenterID,
			})
		}

		// 把每一行数据加入最终数据
		rows = append(rows, row)
	}

	// 第八行开始的数据
	for i := 0; i < 100; i++ {
		row := make([]interface{}, 0)
		// 四列
		row = append(row, excelize.Cell{
			StyleID: styleContentID,
			Value:   "内容",
		}, excelize.Cell{
			StyleID: styleContentID,
			Value:   "内容",
		}, excelize.Cell{
			StyleID: styleContentID,
			Value:   "内容",
		}, excelize.Cell{
			StyleID: styleContentID,
			Value:   "内容",
		})

		rows = append(rows, row)
	}

	mergerCell := []string{"A1:D1", "B2:D2", "B3:D3", "B4:D4", "B5:D5", "B6:D6"}
	err = sw.SetColWidth(1, 4, 20)
	StreamWriterAllRows(sw, rows, mergerCell...)

	sw.Flush()
	// 保存 Excel 文件
	err = ef.SaveAs("output2.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Excel file created successfully.")
}

type TestStruct struct {
	Name      string     `excel:"姓名"`
	DateTime  *time.Time `excel:"时间"`
	StartTime time.Time  `excel:"开始时间"`
	Title     string     `excel:"标题"`
	Age       int64      `excel:"年龄"`
}

func TestParseExcelToSlice(t *testing.T) {
	dir, _ := os.Open("测试.xlsx")
	v2, err := ParseExcelToSlice[TestStruct](dir)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	structs := v2
	for _, testStruct := range structs {
		fmt.Println(testStruct)
	}
}
