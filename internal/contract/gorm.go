package contract

import (
	"github.com/go-sql-driver/mysql"
	"github.com/lackone/gin-scaffold/internal/core"
	"gorm.io/gorm"
	"net"
	"strconv"
	"time"
)

const GORMKey = "app:gorm"

type IGORM interface {
	GetDB(option ...DBOption) (*gorm.DB, error)
}

type DBOption func(container core.IContainer, config *DBConfig) error

type DBConfig struct {
	Username             string `yaml:"username"`               // 用户名
	Password             string `yaml:"password"`               // 密码
	Driver               string `yaml:"driver"`                 // 驱动
	Host                 string `yaml:"host"`                   // 数据库地址
	Port                 int    `yaml:"port"`                   // 端口
	Loc                  string `yaml:"loc"`                    // 时区
	Charset              string `yaml:"charset"`                // 字符集
	ParseTime            bool   `yaml:"parse_time"`             // 是否解析时间
	Protocol             string `yaml:"protocol"`               // 传输协议
	Dsn                  string `yaml:"dsn"`                    // 直接传递dsn，如果传递了，其他关于dsn的配置均无效
	Database             string `yaml:"database"`               // 数据库
	Collation            string `yaml:"collation"`              // 字符序
	WriteTimeout         string `yaml:"write_timeout"`          // 写超时时间
	ReadTimeout          string `yaml:"read_timeout"`           // 读超时时间
	Timeout              string `yaml:"timeout"`                // 连接超时时间
	AllowNativePasswords bool   `yaml:"allow_native_passwords"` // 是否允许nativePassword

	// 连接池配置
	ConnMaxIdle     int    `yaml:"conn_max_idle"`     // 最大空闲连接数
	ConnMaxOpen     int    `yaml:"conn_max_open"`     // 最大连接数
	ConnMaxLifetime string `yaml:"conn_max_lifetime"` // 连接最大生命周期
	ConnMaxIdletime string `yaml:"conn_max_idletime"` // 空闲最大生命周期

	// gorm配置
	*gorm.Config
}

func (c *DBConfig) FormatDSN() (string, error) {
	port := strconv.Itoa(c.Port)
	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return "", err
	}
	readTimeout, err := time.ParseDuration(c.ReadTimeout)
	if err != nil {
		return "", err
	}
	writeTimeout, err := time.ParseDuration(c.WriteTimeout)
	if err != nil {
		return "", err
	}
	location, err := time.LoadLocation(c.Loc)
	if err != nil {
		return "", err
	}
	driverConf := &mysql.Config{
		User:                 c.Username,
		Passwd:               c.Password,
		Net:                  c.Protocol,
		Addr:                 net.JoinHostPort(c.Host, port),
		DBName:               c.Database,
		Collation:            c.Collation,
		Loc:                  location,
		Timeout:              timeout,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
		ParseTime:            c.ParseTime,
		AllowNativePasswords: c.AllowNativePasswords,
	}
	return driverConf.FormatDSN(), nil
}
