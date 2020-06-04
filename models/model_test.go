package models

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"github.com/stretchr/testify/assert"
	"github.com/yuanzhangcai/oxygen/common"
	"github.com/yuanzhangcai/oxygen/tools"
)

var (
	server   = "redis:6379"
	password = "12345678"
	prefix   = ""
)

func init() {
	initConfig()
}

func initConfig() {
	common.CurrRunPath = os.Getenv("CI_PROJECT_DIR")
	if common.CurrRunPath == "" {
		common.CurrRunPath = "/Users/zacyuan/MyWork/oxygen/"
	}

	common.Env = common.EnvTest
	common.LoadConfig()

	str := `
	{
		"db" : {
			"server" : "zacyuan:zacyuan@(mysql:3306)/tds_user_pre?parseTime=true&loc=Local&charset=utf8",
			"write_log" : true
		}
	}`

	s := memory.NewSource(
		memory.WithJSON([]byte(str)),
	)

	_ = config.Load(s)

	// 初始化Redis
	_ = tools.InitRedis(server, password, prefix)
}

func TestInit(t *testing.T) {
	str := `
	{
		"db" : {
			"server" : ""
		}
	}`

	s := memory.NewSource(
		memory.WithJSON([]byte(str)),
	)

	_ = config.Load(s)
	err := Init()
	assert.NotNil(t, err)

	str = `
	{
		"db" : {
			"server" : "www"
		}
	}`

	s = memory.NewSource(
		memory.WithJSON([]byte(str)),
	)

	_ = config.Load(s)
	db = nil
	err = Init()
	assert.NotNil(t, err)

	initConfig()

	err = Init()
	assert.Nil(t, err)
}

func TestExec(t *testing.T) {

	_ = Init()
	db := db.Exec("delete from tbActQual_19_1 where iUin = ?", 99999)
	assert.NotNil(t, db)
}
