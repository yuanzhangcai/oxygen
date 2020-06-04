package oxygen

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/v2/config"
	"github.com/sirupsen/logrus"
	"github.com/yuanzhangcai/oxygen/common"
	"github.com/yuanzhangcai/oxygen/log"
	"github.com/yuanzhangcai/oxygen/models"
	"github.com/yuanzhangcai/oxygen/monitor"
	"github.com/yuanzhangcai/oxygen/services"
	"github.com/yuanzhangcai/oxygen/tools"
)

func init() {
	// 获取程序运行目录信息
	common.GetRunInfo()

	// 获取当前运行环境
	common.GetEnv()

	// 加载配置文件
	common.LoadConfig()

	// 显示版本信息
	common.ShowInfo()

	// 初始化log
	err := log.InitLogrus(nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// 初始化监控
	monitor.Init()

	// 初始化Redis
	if err = tools.InitRedis(config.Get("redis", "server").String(""),
		config.Get("redis", "password").String(""),
		config.Get("redis", "prefix").String("")); err != nil {
		logrus.Fatal(err)
	}

	// 初始化DB
	if err := models.Init(); err != nil {
		logrus.Fatal(err)
	}
}

// Start 开始服务
func Start(setRouter func(router *gin.Engine)) {
	pprof := config.Get("pprof", "server").String("")
	if pprof != "" {
		go func() {
			_ = http.ListenAndServe(pprof, nil) // 非正式环境，开启pprof服务
		}()
	}

	services.Start(setRouter)
}
