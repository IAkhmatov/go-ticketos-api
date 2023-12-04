package service

import (
	"errors"
	"fmt"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type payMasterWebHookUseCase struct {
	orderService domain.OrderService
	logger       *zerolog.Logger
}

// nolint: revive
func NewPayMasterWebHookUseCase(
	orderService domain.OrderService,
	logger *zerolog.Logger,
) (*payMasterWebHookUseCase, error) {
	if orderService == nil {
		return nil, errors.New("NewPayMasterWebHookUseCase: orderService is nil")
	}
	if logger == nil {
		return nil, errors.New("NewPayMasterWebHookUseCase: logger is nil")
	}
	return &payMasterWebHookUseCase{
		orderService: orderService,
		logger:       logger,
	}, nil
}

var _ domain.PayMasterWebHookUseCase = (*payMasterWebHookUseCase)(nil)

func (p payMasterWebHookUseCase) Execute(props domain.PayMasterWebHookUseCaseProps) error {
	p.logger.Info().Any("props", props).Msg("PayMaster webhook handler started")
	orderID, err := uuid.Parse(props.Invoice.OrderNo)
	if err != nil {
		return fmt.Errorf("Execute: can not parse orderNo: %w", err)
	}
	order, err := p.orderService.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("Execute: can not get order: %w", err)
	}
	const ratio = 100
	if props.Amount.Value != float32(order.BuyPrice()*ratio) {
		return errors.New("Execute: incorrect amount value")
	}
	if props.ID != order.Payment.ID {
		return errors.New("Execute: incorrect payment id")
	}
	if props.Status != "Settled" {
		return errors.New("Execute: incorrect payment provider status")
	}
	updateProps := domain.UpdateOrderProps{
		OrderID: orderID,
		Status:  domain.OrderStatusCompleted,
	}
	if _, err = p.orderService.Update(updateProps); err != nil {
		return fmt.Errorf("Execute: can not update order: %w", err)
	}
	p.logger.Info().Any("order", order).Msg("Successful webhook executed")
	return nil
}
