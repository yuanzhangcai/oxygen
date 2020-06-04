package monitor

import (
	"os"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"github.com/stretchr/testify/assert"
	"github.com/yuanzhangcai/oxygen/common"
)

func initConfig() {
	common.CurrRunPath = os.Getenv("CI_PROJECT_DIR")
	if common.CurrRunPath == "" {
		common.CurrRunPath = "/Users/zacyuan/MyWork/oxygen/"
	}

	common.Env = "test"
	common.LoadConfig()

	str := `
	{
		"common" : {
			"etcd_addrs" : ["127.0.0.1:2379"]
		}
	}`

	s := memory.NewSource(
		memory.WithJSON([]byte(str)),
	)

	_ = config.Load(s)
}

func init() {
	time.Sleep(1 * time.Second)
	initConfig()
}

func TestInit(t *testing.T) {
	Init()
	assert.NotNil(t, register)
	assert.NotNil(t, srv)
}

func TestStop(t *testing.T) {
	Init()

	Stop()
	assert.Nil(t, register)
}

func TestMetrics(t *testing.T) {
	Init()

	AddActVisitCount("10", "aa")

	AddCallAIPCount("11", "bb", "11", "13", "121", 0)

	AddQualOperateFailedCount("10", "11", "12", "13", 1)

	SummaryChaosCostTime(234)

	SummaryCallAPICostTime("11", "aa", 349)

	AddURICount("/engine")
}
