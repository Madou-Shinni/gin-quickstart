package initialize

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"log"
)

func init() {
	if err := tools.InitTrans("zh"); err != nil {
		log.Printf("init trans failed, err:%v\n", err)
		return
	}
}
