package database

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type PostgresTestingServiceSuite struct {
	suite.Suite
	postgresTestingService *PostgresTestingService
}

func TestPostgresTestingServiceSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestingServiceSuite))
}

func (suite *PostgresTestingServiceSuite) SetupSuite() {
	suite.postgresTestingService = NewPostgresTestingService()
}

func (suite *PostgresTestingServiceSuite) TearDownSuite() {
	suite.T().Log("tearing down postgres testing service")
	if err := suite.postgresTestingService.Close(); err != nil {
		suite.FailNow("failed to close postgres testing service", "error", err.Error())
	}
}

func (suite *PostgresTestingServiceSuite) TestNewPostgresService() {
	want := &PostgresTestingService{}

	suite.Run("test new postgres testing service", func() {
		suite.Require().NotNil(suite.postgresTestingService)
		suite.Require().Equal(reflect.TypeOf(want), reflect.TypeOf(suite.postgresTestingService))
	})
}

func (suite *PostgresTestingServiceSuite) TestPostgresService_Health() {
	want := map[string]string{
		"message": "testing db is healthy",
	}

	suite.Run("test postgres testing service is healthy", func() {
		suite.Require().Equal(want, suite.postgresTestingService.Health())
	})
}

func (suite *PostgresTestingServiceSuite) TestPostgresService_Instance() {
	want := reflect.TypeOf(&sql.DB{})

	suite.Run("test postgres service has a db instance", func() {
		suite.Require().NotNil(suite.postgresTestingService.Instance())
		suite.Require().Equal(want, reflect.TypeOf(suite.postgresTestingService.Instance()))
	})
}
