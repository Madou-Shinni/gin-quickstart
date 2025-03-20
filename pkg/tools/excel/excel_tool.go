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
}

func NewExcelTool(sheet string) *ExcelTool {
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", sheet)
	sw, err := file.NewStreamWriter(sheet)
	if err != nil {
		return nil
	}

	return &ExcelTool{
		sheet:  sheet,
		file:   file,
		sw:     sw,
		TagCol: make(map[string]string),
	}
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

func (e *ExcelTool) StreamWriteBodyWithMerge(sw *excelize.StreamWriter, d interface{}, mergeConditionIndex int, mergeCols []string) error {
	// 判断参数的有效性
	if sw == nil {
		return errors.New("streamWriter cannot be nil")
	}

	// 判断d的数据类型
	v := reflect.ValueOf(d)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return errors.New("data must be a slice or array")
	}

	// 空数据处理
	dataLen := v.Len()
	if dataLen == 0 {
		return nil // 没有数据，直接返回
	}

	// 预分配数据容量
	data := make([][]interface{}, dataLen)

	// 处理切片中的每个元素
	for i := 0; i < dataLen; i++ {
		item := v.Index(i)

		// 确保元素是指针且可解引用
		if item.Kind() != reflect.Ptr || item.IsNil() {
			return fmt.Errorf("item at index %d is not a valid pointer", i)
		}

		record := item.Elem()
		if record.Kind() != reflect.Struct {
			return fmt.Errorf("item at index %d is not a struct", i)
		}

		// 获取结构体字段数量
		fieldCount := record.NumField()
		row := make([]interface{}, fieldCount)

		// 处理每个字段
		for j := 0; j < fieldCount; j++ {
			field := record.Field(j)
			// 使用封装的函数处理字段值
			row[j] = excelize.Cell{
				StyleID: 0,
				Formula: "",
				Value:   formatValueForExcel(field.Interface()),
			}
		}

		// 保存数据行
		data[i] = row
	}

	// 计算起始行号
	startRow := 2
	if e.remark != "" {
		startRow++
	}

	// 合并单元格处理
	var pendingMerges []struct {
		startRow int
		endRow   int
		column   string
	}

	if len(mergeCols) > 0 && mergeConditionIndex >= 0 && mergeConditionIndex < len(data[0]) {
		// 初始化合并状态
		mergeStartRow := startRow
		lastValue := data[0][mergeConditionIndex]

		// 遍历所有行寻找合并点
		for i := 1; i < dataLen; i++ {
			currentValue := data[i][mergeConditionIndex]
			valueChanged := !reflect.DeepEqual(lastValue, currentValue)
			isLastRow := i == dataLen-1

			// 当值变化或到达最后一行时处理合并
			if valueChanged || isLastRow {
				mergeEndRow := startRow + i - 1
				if isLastRow && !valueChanged {
					mergeEndRow = startRow + i // 包含最后一行
				}

				// 只有当合并范围至少有两行时才进行合并
				if mergeEndRow > mergeStartRow {
					for _, col := range mergeCols {
						pendingMerges = append(pendingMerges, struct {
							startRow int
							endRow   int
							column   string
						}{mergeStartRow, mergeEndRow, col})
					}
				}

				// 重置合并状态
				mergeStartRow = startRow + i
				lastValue = currentValue
			}
		}
	}

	// 写入数据行
	currentRow := startRow
	for i, rowData := range data {
		// 逐行插入数据
		axis, err := excelize.CoordinatesToCellName(1, currentRow)
		if err != nil {
			return fmt.Errorf("failed to convert coordinates for row %d: %w", i+startRow, err)
		}

		if err := sw.SetRow(axis, rowData, excelize.RowOpts{Height: 16}); err != nil {
			return fmt.Errorf("failed to write row %d: %w", i+startRow, err)
		}

		currentRow++
	}

	// 执行所有的单元格合并
	for _, merge := range pendingMerges {
		startCell := fmt.Sprintf("%s%d", merge.column, merge.startRow)
		endCell := fmt.Sprintf("%s%d", merge.column, merge.endRow)

		// 合并单元格
		if err := sw.MergeCell(startCell, endCell); err != nil {
			return fmt.Errorf("failed to merge cells %s:%s: %w", startCell, endCell, err)
		}
	}

	return nil
}

// formatValueForExcel 格式化值以适应Excel的需求
func formatValueForExcel(value interface{}) interface{} {
	if value == nil {
		return ""
	}

	v := reflect.ValueOf(value)

	// 处理指针类型
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		return formatValueForExcel(v.Elem().Interface())
	}

	// 处理时间类型
	if t, ok := value.(time.Time); ok {
		return t.Format("2006-01-02 15:04:05")
	}

	// 返回原始值
	return value
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
