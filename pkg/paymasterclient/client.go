package paymasterclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-ticketos/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type (
	PayMasterClient interface {
		CreateInvoice(props CreateInvoiceRequestDTO) (*CreateInvoiceResponseDTO, error)
	}

	payMasterClient struct {
		cfg    *config.Config
		client *http.Client
		logger *zerolog.Logger
	}
)

// nolint: revive
func NewPayMasterClient(
	cfg *config.Config,
	logger *zerolog.Logger,
) (*payMasterClient, error) {
	if cfg == nil {
		return nil, errors.New("NewPayMasterClient: cfg is nil")
	}
	if logger == nil {
		return nil, errors.New("NewPayMasterClient: logger is nil")
	}
	return &payMasterClient{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.PayMasterTimeout) * time.Second,
		},
		logger: logger,
	}, nil
}

var _ PayMasterClient = (*payMasterClient)(nil)

func (p payMasterClient) CreateInvoice(props CreateInvoiceRequestDTO) (*CreateInvoiceResponseDTO, error) {
	val := validator.New()
	err := val.Struct(props)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice: can not validate input dto: %w", err)
	}
	b, err := json.Marshal(props)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice: can not marshal body: %w", err)
	}
	readerBody := bytes.NewReader(b)
	u := fmt.Sprintf("%s/v2/invoices", payMasterAPIURL)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, u, readerBody)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice: can not build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.cfg.PayMasterAPIKey))
	req.Header.Set("Idempotency-Key", props.Invoice.OrderNo)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice: can not send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("PayMaster response status code = %d", resp.StatusCode)
		bb, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return nil, fmt.Errorf("CreateInvoice: can decode response: %w while: %w", err2, err)
		}
		p.logger.Err(err).Str("response", string(bb)).Int("status code", resp.StatusCode).Msg("problem with pay master")
		return nil, err
	}
	var dto CreateInvoiceResponseDTO
	if err = json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("CreateInvoice: can not decode response: %w", err)
	}
	return &dto, nil
}
