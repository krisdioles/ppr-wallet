package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/krisdioles/ppr-wallet/config"
)

type Bank1Client struct {
	hostname             string
	apikey               string
	disbursementEndpoint string
	httpClient           *http.Client
}

type IBank1Client interface {
	CreateDisbursement(ctx context.Context, requestParam *Bank1CreateDisbursementRequest) (*Bank1CreateDisbursementResponse, error)
}

func NewBank1Client(cfg *config.Bank1Config) IBank1Client {
	return &Bank1Client{
		hostname:             cfg.Hostname,
		apikey:               cfg.APIKey,
		disbursementEndpoint: cfg.DisbursementEndpoint,
		httpClient: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   5 * time.Second,
		},
	}
}

type Bank1CreateDisbursementRequest struct {
	ReferenceID string     `json:"reference_id"`
	Account     AccountObj `json:"account"`
	Amount      AmountObj  `json:"amount"`
}

type AccountObj struct {
	AccountBankCode   string `json:"account_bank_code"`
	AccountNo         string `json:"account_no"`
	AccountHolderName string `json:"account_holder_name"`
}

type AmountObj struct {
	Total    int64  `json:"total"`
	Currency string `json:"currency"`
}

type Bank1CreateDisbursementResponse struct {
	Status  string                              `json:"status"`
	Message string                              `json:"message"`
	Data    Bank1CreateDisbursementResponseData `json:"data"`
}

type Bank1CreateDisbursementResponseData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Bank1CreateDisbursementRequest
	CreatedAt time.Time `json:"created_at"`
}

func (c *Bank1Client) CreateDisbursement(ctx context.Context, requestParam *Bank1CreateDisbursementRequest) (*Bank1CreateDisbursementResponse, error) {
	if requestParam.Amount.Currency == "" {
		requestParam.Amount.Currency = "IDR"
	}

	bodyBytes, _ := json.Marshal(requestParam)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", c.hostname, c.disbursementEndpoint), bytes.NewReader(bodyBytes))
	if err != nil {
		return &Bank1CreateDisbursementResponse{}, err
	}

	req.Header.Set("X-API-Key", c.apikey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return &Bank1CreateDisbursementResponse{}, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return &Bank1CreateDisbursementResponse{}, err
	}

	var resp *Bank1CreateDisbursementResponse
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
