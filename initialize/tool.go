package initialize

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/snowflake"
)

func init() {
	// 初始化雪花算法工具
	snowflake.SnowflakeInit(conf.Conf.MachineID)
}
