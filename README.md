# github.com/Madou-Shinni/gin-quickstart

## 目录结构
```shell
├─api // 接口
│  ├─handle // 接口处理函数类似controller
│  └─routers // 路由
├─cmd // 程序主入口
├─configs // 配置文件
├─docs // 生成的接口文档
├─initialization // 初始化配置信息
├─internal // 该服务所有不对外暴露的代码，通常的业务逻辑都在这下面，使用internal避免错误引用
│  ├─conf // 内部使用的config的结构定义，将文件配置反序列化到结构体
│  ├─data // 业务数据访问 持久层
│  ├─domain // 实体类
│  └─service // 业务层
├─middleware // 中间件
└─pkg // 外部依赖
    ├─constant // 常量
    ├─global // 全局变量
    ├─request // 请求参数
    ├─response // 返回参数
    └─tools // 工具包
        ├─jwt // jwt工具
        ├─letter // 生成随机字母工具
        ├─pagelimit // 将page参数转化为数据库的offset和limit工具
        └─snowflake // 雪花算法工具
```

## 快速开始

### 修改配置文件
```yml
app:
  # 环境：dev or prod
  env: dev
  # 分布式节点(用于生成分布式id应保证各个机器节点不一致)
  machineID: 1
  server-port: 8080
# 数据库
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "123456"
  dbname: "sni-msg"
# redis
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
# jwt
jwt:
  # 过期时间(秒)
  # access_token 过期时间 1小时
  access-expire: 3600
  # refresh_token 过期时间 7天
  refresh-expire: 604800
  # 签名
  issuer: sni
  # 密钥
  secret: yoursecret
```

### 注册路由

在`/api/handle`目录下创建go文件，例如demo.go，定义路由处理函数
```go
package handle

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/gin-gonic/gin"
)

type DemoHandle struct {
	s *service.DemoService
}

func NewDemoHandle(s *service.DemoService) *DemoHandle {
	return &DemoHandle{s: s}
}

// Add 创建Demo
// @Tags     Demo
// @Summary  创建Demo
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.Demo true "创建Demo"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /demo [post]
func (cl *DemoHandle) Add(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindJSON(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Add(demo); err != nil {
		response.Error(c, constant.CODE_ADD_FAILED, constant.CODE_ADD_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Delete 删除Demo
// @Tags     Demo
// @Summary  删除Demo
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.Demo true "删除Demo"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /demo [delete]
func (cl *DemoHandle) Delete(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindJSON(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Delete(demo); err != nil {
		response.Error(c, constant.CODE_DELETE_FAILED, constant.CODE_DELETE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Update 修改Demo
// @Tags     Demo
// @Summary  修改Demo
// @accept   application/json
// @Produce  application/json
// @Param    data body     domain.Demo true "修改Demo"
// @Success  200  {string} string            "{"code":200,"msg":"","data":{}"}"
// @Router   /demo [put]
func (cl *DemoHandle) Update(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindJSON(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	if err := cl.s.Update(demo); err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c)
}

// Find 查询Demo
// @Tags     Demo
// @Summary  查询Demo
// @accept   application/json
// @Produce  application/json
// @Param    data query     domain.Demo true "查询Demo"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /demo [get]
func (cl *DemoHandle) Find(c *gin.Context) {
	var demo domain.Demo
	if err := c.ShouldBindQuery(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.Find(demo)

	if err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c, res)
}

// List 查询Demo列表
// @Tags     Demo
// @Summary  查询Demo列表
// @accept   application/json
// @Produce  application/json
// @Param    data query     domain.Demo true "查询Demo列表"
// @Success  200  {string} string            "{"code":200,"msg":"查询成功","data":{}"}"
// @Router   /demo/list [get]
func (cl *DemoHandle) List(c *gin.Context) {
	var demo domain.PageDemoSearch
	if err := c.ShouldBindQuery(&demo); err != nil {
		response.Error(c, constant.CODE_INVALID_PARAMETER, constant.CODE_INVALID_PARAMETER.Msg())
		return
	}

	res, err := cl.s.List(demo)

	if err != nil {
		response.Error(c, constant.CODE_UPDATE_FAILED, constant.CODE_UPDATE_FAILED.Msg())
		return
	}

	response.Success(c, res)
}
```

在`/api/routers`目录下创建go文件，例如demo.go
```go
package routers

import (
	"github.com/Madou-Shinni/gin-quickstart/api/handle"
	"github.com/Madou-Shinni/gin-quickstart/internal/data"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/gin-gonic/gin"
)

// 注册路由
func DemoRouterRegister(r *gin.Engine) {
	demoGroup := r.Group("demo")
	demoHandle := handle.NewDemoHandle(service.NewDemoService(data.NewDemoRepo()))
	{
		demoGroup.POST("", demoHandle.Add)
		demoGroup.DELETE("", demoHandle.Delete)
		demoGroup.GET("", demoHandle.Find)
		demoGroup.GET("/list", demoHandle.List)
		demoGroup.PUT("", demoHandle.Update)
	}
}
```

在`/initialization/router.go`文件里添加注册路由的信息`routers.DemoRouterRegister(r)`
```go
package initialization

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/api/routers"
	_ "github.com/Madou-Shinni/gin-quickstart/docs"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	// 初始化引擎
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 设置 swagger 访问路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	// 跨域
	//r.Use(cors.Default())

	// 注册路由
	routers.DemoRouterRegister(r)

	r.Run(fmt.Sprintf(":%v", conf.Conf.ServerPort))
}
```

### 建立实体domain层

在`internal/domain/`目录下添加go文件，例如demo.go
```go
package domain

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/request"
	"gorm.io/gorm"
)

type Demo struct {
	gorm.Model
	Did  int64  `json:"did,omitempty" form:"did" gorm:"column:did;comment:唯一标识"`                                        // 唯一标识
	Uid  int64  `json:"uid,omitempty" form:"uid" gorm:"column:uid;index:idx_uid_code,unique;comment:用户id"`              // 用户id
	Code string `json:"code,omitempty" form:"code" gorm:"column:code;index:idx_uid_code,unique;comment:code值(验证请求来源);"` // code
	QPS  int64  `json:"qps,omitempty" form:"qps" gorm:"column:qps;comment:每秒请求量;"`                                      // 每秒请求量
}

type PageDemoSearch struct {
	Demo
	request.PageSearch
}

func (Demo) TableName() string {
	return "demo"
}
```

### 持久层data层

在`internal/data/`目录下添加go文件，例如demo.go，注意：你需要在service中先定义需要用到的方法接口
```go
type DemoRepo interface {
    Create(demo domain.Demo) error
    Delete(demo domain.Demo) error
    Update(demo domain.Demo) error
    Find(demo domain.Demo) (domain.Demo, error)
    List(page domain.PageDemoSearch) ([]domain.Demo, error)
    Count() (int64, error)
}
```
然后在data层实现接口中定义的方法
```go
package data

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/internal/service"
	"github.com/Madou-Shinni/gin-quickstart/pkg/global"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/pagelimit"
)

type DemoRepo struct {
}

func NewDemoRepo() service.DemoRepo {
	return &DemoRepo{}
}

func (s *DemoRepo) Create(demo domain.Demo) error {
	return global.DB.Create(&demo).Error
}

func (s *DemoRepo) Delete(demo domain.Demo) error {
	return global.DB.Delete(&demo).Error
}

func (s *DemoRepo) Update(demo domain.Demo) error {
	return global.DB.Updates(&demo).Error
}

func (s *DemoRepo) Find(demo domain.Demo) (domain.Demo, error) {
	db := global.DB.Model(&domain.Demo{})
	if demo.Did != 0 {
		db = db.Where("cid = ?", demo.Did)
	}

	if demo.Uid != 0 {
		db = db.Where("uid = ?", demo.Uid)
	}

	res := db.First(&demo)

	return demo, res.Error
}

func (s *DemoRepo) List(page domain.PageDemoSearch) ([]domain.Demo, error) {
	var (
		demoList []domain.Demo
		err      error
	)
	// db
	db := global.DB.Model(&domain.Demo{})
	// page
	offset, limit := pagelimit.OffsetLimit(page.PageNum, page.PageSize)

	if page.Did != 0 {
		db = db.Where("cid = ?", page.Did)
	}

	if page.Uid != 0 {
		db = db.Where("uid = ?", page.Uid)
	}

	err = db.Offset(offset).Limit(limit).Find(&demoList).Error

	return demoList, err
}

func (s *DemoRepo) Count() (int64, error) {
	var (
		count int64
		err   error
	)

	err = global.DB.Model(&domain.Demo{}).Count(&count).Error

	return count, err
}

```

### 业务处理service层

在`internal/service/`目录下添加go文件，例如demo.go
```go
package service

import (
	"github.com/Madou-Shinni/gin-quickstart/internal/domain"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/letter"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/snowflake"
	"github.com/Madou-Shinni/go-logger"
	"go.uber.org/zap"
)

// 定义接口
type DemoRepo interface {
	Create(demo domain.Demo) error
	Delete(demo domain.Demo) error
	Update(demo domain.Demo) error
	Find(demo domain.Demo) (domain.Demo, error)
	List(page domain.PageDemoSearch) ([]domain.Demo, error)
	Count() (int64, error)
}

type DemoService struct {
	repo DemoRepo
}

func NewDemoService(repo DemoRepo) *DemoService {
	return &DemoService{repo: repo}
}

func (s *DemoService) Add(demo domain.Demo) error {
	// 1.生成唯一标识
	// 因为我们在全局初始化的时候已经初始化了雪花算法的机器节点
	// 所以我们可以直接使用
	did := snowflake.GenerateID()

	// 2.生成code 20位
	code := letter.GenerateCode(20)

	demo.Did = did
	demo.Code = code

	// 3.持久化入库
	if err := s.repo.Create(demo); err != nil {
		// 4.记录日志
		logger.Error("s.repo.Create(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Delete(demo domain.Demo) error {
	if err := s.repo.Delete(demo); err != nil {
		logger.Error("s.repo.Delete(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Update(demo domain.Demo) error {
	if err := s.repo.Update(demo); err != nil {
		logger.Error("s.repo.Update(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return err
	}

	return nil
}

func (s *DemoService) Find(demo domain.Demo) (domain.Demo, error) {
	res, err := s.repo.Find(demo)

	if err != nil {
		logger.Error("s.repo.Find(demo)", zap.Error(err), zap.Any("domain.Demo", demo))
		return res, err
	}

	return res, nil
}

func (s *DemoService) List(page domain.PageDemoSearch) (response.PageResponse, error) {
	var (
		pageRes response.PageResponse
	)

	data, err := s.repo.List(page)
	if err != nil {
		logger.Error("s.repo.List(demo)", zap.Error(err), zap.Any("domain.PageDemoSearch", page))
		return pageRes, err
	}

	count, err := s.repo.Count()
	if err != nil {
		logger.Error("s.repo.Count()", zap.Error(err), zap.Any("domain.PageDemoSearch", page))
		return pageRes, err
	}

	pageRes.Data = data
	pageRes.Total = count

	return pageRes, nil
}
```

### 启动程序

进入`/cmd/`目录运行main.go文件

### 生成接口文档

进入`/cmd/`目录打开控制台输出命令`swag init --parseDependency --parseInternal --output ../docs`
```go
// 生成swagger文档
// --parseDependency --parseInternal 识别到外部依赖
// --output 文件生成目录
//go:generate swag init --parseDependency --parseInternal --output ../docs
```

## 组件

### 日志组件

我们对zap日志进行了封装处理，以便于更简单的使用日志，如果你想自定义日志的使用可以修改`initialization/logger.go`文件中日志的初始化配置

