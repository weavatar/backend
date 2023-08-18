package id

import (
	"testing"

	"github.com/goravel/framework/testing/mock"
	"github.com/stretchr/testify/suite"
)

type IDTestSuite struct {
	suite.Suite
	id RatID
}

func TestIDTestSuite(t *testing.T) {
	mockConfig := mock.Config()
	mockConfig.On("GetInt", "id.node").Return(0).Once()

	suite.Run(t, &IDTestSuite{
		id: *NewRatID(),
	})

	mockConfig.AssertExpectations(t)
}

func (s *IDTestSuite) TestGenerate() {
	id, err := s.id.Generate()

	s.NoError(err)
	s.NotZero(id)
}
