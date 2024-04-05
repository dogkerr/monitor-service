package rest
import (
	"github.com/prometheus/client_golang/api"
)

type PrometheusAPI struct {
	client *api.Client
}

func NewPrometheusAPI(adress string) (*PrometheusAPI, error) {
	conf := api.Config{
		Address: "http://localhost:9090",
	}
	promeClient, err := api.NewClient(conf);
	if err != nil {
		return nil,  err 
	}
	return &PrometheusAPI{client: &promeClient}, nil
}