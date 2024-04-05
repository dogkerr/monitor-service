package webapi

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"time"
	"dogker/lintang/monitor-service/helper"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"go.uber.org/zap"
)

type PrometheusAPIImpl struct {
	client *api.Client
}

func NewPrometheusAPI(adress string) *PrometheusAPIImpl {
	conf := api.Config{
		Address: adress,
	}
	promeClient, err := api.NewClient(conf)
	if err != nil {
		zap.L().Fatal("error pas init prometheus client", zap.Error(err))
		helper.PanicIfError(err)
	}

	return &PrometheusAPIImpl{client: &promeClient}
}

func (p *PrometheusAPIImpl) GetUserContainerResourceUsageRequest(ctx context.Context, userId string, fromTimeIn string) (domain.Prometheus, error) {
	api := v1.NewAPI(*p.client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	fromTime, err := time.Parse(time.RFC3339, fromTimeIn)
	if err != nil {
		return domain.Prometheus{}, domain.ErrBadParamInput
	}

	defer cancel()
	r := v1.Range{
		Start: time.Now().Add(-time.Duration(fromTime.Second())),
		End:   time.Now(),
		Step:  time.Minute,
	}
	cpuResults, warnings, err := api.QueryRange(ctx, "sum(rate(container_cpu_usage_seconds_total{user_id="+userId+"}[30d]))  * 30 * 24 * 3600 / (12 * 3600)", r, v1.WithTimeout(5*time.Second))

	if err != nil {
		zap.L().Error("Gagal mendapatkan query CPU Usage", zap.Error(err))
		return domain.Prometheus{}, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}
	memoryResults, warnings, err := api.QueryRange(ctx, "avg_over_time(container_memory_usage_bytes{container_label_user_id=~"+userId+"}[30d]) * 30*24*3600 / 3600 / (1024^3)", r, v1.WithTimeout(5*time.Second))

	if err != nil {
		zap.L().Error("Gagal mendapatkan query Memory Usage", zap.Error(err))
		return domain.Prometheus{}, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	networkIngress, warnings, err := api.QueryRange(ctx, "sum(increase(container_network_receive_bytes_total{container_label_user_id=~"+userId+"}[30m])) / 1024", r, v1.WithTimeout(5*time.Second))

	if err != nil {
		zap.L().Error("Gagal mendapatkan query Network Ingress Usage", zap.Error(err))
		return domain.Prometheus{}, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	networkEgress, warnings, err := api.QueryRange(ctx, "sum(increase(container_network_transmit_bytes_total{container_label_user_id=~"+userId+"}[30m])) / 1024", r, v1.WithTimeout(5*time.Second))
	if err != nil {
		zap.L().Error("Gagal mendapatkan query Network Egress Usage", zap.Error(err))
		return domain.Prometheus{}, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	grpcRes := domain.Prometheus{
		CurrentTime:            time.Now(),
		AllCpuUsage:            float32(cpuResults.Type()),
		AllMemoryUsage:         float32(memoryResults.Type()),
		AllNetworkIngressUsage: float32(networkIngress.Type()),
		AllNetworkEgressUsage:  float32(networkEgress.Type()),
		FromTime:               fromTime,
	}
	return grpcRes, nil
}
