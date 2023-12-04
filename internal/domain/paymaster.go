package domain

import "time"

type (
	Amount struct {
		Value    float32
		Currency string
	}

	Invoice struct {
		Description string
		OrderNo     string
	}

	PaymentData struct {
		PaymentMethod          string
		PaymentInstrumentTitle string
	}

	PayMasterWebHookUseCaseProps struct {
		ID          string
		Created     time.Time
		TestMode    bool
		Status      string
		MerchantID  string
		Amount      Amount
		Invoice     Invoice
		PaymentData PaymentData
	}

	PayMasterWebHookUseCase interface {
		Execute(props PayMasterWebHookUseCaseProps) error
	}
)
