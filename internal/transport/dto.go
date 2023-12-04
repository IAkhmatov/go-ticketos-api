package transport

import (
	"time"

	"github.com/google/uuid"
)

type (
	CreateEventRequestDTO struct {
		Name        string    `json:"name" validate:"required"`
		Description *string   `json:"description" validate:"omitempty,min=1"`
		Place       string    `json:"place" validate:"required"`
		AgeRating   int       `json:"ageRating" validate:"required,min=0,max=21"`
		StartAt     time.Time `json:"startAt" validate:"required" time_format:"2006-01-02"`
		EndAt       time.Time `json:"endAt" validate:"required,gtefield=StartAt" time_format:"2006-01-02"`
	}

	EventResponseDTO struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description *string   `json:"description,omitempty"`
		Place       string    `json:"place"`
		AgeRating   int       `json:"ageRating"`
		StartAt     time.Time `json:"startAt"`
		EndAt       time.Time `json:"endAt"`
	}

	CreateOrderRequestDTO struct {
		Name              string   `json:"name" validate:"required,gte=3,lte=150"`
		Email             string   `json:"email" validate:"required,email"`
		Phone             string   `json:"phone" validate:"required,numeric"`
		PromocodeID       *string  `json:"promocodeID" validate:"omitempty,uuid"`
		TicketCategoryIDs []string `json:"ticketCategoryIDs" validate:"required,gte=1,dive,uuid"`
	}

	OrderResponseDTO struct {
		ID        uuid.UUID           `json:"id"`
		Name      string              `json:"name"`
		Email     string              `json:"email"`
		Phone     string              `json:"phone"`
		Payment   *PaymentResponseDTO `json:"payment"`
		Tickets   []TicketResponseDTO `json:"tickets"`
		Status    string              `json:"status"`
		FullPrice uint                `json:"fullPrice"`
		BuyPrice  uint                `json:"buyPrice"`
	}

	PaymentResponseDTO struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}

	TicketResponseDTO struct {
		ID             uuid.UUID                 `json:"id"`
		TicketCategory TicketCategoryResponseDTO `json:"ticketCategory"`
		PromocodeID    *uuid.UUID                `json:"promocodeID"`
		FullPrice      uint                      `json:"fullPrice"`
		BuyPrice       uint                      `json:"buyPrice"`
	}

	TicketCategoryResponseDTO struct {
		ID          uuid.UUID `json:"id"`
		EventID     uuid.UUID `json:"eventID"`
		Price       uint      `json:"price"`
		Name        string    `json:"name"`
		Description *string   `json:"description"`
	}

	Amount struct {
		Value    float32 `json:"value" validate:"required,gte=0"`
		Currency string  `json:"currency" validate:"required"`
	}

	Invoice struct {
		Description string `json:"description"`
		OrderNo     string `json:"orderNo" validate:"required"`
	}

	PaymentData struct {
		PaymentMethod          string `json:"paymentMethod" validate:"required"`
		PaymentInstrumentTitle string `json:"paymentInstrumentTitle" validate:"required"`
	}

	PayMasterWebHookRequestDTO struct {
		ID          string      `json:"id" validate:"required"`
		Created     time.Time   `json:"created" validate:"required"`
		TestMode    bool        `json:"testMode" validate:"required"`
		Status      string      `json:"status" validate:"required"`
		MerchantID  string      `json:"merchantId" validate:"required"`
		Amount      Amount      `json:"amount" validate:"required"`
		Invoice     Invoice     `json:"invoice" validate:"required"`
		PaymentData PaymentData `json:"paymentData" validate:"required"`
	}
)
