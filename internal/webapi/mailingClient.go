package webapi

import (
	"bytes"
	"context"
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gojek/heimdall"
	"github.com/gojek/heimdall/httpclient"
	"go.uber.org/zap"
)

type MailingWebAPI struct {
	MailingURL string
}

func NewWebAPI(cfg *config.Config) *MailingWebAPI {
	return &MailingWebAPI{cfg.HTTPClient.MailingURL}
}

type promeWebhookRes struct {
	Message string `json:"message"`
}

func (api *MailingWebAPI) SendDownSwarmServiceToMailingService(ctx context.Context, ml domain.CommonLabelsMailing) error {
	// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
	backoffInterval := 2 * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := 5 * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	timeout := 1000 * time.Millisecond
	// Create a new client, sets the retry mechanism, and the number of times you would like to retry
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(4),
	)

	body, err := json.Marshal(&ml)
	if err != nil {
		zap.L().Error("json.Marshal (SendDownSwarmServiceToMailingService) (Mailing WebApi)", zap.Error(err))
		return err
	}
	reader := bytes.NewReader(body)
	res, err := client.Post(api.MailingURL, reader, nil)
	if err != nil {
		zap.L().Error("client.Post(api.MailingURL, reader, nil) (SendDownSwarmServiceToMailingService) (MailingWEbAPI)", zap.Error(err))
		return err
	}

	mailingResp := &promeWebhookRes{}
	derr := json.NewDecoder(res.Body).Decode(mailingResp)
	if derr != nil {
		zap.L().Error(" json.NewDecoder(res.Body).Decode(mailingResp) (SendDownSwarmServiceToMailingService) (MailingWEbAPI)", zap.Error(err))
		return err
	}
	zap.L().Info(fmt.Sprintf("send down service to mailing service successs!!!. mailing response: %s", mailingResp.Message))
	return nil
}
