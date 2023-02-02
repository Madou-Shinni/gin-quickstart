package excel

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"reflect"
	"regexp"
)

//head，指定了此结构体字段对应的 Excel 列名。
//type，表示在使用反射进行数据解析时，会将此结构体字段的值作为指定的类型处理。
//select，表示此字段所在的列，包含一个下拉列表，列表中的枚举值由 select 后面的值指定。
//required，表示此字段必须包含非零值，否则在写入 Excel 时会报错。
//omitempty，表示此字段如果是零值，则对应的单元格留空。
//color，指定了列名所在单元格的颜色，通过这个字段，可以为不同的列名设置不同的底色，赋予一些含义，例如，可以将必填的列和选填的列，设置不同的底色。可以通过 Excel 的 RGB 颜色设置窗口，查看不同颜色对应的色号，作为 color 属性的值。
// 解析结果
type Setting struct {
	Head      string
	Type      string
	Select    []string
	Required  bool
	OmitEmpty bool
	Color     string
}

// 解析data中带ex的tag的字段，返回解析后的setting列表
// @param interface{}结构体指针
func ParseExcelTag(data interface{}) []Setting {
	var (
		setting     Setting
		settingList []Setting
	)
	tp := reflect.ValueOf(data).Type().Elem() // 获得结构体的反射Type
	numField := tp.NumField()                 // 获取结构体的字段数量

	for i := 0; i < numField; i++ {
		field := tp.Field(i)       // 获取字段
		tag := field.Tag.Get("ex") // 获取tag中的ex值
		if tag == "" {
			continue
		}
		s := parseFieldTag(setting, tag)
		settingList = append(settingList, s)
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
		}
	}

	return s
}

// 写入第一行标题数据
// @params: *excelize.StreamWriter流写入 interface{}结构体指针
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

// 写入除了标题行的内容数据
// @params: *excelize.StreamWriterexcel流式写入 interface{}切片结构体数据集
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
			record := values.Index(i)
			if record.Kind() == reflect.Struct {
				// 创建一个切片来表示一行数据
				row := make([]interface{}, record.NumField())
				for j := 0; j < record.NumField(); j++ {
					// 遍历结构体中的字段，取出字段值
					field := record.Field(j)
					row[j] = field.Interface()
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
