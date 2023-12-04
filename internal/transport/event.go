package transport

import (
	"errors"

	"go-ticketos/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type EventController[context any] interface {
	Create(c *context) error
}

type eventControllerFiber struct {
	svc     domain.EventService
	adapter EventAdapter
	logger  *zerolog.Logger
}

var _ EventController[fiber.Ctx] = (*eventControllerFiber)(nil)

// nolint: revive
func NewEventController(
	svc domain.EventService,
	logger *zerolog.Logger,
) (*eventControllerFiber, error) {
	if svc == nil {
		return nil, errors.New("NewEventController: svc is nil")
	}
	if logger == nil {
		return nil, errors.New("NewEventController: logger is nil")
	}
	return &eventControllerFiber{
		svc:     svc,
		adapter: NewEventAdapter(),
		logger:  logger,
	}, nil
}

// Create controller for creating events.
func (ec eventControllerFiber) Create(c *fiber.Ctx) error {
	var createDTO CreateEventRequestDTO
	if err := c.BodyParser(&createDTO); err != nil {
		ec.logger.Error().Err(err).Msg("Can not read createDTO")
		return c.Status(fiber.StatusBadRequest).JSON(nil)
	}
	validate := validator.New()
	if err := validate.Struct(createDTO); err != nil {
		ec.logger.Error().Err(err).Any("createDTO", createDTO).Msg("Can not validate createDTO")
		return c.Status(fiber.StatusBadRequest).JSON(nil)
	}
	props := domain.CreateEventProps{
		Name:        createDTO.Name,
		Description: createDTO.Description,
		Place:       createDTO.Place,
		AgeRating:   createDTO.AgeRating,
		StartAt:     createDTO.StartAt,
		EndAt:       createDTO.EndAt,
	}
	event, err := ec.svc.Create(props)
	if err != nil {
		ec.logger.Error().Err(err).Any("createDTO", createDTO).Msg("Can not create event")
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}
	eventDto := ec.adapter.ToResponseDTO(*event)
	return c.Status(fiber.StatusCreated).JSON(eventDto)
}
