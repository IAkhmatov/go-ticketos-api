package transport

import (
	"errors"

	"go-ticketos/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type OrderController[context any] interface {
	Create(c *context) error
}

type orderControllerFiber struct {
	svc     domain.OrderService
	adapter OrderAdapter
	logger  *zerolog.Logger
}

// nolint: revive
func NewOrderController(
	svc domain.OrderService,
	logger *zerolog.Logger,
) (*orderControllerFiber, error) {
	if svc == nil {
		return nil, errors.New("NewOrderController: svc is nil")
	}
	if logger == nil {
		return nil, errors.New("NewOrderController: logger is nil")
	}
	return &orderControllerFiber{
		svc:     svc,
		adapter: NewOrderAdapter(),
		logger:  logger,
	}, nil
}

var _ OrderController[fiber.Ctx] = (*orderControllerFiber)(nil)

func (f orderControllerFiber) Create(c *fiber.Ctx) error {
	var createDTO CreateOrderRequestDTO
	if err := c.BodyParser(&createDTO); err != nil {
		f.logger.Error().Err(err).Msg("Can not read createDTO")
		return c.Status(fiber.StatusBadRequest).JSON(nil)
	}
	validate := validator.New()
	if err := validate.Struct(createDTO); err != nil {
		f.logger.Error().Err(err).Any("createDTO", createDTO).Msg("Can not validate createDTO")
		return c.Status(fiber.StatusBadRequest).JSON(nil)
	}
	props := domain.CreateOrderProps{
		Name:  createDTO.Name,
		Email: createDTO.Email,
		Phone: createDTO.Phone,
	}
	if createDTO.PromocodeID != nil {
		p, err := uuid.Parse(*createDTO.PromocodeID)
		if err != nil {
			f.logger.Error().Err(err).Any("createDTO", createDTO).Msg("Can not parse promocode id")
			return c.Status(fiber.StatusBadRequest).JSON(nil)
		}
		props.PromocodeID = &p
	}
	for _, tcID := range createDTO.TicketCategoryIDs {
		tc, err := uuid.Parse(tcID)
		if err != nil {
			f.logger.Error().Err(err).Any("createDTO", createDTO).Msg("Can not parse ticket category id")
			return c.Status(fiber.StatusBadRequest).JSON(nil)
		}
		props.TicketCategoryIDs = append(props.TicketCategoryIDs, tc)
	}
	o, err := f.svc.Create(props)
	if err != nil {
		f.logger.Error().Err(err).Any("createDTO", createDTO).Msg("Can not create order")
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}
	orderDTO := f.adapter.ToResponseDTO(*o)
	return c.Status(fiber.StatusCreated).JSON(orderDTO)
}
