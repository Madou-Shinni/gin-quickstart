# github.com/Madou-Shinni/gin-quickstart

## 目录结构
```shell
├─api // 接口
│  ├─handle // 接口处理函数类似controller
│  └─routers // 路由
├─cmd // 程序主入口
├─configs // 配置文件
├─docs // 生成的接口文档
├─initialize // 初始化配置信息
├─internal // 该服务所有不对外暴露的代码，通常的业务逻辑都在这下面，使用internal避免错误引用
│  ├─conf // 内部使用的config的结构定义，将文件配置反序列化到结构体
│  ├─data // 业务数据访问 持久层
│  ├─domain // 实体类
│  └─service // 业务层
├─middleware // 中间件
└─pkg // 外部依赖
    ├─constant // 常量
    ├─global // 全局变量
    ├─model // 公共model
    ├─request // 请求参数
    ├─response // 返回参数
    └─tools // 工具包
        ├─jwt // jwt工具
        ├─letter // 生成随机字母工具
        ├─pagelimit // 将page参数转化为数据库的offset和limit工具
        └─snowflake // 雪花算法工具
```

## 快速开始
```bash
git clone https://github.com/Madou-Shinni/gin-quickstart.git

# 初始化
make init

# 依赖整理
go mod tidy

# 运行（默认读取configs目录下的所有yaml配置文件）
go run cmd/main.go [-c] [configPath]
```

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

### 代码生成器
`go install github.com/Madou-Shinni/gctl@latest`
```go
gctl -m Article 自动生成代码
生成代码如下
internal/domain/article.go
internal/service/article.go
internal/data/article.go
api/routers/article.go
api/handle/article.go

GLOBAL OPTIONS:
--module value, -m value  生成模块的名称
--help, -h                show help
```

### 启动程序

进入`/cmd/`目录运行main.go文件

### 生成接口文档
使用swag 1.7.x版本 `go install github.com/swaggo/swag/cmd/swag@v1.7.9`
```go
swag init
```

### 同步api

将api同步到数据库

```shell
make api-sync
```

## 组件

### 日志组件

我们对zap日志进行了封装处理，以便于更简单的使用日志，如果你想自定义日志的使用可以修改`initialize/logger.go`文件中日志的初始化配置

