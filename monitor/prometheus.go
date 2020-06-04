package monitor

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/yuanzhangcai/oxygen/common"
	"github.com/yuanzhangcai/oxygen/tools"
)

const (
	namespace = "tds"
	subsystem = "oxygen"
)

var (
	once     sync.Once
	register *tools.ServicesRegister
	sIP      string // 当前机器IP
	srv      *http.Server

	// actVisitCount 活动访问量
	actVisitCount *prometheus.CounterVec

	// callAIPCount api调用失败
	callAIPCount *prometheus.CounterVec

	// qualOperateFailedCount 资格操作失败次数
	qualOperateFailedCount *prometheus.CounterVec

	// chaosCostTime chaos总耗时情况统计
	chaosCostTime *prometheus.SummaryVec

	// callAPIUsedTime 调用api耗时情况统计
	callAPICostTime *prometheus.SummaryVec

	// uriCount 各uri调用资数
	uriCount *prometheus.CounterVec
)

// SetMetrics 设置监控指标
func SetMetrics() {
	once.Do(func() {
		// 获取本机IP
		sIP = common.GetIntranetIP()

		env := ""
		if common.Env != common.EnvProd {
			env = "_" + common.Env
		}

		// actVisitCount 活动访问量
		actVisitCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "act_visit_count" + env,
				Help:      "act visit count.",
			},
			[]string{"ip", "act_id", "act_name"},
		)

		// callAIPCount api调用失败
		callAIPCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "call_api_count" + env,
				Help:      "call api count.",
			},
			[]string{"ip", "api_id", "name", "act_id", "flow_id", "rule_id", "ret"},
		)

		// qualOperateFailedCount 资格操作失败次数
		qualOperateFailedCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "qual_operate_failed_count" + env,
				Help:      "qual operate failed count.",
			},
			[]string{"ip", "act_id", "flow_id", "pub_qual_id", "rule_id", "type"},
		)

		chaosCostTime = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Subsystem:  subsystem,
				Name:       "cost_time_seconds" + env,
				Help:       "oxygen const time.",
				Objectives: map[float64]float64{0.5: 0.05, 0.7: 0.03, 0.8: 0.02, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"ip"},
		)

		callAPICostTime = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Subsystem:  subsystem,
				Name:       "call_api_cost_time_seconds" + env,
				Help:       "call api const time.",
				Objectives: map[float64]float64{0.5: 0.05, 0.7: 0.03, 0.8: 0.02, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"ip", "api_id", "api_name"},
		)

		// uriCount 总访问量
		uriCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "uri_count" + env,
				Help:      "uri count",
			},
			[]string{"ip", "uri"},
		)

		// 注册监控指标
		prometheus.MustRegister(
			actVisitCount,
			callAIPCount,
			qualOperateFailedCount,
			chaosCostTime,
			callAPICostTime,
			uriCount,
		)
	})
}

// Init 初始化prometheus监控
func Init() {
	if srv != nil {
		return
	}

	// 设置监控指标
	SetMetrics()

	addr := config.Get("monitor", "server").String("")
	if addr != "" { // 开启prometheus监控
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		srv = &http.Server{}
		srv.Addr = addr
		srv.Handler = mux

		go func() {
			// http.Handle("/metrics", promhttp.Handler())
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("listen: %s\n", err)
			}
		}()

		etcdAddrs := config.Get("common", "etcd_addrs").StringSlice([]string{})
		if len(etcdAddrs) > 0 {
			serverName := config.Get("common", "server_name").String("oxygen.zacyuan.com") // 微务服名称
			if common.Env != common.EnvProd {
				serverName += "." + common.Env // 如果当前环境不是正式环境，服务名称添加环境后缀
			}
			serverName += ".metrics"

			register = tools.NewServicesRegister(&tools.RegisterOptions{
				ServerName:    serverName,
				EtcdAddress:   config.Get("common", "etcd_addrs").StringSlice([]string{}),
				ServerAddress: addr,
				Interval:      time.Duration(config.Get("common", "register_interval").Int(15)) * time.Second,
				TTL:           time.Duration(config.Get("common", "register_ttl").Int(30)) * time.Second,
			})

			_ = register.Start()
		}
	}
}

// Stop 停止监控上报
func Stop() {
	if register != nil {
		_ = register.Stop()
		register = nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logrus.Info("Monitor Shutdown Server ...")
	if srv != nil {
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Error("Server Shutdown:", err)
		}
	}
	logrus.Info("Monitor Server exiting")
}

// AddActVisitCount 总访问量加1
func AddActVisitCount(actID, actName string) {
	actVisitCount.WithLabelValues(sIP, actID, actName).Inc()
}

// AddCallAIPCount api接口调用次数统计
func AddCallAIPCount(apiID, apiName, actID, flowID, ruleID string, ret interface{}) {
	callAIPCount.WithLabelValues(sIP, apiID, apiName, actID, flowID, ruleID, common.ToString(ret)).Inc()
}

// AddQualOperateFailedCount 资格操作异常统计
func AddQualOperateFailedCount(actID, flowID, pubQualID, ruleID string, optType int) {
	qualOperateFailedCount.WithLabelValues(sIP, actID, flowID, pubQualID, ruleID, common.ToString(optType)).Inc()
}

// SummaryChaosCostTime 统计接口调用情况
func SummaryChaosCostTime(v float64) {
	chaosCostTime.WithLabelValues(sIP).Observe(v)
}

// SummaryCallAPICostTime 统计接口调用情况
func SummaryCallAPICostTime(apiID, apiName string, v float64) {
	callAPICostTime.WithLabelValues(sIP, apiID, apiName).Observe(v)
}

// AddURICount uri访问量加1
func AddURICount(uri string) {
	uriCount.WithLabelValues(sIP, uri).Inc()
}
