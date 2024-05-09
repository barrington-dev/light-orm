package database

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type DatabaseTestSuite struct {
	suite.Suite
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

type SomeDatabaseService struct {
}

func NewSomeDatabaseService() *SomeDatabaseService {
	return &SomeDatabaseService{}
}

func (s *SomeDatabaseService) Close() error {
	return nil
}

func (s *SomeDatabaseService) Instance() *sql.DB {
	return nil
}

func (s *SomeDatabaseService) Health() map[string]string {
	return map[string]string{
		"message": "testing ok",
	}
}

func (suite *DatabaseTestSuite) TestDatabaseServiceInterface() {
	someDatabaseService := NewSomeDatabaseService()
	i := reflect.TypeOf((*Service)(nil)).Elem()
	suite.True(reflect.TypeOf(someDatabaseService).Implements(i))
}
