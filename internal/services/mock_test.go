package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MockServiceSuite struct {
	suite.Suite
}

func TestMockSuite(t *testing.T) {
	suite.Run(t, new(MockServiceSuite))
}

func (suite *MockServiceSuite) TestMock_OpenMockFile() {
	filename := "/opt/app/api/internal/mocks/users.json"
	suite.T().Logf("filepath %+v\n", filename)
	mockService := NewMockService()

	suite.Run("test can open mock file", func() {
		data, err := mockService.OpenMockFile(filename)

		suite.Require().Nil(err)
		suite.Require().NotNil(data)

		var payload map[string][]interface{}
		err = json.Unmarshal(data, &payload)
		suite.Require().Nil(err)

		suite.Require().True(len(payload["users"]) == 16)
	})
}
