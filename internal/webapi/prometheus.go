package webapi

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PrometheusAPI struct {
	client *api.Client
}

func NewPrometheusAPI(address string) *PrometheusAPI {
	conf := api.Config{
		Address: address,
	}
	promeClient, err := api.NewClient(conf)
	if err != nil {
		zap.L().Fatal("error pas init prometheus client", zap.Error(err))
	}

	return &PrometheusAPI{client: &promeClient}
}

/*
Desc: mendapatkan metrics untuk semua container milik user
*/
func (p *PrometheusAPI) GetUserContainerResourceUsageRequest(ctx context.Context, userID string, fromTimeIn *timestamppb.Timestamp) (*domain.Prometheus, error) {
	promeAPI := v1.NewAPI(*p.client)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	r := v1.Range{
		Start: time.Now().AddDate(0, 0, -30),
		End:   time.Now(),
		Step:  time.Hour,
	}
	seconds := fmt.Sprintf("%f", (time.Since(fromTimeIn.AsTime())).Seconds())

	cpuResults, warnings, err := promeAPI.QueryRange(ctx, "sum(rate(container_cpu_usage_seconds_total{container_label_user_id=~\""+userID+"\"}[1h])) * 100 * "+seconds+" /3600", r, v1.WithTimeout(5*time.Second))
	if err != nil {
		zap.L().Error("Gagal mendapatkan query CPU Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}
	memoryResults, warnings, err := promeAPI.QueryRange(ctx, "sum(avg_over_time(container_memory_usage_bytes{container_label_user_id=~\""+userID+"\"}[1h])) * 30*24*3600 / 3600 / (1024^3)", r, v1.WithTimeout(5*time.Second))

	if err != nil {
		zap.L().Error("Gagal mendapatkan query Memory Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	networkIngress, warnings, err := promeAPI.QueryRange(ctx, "sum(container_network_receive_bytes_total{container_label_user_id=~\""+userID+"\"}) / 1024", r, v1.WithTimeout(5*time.Second))

	if err != nil {
		zap.L().Error("Gagal mendapatkan query Network Ingress Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	networkEgress, warnings, err := promeAPI.QueryRange(ctx, "sum(container_network_transmit_bytes_total{container_label_user_id=~\""+userID+"\"}) / 1024", r, v1.WithTimeout(5*time.Second))
	if err != nil {
		zap.L().Error("Gagal mendapatkan query Network Egress Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	promeRes := domain.Prometheus{
		CurrentTime: timestamppb.Now(),

		FromTime: fromTimeIn,
	}
	switch c := cpuResults.(type) {
	case model.Matrix:
		for _, s := range c {
			promeRes.AllCPUUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad

	switch c := memoryResults.(type) {
	case model.Matrix:
		for _, s := range c {
			promeRes.AllMemoryUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad
	switch c := networkIngress.(type) {
	case model.Matrix:
		for _, s := range c {
			promeRes.AllNetworkIngressUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad
	switch c := networkEgress.(type) {
	case model.Matrix:
		for _, s := range c {
			promeRes.AllNetworkEgressUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad
	return &promeRes, nil
}

// TODO: bikin Get Metrics Usage Per User
/*
	Desc: mendapatkan metrics untuk container dg serviceID tertentu
*/
func (p *PrometheusAPI) GetMetricsByServiceID(ctx context.Context, serviceID string, fromTimeIn *timestamppb.Timestamp) (*domain.Metric, error) {
	promeAPI := v1.NewAPI(*p.client)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	r := v1.Range{
		Start: time.Now().AddDate(0, 0, -30),
		End:   time.Now(),
		Step:  time.Hour,
	}
	seconds := fmt.Sprintf("%f", (time.Since(fromTimeIn.AsTime())).Seconds())

	cpuResults, warnings, err := promeAPI.QueryRange(ctx, "sum(rate(container_cpu_usage_seconds_total{container_label_com_docker_swarm_service_id=~\""+serviceID+".*\"}[1h])) * 100 * "+seconds+" /3600", r, v1.WithTimeout(5*time.Second))
	if err != nil {
		zap.L().Error("Gagal mendapatkan query CPU Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}
	memoryResults, warnings, err := promeAPI.QueryRange(ctx, "sum(avg_over_time(container_memory_usage_bytes{container_label_com_docker_swarm_service_id=~\""+serviceID+".*\"}[1h])) * 30*24*3600 / 3600 / (1024^3)", r, v1.WithTimeout(5*time.Second))

	if err != nil {
		zap.L().Error("Gagal mendapatkan query Memory Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	networkIngress, warnings, err := promeAPI.QueryRange(ctx, "sum(container_network_receive_bytes_total{container_label_com_docker_swarm_service_id=~\""+serviceID+".*\"}) / 1024", r, v1.WithTimeout(5*time.Second))

	if err != nil {
		zap.L().Error("Gagal mendapatkan query Network Ingress Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	networkEgress, warnings, err := promeAPI.QueryRange(ctx, "sum(container_network_transmit_bytes_total{container_label_com_docker_swarm_service_id=~\""+serviceID+".*\"}) / 1024", r, v1.WithTimeout(5*time.Second))
	if err != nil {
		zap.L().Error("Gagal mendapatkan query Network Egress Usage", zap.Error(err))
		return nil, err
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warnings: pas query CPU Usage\n")
	}

	var metric domain.Metric
	switch c := cpuResults.(type) {
	case model.Matrix:
		for _, s := range c {
			metric.CpuUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad
	switch c := memoryResults.(type) {
	case model.Matrix:
		for _, s := range c {
			metric.MemoryUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad
	switch c := networkIngress.(type) {
	case model.Matrix:
		for _, s := range c {
			metric.NetworkIngressUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad
	switch c := networkEgress.(type) {
	case model.Matrix:
		for _, s := range c {
			metric.NetworkEgressUsage = float32(s.Values[len(s.Values)-1].Value)
		}
	} //nolint: gocritic // asdsad
	return &metric, nil
}
