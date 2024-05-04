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
	client       *goapi.GrafanaHTTPAPI
	fileLocation string
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

	return &GrafanaAPI{client: grafanaClient, fileLocation: cfg.Grafana.GrafanaFileLoc}
}

// genrate monitoring dashboard prometheus untuk setiap docker swarm service milik user
func (g *GrafanaAPI) CreateMonitorDashboard(ctx context.Context, userID string) (*domain.Dashboard, error) {
	// get prometheus datasource id
	getDt, err := g.client.Datasources.GetDataSources()
	if err != nil {
		zap.L().Error("datasources not found", zap.Error(err))
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
	jsonFile, err := os.Open(path + g.fileLocation) // config/docker_prometheus_template.json buat didalem docker container
	// if we os.Open returns an error then handle it
	if err != nil {
		zap.L().Error("grafana monitor config template file not found", zap.String("path", g.fileLocation), zap.Error(err))
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		zap.L().Error("io.ReadAll(jsonFile) [grafana monitoring config file]", zap.Error(err))
		return nil, err
	}
	var grafanaConfig domain.GrafanaMonitorConfig

	err = json.Unmarshal(byteValue, &grafanaConfig)
	if err != nil {
		zap.L().Error("json.Unmarshal", zap.Error(err)) // jangan di return karena emang banyak field json yang ga sesuai struct , tapi tetep bisa generate dashboard
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

	makeDB, err := g.client.Dashboards.PostDashboard(&models.SaveDashboardCommand{
		Dashboard: grafanaConfig.Dashboard,
	})
	if err != nil {
		zap.L().Error("cant create new dashboard", zap.Error(err), zap.String("userId", userID))
		return nil, err
	}

	dbUID := *makeDB.GetPayload().UID
	return &domain.Dashboard{
		Uid:   dbUID,
		Owner: userID,
		Type:  "monitor",
	}, nil
}

// generate random string uid,title user dashboard
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

type lokiQueryVariable struct {
	Label  string `json:"label"`
	RefID  string `json:"refId"`
	Stream string `json:"stream"`
	Tipe   int    `json:"type"`
}

// create logs dashboard untuk semua docker swarm  service milik user
func (g *GrafanaAPI) CreateLogsDashboard(ctx context.Context, userID string) (*domain.Dashboard, error) {
	// get loki datasource id
	getDt, err := g.client.Datasources.GetDataSources()
	if err != nil {
		zap.L().Error("datasources not found", zap.Error(err))
		return nil, err
	}
	datasources := getDt.Payload
	datasourceID := ""

	for _, dt := range datasources {
		if dt.Type == "loki" {
			zap.L().Debug("dapet loki data source uid", zap.String("loki_uid", dt.UID))
			datasourceID = dt.UID
		}
	}

	path, _ := os.Getwd()
	jsonFile, err := os.Open(path + "/config/loki_logs_per_user_template.json") // /../config kalau debug, /config kalau gak debug

	if err != nil {
		zap.L().Error("grafana logs dashboard config template file not found", zap.String("path", "/config/loki_logs_per_user_template.json"), zap.Error(err))
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		zap.L().Error("io.ReadAll(jsonFile) [grafana logs config file]", zap.Error(err))
		return nil, err
	}

	var grafanaConfig domain.GrafanaLogsDashboard
	err = json.Unmarshal(byteValue, &grafanaConfig)
	if err != nil {
		zap.L().Debug("json.Unmarshal") // jangan di return karena emang banyak field json yang ga sesuai struct, tapi tetep bisa generate dashboard
	}

	for _, panel := range grafanaConfig.Dashboard.Panels {
		panel.Targets[0].Datasource.UID = datasourceID
	}
	grafanaConfig.Dashboard.Panels[0].Datasource.UID = datasourceID
	grafanaConfig.Dashboard.Panels[1].Datasource.UID = datasourceID

	list := grafanaConfig.Dashboard.Templating.List
	list[2].Datasource.UID = datasourceID
	// var queryNameloki interface{} = lokiQueryVariable{
	// 	label: "container_name",
	// 	refId: "LokiVariableQueryEditor-VariableQuery",
	// 	stream: "{userId=\"" + userID + "\"}",
	// 	tipe: 1,
	// }

	// list[2].Query = queryNameloki.(lokiQueryVariable)
	list[2].Query = lokiQueryVariable{
		Label:  "container_name",
		RefID:  "LokiVariableQueryEditor-VariableQuery",
		Stream: "{userId=\"" + userID + "\"}",
		Tipe:   1,
	} // bug disni gak ketulis ke dashboardnya

	list[2].Datasource.UID = datasourceID

	grafanaConfig.Dashboard.Time.From = "now-5d"

	randomString := generateRandomString(8)
	grafanaConfig.Dashboard.Title = randomString
	grafanaConfig.Dashboard.UID = randomString
	grafanaConfig.Dashboard.Style = "light"

	zap.L().Debug("new panel datasource ID: ", zap.String("panel.Datasource.UID ", grafanaConfig.Dashboard.Panels[0].Datasource.UID))
	zap.L().Debug("new panel datasource ID2: ", zap.String("panel.Datasource.UID ", grafanaConfig.Dashboard.Panels[1].Datasource.UID))

	makeLogsDB, err := g.client.Dashboards.PostDashboard(&models.SaveDashboardCommand{
		Dashboard: grafanaConfig.Dashboard,
	})

	if err != nil {
		zap.L().Error("cant create new dashboard", zap.Error(err), zap.String("userId", userID))
		return nil, err
	}
	dbUID := *makeLogsDB.GetPayload().UID
	return &domain.Dashboard{
		Uid:   dbUID,
		Owner: userID,
		Type:  "log",
	}, nil
}
