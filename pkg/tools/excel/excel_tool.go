package excel

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"time"
)

type ExcelTool struct {
	file  *excelize.File
	sw    *excelize.StreamWriter
	model interface{} // 结构体(head)
	list  interface{} // 数据(body)
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

func (e *ExcelTool) Flush() error {
	err := StreamWriteHead(e.sw, e.model)
	if err != nil {
		return err
	}
	err = StreamWriteBody(e.sw, e.list)
	if err != nil {
		return err
	}
	return e.sw.Flush()
}

func (e *ExcelTool) WriteToBuffer() (*bytes.Buffer, error) {
	return e.file.WriteToBuffer()
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
