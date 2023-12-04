package transport_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"go-ticketos/internal/config"
	"go-ticketos/internal/domain"
	"go-ticketos/internal/service"
	"go-ticketos/internal/storage"
	"go-ticketos/internal/transport"
	"go-ticketos/pkg/log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type eventControllerTestSuite struct {
	suite.Suite
	a          *assert.Assertions
	db         *sqlx.DB
	svc        domain.EventService
	controller transport.EventController[fiber.Ctx]
	app        *fiber.App
}

func TestEventControllerTestSuite(t *testing.T) {
	suite.Run(t, &eventControllerTestSuite{})
}

func (s *eventControllerTestSuite) SetupSuite() {
	s.a = assert.New(s.T())

	cfg, err := config.NewConfig()
	s.a.NoError(err)
	db, err := storage.NewSqlxDB(cfg.TestDBConnectString)
	s.a.NoError(err)
	err = storage.RecreateSchema(db)
	s.a.NoError(err)
	err = storage.CreateTables(db)
	s.a.NoError(err)
	s.db = db

	repo, err := storage.NewEventRepo(db)
	s.a.NoError(err)

	svc, err := service.NewEventService(repo)
	s.a.NoError(err)
	s.svc = svc

	controller, err := transport.NewEventController(s.svc, log.NewLogger(cfg))
	s.a.NoError(err)
	s.controller = controller

	s.app = fiber.New()
	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/event", s.controller.Create)
}

func (s *eventControllerTestSuite) TearDownSuite() {
	err := storage.RecreateSchema(s.db)
	s.a.NoError(err)
}

func (s *eventControllerTestSuite) TearDownTest() {
	err := storage.CleanAllTables(s.db)
	s.a.NoError(err)
}

func (s *eventControllerTestSuite) TestCreate() {
	strPointer := func(s string) *string {
		return &s
	}
	tests := []struct {
		name         string
		method       string
		url          string
		body         io.Reader
		expectedCode int
	}{
		{
			name:   "incorrect age rating",
			method: "POST",
			url:    "/api/v1/event",
			body: s.generateBodyString(transport.CreateEventRequestDTO{
				Name:      "test1",
				Place:     "test2",
				AgeRating: 222,
				StartAt:   time.Now().UTC(),
				EndAt:     time.Now().UTC().Add(1 * time.Hour),
			}),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "incorrect description",
			method: "POST",
			url:    "/api/v1/event",
			body: s.generateBodyString(transport.CreateEventRequestDTO{
				Name:        "test1",
				Place:       "test2",
				Description: strPointer(""),
				AgeRating:   4,
				StartAt:     time.Now().UTC(),
				EndAt:       time.Now().UTC().Add(1 * time.Hour),
			}),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "incorrect end at",
			method: "POST",
			url:    "/api/v1/event",
			body: s.generateBodyString(transport.CreateEventRequestDTO{
				Name:        "test1",
				Place:       "test2",
				Description: strPointer("1"),
				AgeRating:   4,
				StartAt:     time.Now().UTC(),
				EndAt:       time.Now().UTC().Add(-1 * time.Hour),
			}),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "correct with description",
			method: "POST",
			url:    "/api/v1/event",
			body: s.generateBodyString(transport.CreateEventRequestDTO{
				Name:        "test1",
				Place:       "test2",
				Description: strPointer("1"),
				AgeRating:   4,
				StartAt:     time.Now().UTC(),
				EndAt:       time.Now().UTC().Add(1 * time.Hour),
			}),
			expectedCode: http.StatusCreated,
		},
		{
			name:   "correct without description",
			method: "POST",
			url:    "/api/v1/event",
			body: s.generateBodyString(transport.CreateEventRequestDTO{
				Name:      "test1",
				Place:     "test2",
				AgeRating: 4,
				StartAt:   time.Now().UTC(),
				EndAt:     time.Now().UTC().Add(1 * time.Hour),
			}),
			expectedCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			s.a.NoError(err)
			defer req.Body.Close()
			req.Header.Set("Content-Type", "application/json")
			resp, err := s.app.Test(req, 10000)
			s.a.NoError(err)
			s.Equalf(tt.expectedCode, resp.StatusCode, tt.name)
			s.TearDownTest()
		})
	}
}

func (s *eventControllerTestSuite) generateBodyString(createEventRequestDTO transport.CreateEventRequestDTO) io.Reader {
	s.T().Helper()
	b, err := json.Marshal(createEventRequestDTO)
	s.a.NoError(err)
	return bytes.NewReader(b)
}
