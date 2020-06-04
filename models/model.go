package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/config"
	"github.com/sirupsen/logrus"
	"github.com/yuanzhangcai/oxygen/common"
	"github.com/yuanzhangcai/oxygen/tools"
)

var (
	db                 *gorm.DB
	redis              *tools.Redis
	engineConfigPrefix string
	engineConfigEnv    string
)

// Model 数据库操作组件基类
type Model struct {
}

// Exec 执行sql语句
func (c *Model) Exec(sql string, values ...interface{}) *gorm.DB {
	return db.Exec(sql, values...)
}

type dbLogger struct {
}

func (c *dbLogger) Print(v ...interface{}) {
	logrus.Info(v...)
}

// Init 初始化顾
func Init() error {
	var err error
	dbInfo := config.Get("db", "server").String("")
	if dbInfo == "" {
		return fmt.Errorf("没有获取到数据库配置")
	}

	if db != nil {
		return nil
	}

	engineConfigPrefix = config.Get("engine_config", "prefix").String("")

	engineConfigEnv = "prod"
	if common.Env != common.EnvProd {
		engineConfigEnv = "test"
	}

	// 初始化连接
	db, err = gorm.Open("mysql", dbInfo)
	if err != nil {
		logrus.Error("数据库初始化失败。错误信息：" + err.Error())
		db = nil
		return err
	}

	// 取消DB复数
	db.SingularTable(true)

	if config.Get("db", "write_log").Bool(false) {
		// 设置sql语句输出到日志文件中
		db.LogMode(true)
		logger := &dbLogger{}
		db.SetLogger(logger)
	}

	redis = tools.GetRedis()

	return err
}

// GormTime Grom datetime类型
type GormTime struct {
	time.Time
}

// MarshalJSON 数据序列化
func (t GormTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(common.YMDHIS))
	return []byte(formatted), nil
}

// UnmarshalJSON json反序列表
func (t *GormTime) UnmarshalJSON(data []byte) error {
	tm, err := time.Parse(common.YMDHIS, string(data[1:len(data)-1]))
	if err != nil {
		return nil
	}
	*t = GormTime{Time: tm}
	return nil
}

// Value 返回datetime值
func (t GormTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 设置datetime值
func (t *GormTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = GormTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
