package transport_test

import (
	"testing"

	"go-ticketos/internal/domain"
	"go-ticketos/internal/transport"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ticketCategoryAdapterTestSuite struct {
	suite.Suite
	a       *assert.Assertions
	adapter transport.TicketCategoryAdapter
}

func TestTicketCategoryAdapterTestSuite(t *testing.T) {
	suite.Run(t, &ticketCategoryAdapterTestSuite{})
}

func (s *ticketCategoryAdapterTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
	s.adapter = transport.NewTicketCategoryAdapter()
}

func (s *ticketCategoryAdapterTestSuite) TestToResponseDTO() {
	d := "tt"
	tc := domain.NewTicketCategory(uuid.New(), 10000, "test", &d)

	tcDto := s.adapter.ToResponseDTO(tc)

	s.a.Equal(tc.ID, tcDto.ID)
	s.a.Equal(tc.EventID, tcDto.EventID)
	s.a.Equal(tc.Description, tcDto.Description)
	s.a.Equal(tc.Price, tcDto.Price)
	s.a.Equal(tc.Name, tcDto.Name)
}
