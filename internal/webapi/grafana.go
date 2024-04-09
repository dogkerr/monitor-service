package webapi

import (
	"context"
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/domain"
	"encoding/json"
	"io"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
	"go.uber.org/zap"
)

type GrafanaAPI struct {
	client *goapi.GrafanaHTTPAPI
}

func NewGrafanaAPI(cfg *config.Config) *GrafanaAPI {
	grafanaCfg := goapi.TransportConfig{
		Host:      cfg.Grafana.URLGrafana,
		BasePath:  "/api",
		Schemes:   []string{"http"},
		APIKey:    cfg.Grafana.Apikey,
		OrgID:     1,
		BasicAuth: url.UserPassword("admin", "password"),
	}

	grafanaClient := goapi.NewHTTPClientWithConfig(strfmt.Default, &grafanaCfg)

	return &GrafanaAPI{grafanaClient}
}

func (g *GrafanaAPI) CreateDashboard(ctx context.Context, userID string) (*domain.Dashboard, error) {
	// get prometheus datasource id
	getDt, err := g.client.Datasources.GetDataSources()
	if err != nil {
		zap.L().Error("datasources not found")
		return nil, err
	}
	datasources := getDt.Payload
	datasourceID := ""
	for _, dt := range datasources {
		if dt.Type == "prometheus" {
			zap.L().Debug("get prometheus data source", zap.String("datasource_id", dt.UID))
			datasourceID = dt.UID
		}
	}
	path, _ := os.Getwd()
	// Open our jsonFile
	jsonFile, err := os.Open(path + "./config/docker-quest-prometheus.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		zap.L().Error("grafana config file not found", zap.String("path", "../config/docker-quest-prometheus.json"))
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		zap.L().Error("file not found")
		return nil, err
	}
	var grafanaConfig domain.GrafanaConfig

	err = json.Unmarshal(byteValue, &grafanaConfig)
	if err != nil {
		zap.L().Error("json.Unmarshal")
		// return nil, err
	}

	// mengubah isi config grafana sesuai dg userId yg meminta request
	grafanaConfig.Dashboard.Inputs[0].PluginID = datasourceID
	for i := 0; i < len(grafanaConfig.Dashboard.Panels); i++ {
		grafanaConfig.Dashboard.Panels[i].Datasource = datasourceID
	}

	for i := 0; i < len(grafanaConfig.Dashboard.Templating.List); i++ {
		grafanaConfig.Dashboard.Templating.List[i].Datasource = datasourceID
	}

	// idxCpuUsage := 14
	// idxSentNetworkTraffic := 13
	// idxRcvdNetworkTraffic := 12
	// idxSwap := 15
	// idxMemUsage :=16
	// idxLimitMemory := 17
	// idxUsageMemory := 18
	// idxRemainingMemory=19
	user := userID
	cpuUsageByUser := "sum(rate(container_cpu_usage_seconds_total{container_label_user_id=~\"" + user + "\"}[$interval])) by (name) * 100"
	grafanaConfig.Dashboard.Panels[14].Targets[0].Expr = cpuUsageByUser

	sentNetworkTrafficByUser := "sum(rate(container_network_transmit_bytes_total{container_label_user_id=~\"" + user + " \"}[$interval])) by (name)"

	grafanaConfig.Dashboard.Panels[13].Targets[0].Expr = sentNetworkTrafficByUser
	grafanaConfig.Dashboard.Panels[12].Targets[0].Expr = "sum(rate(container_network_receive_bytes_total{container_label_user_id=~\"" + user + "\"}[$interval])) by (name)"

	grafanaConfig.Dashboard.Panels[15].Targets[0].Expr = "sum(container_memory_swap{container_label_user_id=~\"" + user + "\"}) by (name)"

	grafanaConfig.Dashboard.Panels[16].Targets[0].Expr = "sum(container_memory_rss{container_label_user_id=~\"" + user + "\"}) by (name)"
	grafanaConfig.Dashboard.Panels[16].Targets[1].Expr = "container_memory_usage_bytes{container_label_user_id=~\"" + user + "\"}"

	grafanaConfig.Dashboard.Panels[17].Targets[0].Expr = "sum(container_spec_memory_limit_bytes{container_label_user_id=~\"" + user + "\"} - container_memory_usage_bytes{container_label_user_id=~\"" + user + "\"}) by (name) "

	grafanaConfig.Dashboard.Panels[18].Targets[2].Expr = " container_memory_usage_bytes{container_label_user_id=~\"" + user + "\"} "
	grafanaConfig.Dashboard.Panels[1].Targets[0].Expr = "count(rate(container_last_seen{container_label_user_id=~\"" + user + "\"}[$interval]))"

	grafanaConfig.Dashboard.Panels = grafanaConfig.Dashboard.Panels[:len(grafanaConfig.Dashboard.Panels)-1]

	randomString := generateRandomString(8)
	grafanaConfig.Dashboard.Title = randomString
	grafanaConfig.Dashboard.UID = randomString
	grafanaConfig.Dashboard.Style = "light"

	// file, _ := json.MarshalIndent(grafanaConfig, "", " ")

	// newDb, err := os.Open(path + "/" + randomString + ".json")
	// file, _ := json.Marshal(grafanaConfig)
	makeDB, err := g.client.Dashboards.PostDashboard(&models.SaveDashboardCommand{
		Dashboard: grafanaConfig.Dashboard,
	})
	if err != nil {
		zap.L().Error("cant create new dashboard", zap.String("userId", userID))
		return nil, err
	}

	dbUID := *makeDB.GetPayload().UID
	return &domain.Dashboard{
		Uid:   dbUID,
		Owner: userID,
		Type:  "monitor",
	}, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed) //nolint: gosec // asda
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
} //nolint: gosec // asda
