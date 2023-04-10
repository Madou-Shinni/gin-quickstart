package main

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/str"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
)

type Temp struct {
	Module      string // 模块名
	ModuleLower string
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "module",
				Aliases:  []string{"m"},
				Usage:    "生成模块的名称",
				Required: true,
			},
		},
		Action: gen,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// 生成代码
func gen(c *cli.Context) error {
	var err error
	var wg sync.WaitGroup
	var tempSlice []string

	s := c.String("module")

	// 定义变量
	//k v
	data := Temp{
		Module:      s,
		ModuleLower: strings.ToLower(s[:1]) + s[1:],
	}

	// 遍历模板文件
	dir, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("read dir failed: %s", err.Error())
		return err
	}

	for _, entry := range dir {
		if strings.Contains(entry.Name(), "txt") {
			// 将模板文件名加入列表
			tempSlice = append(tempSlice, entry.Name())
		}
	}

	wg.Add(5)
	for i := 0; i < len(tempSlice); i++ {
		// 启动5个goroutine生成不同的模板文件
		go func(i int) {
			defer wg.Done()

			var t *template.Template
			var f *os.File

			// 解析模板文件
			t, err = template.ParseFiles(tempSlice[i])
			if err != nil {
				return
			}

			// 写出文件
			err = writeOutput(s, tempSlice[i], data, f, t)
			if err != nil {
				return
			}

		}(i)
	}

	wg.Wait()

	if err != nil {
		return err
	}

	fmt.Println("gen code success")

	return nil
}

// 写出文件
func writeOutput(module string, sliceItem string, data Temp, f *os.File, t *template.Template) error {
	dirname := strings.Split(sliceItem, "_")[0]
	outputDir := "."

	switch dirname {
	case "data":
		outputDir = "../../internal/data"
	case "domain":
		outputDir = "../../internal/domain"
	case "service":
		outputDir = "../../internal/service"
	case "handle":
		outputDir = "../../api/handle"
	case "route":
		outputDir = "../../api/routers"
	default:
		outputDir = "./gen_code"
	}

	// 创建文件夹If path is already a directory, MkdirAll does nothing and returns nil.
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	// 将驼峰式字符串转换为下划线式字符串
	module = str.CamelToSnake(module)

	// 渲染模板并将结果写入文件
	f, err = os.Create(fmt.Sprintf("%s/%s.go", outputDir, module))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
