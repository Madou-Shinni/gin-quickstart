package excel

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type ExcelTool struct {
	sheet               string
	file                *excelize.File
	sw                  *excelize.StreamWriter
	model               interface{}       // 结构体(head)
	list                interface{}       // 数据(body)
	mergeConditionIndex int               // 合并条件（列下标从0开始，例如第一列相同则输入0）
	mergeCols           []string          // 合并列
	remark              string            // 备注(A1单元格)
	TagCol              map[string]string // 结构体标签对应的列
	formatBool          map[bool]string   // bool格式化
}

func NewExcelTool(sheet string) *ExcelTool {
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", sheet)
	sw, err := file.NewStreamWriter(sheet)
	if err != nil {
		return nil
	}

	return &ExcelTool{
		sheet:      sheet,
		file:       file,
		sw:         sw,
		TagCol:     make(map[string]string),
		formatBool: map[bool]string{true: "是", false: "否"},
	}
}

func (e *ExcelTool) FormatBool(m map[bool]string) *ExcelTool {
	for b, s := range m {
		e.formatBool[b] = s
	}
	return e
}

func (e *ExcelTool) WriteHead(data interface{}) *ExcelTool {
	if e.model == nil {
		e.model = data
	}
	return e
}

func (e *ExcelTool) WriteBody(list interface{}) *ExcelTool {
	e.list = list
	return e
}

// Remark 备注
// mergeEndCell 合并结束单元格
func (e *ExcelTool) Remark(remark string) *ExcelTool {
	e.remark = remark
	return e
}

// MergeCols 合并列
// headCondition 合并条件 excel标签
// mergeHeads 合并列 被合并的excel标签
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
	if err := e.setRemark(); err != nil {
		return err
	}
	if err := e.StreamWriteHead(e.sw, e.model); err != nil {
		return err
	}
	if e.list != nil {
		err := e.StreamWriteBodyWithMerge(e.sw, e.list, e.mergeConditionIndex, e.mergeCols)
		if err != nil {
			return err
		}
	}
	return e.sw.Flush()
}

func (e *ExcelTool) setRemark() error {
	if e.remark == "" {
		return nil
	}
	if err := e.sw.MergeCell("A1", "Z1"); err != nil {
		return err
	}

	// 文字靠左顶部对其
	style := &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "top",
		},
	}
	h := 16
	newStyle, err := e.file.NewStyle(style)
	if err != nil {
		return err
	}

	// 查询remark有多少个换行符
	h = h * strings.Count(e.remark, "\n")
	if h > 100 {
		h = 100
	}

	return e.sw.SetRow(fmt.Sprintf("A%d", 1), []interface{}{e.remark}, excelize.RowOpts{Height: float64(h), StyleID: newStyle})
}

func (e *ExcelTool) StreamWriteHead(sw *excelize.StreamWriter, data interface{}) error {
	row := 1
	if e.remark != "" {
		// 第一行备注
		row++
	}

	settingSlice := ParseExcelTag(data)
	rows := make([]interface{}, len(settingSlice)) // 创建一个切片，表示一行数据
	for i := range settingSlice {
		rows[i] = excelize.Cell{
			Value: settingSlice[i].Head,
		}
	}
	// 列名都是从列号1开始；行号从1开始
	axis, err := excelize.CoordinatesToCellName(1, row)
	if err != nil {
		return err
	}
	// 流式写入行，并指定高度
	return sw.SetRow(axis, rows, excelize.RowOpts{Height: 16})
}

func (e *ExcelTool) formatValue(v reflect.Value) (any, error) {
	// 处理指针类型
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil // 或返回 "", 0，视你的使用场景而定
		}
		v = v.Elem() // 解引用
	}

	switch v.Kind() {
	case reflect.Bool:
		if e.formatBool != nil {
			return e.formatBool[v.Bool()], nil
		}
	default:
		return v.Interface(), nil
	}

	return nil, errors.New("not support")
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
					cellValue, err := e.formatValue(field)
					if err != nil {
						return err
					}
					row[j] = excelize.Cell{
						StyleID: 0,
						Formula: "",
						Value:   cellValue,
					}
				}
				// 将每一行数据保存到二维数组中
				data[i] = row
			}
		}

		row := 2
		if e.remark != "" {
			row++
		}

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

func (e *ExcelTool) Model(m interface{}) *ExcelTool {
	e.model = m
	tp := reflect.ValueOf(e.model).Type().Elem() // 获得结构体的反射Type
	numField := tp.NumField()                    // 获取结构体的字段数量
	for i := 0; i < numField; i++ {
		field := tp.Field(i)                 // 获取字段
		tag := field.Tag.Get("excel")        // 获取tag中的ex值
		e.TagCol[tag] = indexToColumnName(i) // 这是tag对应的列名
	}

	return e
}

func (e *ExcelTool) SetDropList(tagMap map[string][]string) error {
	for t, list := range tagMap {
		if _, ok := e.TagCol[t]; !ok {
			return fmt.Errorf("tag %s not found", t)
		}

		dvRange := excelize.NewDataValidation(true)
		s := e.TagCol[t]
		dvRange.SetSqref(fmt.Sprintf("%s2:%s65535", s, s))
		err := dvRange.SetDropList(list)
		if err != nil {
			return err
		}
		err = e.file.AddDataValidation(e.sheet, dvRange)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ExcelTool) SetDropListPro(tagMap map[string][]string) error {
	for t, list := range tagMap {
		if _, ok := e.TagCol[t]; !ok {
			return fmt.Errorf("tag %s not found", t)
		}

		// 创建一个工作表用于存储选项(需要删除()，不然下拉列表会出现问题 )
		re := regexp.MustCompile(`\(.+?\)`) // 匹配括号及其中的内容
		result := re.ReplaceAllString(t, "")
		optionsSheet := result
		_, err := e.file.NewSheet(optionsSheet)
		if err != nil {
			return err
		}
		for i, s := range list {
			cell, _ := excelize.CoordinatesToCellName(1, i+1)
			err = e.file.SetCellValue(optionsSheet, cell, s)
			if err != nil {
				return err
			}
		}

		dvRange := excelize.NewDataValidation(true)
		s := e.TagCol[t]
		dvRange.SetSqref(fmt.Sprintf("%s2:%s65535", s, s))
		dvRange.SetSqrefDropList(fmt.Sprintf("%s!$A$1:$A$%d", optionsSheet, len(list)))

		err = e.file.AddDataValidation(e.sheet, dvRange)
		if err != nil {
			return err
		}
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
	// 检查输入是否为指针类型
	value := reflect.ValueOf(ptr)
	if value.Kind() != reflect.Ptr {
		return fmt.Sprintf("%v", ptr)
	}

	// 检查指针是否为空
	if value.IsNil() {
		return ""
	}

	// 获取指针指向的值
	elem := value.Elem()
	if !elem.IsValid() {
		return "" // 处理无效值的情况
	}

	// 特殊处理时间类型
	if value.Type() == reflect.TypeOf(&time.Time{}) {
		return elem.Interface().(time.Time).Format("2006-01-02 15:04:05")
	}

	// 处理自定义类型
	if elem.Kind() == reflect.Struct {
		if stringer, ok := elem.Interface().(fmt.Stringer); ok {
			return stringer.String() // 如果实现了 fmt.Stringer 接口，调用其 String 方法
		}
	}

	// 处理其他类型
	return fmt.Sprintf("%v", elem.Interface())
}
