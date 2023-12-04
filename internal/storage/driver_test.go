package storage_test

import (
	"testing"

	"go-ticketos/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type toolsTestSuite struct {
	suite.Suite
	a *assert.Assertions
}

func TestToolsTestSuite(t *testing.T) {
	suite.Run(t, &toolsTestSuite{})
}

func (s *toolsTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
}

func (s *toolsTestSuite) TestToNamedInsert() {
	stmt := "INSERT INTO test %s VALUES %s"
	rows := []string{
		"test.id",
		"test.name",
		"test.description",
	}

	stmt = storage.ToNamedInsert(stmt, rows)

	s.a.Equal("INSERT INTO test (id, name, description) VALUES (:test.id, :test.name, :test.description)", stmt)
}

func (s *toolsTestSuite) TestToNamedUpdate() {
	stmt := "UPDATE test SET %s WHERE 1=1"
	rows := []string{
		"test.id",
		"test.name",
		"test.description",
	}

	stmt = storage.ToNamedUpdate(stmt, rows)

	s.a.Equal("UPDATE test SET id=:test.id, name=:test.name, description=:test.description WHERE 1=1", stmt)
}

func (s *toolsTestSuite) TestToAlias() {
	stmt := "SELECT %s FROM test"
	rows := []string{
		"id",
		"name",
		"description",
	}

	stmt = storage.ToSelect(stmt, rows)

	s.a.Equal("SELECT id \"id\", name \"name\", description \"description\" FROM test", stmt)
}
