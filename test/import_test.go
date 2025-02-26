package test

import (
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"testing"
	"time"
)

func TestImportDemo(t *testing.T) {
	host := "http://localhost:8081"
	api := "/dataImport/import"
	jsonData := `
	{
		"category": "demo",
		"data": [{"name": "名称1", "age": 10, "birth_day": "2025-01-01 00:00:00"}, {"name": "名称2", "age": 15, "birth_day": "2025-01-01 00:00:00"}],
		"file_url": "http://cznu.gy/cfrx",
		"filename": "批量导入demo模板.xlsx"
	}
	`
	// 模拟错误
	//errJsonData := `
	//{
	//	"category": "demo",
	//	"data": [{"name": "名称1", "age": 10, "birth_day": "2025-01-01 00:00:00"}, {"name": "error", "age": 15, "birth_day": "2025-01-01 00:00:00"}],
	//	"file_url": "http://cznu.gy/cfrx",
	//	"filename": "批量导入demo模板.xlsx"
	//}
	//`

	headers := map[string]string{
		"accept": "application/json",
	}
	data := make(map[string]interface{})
	json.Unmarshal([]byte(jsonData), &data)

	v2, i, err := tools.NewRequestV2(tools.POST, 10*time.Second, host+api, data, headers)
	if err != nil {
		t.Fatalf("err: %v, i: %v", err, i)
	}

	t.Logf("v2: %v", string(v2))
}
