package excel

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"mime/multipart"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// head，指定了此结构体字段对应的 Excel 列名。
// type，表示在使用反射进行数据解析时，会将此结构体字段的值作为指定的类型处理。
// select，表示此字段所在的列，包含一个下拉列表，列表中的枚举值由 select 后面的值指定。
// required，表示此字段必须包含非零值，否则在写入 Excel 时会报错。
// omitempty，表示此字段如果是零值，则对应的单元格留空。
// color，指定了列名所在单元格的颜色，通过这个字段，可以为不同的列名设置不同的底色，赋予一些含义，例如，可以将必填的列和选填的列，设置不同的底色。可以通过 Excel 的 RGB 颜色设置窗口，查看不同颜色对应的色号，作为 color 属性的值。
// 字段的解析结果
type Setting struct {
	Head      string
	Type      string
	Select    []string
	Required  bool
	OmitEmpty bool
	Color     string
}

// 解析data中带ex的tag的字段，返回解析后的setting列表
// param interface{}结构体指针
func ParseExcelTag(data interface{}) []Setting {
	var (
		setting     Setting
		settingList []Setting
	)
	tp := reflect.ValueOf(data).Type().Elem() // 获得结构体的反射Type
	numField := tp.NumField()                 // 获取结构体的字段数量

	for i := 0; i < numField; i++ {
		field := tp.Field(i)          // 获取字段
		tag := field.Tag.Get("excel") // 获取tag中的ex值
		if tag == "" {
			continue
		}
		//s := parseFieldTag(setting, tag)
		setting.Head = tag
		settingList = append(settingList, setting)
	}
	return settingList
}

// 解析tag到setting里面，返回setting
func parseFieldTag(s Setting, tag string) Setting {
	re := regexp.MustCompile(`(\w+):([^;]+)(;|$)`)
	attrs := re.FindAllStringSubmatch(tag, -1)

	for _, attr := range attrs {
		key := attr[1]
		value := attr[2]

		switch key {
		case "head":
			s.Head = value
		case "type":
			s.Type = value
		case "required":
			s.Required = true
		case "omitempty":
			s.OmitEmpty = true
		case "color":
			s.Color = value
		case "select":
			items := strings.Split(value, ",")
			s.Select = append(s.Select, items...)
		}
	}

	return s
}

// 写入第一行标题数据，并给指定列添加数据校验
// params: f *excelize.File写入 data interface{}结构体指针 map下拉选项 string列索引(A B C)，[]string选项
func WriteHead(f *excelize.File, data interface{}, dataValidation map[string][]string) error {
	var err error
	sheet, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}

	settingSlice := ParseExcelTag(data)
	row := make([]interface{}, len(settingSlice)) // 创建一个切片，表示一行数据
	for i := range settingSlice {
		row[i] = settingSlice[i].Head
	}
	axis, err := excelize.CoordinatesToCellName(1, 1)
	err = f.SetSheetRow("Sheet1", axis, &row)

	for s := range dataValidation {
		// 创建下拉选项列表
		dv := excelize.NewDataValidation(true)
		dv.SetSqref(fmt.Sprintf("%s1:%s1048576", s, s)) // 设置为整个 %s 列
		err = dv.SetDropList(dataValidation[s])
		err = f.AddDataValidation("Sheet1", dv)
	}

	// 设置活动工作表为 Sheet1
	f.SetActiveSheet(sheet)
	return err
}

// 写入第一行标题数据
// params: *excelize.StreamWriter流写入 interface{}结构体指针
func StreamWriteHead(sw *excelize.StreamWriter, data interface{}) error {
	settingSlice := ParseExcelTag(data)
	rows := make([]interface{}, len(settingSlice)) // 创建一个切片，表示一行数据
	for i := range settingSlice {
		rows[i] = excelize.Cell{
			Value: settingSlice[i].Head,
		}
	}
	// 列名都是从列号1开始；行号从1开始
	axis, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		return err
	}
	// 流式写入行，并指定高度
	return sw.SetRow(axis, rows, excelize.RowOpts{Height: 16})
}

// 写入除了标题行的内容数据，按结构体属性顺序写入
// params: *excelize.StreamWriterexcel流式写入 interface{}切片结构体指针数据集
func StreamWriteBody(sw *excelize.StreamWriter, d interface{}) error {
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
		for i := range data {
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

// StreamWriterAllRows 流式写出数据集，合并单元格，自定义样式
// param: sw *excelize.StreamWriter 流式写入器
// param: rows [][]interface{} 所有行数据 []interface{} 一行所有列数据
// param: mergeCell 需要合并的单元格 "A1:B1"
func StreamWriterAllRows(sw *excelize.StreamWriter, content [][]interface{}, mergeCell ...string) {
	// 合并单元格
	for i := 0; i < len(mergeCell); i++ {
		if !strings.Contains(mergeCell[i], ":") {
			fmt.Printf("Invalid parameter: %s", mergeCell[i])
			return
		}
		split := strings.Split(mergeCell[i], ":")
		sw.MergeCell(split[0], split[1])
	}

	// 生成内容
	for i := range content {
		cell, err := excelize.CoordinatesToCellName(1, i+1)
		err = sw.SetRow(cell, content[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ParseExcelToSlice 解析excel中全部数据，返回[]T，与字段标签excel匹配的标题 列的内容为字段值
// 目前支持类型string int int32 int64 *int *int32 *int64 time.time *time.Time
//
//	type TestStruct struct {
//		Name      string     `excel:"姓名"`
//		DateTime  *time.Time `excel:"时间"`
//		StartTime time.Time  `excel:"开始时间"`
//		}
//
// +---------+---------------------+---------------------+
// |   姓名   |         时间         |        开始时间      |
// +---------+---------------------+---------------------+
// | 李嘉图   | 2000-01-20 00:00:00 | 2000-01-20 00:00:00 |
// | Ricardo | 2000-01-20 00:00:00 | 2000-01-20 00:00:00 |
// +---------+---------------------+---------------------+
func ParseExcelToSlice[T any](file multipart.File) ([]T, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	var slice []T
	var headerMap map[string]int // 列标题与索引的映射

	// 遍历所有工作表
	for _, s := range f.GetSheetList() {
		rows, err := f.GetRows(s)
		if err != nil {
			return nil, err
		}

		// 遍历所有行
		for i, row := range rows {
			if i == 0 {
				// 取标题行，记录标题及其索引
				headerMap = make(map[string]int)
				for j, colCell := range row {
					headerMap[colCell] = j
				}
				continue
			}

			// 通过反射获取传入的结构体字段的值和tag标签，设置值
			var structType T
			val := reflect.ValueOf(&structType).Elem()
			tp := val.Type()
			if val.Kind() == reflect.Struct {
				// 获取所有的字段和excel tag
				for k := 0; k < val.NumField(); k++ {
					field := val.Field(k)
					structField := tp.Field(k)
					excelTag := structField.Tag.Get("excel")
					// 根据指定tag标签来设置值
					if colIndex, ok := headerMap[excelTag]; ok {
						switch field.Type().Kind() {
						case reflect.String:
							val.Field(k).SetString(row[colIndex])
						case reflect.Int, reflect.Int32, reflect.Int64:
							intVal, _ := strconv.Atoi(row[colIndex])
							val.Field(k).SetInt(int64(intVal))
						case reflect.TypeOf(time.Time{}).Kind():
							parse, err := time.Parse("2006-01-02 15:04:05", row[colIndex])
							if err != nil {
								return nil, fmt.Errorf("err: %v\n", err)
							}
							val.Field(k).Set(reflect.ValueOf(parse))
						case reflect.Pointer: // 指针指向类型
							switch field.Type().Elem().Kind() {
							case reflect.TypeOf(time.Time{}).Kind():
								parse, err := time.Parse("2006-01-02 15:04:05", row[colIndex])
								if err != nil {
									return nil, fmt.Errorf("err: %v\n", err)
								}
								val.Field(k).Set(reflect.ValueOf(&parse))
							case reflect.Int, reflect.Int32, reflect.Int64:
								intVal, _ := strconv.Atoi(row[colIndex])
								val.Field(k).Set(reflect.ValueOf(&intVal))
							}
						}
					}
				}
			}
			if !reflect.ValueOf(structType).IsZero() {
				slice = append(slice, structType)
			}
		}
	}

	return slice, nil
}
