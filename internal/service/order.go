package service

import (
	"errors"
	"fmt"
	"time"

	"go-ticketos/internal/config"
	"go-ticketos/internal/domain"
	"go-ticketos/pkg/paymasterclient"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type orderService struct {
	orderRepo          domain.OrderRepo
	ticketCategoryRepo domain.TicketCategoryRepo
	promocodeRepo      domain.PromocodeRepo
	payMasterClient    paymasterclient.PayMasterClient
	cfg                *config.Config
	logger             *zerolog.Logger
}

// nolint: revive
func NewOrderService(
	orderRepo domain.OrderRepo,
	ticketCategoryRepo domain.TicketCategoryRepo,
	promocodeRepo domain.PromocodeRepo,
	payMasterClient paymasterclient.PayMasterClient,
	cfg *config.Config,
	logger *zerolog.Logger,
) (*orderService, error) {
	if orderRepo == nil {
		return nil, errors.New("NewEventService: orderRepo in nil")
	}
	if ticketCategoryRepo == nil {
		return nil, errors.New("NewEventService: ticketCategoryRepo in nil")
	}
	if promocodeRepo == nil {
		return nil, errors.New("NewEventService: promocodeRepo in nil")
	}
	if payMasterClient == nil {
		return nil, errors.New("NewOrderService: payMasterClient is nil")
	}
	if cfg == nil {
		return nil, errors.New("NewOrderService: config is nil")
	}
	if logger == nil {
		return nil, errors.New("NewOrderService: logger is nil")
	}
	return &orderService{
		orderRepo:          orderRepo,
		ticketCategoryRepo: ticketCategoryRepo,
		promocodeRepo:      promocodeRepo,
		payMasterClient:    payMasterClient,
		cfg:                cfg,
		logger:             logger,
	}, nil
}

var _ domain.OrderService = (*orderService)(nil)

func (o orderService) Create(props domain.CreateOrderProps) (*domain.Order, error) {
	o.logger.Debug().Any("props", props).Msg("Order update started")
	var tickets []domain.Ticket
	var promocode domain.Promocode
	// find promocode
	if props.PromocodeID != nil {
		promocodeNew, err := o.promocodeRepo.GetByID(*props.PromocodeID)
		if err != nil {
			return nil, fmt.Errorf("Create: can not get promocode: %w", err)
		}
		promocode = *promocodeNew
	}

	// create tickets
	// promocode can append only on tickets category that can uses this promocode
	for _, tid := range props.TicketCategoryIDs {
		tc, err := o.ticketCategoryRepo.GetByID(tid)
		if err != nil {
			return nil, fmt.Errorf("Create: can not get ticket category: %w", err)
		}
		var promocodeForTicket *domain.Promocode
		if props.PromocodeID != nil {
			if tc.CanUsePromocode(promocode) {
				promocodeForTicket = &promocode
			}
		}
		ticket, err := domain.NewTicket(*tc, promocodeForTicket)
		if err != nil {
			return nil, fmt.Errorf("Create: can not create ticket: %w", err)
		}
		tickets = append(tickets, *ticket)
	}

	// create order
	order, err := domain.NewOrder(
		props.Name,
		props.Email,
		props.Phone,
		tickets,
	)
	if err != nil {
		return nil, fmt.Errorf("Create: can not create order: %w", err)
	}
	if err = o.orderRepo.Create(*order); err != nil {
		return nil, fmt.Errorf("Create: can not save order: %w", err)
	}

	// create invoice in payment system
	createInvoiceProps := paymasterclient.CreateInvoiceRequestDTO{
		MerchantID: o.cfg.PayMasterMerchantID,
		TestMode:   true,
		Invoice: paymasterclient.Invoice{
			Description: "Tickets",
			OrderNo:     order.ID.String(),
			Expires:     time.Now().UTC().Add(time.Duration(o.cfg.OrderTTL) * time.Minute),
		},
		Amount: paymasterclient.Amount{
			Value:    float32(order.BuyPrice()),
			Currency: paymasterclient.CurrencyRUB,
		},
		PaymentMethod: paymasterclient.PaymentMethodBankCard,
		Protocol: paymasterclient.Protocol{
			ReturnURL:   o.cfg.PayMasterReturnURL,
			CallbackURL: o.cfg.PayMasterCallbackURL,
		},
		Customer: paymasterclient.Customer{
			Email: order.Email,
			Phone: order.Phone,
		},
	}
	resp, err := o.payMasterClient.CreateInvoice(createInvoiceProps)
	if err != nil {
		return nil, fmt.Errorf("Create: can not create payment in payment system: %w", err)
	}
	o.logger.Info().Str("payment url", resp.URL).Str("payment id", resp.PaymentID).Msg("Created payment in pp system")

	// update order
	updateProps := domain.UpdateOrderProps{
		OrderID:    order.ID,
		Status:     domain.OrderStatusAwaitingPayment,
		PaymentID:  &resp.PaymentID,
		PaymentURL: &resp.URL,
	}
	order, err = o.Update(updateProps)
	if err != nil {
		return nil, fmt.Errorf("Create: can not udpate order: %w", err)
	}
	return order, nil
}

func (o orderService) Update(props domain.UpdateOrderProps) (*domain.Order, error) {
	o.logger.Debug().Any("props", props).Msg("Order update started")
	val := validator.New()
	if err := val.Struct(props); err != nil {
		return nil, fmt.Errorf("Update: can not validate props: %w", err)
	}
	order, err := o.orderRepo.GetByID(props.OrderID)
	if err != nil {
		return nil, fmt.Errorf("Update: can not get order: %w", err)
	}
	// nolint: exhaustive
	switch props.Status {
	case domain.OrderStatusAwaitingPayment:
		if err = order.UpdateStatus(props.Status); err != nil {
			return nil, fmt.Errorf("Update: can not update order to this status: %w", err)
		}
		payment := domain.Payment{
			ID:  *props.PaymentID,
			URL: *props.PaymentURL,
		}
		order.UpdatePayment(payment)
	case domain.OrderStatusCompleted:
		if err = order.UpdateStatus(props.Status); err != nil {
			return nil, fmt.Errorf("Update: can not update order to this status: %w", err)
		}
	default:
		return nil, errors.New("Update: can not update order to this status")
	}
	if err = o.orderRepo.Update(*order); err != nil {
		return nil, fmt.Errorf("Update: can not update order in db: %w", err)
	}
	o.logger.Debug().Any("order", order).Msg("Successful order update")
	return order, nil
}

func (o orderService) GetByID(id uuid.UUID) (*domain.Order, error) {
	order, err := o.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetByID: can not get order: %w", err)
	}
	return order, nil
}
