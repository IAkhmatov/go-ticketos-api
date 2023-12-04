//go:generate go-enum --marshal --mustparse

package paymasterclient

import (
	"time"
)

const payMasterAPIURL string = "https://paymaster.ru/api"

type (
	// ENUM(BankCard).
	PaymentMethod string
	// ENUM(RUB).
	Currency string

	Invoice struct {
		Description string    `json:"description" validate:"required"`
		OrderNo     string    `json:"orderNo" validate:"required"`
		Expires     time.Time `json:"expires" validate:"required,gte"`
	}

	Amount struct {
		Value    float32  `json:"value" validate:"required,gte=1"`
		Currency Currency `json:"currency" validate:"required,oneof=RUB"`
	}

	Protocol struct {
		ReturnURL   string `json:"returnURL" validate:"required,url"`
		CallbackURL string `json:"callbackURL" validate:"required,url"`
	}

	Customer struct {
		Email string `json:"email" validate:"required,email"`
		Phone string `json:"phone" validate:"required,gte=11"`
	}

	CreateInvoiceRequestDTO struct {
		MerchantID    string        `json:"merchantId" validate:"required"`
		TestMode      bool          `json:"testMode" validate:"required"`
		Invoice       Invoice       `json:"invoice" validate:"required"`
		Amount        Amount        `json:"amount" validate:"required"`
		PaymentMethod PaymentMethod `json:"paymentMethod" validate:"required,oneof=BankCard"`
		Protocol      Protocol      `json:"protocol" validate:"required"`
		Customer      Customer      `json:"customer" validate:"required"`
	}

	CreateInvoiceResponseDTO struct {
		PaymentID string `json:"paymentId"`
		URL       string `json:"url"`
	}
)
