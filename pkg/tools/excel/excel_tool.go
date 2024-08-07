package excel

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"time"
)

type ExcelTool struct {
	file                *excelize.File
	sw                  *excelize.StreamWriter
	model               interface{} // 结构体(head)
	list                interface{} // 数据(body)
	mergeConditionIndex int         // 合并条件（列下标从0开始，例如第一列相同则输入0）
	mergeCols           []string    // 合并列
}

func NewExcelTool(sheet string) *ExcelTool {
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", sheet)
	sw, err := file.NewStreamWriter(sheet)
	if err != nil {
		return nil
	}

	return &ExcelTool{
		file: file,
		sw:   sw,
	}
}

func (e *ExcelTool) WriteHead(data interface{}) *ExcelTool {
	e.model = data
	return e
}

func (e *ExcelTool) WriteBody(list interface{}) *ExcelTool {
	e.list = list
	return e
}

func (e *ExcelTool) MergeCols(headCondition string, mergeHeads ...string) *ExcelTool {
	e.mergeConditionIndex = e.headToMergeConditionIndex(headCondition)
	e.mergeCols = e.headToCols(mergeHeads...)
	return e
}

func (e *ExcelTool) headToCols(heads ...string) []string {
	var cols []string
	var headMap = make(map[string]struct{})
	for _, head := range heads {
		headMap[head] = struct{}{}
	}
	tp := reflect.ValueOf(e.model).Type().Elem() // 获得结构体的反射Type
	numField := tp.NumField()                    // 获取结构体的字段数量
	for i := 0; i < numField; i++ {
		field := tp.Field(i)          // 获取字段
		tag := field.Tag.Get("excel") // 获取tag中的ex值
		if _, ok := headMap[tag]; ok {
			cols = append(cols, indexToColumnName(i))
		}
	}
	return cols
}

// Convert an integer to an Excel column name
func indexToColumnName(index int) string {
	columnName := ""
	for index >= 0 {
		columnName = string(rune('A'+(index%26))) + columnName
		index = index/26 - 1
	}
	return columnName
}

func (e *ExcelTool) headToMergeConditionIndex(head string) int {
	tp := reflect.ValueOf(e.model).Type().Elem() // 获得结构体的反射Type
	numField := tp.NumField()                    // 获取结构体的字段数量
	for i := 0; i < numField; i++ {
		field := tp.Field(i)          // 获取字段
		tag := field.Tag.Get("excel") // 获取tag中的ex值
		if tag == head {
			return i
		}
	}
	return 0
}

func (e *ExcelTool) Flush() error {
	err := StreamWriteHead(e.sw, e.model)
	if err != nil {
		return err
	}
	err = e.StreamWriteBodyWithMerge(e.sw, e.list, e.mergeConditionIndex, e.mergeCols)
	if err != nil {
		return err
	}
	return e.sw.Flush()
}

func (e *ExcelTool) StreamWriteBodyWithMerge(sw *excelize.StreamWriter, d interface{}, mergeConditionIndex int, mergeCols []string) error {
	// 判断d的数据类型
	switch reflect.TypeOf(d).Kind() {
	case reflect.Slice, reflect.Array:
		// 是切片或者数组
		values := reflect.ValueOf(d)
		// 创建一个二维数组的数据集，用来存放最终数据集
		data := make([][]interface{}, values.Len())
		for i := 0; i < values.Len(); i++ {
			// 取出切片中的每个结构体，利用反射获取值
			record := values.Index(i).Elem()
			if record.Kind() == reflect.Struct {
				// 创建一个切片来表示一行数据
				row := make([]interface{}, record.NumField())
				for j := 0; j < record.NumField(); j++ {
					// 遍历结构体中的字段，取出字段值
					field := record.Field(j)
					row[j] = excelize.Cell{
						StyleID: 0,
						Formula: "",
						Value:   field.Interface(),
					}
				}
				// 将每一行数据保存到二维数组中
				data[i] = row
			}
		}

		row := 2
		var mergeStarted bool
		var mergeStartedRow, mergeEndedRow int
		for i := range data {
			if i < len(data)-1 {
				if !mergeStarted && data[i][mergeConditionIndex] == data[i+1][mergeConditionIndex] {
					// 合并单元格 开始
					mergeStarted = true
					mergeStartedRow = row
				}
				if mergeStarted && (data[i][mergeConditionIndex] != data[i+1][mergeConditionIndex]) {
					// 后一行合并条件不同 合并结束
					mergeEndedRow = row
					mergeStarted = false
					for _, col := range mergeCols {
						sw.MergeCell(fmt.Sprintf("%s%d", col, mergeStartedRow), fmt.Sprintf("%s%d", col, mergeEndedRow))
					}
				}
				if mergeStarted && (data[i][mergeConditionIndex] == data[i+1][mergeConditionIndex]) && i+1 == len(data)-1 {
					// 后是最后一行并且和前数据相同 合并结束
					mergeEndedRow = row + 1
					mergeStarted = false
					for _, col := range mergeCols {
						sw.MergeCell(fmt.Sprintf("%s%d", col, mergeStartedRow), fmt.Sprintf("%s%d", col, mergeEndedRow))
					}
				}
			}
			// 逐行插入数据 将数据写入excel
			// 数据都是从列号1开始；行号从2开始，因为第一行为标题行
			axis, err := excelize.CoordinatesToCellName(1, row)
			if err != nil {
				return err
			}
			if err := sw.SetRow(axis, data[i], excelize.RowOpts{Height: 16}); err != nil {
				return err
			}
			row++
		}
	default:
		// 不支持改数据类型
		return errors.New("resolution of this data type is not supported")
	}

	return nil
}

func (e *ExcelTool) WriteToBuffer() (*bytes.Buffer, error) {
	return e.file.WriteToBuffer()
}

func (e *ExcelTool) SaveAs(filename string) error {
	return e.file.SaveAs(filename)
}

// ReturnEmptyIfPointerIsNil 空指针则返回""
func ReturnEmptyIfPointerIsNil(ptr interface{}) string {
	value := reflect.ValueOf(ptr)
	if value.Kind() != reflect.Ptr {
		return fmt.Sprintf("%v", ptr)
	}

	// 检查指针是否为空
	if value.IsNil() {
		return ""
	}

	// 使用 Elem() 获取指针指向的值
	if value.Elem().Kind() == reflect.Invalid {
		return "" // 处理空值的情况
	}

	// 处理时间指向类型
	if value.Type() == reflect.TypeOf(&time.Time{}) {
		return fmt.Sprintf("%v", reflect.ValueOf(ptr).Elem().Interface().(time.Time).Format("2006-01-02 15:04:05"))
	}

	return fmt.Sprintf("%v", reflect.ValueOf(ptr).Elem().Interface())
}
