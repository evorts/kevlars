/**
 * @Author: steven
 * @Description:
 * @File: log_test
 * @Date: 18/12/23 07.44
 */

package audit

import (
	"github.com/evorts/kevlars/db"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestSuite struct {
	suite.Suite
	db db.Manager
}

func (ts *TestSuite) SetupTest() {
	ts.db = db.NewWithMock()
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
