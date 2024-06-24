package conf

var Conf = new(ProfileInfo)

type ProfileInfo struct {
	*App          `mapstructure:"app"`
	*MysqlConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
	*JwtConfig    `mapstructure:"jwt"`
	*UploadConfig `mapstructure:"upload"`
}

// 系统配置
type App struct {
	Env        string `mapstructure:"env"`
	MachineID  int64  `mapstructure:"machineID"`
	ServerPort int    `mapstructure:"server-port"`
	LogFile    string `mapstructure:"log-file"`
}

// mysql配置
type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
}

// redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// jwt配置
type JwtConfig struct {
	AccessExpire  int64  `mapstructure:"access-expire"`
	RefreshExpire int64  `mapstructure:"refresh-expire"`
	Issuer        string `mapstructure:"issuer"`
	Secret        string `mapstructure:"secret"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Dir string `mapstructure:"dir"`
}
