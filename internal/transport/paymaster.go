package transport

import (
	"errors"

	"go-ticketos/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type PayMasterController[context any] interface {
	WebHook(c *context) error
}

type payMasterControllerFiber struct {
	useCase domain.PayMasterWebHookUseCase
	logger  *zerolog.Logger
}

// nolint: revive
func NewPayMasterController(
	useCase domain.PayMasterWebHookUseCase,
	logger *zerolog.Logger,
) (*payMasterControllerFiber, error) {
	if useCase == nil {
		return nil, errors.New("NewPayMasterController: use case is nil")
	}
	if logger == nil {
		return nil, errors.New("NewPayMasterController: logger is nil")
	}
	return &payMasterControllerFiber{
		useCase: useCase,
		logger:  logger,
	}, nil
}

var _ PayMasterController[fiber.Ctx] = (*payMasterControllerFiber)(nil)

func (p payMasterControllerFiber) WebHook(c *fiber.Ctx) error {
	var requestDTO PayMasterWebHookRequestDTO
	if err := c.BodyParser(&requestDTO); err != nil {
		p.logger.Error().Err(err).Msg("Can not read request dto")
		return c.Status(fiber.StatusBadRequest).JSON(nil)
	}
	validate := validator.New()
	if err := validate.Struct(requestDTO); err != nil {
		p.logger.Error().Err(err).Any("requestDTO", requestDTO).Msg("Can not validate requestDTO")
		return c.Status(fiber.StatusBadRequest).JSON(nil)
	}
	p.logger.Info().Any("requestDTO", requestDTO).Msg("New webhook from pay master")

	props := domain.PayMasterWebHookUseCaseProps{
		ID:         requestDTO.ID,
		Created:    requestDTO.Created,
		TestMode:   requestDTO.TestMode,
		Status:     requestDTO.Status,
		MerchantID: requestDTO.MerchantID,
		Amount: domain.Amount{
			Value:    requestDTO.Amount.Value,
			Currency: requestDTO.Amount.Currency,
		},
		Invoice: domain.Invoice{
			Description: requestDTO.Invoice.Description,
			OrderNo:     requestDTO.Invoice.OrderNo,
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          requestDTO.PaymentData.PaymentMethod,
			PaymentInstrumentTitle: requestDTO.PaymentData.PaymentInstrumentTitle,
		},
	}
	if err := p.useCase.Execute(props); err != nil {
		p.logger.Error().Err(err).Msg("Can not process webhook ")
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}
	p.logger.Info().Any("requestDTO", requestDTO).Msg("Successful webhook processing")
	return c.Status(fiber.StatusOK).JSON(nil)
}
