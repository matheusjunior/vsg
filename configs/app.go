package configs

const (
	appName     = "vgs"
	host        = "localhost:8080"
	metricsHost = "localhost:8091"
)

func GetAppName() string {
	return appName
}

func GetHost() string {
	return host
}

func GetMetricsHost() string {
	return metricsHost
}

func GetVoucherDuration() int32 {
	return 60 * 10 // 10 minutes
}
