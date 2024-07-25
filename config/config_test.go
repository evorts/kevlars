/**
 * @Author: steven
 * @Description:
 * @File: config
 * @Date: 21/07/24 06.17
 */

package config

import (
	mockConfig "github.com/evorts/kevlars/mocks/config"
	"github.com/stretchr/testify/suite"
	"testing"
)

type configTestSuite struct {
	suite.Suite

	lfp *mockConfig.Provider // local file provider
	rgp *mockConfig.Provider // gsm provider
	rcp *mockConfig.Provider // consul provider
	rdp *mockConfig.Provider // database provider
}

func (ts *configTestSuite) SetupSuite() {
	ts.lfp = mockConfig.NewProvider(ts.T())
	ts.rgp = mockConfig.NewProvider(ts.T())
	ts.rcp = mockConfig.NewProvider(ts.T())
	ts.rdp = mockConfig.NewProvider(ts.T())
}

func (ts *configTestSuite) TearDownTest() {
	//TODO implement me
}

func (ts *configTestSuite) TearDownSubTest() {
	//TODO implement me
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}
