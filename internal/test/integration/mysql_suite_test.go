//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MysqlRepositoryTestSuite struct {
	suite.Suite
}

func TestPetMysqlRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlRepositoryTestSuite))
}

func (p *MysqlRepositoryTestSuite) TearDownSuite() {
	//TODO implement me
	panic("implement me")
}

func (p *MysqlRepositoryTestSuite) SetupSuite() {
	//TODO implement me
	panic("implement me")
}

func (p *MysqlRepositoryTestSuite) SetupTest() {
	//TODO implement me
	panic("implement me")
}
